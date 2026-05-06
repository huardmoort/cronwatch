package metrics_test

import (
	"testing"

	"github.com/cronwatch/cronwatch/internal/metrics"
)

// TestCollector_FullLifecycle simulates a short daemon run and verifies that
// all counters accumulate correctly across multiple check cycles.
func TestCollector_FullLifecycle(t *testing.T) {
	col := metrics.New()

	// Simulate 3 check cycles, each finding one missed job.
	for i := 0; i < 3; i++ {
		col.IncChecks()
		col.IncMissed()
		col.IncAlerts()
	}

	// Simulate 5 heartbeats arriving from jobs.
	for i := 0; i < 5; i++ {
		col.IncHeartbeats()
	}

	snap := col.Snapshot()

	if snap.ChecksTotal != 3 {
		t.Errorf("ChecksTotal: want 3, got %d", snap.ChecksTotal)
	}
	if snap.MissedJobs != 3 {
		t.Errorf("MissedJobs: want 3, got %d", snap.MissedJobs)
	}
	if snap.AlertsSent != 3 {
		t.Errorf("AlertsSent: want 3, got %d", snap.AlertsSent)
	}
	if snap.HeartbeatsRecv != 5 {
		t.Errorf("HeartbeatsRecv: want 5, got %d", snap.HeartbeatsRecv)
	}
	if col.Uptime() <= 0 {
		t.Error("expected positive uptime after lifecycle test")
	}
}
