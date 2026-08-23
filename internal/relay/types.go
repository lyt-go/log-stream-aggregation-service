package relay

import "context"

type Job struct {
	ID   string
	Fail bool
}

type Receipt struct{ ID string }

type session struct {
	ctx    context.Context
	cancel context.CancelFunc
}
