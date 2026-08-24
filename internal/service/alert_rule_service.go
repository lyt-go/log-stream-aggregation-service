package service

import (
	"sort"
	"time"

	"logaggregation/internal/model"
	"logaggregation/pkg/idgen"
)

func (s *Service) CreateAlertRule(input model.AlertRule) (*model.AlertRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.store.GetStream(input.StreamID); err != nil {
		return nil, model.NewValidationError("stream_id", "所属日志流不存在")
	}
	now := time.Now()
	a := &model.AlertRule{
		ID:        idgen.Hex(),
		StreamID:  input.StreamID,
		Level:     input.Level,
		Pattern:   input.Pattern,
		Status:    input.Status,
		CreatedAt: now,
	}
	if err := s.store.CreateAlertRule(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) GetAlertRule(id string) (*model.AlertRule, error) {
	return s.store.GetAlertRule(id)
}

func (s *Service) ListAlertRules(filter model.AlertRuleFilter, page, size int) ([]*model.AlertRule, int, error) {
	all := s.store.ListAlertRules()
	matched := make([]*model.AlertRule, 0, len(all))
	for _, a := range all {
		if filter.Match(a) {
			matched = append(matched, a)
		}
	}
	sort.Slice(matched, func(i, j int) bool {
		return matched[i].CreatedAt.After(matched[j].CreatedAt)
	})
	total := len(matched)
	start := (page - 1) * size
	if start >= total {
		return []*model.AlertRule{}, total, nil
	}
	end := start + size
	if end > total {
		end = total
	}
	return matched[start:end], total, nil
}

func (s *Service) UpdateAlertRule(id string, input model.AlertRule) (*model.AlertRule, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}
	a, err := s.store.GetAlertRule(id)
	if err != nil {
		return nil, err
	}
	if a.StreamID != input.StreamID {
		if _, err := s.store.GetStream(input.StreamID); err != nil {
			return nil, model.NewValidationError("stream_id", "所属日志流不存在")
		}
	}
	a.StreamID = input.StreamID
	a.Level = input.Level
	a.Pattern = input.Pattern
	if err := s.store.UpdateAlertRule(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) DeleteAlertRule(id string) error {
	return s.store.DeleteAlertRule(id)
}

func (s *Service) TransitionAlertRuleStatus(id string, toStatus string) (*model.AlertRule, error) {
	a, err := s.store.GetAlertRule(id)
	if err != nil {
		return nil, err
	}
	if !model.CanTransitionAlertRule(a.Status, toStatus) {
		return nil, model.NewValidationError("status", "状态流转不合法")
	}
	a.Status = toStatus
	if err := s.store.UpdateAlertRule(a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *Service) BatchUpdateAlertRuleStatus(ids []string, toStatus string) error {
	if len(ids) == 0 {
		return model.NewValidationError("ids", "ID 列表不能为空")
	}
	for _, id := range ids {
		if _, err := s.TransitionAlertRuleStatus(id, toStatus); err != nil {
			return err
		}
	}
	return nil
}
