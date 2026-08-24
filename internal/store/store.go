// Package store 定义数据访问接口与内存实现。
package store

import (
	"errors"

	"logaggregation/internal/model"
)

var (
	ErrNotFound = errors.New("记录不存在")
	ErrConflict = errors.New("记录已存在或状态冲突")
)

type Store interface {
	CreateStream(s *model.Stream) error
	GetStream(id string) (*model.Stream, error)
	GetStreamByName(name string) (*model.Stream, error)
	ListStreams() []*model.Stream
	UpdateStream(s *model.Stream) error
	DeleteStream(id string) error

	CreateCollector(c *model.Collector) error
	GetCollector(id string) (*model.Collector, error)
	GetCollectorByName(name string) (*model.Collector, error)
	ListCollectors() []*model.Collector
	UpdateCollector(c *model.Collector) error
	DeleteCollector(id string) error

	CreateLogEntry(e *model.LogEntry) error
	GetLogEntry(id string) (*model.LogEntry, error)
	ListLogEntries() []*model.LogEntry
	UpdateLogEntry(e *model.LogEntry) error
	DeleteLogEntry(id string) error
	DeleteLogEntriesByIDs(ids []string) error

	CreateSearchRecord(sr *model.SearchRecord) error
	GetSearchRecord(id string) (*model.SearchRecord, error)
	ListSearchRecords() []*model.SearchRecord
	UpdateSearchRecord(sr *model.SearchRecord) error
	DeleteSearchRecord(id string) error

	CreateAlertRule(a *model.AlertRule) error
	GetAlertRule(id string) (*model.AlertRule, error)
	ListAlertRules() []*model.AlertRule
	UpdateAlertRule(a *model.AlertRule) error
	DeleteAlertRule(id string) error
}
