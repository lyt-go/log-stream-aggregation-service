package decodercache

import (
	"sync"

	"logaggregation/internal/decoderregistry"
)

type Cache struct {
	mu     sync.Mutex
	values map[string]decoderregistry.Decoder
}

func New() *Cache {
	return &Cache{values: make(map[string]decoderregistry.Decoder)}
}

// Get returns the cached decoder for route, resolving via resolve when not yet
// cached. A nil resolve result is NOT memoized as a permanent negative: if the
// decoder is registered later, the next Get re-resolves and picks it up.
// Negative-result caching would let a transiently-absent decoder (e.g. a
// typed-nil placeholder) poison every subsequent call for that route.
func (c *Cache) Get(route string, resolve func(string) decoderregistry.Decoder) decoderregistry.Decoder {
	c.mu.Lock()
	defer c.mu.Unlock()
	if value := c.values[route]; value != nil {
		return value
	}
	value := resolve(route)
	if value != nil {
		c.values[route] = value
	}
	return value
}

// Invalidate drops the cached entry for route. Callers that mutate the registry
// (replacing a decoder) use this to ensure the next Get reflects the new state.
func (c *Cache) Invalidate(route string) {
	c.mu.Lock()
	delete(c.values, route)
	c.mu.Unlock()
}
