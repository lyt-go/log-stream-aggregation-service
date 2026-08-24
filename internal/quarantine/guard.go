package quarantine

import (
	"errors"
	"fmt"

	"logaggregation/internal/telemetry"
)

var ErrRouteQuarantined = errors.New("decode route quarantined")

type Guard struct {
	blocked map[string]bool
}

func New() *Guard {
	return &Guard{blocked: make(map[string]bool)}
}

func (g *Guard) Run(route string, reporter telemetry.Reporter, decode func() ([]string, error)) (entries []string, err error) {
	if g.blocked[route] {
		return nil, ErrRouteQuarantined
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			g.blocked[route] = true
			reporter.Report(route)
			entries = nil
			err = fmt.Errorf("decoder panic: %v", recovered)
		}
	}()
	return decode()
}
