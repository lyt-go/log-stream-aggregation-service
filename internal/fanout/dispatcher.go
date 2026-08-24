package fanout

import (
	"context"
	"fmt"
	"sync"

	"logaggregation/internal/ackstate"
	"logaggregation/internal/envelopepool"
	"logaggregation/internal/routequeue"
)

type Dispatcher struct {
	mu    sync.Mutex
	pool  *envelopepool.Pool
	queue *routequeue.Queue
	acks  *ackstate.Registry
}

func New() *Dispatcher {
	return &Dispatcher{
		pool: envelopepool.New(), queue: &routequeue.Queue{},
		acks: ackstate.New(),
	}
}

// Dispatch 把一条消息发往指定路由并返回送达回执。
//
// 状态隔离原则：每次调用都基于「本次」的上下文与入参独立判定结果，
// 不继承前一次（尤其是前一条消息、上一个 ctx）的失败状态或残留对象。
// 因此：
//   - 不再维护按路由长期保留的失败缓存（旧的 lastErr[route] 会让
//     后一条新消息继承前一条的失败），每条消息只看它自己的 ack 状态；
//   - 失败只登记到「本条」消息，绝不污染整条路由；
//   - 信封在确认送达之前不归还，避免被并发复用、把新 payload
//     写进还在队列里、尚未 Pop 出来的旧信封。
func (d *Dispatcher) Dispatch(ctx context.Context, id, route, payload string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()

	// 先校验本次上下文，避免对一条注定要失败的请求做入队、出队、
	// 登记失败等副作用——副作用一旦落到共享状态上，就会污染后续
	// 无关消息。
	if err := ctx.Err(); err != nil {
		// 只登记「本条」消息失败，绝不写入路由级别的失败标记。
		d.acks.MarkFailed(route, id)
		return "", err
	}

	// 检查本条消息是否可发送。失败记录按消息单次有效，新消息天然可发送。
	if !d.acks.CanSend(route, id) {
		return "", fmt.Errorf("message %s on route %s is blocked", id, route)
	}

	// 构造本次消息的信封。Acquire 已保证全量覆写字段，不会携带
	// 前一次复用对象的残留 Route/Payload。
	env := d.pool.Acquire(id, route, payload)
	d.queue.Push(env)

	// 在确认送达之前不归还对象：若先 Release，下一次 Acquire 可能
	// 就地复用这个仍在队列里的信封，把新 payload 覆写上去，导致
	// 本次 Pop 拿到的是别人的 payload。Pop 返回的是值拷贝，
	// 之后即便对象被复用也不影响已取出的结果。
	item := d.queue.Pop()
	defer d.pool.Release(env)

	// 送达前再做一次上下文检查；失败同样只记到本条消息。
	if err := ctx.Err(); err != nil {
		d.acks.MarkFailed(route, id)
		return "", err
	}

	return "delivered:" + item.ID + ":" + item.Payload, nil
}
