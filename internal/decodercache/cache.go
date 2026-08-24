package decodercache

import (
	"sync"

	"logaggregation/internal/decoderregistry"
)

type Cache struct {
	mu      sync.Mutex
	values  map[string]decoderregistry.Decoder
	missing map[string]bool
}

func New() *Cache {
	return &Cache{values: make(map[string]decoderregistry.Decoder), missing: make(map[string]bool)}
}

func (c *Cache) Get(route string, resolve func(string) decoderregistry.Decoder) decoderregistry.Decoder {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.missing[route] {
		return nil
	}
	if value := c.values[route]; value != nil {
		return value
	}
	value := resolve(route)
	if value == nil {
		c.missing[route] = true
	} else {
		c.values[route] = value
	}
	return value
}
