package store

import (
	"logaggregation/internal/model"
)

func (s *MemoryStore) CreateLogEntry(e *model.LogEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.logEntries[e.ID] = e
	return nil
}

func (s *MemoryStore) GetLogEntry(id string) (*model.LogEntry, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.logEntries[id]
	if !ok {
		return nil, ErrNotFound
	}
	return e, nil
}

func (s *MemoryStore) ListLogEntries() []*model.LogEntry {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.LogEntry, 0, len(s.logEntries))
	for _, e := range s.logEntries {
		list = append(list, e)
	}
	return list
}

func (s *MemoryStore) UpdateLogEntry(e *model.LogEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.logEntries[e.ID]; !ok {
		return ErrNotFound
	}
	s.logEntries[e.ID] = e
	return nil
}

func (s *MemoryStore) DeleteLogEntry(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.logEntries[id]; !ok {
		return ErrNotFound
	}
	delete(s.logEntries, id)
	return nil
}

func (s *MemoryStore) DeleteLogEntriesByIDs(ids []string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, id := range ids {
		delete(s.logEntries, id)
	}
	return nil
}
