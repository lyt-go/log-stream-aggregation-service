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
	if err := p.lastErr[route]; err != nil {
		return "", err
	}
	token, err := p.leases.Begin(route, batchID)
	if err != nil {
		return "", err
	}
	value, err := p.worker.Execute(payload)
	if err != nil {
		p.commits.MarkPanic(route, batchID)
		p.lastErr[route] = err
		return "", err
	}
	p.leases.End(token)
	result, err := p.commits.Commit(route, batchID, value)
	if err != nil {
		return "", fmt.Errorf("commit batch: %w", err)
	}
	return result, nil
}
