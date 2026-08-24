package routequeue

import (
	"sync"

	"logaggregation/internal/envelopepool"
)

type Queue struct {
	mu    sync.Mutex
	items []*envelopepool.Envelope
}

func (q *Queue) Push(env *envelopepool.Envelope) {
	q.mu.Lock()
	q.items = append(q.items, env)
	q.mu.Unlock()
}

func (q *Queue) Pop() envelopepool.Envelope {
	q.mu.Lock()
	defer q.mu.Unlock()
	item := q.items[0]
	q.items = q.items[1:]
	return *item
}
