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
	lastErr  map[string]error
}

func New() *Pipeline {
	return &Pipeline{registry: decoderregistry.New(), cache: decodercache.New(), guard: panicguard.New(), lastErr: make(map[string]error)}
}

func (p *Pipeline) RegisterNil(route string) { p.registry.RegisterNil(route) }

func (p *Pipeline) Register(route, prefix string) { p.registry.Register(route, prefix) }

func (p *Pipeline) Decode(route, payload string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if err := p.lastErr[route]; err != nil {
		return "", err
	}
	decoder := p.cache.Get(route, p.registry.Resolve)
	value, err := p.guard.Run(route, func() string {
		if decoder == nil {
			panic("missing decoder")
		}
		return decoder.Decode(payload)
	})
	if err != nil {
		p.lastErr[route] = err
		return "", fmt.Errorf("decode route: %w", err)
	}
	return value, nil
}
