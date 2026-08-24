package delivery

import (
	"logaggregation/internal/checkpoint"
	"logaggregation/internal/outsink"
)

type Batch struct {
	ID      string
	Stream  string
	Cursor  int
	Entries []string
}

type Coordinator struct {
	checkpoints *checkpoint.Store
	sink        *outsink.Sink
}

func New(checkpoints *checkpoint.Store, sink *outsink.Sink) *Coordinator {
	return &Coordinator{checkpoints: checkpoints, sink: sink}
}

func (c *Coordinator) Deliver(batch Batch) error {
	tx := c.checkpoints.Begin(batch.Stream)
	tx.Stage(batch.Cursor)
	if err := c.sink.Send(batch.ID, batch.Entries); err != nil {
		return err
	}
	tx.Commit()
	return nil
}
