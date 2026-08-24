package inflight

import (
	"errors"
	"sync"
)

var ErrBusy = errors.New("ingestion already has an active batch")

type Tracker struct {
	mu     sync.Mutex
	active int
}

type Lease struct {
	tracker   *Tracker
	committed bool
}

func (t *Tracker) Begin() (*Lease, error) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.active != 0 {
		return nil, ErrBusy
	}
	t.active++
	return &Lease{tracker: t}, nil
}

func (l *Lease) Commit() { l.committed = true }

func (l *Lease) Done() {
	if !l.committed {
		return
	}
	l.tracker.mu.Lock()
	l.tracker.active--
	l.tracker.mu.Unlock()
}

func (t *Tracker) Active() int {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.active
}
