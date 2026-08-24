package checkpoint

import "sync"

type Store struct {
	mu      sync.Mutex
	cursors map[string]int
}

func NewStore() *Store {
	return &Store{cursors: make(map[string]int)}
}

func (s *Store) Cursor(stream string) int {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.cursors[stream]
}

func (s *Store) Begin(stream string) *Tx {
	return &Tx{store: s, stream: stream}
}
