package alertlog_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/alertlog"
)

func TestAlertLog_MultipleJobs_IndependentEntries(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "alerts.json")

	log, err := alertlog.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	now := time.Now().UTC()

	if err := log.Record("job-a", "missed run"); err != nil {
		t.Fatalf("Record job-a: %v", err)
	}
	if err := log.Record("job-b", "exit code 1"); err != nil {
		t.Fatalf("Record job-b: %v", err)
	}
	if err := log.Record("job-a", "missed run again"); err != nil {
		t.Fatalf("Record job-a second: %v", err)
	}

	entries := log.Entries()
	if len(entries) != 3 {
		t.Fatalf("expected 3 entries, got %d", len(entries))
	}

	if entries[0].Job != "job-a" {
		t.Errorf("entries[0].Job = %q, want job-a", entries[0].Job)
	}
	if entries[1].Job != "job-b" {
		t.Errorf("entries[1].Job = %q, want job-b", entries[1].Job)
	}
	if entries[2].Job != "job-a" {
		t.Errorf("entries[2].Job = %q, want job-a", entries[2].Job)
	}

	for _, e := range entries {
		if e.At.Before(now.Add(-time.Second)) {
			t.Errorf("entry timestamp %v is unexpectedly old", e.At)
		}
	}

	// Reload from disk and verify persistence.
	log2, err := alertlog.New(path)
	if err != nil {
		t.Fatalf("New (reload): %v", err)
	}
	loaded := log2.Entries()
	if len(loaded) != 3 {
		t.Fatalf("reloaded: expected 3 entries, got %d", len(loaded))
	}

	// Cleanup.
	_ = os.Remove(path)
}

func TestAlertLog_EmptyFile_Reload(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "alerts.json")

	log, err := alertlog.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if len(log.Entries()) != 0 {
		t.Errorf("expected empty log, got %d entries", len(log.Entries()))
	}
}
