package relay

import "context"

func New(ctx context.Context) *session {
	ctx, cancel := context.WithCancel(ctx)
	return &session{ctx: ctx, cancel: cancel}
}

func (s *session) begin() context.Context { return s.ctx }
