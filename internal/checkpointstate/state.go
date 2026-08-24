package checkpointstate

import "sync"

type State struct {
	mu      sync.Mutex
	flushed map[string]bool
}

func New() *State { return &State{flushed: make(map[string]bool)} }

func (s *State) ShouldFlush(stream string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return !s.flushed[stream]
}

func (s *State) MarkAttempt(stream string) {
	s.mu.Lock()
	s.flushed[stream] = true
	s.mu.Unlock()
}

func (s *State) MarkSuccess(stream string) {
	s.mu.Lock()
	s.flushed[stream] = true
	s.mu.Unlock()
}
