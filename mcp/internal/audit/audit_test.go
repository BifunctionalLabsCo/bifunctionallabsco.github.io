package audit

import (
	"os"
	"path/filepath"
	"testing"
)

func TestRecordCreatesPrivateJSONL(t *testing.T) {
	path := filepath.Join(t.TempDir(), "nested", "audit.jsonl")
	logger := New(path)
	if err := logger.Record(Event{Operation: "test", Outcome: "success"}); err != nil {
		t.Fatal(err)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := info.Mode().Perm(); got != 0o600 {
		t.Fatalf("mode = %o, want 600", got)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) == 0 || b[len(b)-1] != '\n' {
		t.Fatalf("expected JSONL record, got %q", b)
	}
}
