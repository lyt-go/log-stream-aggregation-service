package batchlease

import (
	"fmt"
	"sync"
)

type Registry struct {
	mu     sync.Mutex
	active map[string]string
}

func New() *Registry { return &Registry{active: make(map[string]string)} }

func (r *Registry) Begin(route, batchID string) (string, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if current := r.active[route]; current != "" {
		return "", fmt.Errorf("route %s still leased by %s", route, current)
	}
	r.active[route] = batchID
	return route + "/" + batchID, nil
}

func (r *Registry) End(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()
}
