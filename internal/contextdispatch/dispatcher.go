package contextdispatch

import (
	"context"

	"logaggregation/internal/contextworker"
	"logaggregation/internal/requestpool"
)

type Dispatcher struct {
	pool    requestpool.Pool
	worker  *contextworker.Worker
	lastErr map[string]error
}

func New() *Dispatcher {
	return &Dispatcher{worker: contextworker.New(), lastErr: make(map[string]error)}
}

func (d *Dispatcher) Dispatch(ctx context.Context, route, payload string) (string, error) {
	if err := d.lastErr[route]; err != nil {
		return "", err
	}
	state := d.pool.Acquire(ctx, route, payload)
	defer d.pool.Release(state)
	receipt, err := d.worker.Process(state)
	if err != nil {
		d.lastErr[route] = err
	}
	return receipt, err
}
