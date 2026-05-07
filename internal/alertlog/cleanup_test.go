package alertlog

import (
	"testing"
	"time"
)

func TestCleanup_RemovesOldEntries(t *testing.T) {
	path := tempPath(t)
	l, err := New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	now := time.Now()
	old := now.Add(-48 * time.Hour)

	// inject entries directly to control timestamps
	l.entries = []Entry{
		{Job: "job-old", Message: "missed", At: old},
		{Job: "job-new", Message: "missed", At: now},
	}
	if err := l.persist(); err != nil {
		t.Fatalf("persist: %v", err)
	}

	pruned, err := l.Cleanup(24 * time.Hour)
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	if pruned != 1 {
		t.Errorf("expected 1 pruned, got %d", pruned)
	}

	entries := l.All("")
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry remaining, got %d", len(entries))
	}
	if entries[0].Job != "job-new" {
		t.Errorf("expected job-new, got %s", entries[0].Job)
	}
}

func TestCleanup_NothingToRemove(t *testing.T) {
	path := tempPath(t)
	l, _ := New(path)

	l.entries = []Entry{
		{Job: "job-a", Message: "ok", At: time.Now()},
	}
	_ = l.persist()

	pruned, err := l.Cleanup(24 * time.Hour)
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	if pruned != 0 {
		t.Errorf("expected 0 pruned, got %d", pruned)
	}
	if len(l.All("")) != 1 {
		t.Error("entry should still be present")
	}
}

func TestCleanup_ZeroRetentionIsNoop(t *testing.T) {
	path := tempPath(t)
	l, _ := New(path)

	l.entries = []Entry{
		{Job: "job-x", Message: "missed", At: time.Now().Add(-999 * time.Hour)},
	}
	_ = l.persist()

	pruned, err := l.Cleanup(0)
	if err != nil {
		t.Fatalf("Cleanup: %v", err)
	}
	if pruned != 0 {
		t.Errorf("expected 0 pruned with zero retention, got %d", pruned)
	}
}

func TestOlderThan_FiltersCorrectly(t *testing.T) {
	path := tempPath(t)
	l, _ := New(path)

	now := time.Now()
	l.entries = []Entry{
		{Job: "old", At: now.Add(-2 * time.Hour)},
		{Job: "new", At: now.Add(time.Minute)},
	}

	out := l.OlderThan(now)
	if len(out) != 1 {
		t.Fatalf("expected 1, got %d", len(out))
	}
	if out[0].Job != "old" {
		t.Errorf("expected old, got %s", out[0].Job)
	}
}
