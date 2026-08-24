package batchpipeline

import (
	"strings"
	"testing"
)

func TestNextBatchCommitsAfterWorkerRecovery(t *testing.T) {
	pipeline := New()
	if result, err := pipeline.Process("ingest-a", "batch-old", "panic"); result != "" || err == nil || !strings.Contains(err.Error(), "worker panic") {
		t.Fatalf("panic batch = result %q error %v, want empty result and recovered worker panic", result, err)
	}

	result, err := pipeline.Process("ingest-a", "batch-new", "healthy")
	if err != nil {
		t.Fatalf("healthy batch after recovered panic returned %v", err)
	}
	if result != "committed:batch-new:decoded:healthy" {
		t.Fatalf("healthy batch result = %q, want committed:batch-new:decoded:healthy", result)
	}
}
