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

// Pending returns a defensive copy of the in-flight retry snapshot, or nil if
// there is no pending batch. The copy lets callers hand the slice to the sink
// without aliasing the queue's internal storage.
func (q *Queue) Pending() []string {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.pending) == 0 {
		return nil
	}
	cp := make([]string, len(q.pending))
	copy(cp, q.pending)
	return cp
}

// Fail records that the current flush attempt failed but keeps the pending
// snapshot so the batch can be retried as-is. Dropping it here would force the
// next flush to re-snapshot the buffer (which may have since received fresh
// items), mixing new data into the retry batch.
func (q *Queue) Fail() {}

// Ack clears the pending snapshot after a flush succeeds; the batch has been
// durably written to the sink and is no longer eligible for retry.
func (q *Queue) Ack() {
	q.mu.Lock()
	q.pending = nil
	q.mu.Unlock()
}
