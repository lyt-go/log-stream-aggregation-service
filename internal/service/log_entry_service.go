package service

import (
	"sort"
	"time"

	"logaggregation/internal/model"
	"logaggregation/pkg/idgen"
)

func (s *Service) CreateLogEntry(input model.LogEntry) (*model.LogEntry, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetStream(input.StreamID); err != nil {
		return nil, model.NewValidationError("stream_id", "所属日志流不存在")
	}
	now := time.Now()
	e := &model.LogEntry{
		ID:        idgen.Hex(),
		StreamID:  input.StreamID,
		Level:     input.Level,
		Message:   input.Message,
		Tags:      input.Tags,
		Timestamp: now,
	}
	if err := s.store.CreateLogEntry(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) GetLogEntry(id string) (*model.LogEntry, error) {
	return s.store.GetLogEntry(id)
}

func (s *Service) ListLogEntries(filter model.LogEntryFilter, page, size int) ([]*model.LogEntry, int, error) {
	all := s.store.ListLogEntries()
	matched := make([]*model.LogEntry, 0, len(all))
	for _, e := range all {
		if filter.Match(e) {
			matched = append(matched, e)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].Timestamp.After(matched[j].Timestamp)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.LogEntry{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateLogEntry(id string, input model.LogEntry) (*model.LogEntry, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	e, err := s.store.GetLogEntry(id)
	if err != nil {
		return nil, err
	}
	if e.StreamID != input.StreamID {
		if _, err := s.store.GetStream(input.StreamID); err != nil {
			return nil, model.NewValidationError("stream_id", "所属日志流不存在")
		}
	}
	e.StreamID = input.StreamID
	e.Level = input.Level
	e.Message = input.Message
	e.Tags = input.Tags
	if err := s.store.UpdateLogEntry(e); err != nil {
		return nil, err
	}
	return e, nil
}

func (s *Service) DeleteLogEntry(id string) error {
	return s.store.DeleteLogEntry(id)
}

func (s *Service) BatchDeleteLogEntries(ids []string) error {
	if len(ids) == 0 {
		return model.NewValidationError("ids", "ID 列表不能为空")
	}
	return s.store.DeleteLogEntriesByIDs(ids)
}
