package model

import (
	"strings"
	"time"
)

const (
	StreamStatusActive = "active"
	StreamStatusPaused = "paused"
)

const (
	StreamFormatJSON = "json"
	StreamFormatText = "text"
)

type Stream struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Source    string    `json:"source"`
	Format    string    `json:"format"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *Stream) Validate() error {
	s.Name = strings.TrimSpace(s.Name)
	s.Source = strings.TrimSpace(s.Source)
	if s.Name == "" {
		return NewValidationError("name", "日志流名称不能为空")
	}
	if s.Source == "" {
		return NewValidationError("source", "日志源不能为空")
	}
	if s.Format == "" {
		s.Format = StreamFormatJSON
	}
	if s.Format != StreamFormatJSON && s.Format != StreamFormatText {
		return NewValidationError("format", "格式只能是 json 或 text")
	}
	if s.Status == "" {
		s.Status = StreamStatusActive
	}
	if s.Status != StreamStatusActive && s.Status != StreamStatusPaused {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type StreamFilter struct {
	Status  string
	Format  string
	Keyword string
}

func (f StreamFilter) Match(s *Stream) bool {
	if f.Status != "" && s.Status != f.Status {
		return false
	}
	if f.Format != "" && s.Format != f.Format {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(s.Name), k) &&
			!strings.Contains(strings.ToLower(s.Source), k) {
			return false
		}
	}
	return true
}

var streamTransitions = map[string]map[string]bool{
	StreamStatusActive: {StreamStatusPaused: true},
	StreamStatusPaused: {StreamStatusActive: true},
}

func CanTransitionStream(from, to string) bool {
	if m, ok := streamTransitions[from]; ok {
		return m[to]
	}
	return false
}
