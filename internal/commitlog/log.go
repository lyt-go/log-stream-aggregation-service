package commitlog

import (
	"fmt"
	"sync"
)

type Log struct {
	mu      sync.Mutex
	aborted map[string]bool
}

func New() *Log { return &Log{aborted: make(map[string]bool)} }

func (l *Log) MarkPanic(route, batchID string) {
	l.mu.Lock()
	l.aborted[route] = true
	l.mu.Unlock()
}

func (l *Log) Commit(route, batchID, value string) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.aborted[route] {
		return "", fmt.Errorf("route %s transaction is aborted", route)
	}
	return "committed:" + batchID + ":" + value, nil
}
