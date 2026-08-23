package relay

func (s *session) Deliver(jobs []Job) ([]Receipt, error) {
	ctx := s.begin()
	return deliver(ctx, jobs, s.cancel)
}
