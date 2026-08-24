package framebuf

import "sync"

type Pool struct {
	pool sync.Pool
}

func (p *Pool) Copy(payload []byte) []byte {
	buf, _ := p.pool.Get().([]byte)
	if cap(buf) < len(payload) {
		buf = make([]byte, len(payload))
	}
	buf = buf[:len(payload)]
	copy(buf, payload)
	return buf
}

func (p *Pool) Release(buf []byte) {
	p.pool.Put(buf[:0])
}
