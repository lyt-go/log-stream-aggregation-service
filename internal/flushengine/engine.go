package flushengine

import (
	"sync"

	"logaggregation/internal/checkpointstate"
	"logaggregation/internal/flushqueue"
	"logaggregation/internal/windowbuffer"
)

type Sink interface{ Write(string, []string) error }

type Engine struct {
	mu     sync.Mutex
	buffer *windowbuffer.Buffer
	queue  *flushqueue.Queue
	state  *checkpointstate.State
	sink   Sink
}

func New(sink Sink) *Engine {
	return &Engine{
		buffer: &windowbuffer.Buffer{},
		queue:  &flushqueue.Queue{},
		state:  checkpointstate.New(),
		sink:   sink,
	}
}

func (e *Engine) Append(value string) { e.buffer.Append(value) }

// Flush writes the current window's batch for stream to the sink.
//
// If a previous flush failed, the failed snapshot (e.g. [old]) is retried
// verbatim; items appended after the failure (e.g. fresh) stay in the buffer
// for the next window and are never folded into the retry batch. A transient
// sink failure must not permanently block retries: each Flush re-invokes the
// sink rather than returning a stale cached error. The checkpoint advances
// only after the sink acknowledges the write.
func (e *Engine) Flush(stream string) error {
	e.mu.Lock()
	defer e.mu.Unlock()

	// Retry an in-flight failed snapshot first; otherwise drain the current
	// window. Either path leaves `batch` as the exact slice to write.
	batch := e.queue.Pending()
	if len(batch) == 0 {
		batch = e.buffer.Snapshot()
		if len(batch) == 0 {
			// Nothing to flush and nothing to retry: leave the checkpoint and
			// the sink untouched.
			return nil
		}
		e.queue.Enqueue(batch)
		e.buffer.Reset()
	}

	if err := e.sink.Write(stream, batch); err != nil {
		// Keep the snapshot pending so the next Flush retries the same batch;
		// do not advance the checkpoint.
		e.queue.Fail()
		return err
	}
	e.queue.Ack()
	e.state.MarkSuccess(stream)
	return nil
}
