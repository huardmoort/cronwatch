package store_test

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/yourorg/cronwatch/internal/store"
)

func tempPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "runs.json")
}

func TestNew_EmptyWhenMissing(t *testing.T) {
	s, err := store.New(tempPath(t))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	_, found := s.LastRun("myjob")
	if found {
		t.Error("expected no run for new store")
	}
}

func TestRecord_AndLastRun(t *testing.T) {
	p := tempPath(t)
	s, _ := store.New(p)

	run := store.JobRun{
		JobName:   "backup",
		StartedAt: time.Now().UTC().Truncate(time.Second),
		Success:   true,
	}
	if err := s.Record(run); err != nil {
		t.Fatalf("Record failed: %v", err)
	}

	got, found := s.LastRun("backup")
	if !found {
		t.Fatal("expected to find last run")
	}
	if got.JobName != run.JobName || !got.StartedAt.Equal(run.StartedAt) || got.Success != run.Success {
		t.Errorf("got %+v, want %+v", got, run)
	}
}

func TestRecord_Persisted(t *testing.T) {
	p := tempPath(t)
	s, _ := store.New(p)
	_ = s.Record(store.JobRun{JobName: "sync", StartedAt: time.Now(), Success: false, Note: "exit 1"})

	// reload from disk
	s2, err := store.New(p)
	if err != nil {
		t.Fatalf("reload failed: %v", err)
	}
	_, found := s2.LastRun("sync")
	if !found {
		t.Error("run not persisted to disk")
	}
}

func TestNew_InvalidJSON(t *testing.T) {
	p := tempPath(t)
	_ = os.WriteFile(p, []byte("not json"), 0o644)
	_, err := store.New(p)
	if err == nil {
		t.Error("expected error for invalid JSON")
	}
}

func TestLastRun_ReturnsLatest(t *testing.T) {
	s, _ := store.New(tempPath(t))
	t1 := time.Now().Add(-time.Hour)
	t2 := time.Now()
	_ = s.Record(store.JobRun{JobName: "job", StartedAt: t1, Success: false})
	_ = s.Record(store.JobRun{JobName: "job", StartedAt: t2, Success: true})

	got, _ := s.LastRun("job")
	if !got.Success {
		t.Error("expected latest (successful) run to be returned")
	}
}
