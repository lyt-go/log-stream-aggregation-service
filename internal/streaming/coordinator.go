package streaming

import (
	"logaggregation/internal/fanout"
	"logaggregation/internal/framearchive"
	"logaggregation/internal/framebuf"
)

type Coordinator struct {
	pool    framebuf.Pool
	queue   fanout.Queue
	archive framearchive.Archive
}

func New() *Coordinator {
	return &Coordinator{}
}

func (c *Coordinator) Submit(payload []byte) {
	buf := c.pool.Copy(payload)
	c.queue.Push(buf)
	c.pool.Release(buf)
}

func (c *Coordinator) Pending() [][]byte {
	return c.queue.Peek()
}

func (c *Coordinator) Flush() {
	for _, frame := range c.queue.Drain() {
		c.archive.Append(frame)
	}
}

func (c *Coordinator) Archived() [][]byte {
	return c.archive.Frames()
}
