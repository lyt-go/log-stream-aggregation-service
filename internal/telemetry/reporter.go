package telemetry

import "sync"

type Reporter interface {
	Report(route string)
}

type Counter struct {
	mu     sync.Mutex
	counts map[string]int
}

func New(enabled bool) Reporter {
	if !enabled {
		var counter *Counter
		return counter
	}
	return &Counter{counts: make(map[string]int)}
}

func (c *Counter) Report(route string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.counts[route]++
}
