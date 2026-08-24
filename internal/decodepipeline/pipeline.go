package decodepipeline

import (
	"fmt"
	"sync"

	"logaggregation/internal/decodercache"
	"logaggregation/internal/decoderregistry"
	"logaggregation/internal/panicguard"
)

type Pipeline struct {
	mu       sync.Mutex
	registry *decoderregistry.Registry
	cache    *decodercache.Cache
	guard    *panicguard.Guard
}

func New() *Pipeline {
	return &Pipeline{registry: decoderregistry.New(), cache: decodercache.New(), guard: panicguard.New()}
}

func (p *Pipeline) RegisterNil(route string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.registry.RegisterNil(route)
	p.cache.Invalidate(route)
}

func (p *Pipeline) Register(route, prefix string) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.registry.Register(route, prefix)
	// Drop any stale cached entry so the next Decode picks up the new decoder
	// (it may previously have been a typed-nil placeholder that resolved to nil).
	p.cache.Invalidate(route)
}

func (p *Pipeline) Decode(route, payload string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	decoder := p.cache.Get(route, p.registry.Resolve)
	value, err := p.guard.Run(route, func() string {
		if decoder == nil {
			panic("missing decoder")
		}
		return decoder.Decode(payload)
	})
	if err != nil {
		// The error (panic or otherwise) is scoped to this call only; a
		// subsequent call after the decoder is fixed must succeed.
		return "", fmt.Errorf("decode route: %w", err)
	}
	return value, nil
}
