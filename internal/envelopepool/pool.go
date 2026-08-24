package envelopepool

import "sync"

type Envelope struct {
	ID      string
	Route   string
	Payload string
}

type Pool struct{ pool sync.Pool }

func New() *Pool { return &Pool{} }

func (p *Pool) Acquire(id, route, payload string) *Envelope {
	v := p.pool.Get()
	if v == nil {
		return &Envelope{ID: id, Route: route, Payload: payload}
	}
	env := v.(*Envelope)
	env.ID = id
	return env
}

func (p *Pool) Release(env *Envelope) { p.pool.Put(env) }
