package store

import (
	"logaggregation/internal/model"
)

func (s *MemoryStore) CreateSearchRecord(sr *model.SearchRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.searchRecords[sr.ID] = sr
	return nil
}

func (s *MemoryStore) GetSearchRecord(id string) (*model.SearchRecord, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	sr, ok := s.searchRecords[id]
	if !ok {
		return nil, ErrNotFound
	}
	return sr, nil
}

func (s *MemoryStore) ListSearchRecords() []*model.SearchRecord {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.SearchRecord, 0, len(s.searchRecords))
	for _, sr := range s.searchRecords {
		list = append(list, sr)
	}
	return list
}

func (s *MemoryStore) UpdateSearchRecord(sr *model.SearchRecord) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.searchRecords[sr.ID]; !ok {
		return ErrNotFound
	}
	s.searchRecords[sr.ID] = sr
	return nil
}

func (s *MemoryStore) DeleteSearchRecord(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.searchRecords[id]; !ok {
		return ErrNotFound
	}
	delete(s.searchRecords, id)
	return nil
}
