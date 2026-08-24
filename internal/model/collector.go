package model

import (
	"strings"
	"time"
)

const (
	CollectorStatusActive  = "active"
	CollectorStatusStopped = "stopped"
)

type Collector struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	StreamID  string    `json:"stream_id"`
	Endpoint  string    `json:"endpoint"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (c *Collector) Validate() error {
	c.Name = strings.TrimSpace(c.Name)
	c.Endpoint = strings.TrimSpace(c.Endpoint)
	c.StreamID = strings.TrimSpace(c.StreamID)
	if c.Name == "" {
		return NewValidationError("name", "采集器名称不能为空")
	}
	if c.StreamID == "" {
		return NewValidationError("stream_id", "所属日志流不能为空")
	}
	if c.Endpoint == "" {
		return NewValidationError("endpoint", "采集端点不能为空")
	}
	if c.Status == "" {
		c.Status = CollectorStatusActive
	}
	if c.Status != CollectorStatusActive && c.Status != CollectorStatusStopped {
		return NewValidationError("status", "状态不合法")
	}
	return nil
}

type CollectorFilter struct {
	Status   string
	StreamID string
	Keyword  string
}

func (f CollectorFilter) Match(c *Collector) bool {
	if f.Status != "" && c.Status != f.Status {
		return false
	}
	if f.StreamID != "" && c.StreamID != f.StreamID {
		return false
	}
	if f.Keyword != "" {
		k := strings.ToLower(strings.TrimSpace(f.Keyword))
		if k != "" && !strings.Contains(strings.ToLower(c.Name), k) &&
			!strings.Contains(strings.ToLower(c.Endpoint), k) {
			return false
		}
	}
	return true
}

var collectorTransitions = map[string]map[string]bool{
	CollectorStatusActive:  {CollectorStatusStopped: true},
	CollectorStatusStopped: {CollectorStatusActive: true},
}

func CanTransitionCollector(from, to string) bool {
	if m, ok := collectorTransitions[from]; ok {
		return m[to]
	}
	return false
}
