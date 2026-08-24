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

// aborted 按 batchID 而非 route 记录：某个批次 panic 仅标记该批次中止，
// 不得阻断同一 route 上其它批次的提交。
func (l *Log) MarkPanic(route, batchID string) {
	l.mu.Lock()
	l.aborted[route+"|"+batchID] = true
	l.mu.Unlock()
}

func (l *Log) Commit(route, batchID, value string) (string, error) {
	l.mu.Lock()
	defer l.mu.Unlock()
	key := route + "|" + batchID
	if l.aborted[key] {
		return "", fmt.Errorf("batch %s on route %s transaction is aborted", batchID, route)
	}
	return "committed:" + batchID + ":" + value, nil
}
