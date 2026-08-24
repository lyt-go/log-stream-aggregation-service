package model

import (
	"strings"
	"time"
)

const (
	LogLevelDebug = "debug"
	LogLevelInfo  = "info"
	LogLevelWarn  = "warn"
	LogLevelError = "error"
)

type LogEntry struct {
	ID        string    `json:"id"`
	StreamID  string    `json:"stream_id"`
	Level     string    `json:"level"`
	Message   string    `json:"message"`
	Tags      []string  `json:"tags"`
	Timestamp time.Time `json:"timestamp"`
}

func (e *LogEntry) Validate() error {
	e.StreamID = strings.TrimSpace(e.StreamID)
	e.Message = strings.TrimSpace(e.Message)
	if e.StreamID == "" {
		return NewValidationError("stream_id", "所属日志流不能为空")
	}
	if e.Message == "" {
		return NewValidationError("message", "日志消息不能为空")
	}
	if e.Level == "" {
		e.Level = LogLevelInfo
	}
	if e.Level != LogLevelDebug && e.Level != LogLevelInfo && e.Level != LogLevelWarn && e.Level != LogLevelError {
		return NewValidationError("level", "日志级别不合法")
	}
	return nil
}

type LogEntryFilter struct {
	StreamID string
	Level    string
	Keyword  string
	Tag      string
}

func (f LogEntryFilter) Match(e *LogEntry) bool {
	if f.StreamID != "" && e.StreamID != f.StreamID {
		return false
	}
	if f.Level != "" && e.Level != f.Level {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(e.Message), k) {
			return false
		}
	}
	if f.Tag != "" {
		found := false
		for _, t := range e.Tags {
			if t == f.Tag {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}
