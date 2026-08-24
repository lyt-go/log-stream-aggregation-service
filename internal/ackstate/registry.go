package ackstate

import "sync"

type Registry struct {
	mu     sync.Mutex
	failed map[string]bool
}

func New() *Registry { return &Registry{failed: make(map[string]bool)} }

func (r *Registry) MarkFailed(route, messageID string) {
	r.mu.Lock()
	r.failed[route] = true
	r.mu.Unlock()
}

func (r *Registry) CanSend(route, messageID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return !r.failed[route]
}
