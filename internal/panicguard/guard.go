package panicguard

import (
	"fmt"
	"sync"
)

type Guard struct {
	mu      sync.Mutex
	blocked map[string]bool
}

func New() *Guard { return &Guard{blocked: make(map[string]bool)} }

func (g *Guard) Run(route string, fn func() string) (value string, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	if g.blocked[route] {
		return "", fmt.Errorf("route %s remains blocked after panic", route)
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			g.blocked[route] = true
			value = ""
			err = fmt.Errorf("decoder panic: %v", recovered)
		}
	}()
	return fn(), nil
}
