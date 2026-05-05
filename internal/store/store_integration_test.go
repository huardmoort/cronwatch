package store_test

import (
	"testing"
	"time"

	"github.com/yourorg/cronwatch/internal/store"
)

// TestMultipleJobs verifies that LastRun isolates records by job name.
func TestMultipleJobs(t *testing.T) {
	s, _ := store.New(tempPath(t))

	now := time.Now().UTC()
	_ = s.Record(store.JobRun{JobName: "alpha", StartedAt: now.Add(-2 * time.Hour), Success: true})
	_ = s.Record(store.JobRun{JobName: "beta", StartedAt: now.Add(-1 * time.Hour), Success: false})
	_ = s.Record(store.JobRun{JobName: "alpha", StartedAt: now, Success: false})

	alpha, foundA := s.LastRun("alpha")
	beta, foundB := s.LastRun("beta")
	_, foundC := s.LastRun("gamma")

	if !foundA {
		t.Error("alpha not found")
	}
	if alpha.Success {
		t.Error("expected latest alpha run to be failed")
	}
	if !foundB {
		t.Error("beta not found")
	}
	if beta.Success {
		t.Error("expected beta run to be failed")
	}
	if foundC {
		t.Error("gamma should not exist")
	}
}
