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
	// 归还的对象可能携带前一次调度的残留字段，必须全量覆写，
	// 否则新消息会继承上一次的 Route/Payload。
	env.ID = id
	env.Route = route
	env.Payload = payload
	return env
}

// Release 把对象归还到池中。归还后调用方不得再持有或读写该对象，
// 因为下一次 Acquire 会原地复用并覆写它。
func (p *Pool) Release(env *Envelope) { p.pool.Put(env) }
