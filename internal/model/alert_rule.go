package model

import (
	"strings"
	"time"
)

const (
	AlertRuleStatusActive   = "active"
	AlertRuleStatusDisabled = "disabled"
)

type AlertRule struct {
	ID        string    `json:"id"`
	StreamID  string    `json:"stream_id"`
	Level     string    `json:"level"`
	Pattern   string    `json:"pattern"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (a *AlertRule) Validate() error {
	a.StreamID = strings.TrimSpace(a.StreamID)
	a.Pattern = strings.TrimSpace(a.Pattern)
	if a.StreamID == "" {
		return NewValidationError("stream_id", "所属日志流不能为空")
	}
	if a.Pattern == "" {
		return NewValidationError("pattern", "匹配规则不能为空")
	}
	if a.Level == "" {
		a.Level = LogLevelError
	}
	if a.Level != LogLevelDebug && a.Level != LogLevelInfo && a.Level != LogLevelWarn && a.Level != LogLevelError {
		return NewValidationError("level", "日志级别不合法")
	}
	if a.Status == "" {
		a.Status = AlertRuleStatusActive
	}
	if a.Status != AlertRuleStatusActive && a.Status != AlertRuleStatusDisabled {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type AlertRuleFilter struct {
	Status   string
	StreamID string
	Level    string
	Keyword  string
}

func (f AlertRuleFilter) Match(a *AlertRule) bool {
	if f.Status != "" && a.Status != f.Status {
		return false
	}
	if f.StreamID != "" && a.StreamID != f.StreamID {
		return false
	}
	if f.Level != "" && a.Level != f.Level {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(a.Pattern), k) {
			return false
		}
	}
	return true
}

var alertRuleTransitions = map[string]map[string]bool{
	AlertRuleStatusActive:   {AlertRuleStatusDisabled: true},
	AlertRuleStatusDisabled: {AlertRuleStatusActive: true},
}

func CanTransitionAlertRule(from, to string) bool {
	if m, ok := alertRuleTransitions[from]; ok {
		return m[to]
	}
	return false
}
