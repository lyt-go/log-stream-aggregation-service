package outsink

import (
	"errors"
	"sync"
)

var ErrTemporary = errors.New("temporary sink failure")

type Sink struct {
	mu       sync.Mutex
	failNext map[string]bool
	reserved map[string]bool
	entries  map[string][]string
}

func New() *Sink {
	return &Sink{
		failNext: make(map[string]bool),
		reserved: make(map[string]bool),
		entries:  make(map[string][]string),
	}
}

func (s *Sink) FailNext(batchID string) {
	s.mu.Lock()
	s.failNext[batchID] = true
	s.mu.Unlock()
}

func (s *Sink) Send(batchID string, entries []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.reserved[batchID] {
		return nil
	}
	s.reserved[batchID] = true
	if s.failNext[batchID] {
		delete(s.failNext, batchID)
		return ErrTemporary
	}
	s.entries[batchID] = append([]string(nil), entries...)
	return nil
}

func (s *Sink) Entries(batchID string) []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return append([]string(nil), s.entries[batchID]...)
}
