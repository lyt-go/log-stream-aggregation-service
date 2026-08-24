package substate

import "sync"

type Registry struct {
	mu   sync.Mutex
	done map[string]bool
}

func New() *Registry { return &Registry{done: make(map[string]bool)} }

func (r *Registry) MarkClosed(topic, subscriptionID string) {
	r.mu.Lock()
	r.done[topic] = true
	r.mu.Unlock()
}

func (r *Registry) CanSubscribe(topic string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	return !r.done[topic]
}
