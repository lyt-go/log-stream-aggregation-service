package decoderregistry

import "sync"

type Decoder interface{ Decode(string) string }

type TextDecoder struct{ Prefix string }

func (d *TextDecoder) Decode(payload string) string { return d.Prefix + payload }

type Registry struct {
	mu      sync.Mutex
	entries map[string]*TextDecoder
}

func New() *Registry { return &Registry{entries: make(map[string]*TextDecoder)} }

func (r *Registry) RegisterNil(route string) {
	r.mu.Lock()
	r.entries[route] = nil
	r.mu.Unlock()
}

func (r *Registry) Register(route, prefix string) {
	r.mu.Lock()
	if _, exists := r.entries[route]; !exists {
		r.entries[route] = &TextDecoder{Prefix: prefix}
	}
	r.mu.Unlock()
}

func (r *Registry) Resolve(route string) Decoder {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.entries[route]
}
