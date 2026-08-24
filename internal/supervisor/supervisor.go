package supervisor

import (
	"errors"
	"fmt"
)

var ErrStopped = errors.New("parser supervisor stopped")

type Supervisor struct {
	stopped bool
}

func (s *Supervisor) Run(parse func() ([]string, error)) (entries []string, err error) {
	if s.stopped {
		return nil, ErrStopped
	}
	defer func() {
		if recovered := recover(); recovered != nil {
			s.stopped = true
			entries = nil
			err = fmt.Errorf("parser panic recovered: %v", recovered)
		}
	}()
	return parse()
}
