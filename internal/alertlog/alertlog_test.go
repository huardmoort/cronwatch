package alertlog_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/alertlog"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "alerts.json")
}

func TestNew_EmptyWhenMissing(t *testing.T) {
	l, err := alertlog.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got := l.Entries(); len(got) != 0 {
		t.Fatalf("expected 0 entries, got %d", len(got))
	}
}

func TestRecord_AppendsEntry(t *testing.T) {
	path := tempPath(t)
	l, _ := alertlog.New(path)

	if err := l.Record("backup", "missed run"); err != nil {
		t.Fatalf("Record error: %v", err)
	}

	entries := l.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Job != "backup" {
		t.Errorf("expected job 'backup', got %q", entries[0].Job)
	}
	if entries[0].Reason != "missed run" {
		t.Errorf("expected reason 'missed run', got %q", entries[0].Reason)
	}
	if entries[0].FiredAt.IsZero() {
		t.Error("FiredAt should not be zero")
	}
}

func TestRecord_Persisted(t *testing.T) {
	path := tempPath(t)
	l, _ := alertlog.New(path)
	_ = l.Record("sync", "timeout")

	// Reload from disk
	l2, err := alertlog.New(path)
	if err != nil {
		t.Fatalf("reload error: %v", err)
	}
	entries := l2.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 persisted entry, got %d", len(entries))
	}
	if entries[0].Job != "sync" {
		t.Errorf("expected job 'sync', got %q", entries[0].Job)
	}
}

func TestNew_InvalidJSON(t *testing.T) {
	path := tempPath(t)
	_ = os.WriteFile(path, []byte("not json"), 0o644)
	_, err := alertlog.New(path)
	if err == nil {
		t.Fatal("expected error for invalid JSON")
	}
}

func TestRecord_MultipleEntries(t *testing.T) {
	path := tempPath(t)
	l, _ := alertlog.New(path)

	_ = l.Record("jobA", "missed")
	before := time.Now()
	_ = l.Record("jobB", "failed")

	entries := l.Entries()
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[1].FiredAt.Before(before) && !entries[1].FiredAt.Equal(before) {
		t.Error("second entry FiredAt should not be before first record call")
	}
}
