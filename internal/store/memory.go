package store

import (
	"sync"

	"logaggregation/internal/model"
)

type MemoryStore struct {
	mu             sync.RWMutex
	streams        map[string]*model.Stream
	collectors     map[string]*model.Collector
	logEntries     map[string]*model.LogEntry
	searchRecords  map[string]*model.SearchRecord
	alertRules     map[string]*model.AlertRule
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		streams:       make(map[string]*model.Stream),
		collectors:    make(map[string]*model.Collector),
		logEntries:    make(map[string]*model.LogEntry),
		searchRecords: make(map[string]*model.SearchRecord),
		alertRules:    make(map[string]*model.AlertRule),
	}
}

var _ Store = (*MemoryStore)(nil)
