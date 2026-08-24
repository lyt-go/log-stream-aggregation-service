package service

import (
	"sort"
	"time"

	"logaggregation/internal/model"
	"logaggregation/pkg/idgen"
)

func (s *Service) CreateCollector(input model.Collector) (*model.Collector, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetStream(input.StreamID); err != nil {
		return nil, model.NewValidationError("stream_id", "所属日志流不存在")
	}
	now := time.Now()
	c := &model.Collector{
		ID:        idgen.Hex(),
		Name:      input.Name,
		StreamID:  input.StreamID,
		Endpoint:  input.Endpoint,
		Status:    input.Status,
		CreatedAt: now,
	}
	if err := s.store.CreateCollector(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) GetCollector(id string) (*model.Collector, error) {
	return s.store.GetCollector(id)
}

func (s *Service) ListCollectors(filter model.CollectorFilter, page, size int) ([]*model.Collector, int, error) {
	all := s.store.ListCollectors()
	matched := make([]*model.Collector, 0, len(all))
	for _, c := range all {
		if filter.Match(c) {
			matched = append(matched, c)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Collector{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateCollector(id string, input model.Collector) (*model.Collector, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	c, err := s.store.GetCollector(id)
	if err != nil {
		return nil, err
	}
	if c.StreamID != input.StreamID {
		if _, err := s.store.GetStream(input.StreamID); err != nil {
			return nil, model.NewValidationError("stream_id", "所属日志流不存在")
		}
	}
	c.Name = input.Name
	c.StreamID = input.StreamID
	c.Endpoint = input.Endpoint
	if err := s.store.UpdateCollector(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *Service) DeleteCollector(id string) error {
	return s.store.DeleteCollector(id)
}

func (s *Service) TransitionCollectorStatus(id string, toStatus string) (*model.Collector, error) {
	c, err := s.store.GetCollector(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionCollector(c.Status, toStatus) {
		return nil, model.NewValidationError("status", "状态流转不合法")
	}
	c.Status = toStatus
	if err := s.store.UpdateCollector(c); err != nil {
		return nil, err
	}
	return c, nil
}
