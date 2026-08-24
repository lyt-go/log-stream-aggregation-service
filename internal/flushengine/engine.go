package flushengine

import (
	"sync"

	"logaggregation/internal/checkpointstate"
	"logaggregation/internal/flushqueue"
	"logaggregation/internal/windowbuffer"
)

type Sink interface{ Write(string, []string) error }

type Engine struct {
	mu      sync.Mutex
	buffer  *windowbuffer.Buffer
	queue   *flushqueue.Queue
	state   *checkpointstate.State
	sink    Sink
	lastErr map[string]error
}

func New(sink Sink) *Engine {
	return &Engine{buffer: &windowbuffer.Buffer{}, queue: &flushqueue.Queue{}, state: checkpointstate.New(), sink: sink, lastErr: make(map[string]error)}
}

func (e *Engine) Append(value string) { e.buffer.Append(value) }

func (e *Engine) Flush(stream string) error {
	e.mu.Lock()
	defer e.mu.Unlock()
	if err := e.lastErr[stream]; err != nil {
		return err
	}
	if !e.state.ShouldFlush(stream) {
		return nil
	}
	batch := e.queue.Pending()
	if len(batch) == 0 {
		batch = e.buffer.Snapshot()
		e.queue.Enqueue(batch)
		e.buffer.Reset()
	}
	e.state.MarkAttempt(stream)
	if err := e.sink.Write(stream, batch); err != nil {
		e.queue.Fail()
		e.lastErr[stream] = err
		return err
	}
	e.queue.Ack()
	e.state.MarkSuccess(stream)
	return nil
}
