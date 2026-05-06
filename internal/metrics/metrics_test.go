package metrics_test

import (
	"sync"
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/metrics"
)

func TestNew_ZeroValues(t *testing.T) {
	c := metrics.New()
	snap := c.Snapshot()
	if snap.ChecksTotal != 0 || snap.AlertsSent != 0 ||
		snap.HeartbeatsRecv != 0 || snap.MissedJobs != 0 {
		t.Fatalf("expected all zeros, got %+v", snap)
	}
}

func TestIncrements(t *testing.T) {
	c := metrics.New()
	c.IncChecks()
	c.IncChecks()
	c.IncAlerts()
	c.IncHeartbeats()
	c.IncMissed()
	c.IncMissed()

	snap := c.Snapshot()
	if snap.ChecksTotal != 2 {
		t.Errorf("ChecksTotal: want 2, got %d", snap.ChecksTotal)
	}
	if snap.AlertsSent != 1 {
		t.Errorf("AlertsSent: want 1, got %d", snap.AlertsSent)
	}
	if snap.HeartbeatsRecv != 1 {
		t.Errorf("HeartbeatsRecv: want 1, got %d", snap.HeartbeatsRecv)
	}
	if snap.MissedJobs != 2 {
		t.Errorf("MissedJobs: want 2, got %d", snap.MissedJobs)
	}
}

func TestSnapshot_CapturedAt(t *testing.T) {
	before := time.Now()
	c := metrics.New()
	snap := c.Snapshot()
	after := time.Now()

	if snap.CapturedAt.Before(before) || snap.CapturedAt.After(after) {
		t.Errorf("CapturedAt %v not in expected range [%v, %v]", snap.CapturedAt, before, after)
	}
}

func TestUptime_Positive(t *testing.T) {
	c := metrics.New()
	time.Sleep(2 * time.Millisecond)
	if c.Uptime() <= 0 {
		t.Error("expected positive uptime")
	}
}

func TestConcurrentIncrements(t *testing.T) {
	c := metrics.New()
	const goroutines = 50
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			c.IncChecks()
			c.IncAlerts()
		}()
	}
	wg.Wait()
	snap := c.Snapshot()
	if snap.ChecksTotal != goroutines {
		t.Errorf("ChecksTotal: want %d, got %d", goroutines, snap.ChecksTotal)
	}
	if snap.AlertsSent != goroutines {
		t.Errorf("AlertsSent: want %d, got %d", goroutines, snap.AlertsSent)
	}
}
