package service

import (
	"sort"
	"time"

	"logaggregation/internal/model"
	"logaggregation/pkg/idgen"
)

func (s *Service) CreateStream(input model.Stream) (*model.Stream, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	st := &model.Stream{
		ID:        idgen.Hex(),
		Name:      input.Name,
		Source:    input.Source,
		Format:    input.Format,
		Status:    input.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.store.CreateStream(st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) GetStream(id string) (*model.Stream, error) {
	return s.store.GetStream(id)
}

func (s *Service) ListStreams(filter model.StreamFilter, page, size int) ([]*model.Stream, int, error) {
	all := s.store.ListStreams()
	matched := make([]*model.Stream, 0, len(all))
	for _, st := range all {
		if filter.Match(st) {
			matched = append(matched, st)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.Stream{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateStream(id string, input model.Stream) (*model.Stream, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	st, err := s.store.GetStream(id)
	if err != nil {
		return nil, err
	}
	st.Name = input.Name
	st.Source = input.Source
	st.Format = input.Format
	st.UpdatedAt = time.Now()
	if err := s.store.UpdateStream(st); err != nil {
		return nil, err
	}
	return st, nil
}

func (s *Service) DeleteStream(id string) error {
	return s.store.DeleteStream(id)
}

func (s *Service) TransitionStreamStatus(id string, toStatus string) (*model.Stream, error) {
	st, err := s.store.GetStream(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionStream(st.Status, toStatus) {
		return nil, model.NewValidationError("status", "状态流转不合法")
	}
	st.Status = toStatus
	st.UpdatedAt = time.Now()
	if err := s.store.UpdateStream(st); err != nil {
		return nil, err
	}
	return st, nil
}
