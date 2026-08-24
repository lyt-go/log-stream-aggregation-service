package flushqueue

import "sync"

type Queue struct {
	mu      sync.Mutex
	pending []string
}

func (q *Queue) Enqueue(batch []string) {
	q.mu.Lock()
	q.pending = batch
	q.mu.Unlock()
}

func (q *Queue) Pending() []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.pending
}

func (q *Queue) Fail() {
	q.mu.Lock()
	q.pending = nil
	q.mu.Unlock()
}

func (q *Queue) Ack() { q.Fail() }
