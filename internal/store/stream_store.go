package store

import (
	"strings"
	"time"

	"logaggregation/internal/model"
)

func (s *MemoryStore) CreateStream(st *model.Stream) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.streams {
		if strings.EqualFold(exist.Name, st.Name) {
			return ErrConflict
		}
	}
	s.streams[st.ID] = st
	return nil
}

func (s *MemoryStore) GetStream(id string) (*model.Stream, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	st, ok := s.streams[id]
	if !ok {
		return nil, ErrNotFound
	}
	return st, nil
}

func (s *MemoryStore) GetStreamByName(name string) (*model.Stream, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, st := range s.streams {
		if strings.EqualFold(st.Name, name) {
			return st, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListStreams() []*model.Stream {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Stream, 0, len(s.streams))
	for _, st := range s.streams {
		list = append(list, st)
	}
	return list
}

func (s *MemoryStore) UpdateStream(st *model.Stream) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.streams[st.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.streams {
		if exist.ID != st.ID && strings.EqualFold(exist.Name, st.Name) {
			return ErrConflict
		}
	}
	st.UpdatedAt = time.Now()
	s.streams[st.ID] = st
	return nil
}

func (s *MemoryStore) DeleteStream(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.streams[id]; !ok {
		return ErrNotFound
	}
	delete(s.streams, id)
	return nil
}
