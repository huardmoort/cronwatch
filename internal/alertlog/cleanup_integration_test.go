package alertlog_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/your-org/cronwatch/internal/alertlog"
)

func TestCleanup_PersistenceAfterPrune(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "alerts.json")

	l, err := alertlog.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// Record two entries via public API, then manipulate via reload.
	if err := l.Record("job-keep", "missed run"); err != nil {
		t.Fatalf("Record: %v", err)
	}

	// Reload and verify cleanup removes nothing (entry is fresh).
	l2, err := alertlog.New(path)
	if err != nil {
		t.Fatalf("New reload: %v", err)
	}

	pruned, err := l2.Cleanup(24 * time.Hour)
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	if pruned != 0 {
		t.Errorf("expected 0 pruned for fresh entry, got %d", pruned)
	}

	// Verify file still exists and has content.
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat: %v", err)
	}
	if info.Size() == 0 {
		t.Error("log file should not be empty after cleanup of fresh entries")
	}

	// Entries still readable.
	l3, _ := alertlog.New(path)
	if len(l3.All("")) != 1 {
		t.Errorf("expected 1 entry after reload, got %d", len(l3.All("")))
	}
}

func TestCleanup_EmptyLogAfterFullPrune(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "alerts.json")

	l, _ := alertlog.New(path)

	// Record one entry, then prune with tiny retention so it's removed.
	_ = l.Record("job-gone", "missed")

	// Use negative retention trick: prune anything older than "now + 1s"
	// achieved by using a very small retention that excludes even fresh entries.
	// We instead use OlderThan to confirm the entry exists, then use
	// a future-biased cutoff via a negative retention workaround:
	// directly test with 1 nanosecond after sleeping.
	time.Sleep(2 * time.Millisecond)
	pruned, err := l.Cleanup(1 * time.Nanosecond)
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	if pruned != 1 {
		t.Errorf("expected 1 pruned, got %d", pruned)
	}

	l2, _ := alertlog.New(path)
	if len(l2.All("")) != 0 {
		t.Errorf("expected 0 entries after full prune, got %d", len(l2.All("")))
	}
}
