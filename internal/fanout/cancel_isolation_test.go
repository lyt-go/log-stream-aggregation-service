package fanout

import (
	"context"
	"errors"
	"testing"
)

func TestFreshEnvelopeSurvivesPriorRouteCancellation(t *testing.T) {
	dispatcher := New()
	cancelled, cancel := context.WithCancel(context.Background())
	cancel()
	if receipt, err := dispatcher.Dispatch(cancelled, "msg-old", "edge-east", "stale"); receipt != "" || !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled delivery = receipt %q error %v, want empty receipt and context canceled", receipt, err)
	}

	receipt, err := dispatcher.Dispatch(context.Background(), "msg-new", "edge-east", "fresh")
	if err != nil {
		t.Fatalf("fresh delivery after cancellation returned %v", err)
	}
	if receipt != "delivered:msg-new:fresh" {
		t.Fatalf("fresh delivery receipt = %q, want delivered:msg-new:fresh", receipt)
	}
}
