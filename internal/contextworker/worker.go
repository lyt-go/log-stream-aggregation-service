package contextworker

import (
	"errors"

	"logaggregation/internal/requestpool"
)

var ErrRouteBlocked = errors.New("worker route blocked")

type Worker struct {
	blocked map[string]bool
}

func New() *Worker { return &Worker{blocked: make(map[string]bool)} }

func (w *Worker) Process(state *requestpool.State) (string, error) {
	if w.blocked[state.Route] {
		return "", ErrRouteBlocked
	}
	select {
	case <-state.Context.Done():
		w.blocked[state.Route] = true
		return "", state.Context.Err()
	default:
	}
	return "accepted:" + state.Payload, nil
}
