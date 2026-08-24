package decodepipeline

import (
	"strings"
	"testing"
)

func TestHealthyCodecTakesOverAfterTypedNilRecovery(t *testing.T) {
	pipeline := New()
	pipeline.RegisterNil("json-west")
	if value, err := pipeline.Decode("json-west", "old"); value != "" || err == nil || !strings.Contains(err.Error(), "decoder panic") {
		t.Fatalf("typed-nil decode = value %q error %v, want empty value and recovered decoder panic", value, err)
	}

	pipeline.Register("json-west", "decoded:")
	value, err := pipeline.Decode("json-west", "fresh")
	if err != nil {
		t.Fatalf("replacement decoder after panic returned %v", err)
	}
	if value != "decoded:fresh" {
		t.Fatalf("replacement decoder value = %q, want decoded:fresh", value)
	}
}
