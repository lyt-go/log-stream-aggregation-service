package routequeue

import (
	"sync"

	"logaggregation/internal/envelopepool"
)

type Queue struct {
	mu    sync.Mutex
	items []*envelopepool.Envelope
}

func (q *Queue) Push(env *envelopepool.Envelope) {
	q.mu.Lock()
	q.items = append(q.items, env)
	q.mu.Unlock()
}

// Pop 弹出队首元素，返回其值拷贝。
//
// 返回值拷贝是为了与 envelopepool 的复用语义解耦：调用方拿到的是
// 一份独立的快照，即便原始 envelope 随后被归还并复用，也不会影响
// 已经 Pop 出来的结果。
func (q *Queue) Pop() envelopepool.Envelope {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return envelopepool.Envelope{}
	}
	item := q.items[0]
	q.items = q.items[1:]
	return *item
}
