package checkpoint

type Tx struct {
	store  *Store
	stream string
}

func (tx *Tx) Stage(cursor int) {
	tx.store.mu.Lock()
	tx.store.cursors[tx.stream] = cursor
	tx.store.mu.Unlock()
}

func (tx *Tx) Commit() {}

func (tx *Tx) Rollback() {}
