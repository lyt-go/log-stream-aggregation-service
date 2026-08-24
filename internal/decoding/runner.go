package decoding

import (
	"logaggregation/internal/decodeplugin"
	"logaggregation/internal/quarantine"
	"logaggregation/internal/telemetry"
)

type Runner struct {
	plugin   decodeplugin.Plugin
	guard    *quarantine.Guard
	reporter telemetry.Reporter
	lastErr  map[string]error
}

func New(reportingEnabled bool) *Runner {
	return &Runner{
		guard:    quarantine.New(),
		reporter: telemetry.New(reportingEnabled),
		lastErr:  make(map[string]error),
	}
}

func (r *Runner) Process(route, payload string) ([]string, error) {
	if err := r.lastErr[route]; err != nil {
		return nil, err
	}
	entries, err := r.guard.Run(route, r.reporter, func() ([]string, error) {
		return r.plugin.Decode(payload)
	})
	if err != nil {
		r.lastErr[route] = err
	}
	return entries, err
}
