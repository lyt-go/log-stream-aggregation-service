package panicguard

import (
	"fmt"
	"sync"
)

type Guard struct {
	mu sync.Mutex
}

func New() *Guard { return &Guard{} }

// Run executes fn, recovering any panic and converting it to an error scoped to
// THIS call only. A panic must not durably block the route: the decoder may be
// replaced between requests, so each call gets a fresh chance to run.
func (g *Guard) Run(route string, fn func() string) (value string, err error) {
	g.mu.Lock()
	defer g.mu.Unlock()
	defer func() {
		if recovered := recover(); recovered != nil {
			value = ""
			err = fmt.Errorf("decoder panic: %v", recovered)
		}
	}()
	return fn(), nil
}
