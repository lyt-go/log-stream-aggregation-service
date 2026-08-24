package service

import (
	"sort"
	"time"

	"logaggregation/internal/model"
	"logaggregation/pkg/idgen"
)

func (s *Service) CreateSearchRecord(input model.SearchRecord) (*model.SearchRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	now := time.Now()
	sr := &model.SearchRecord{
		ID:          idgen.Hex(),
		Query:       input.Query,
		ResultCount: input.ResultCount,
		DurationMs:  input.DurationMs,
		CreatedAt:   now,
	}
	if err := s.store.CreateSearchRecord(sr); err != nil {
		return nil, err
	}
	return sr, nil
}

func (s *Service) GetSearchRecord(id string) (*model.SearchRecord, error) {
	return s.store.GetSearchRecord(id)
}

func (s *Service) ListSearchRecords(filter model.SearchRecordFilter, page, size int) ([]*model.SearchRecord, int, error) {
	all := s.store.ListSearchRecords()
	matched := make([]*model.SearchRecord, 0, len(all))
	for _, sr := range all {
		if filter.Match(sr) {
			matched = append(matched, sr)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.SearchRecord{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateSearchRecord(id string, input model.SearchRecord) (*model.SearchRecord, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	sr, err := s.store.GetSearchRecord(id)
	if err != nil {
		return nil, err
	}
	sr.Query = input.Query
	sr.ResultCount = input.ResultCount
	sr.DurationMs = input.DurationMs
	if err := s.store.UpdateSearchRecord(sr); err != nil {
		return nil, err
	}
	return sr, nil
}

func (s *Service) DeleteSearchRecord(id string) error {
	return s.store.DeleteSearchRecord(id)
}
