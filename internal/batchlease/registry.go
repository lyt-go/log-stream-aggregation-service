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

// End 释放 token 对应的租约。租约必须按批次释放，
// 否则同一 route 上的失败批次会永久占用租约，拦截后续批次。
func (r *Registry) End(token string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for route, bid := range r.active {
		if route+"/"+bid == token {
			delete(r.active, route)
			return
		}
	}
}
