package store

import (
	"strings"

	"logaggregation/internal/model"
)

func (s *MemoryStore) CreateCollector(c *model.Collector) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	for _, exist := range s.collectors {
		if strings.EqualFold(exist.Name, c.Name) {
			return ErrConflict
		}
	}
	s.collectors[c.ID] = c
	return nil
}

func (s *MemoryStore) GetCollector(id string) (*model.Collector, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c, ok := s.collectors[id]
	if !ok {
		return nil, ErrNotFound
	}
	return c, nil
}

func (s *MemoryStore) GetCollectorByName(name string) (*model.Collector, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, c := range s.collectors {
		if strings.EqualFold(c.Name, name) {
			return c, nil
		}
	}
	return nil, ErrNotFound
}

func (s *MemoryStore) ListCollectors() []*model.Collector {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]*model.Collector, 0, len(s.collectors))
	for _, c := range s.collectors {
		list = append(list, c)
	}
	return list
}

func (s *MemoryStore) UpdateCollector(c *model.Collector) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.collectors[c.ID]; !ok {
		return ErrNotFound
	}
	for _, exist := range s.collectors {
		if exist.ID != c.ID && strings.EqualFold(exist.Name, c.Name) {
			return ErrConflict
		}
	}
	s.collectors[c.ID] = c
	return nil
}

func (s *MemoryStore) DeleteCollector(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.collectors[id]; !ok {
		return ErrNotFound
	}
	delete(s.collectors, id)
	return nil
}
