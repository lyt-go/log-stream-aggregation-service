package ingestion

import (
	"logaggregation/internal/batchparser"
	"logaggregation/internal/inflight"
	"logaggregation/internal/supervisor"
)

type Engine struct {
	parser     batchparser.Parser
	tracker    inflight.Tracker
	supervisor supervisor.Supervisor
	lastErr    error
}

func New() *Engine { return &Engine{} }

func (e *Engine) Process(payload string) ([]string, error) {
	if e.lastErr != nil {
		return nil, e.lastErr
	}
	lease, err := e.tracker.Begin()
	if err != nil {
		return nil, err
	}
	entries, err := e.supervisor.Run(func() ([]string, error) {
		return e.parser.Parse(payload)
	})
	if err != nil {
		e.lastErr = err
		return nil, err
	}
	lease.Commit()
	lease.Done()
	return entries, nil
}

func (e *Engine) ActiveBatches() int { return e.tracker.Active() }
