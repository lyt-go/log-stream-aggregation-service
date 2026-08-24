package eventchannel

import (
	"fmt"
	"sync"
)

type Subscription struct {
	ID     string
	Events <-chan string
	Cancel func()
}

type stream struct {
	ch     chan string
	closed bool
}

type Hub struct {
	mu      sync.Mutex
	nextID  int
	streams map[string]*stream
}

func New() *Hub { return &Hub{streams: make(map[string]*stream)} }

func (h *Hub) Subscribe(topic string) *Subscription {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := h.streams[topic]
	if s == nil {
		s = &stream{ch: make(chan string, 1)}
		h.streams[topic] = s
	}
	h.nextID++
	id := fmt.Sprintf("sub-%d", h.nextID)
	return &Subscription{ID: id, Events: s.ch, Cancel: func() {
		h.mu.Lock()
		defer h.mu.Unlock()
		if !s.closed {
			close(s.ch)
			s.closed = true
		}
	}}
}

func (h *Hub) Publish(topic, event string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	s := h.streams[topic]
	if s == nil || s.closed {
		return
	}
	s.ch <- event
}
