package ackstate

import "sync"

// Registry 追踪每条消息的确认状态。
//
// 注意：失败状态按消息（route+messageID）记录，且单次有效——
// 一条消息失败只阻断它自己的重发，不能永久污染整条路由，
// 否则后续走同一路由的新消息会继承上次的失败标记。
type Registry struct {
	mu     sync.Mutex
	failed map[string]bool
}

func New() *Registry { return &Registry{failed: make(map[string]bool)} }

// key 为 route 与 messageID 的组合，避免不同路由上的同名消息互相误伤。
func (r *Registry) key(route, messageID string) string {
	return route + "\x00" + messageID
}

// MarkFailed 标记「这一条」消息失败，仅阻断该消息的重发。
func (r *Registry) MarkFailed(route, messageID string) {
	r.mu.Lock()
	r.failed[r.key(route, messageID)] = true
	r.mu.Unlock()
}

// CanSend 判断指定消息是否可发送。失败记录单次有效，查询后即清除，
// 这样同一消息重发时只看本次结果，新消息更是天然为可发送。
func (r *Registry) CanSend(route, messageID string) bool {
	r.mu.Lock()
	defer r.mu.Unlock()
	k := r.key(route, messageID)
	if r.failed[k] {
		delete(r.failed, k)
		return false
	}
	return true
}
