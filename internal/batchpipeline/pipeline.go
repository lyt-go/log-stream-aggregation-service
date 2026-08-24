package batchpipeline

import (
	"fmt"
	"sync"

	"logaggregation/internal/batchlease"
	"logaggregation/internal/commitlog"
	"logaggregation/internal/panicworker"
)

type Pipeline struct {
	mu      sync.Mutex
	leases  *batchlease.Registry
	worker  *panicworker.Worker
	commits *commitlog.Log
	lastErr map[string]error
}

func New() *Pipeline {
	return &Pipeline{leases: batchlease.New(), worker: &panicworker.Worker{}, commits: commitlog.New(), lastErr: make(map[string]error)}
}

func (p *Pipeline) Process(route, batchID, payload string) (string, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// lastErr 按 (route,batchID) 隔离：失败批次的错误只对该批次自身（含重试）生效，
	// 不得泄露到同一 route 上的其它批次，也不得跨 route 误伤同 batchID 的批次。
	key := route + "|" + batchID
	if err := p.lastErr[key]; err != nil {
		return "", err
	}
	token, err := p.leases.Begin(route, batchID)
	if err != nil {
		return "", err
	}
	// 无论成功还是失败都必须释放租约，否则失败批次的租约会拦截后续批次。
	defer p.leases.End(token)
	value, err := p.worker.Execute(payload)
	if err != nil {
		p.commits.MarkPanic(route, batchID)
		p.lastErr[key] = err
		return "", err
	}
	result, err := p.commits.Commit(route, batchID, value)
	if err != nil {
		return "", fmt.Errorf("commit batch: %w", err)
	}
	return result, nil
}
