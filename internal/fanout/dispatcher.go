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
	mu      sync.Mutex
	pool    *envelopepool.Pool
	queue   *routequeue.Queue
	acks    *ackstate.Registry
	lastErr map[string]error
}

func New() *Dispatcher {
	return &Dispatcher{
		pool: envelopepool.New(), queue: &routequeue.Queue{},
		acks: ackstate.New(), lastErr: make(map[string]error),
	}
}

func (d *Dispatcher) Dispatch(ctx context.Context, id, route, payload string) (string, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.lastErr[route]; err != nil {
		return "", err
	}
	if !d.acks.CanSend(route, id) {
		return "", fmt.Errorf("route %s is blocked", route)
	}
	env := d.pool.Acquire(id, route, payload)
	d.queue.Push(env)
	d.pool.Release(env)
	item := d.queue.Pop()
	if err := ctx.Err(); err != nil {
		d.acks.MarkFailed(route, id)
		d.lastErr[route] = err
		return "", err
	}
	return "delivered:" + item.ID + ":" + item.Payload, nil
}
