package flushengine

import (
	"errors"
	"reflect"
	"testing"
)

type flakySink struct {
	calls   int
	batches [][]string
}

func (s *flakySink) Write(stream string, batch []string) error {
	s.calls++
	s.batches = append(s.batches, append([]string(nil), batch...))
	if s.calls == 1 {
		return errors.New("temporary sink failure")
	}
	return nil
}

func TestFailedFlushRetriesOriginalSnapshotBeforeNewWindow(t *testing.T) {
	sink := &flakySink{}
	engine := New(sink)
	engine.Append("old")
	if err := engine.Flush("metrics-east"); err == nil || err.Error() != "temporary sink failure" {
		t.Fatalf("first flush error = %v, want temporary sink failure", err)
	}

	engine.Append("fresh")
	if err := engine.Flush("metrics-east"); err != nil {
		t.Fatalf("retry after temporary failure returned %v", err)
	}
	if sink.calls != 2 {
		t.Fatalf("sink calls = %d, want 2", sink.calls)
	}
	if !reflect.DeepEqual(sink.batches[1], []string{"old"}) {
		t.Fatalf("retried batch = %v, want [old] before the fresh window", sink.batches[1])
	}
}
