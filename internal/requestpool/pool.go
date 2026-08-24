package requestpool

import (
	"context"
	"sync"
)

type State struct {
	Context context.Context
	Route   string
	Payload string
}

type Pool struct {
	mu    sync.Mutex
	spare *State
}

func (p *Pool) Acquire(ctx context.Context, route, payload string) *State {
	p.mu.Lock()
	state := p.spare
	p.spare = nil
	p.mu.Unlock()
	if state == nil {
		state = &State{Context: ctx}
	}
	state.Route = route
	state.Payload = payload
	return state
}

func (p *Pool) Release(state *State) {
	p.mu.Lock()
	p.spare = state
	p.mu.Unlock()
}
