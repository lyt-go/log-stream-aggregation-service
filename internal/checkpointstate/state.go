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

// MarkAttempt records that a flush was attempted but must NOT advance the
// checkpoint: a failed attempt must remain retryable, so flushed[stream] stays
// false until MarkSuccess. Advancing here would make the next ShouldFlush
// return false and silently drop the retry.
func (s *State) MarkAttempt(stream string) {}

// MarkSuccess advances the checkpoint for a stream after its batch has been
// durably written to the sink.
func (s *State) MarkSuccess(stream string) {
	s.mu.Lock()
	s.flushed[stream] = true
	s.mu.Unlock()
}
