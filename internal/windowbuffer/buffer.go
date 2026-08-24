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
	return b.items
}

func (b *Buffer) Reset() {
	b.mu.Lock()
	b.items = b.items[:0]
	b.mu.Unlock()
}
