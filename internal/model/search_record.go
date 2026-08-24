package model

import (
	"strings"
	"time"
)

type SearchRecord struct {
	ID          string    `json:"id"`
	Query       string    `json:"query"`
	ResultCount int       `json:"result_count"`
	DurationMs  int       `json:"duration_ms"`
	CreatedAt   time.Time `json:"created_at"`
}

func (s *SearchRecord) Validate() error {
	s.Query = strings.TrimSpace(s.Query)
	if s.Query == "" {
		return NewValidationError("query", "检索关键词不能为空")
	}
	if s.ResultCount < 0 {
		return NewValidationError("result_count", "结果数不能为负数")
	}
	if s.DurationMs < 0 {
		return NewValidationError("duration_ms", "耗时不能为负数")
	}
	return nil
}

type SearchRecordFilter struct {
	Keyword string
}

func (f SearchRecordFilter) Match(s *SearchRecord) bool {
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Query), k) {
			return false
		}
	}
	return true
}
