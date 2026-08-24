package windowbuffer

import "sync"

type Buffer struct {
	mu    sync.Mutex
	items []string
}

func (b *Buffer) Append(value string) {
	b.mu.Lock()
	b.items = append(b.items, value)
	b.mu.Unlock()
}

func (b *Buffer) Snapshot() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	// Defensive copy: callers (e.g. the flush queue) may retain the returned
	// slice across subsequent Append/Reset calls, which reuse the backing
	// array. Without a copy, a later Append can mutate an enqueued snapshot
	// in place (e.g. turning [old] into [fresh]).
	cp := make([]string, len(b.items))
	copy(cp, b.items)
	return cp
}

func (b *Buffer) Reset() {
	b.mu.Lock()
	b.items = b.items[:0]
	b.mu.Unlock()
}
