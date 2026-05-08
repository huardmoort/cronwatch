package alertlog_test

import (
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/alertlog"
)

func TestThrottle_ZeroIntervalAlwaysAllows(t *testing.T) {
	p := alertlog.NewThrottlePolicy(0)
	p.Record("job-a")
	if !p.Allow("job-a") {
		t.Fatal("expected Allow to return true with zero interval")
	}
}

func TestThrottle_FirstCallAllowed(t *testing.T) {
	p := alertlog.NewThrottlePolicy(5 * time.Minute)
	if !p.Allow("job-a") {
		t.Fatal("expected first Allow to return true")
	}
}

func TestThrottle_SecondCallSuppressed(t *testing.T) {
	p := alertlog.NewThrottlePolicy(5 * time.Minute)
	p.Record("job-a")
	if p.Allow("job-a") {
		t.Fatal("expected Allow to be suppressed after recent Record")
	}
}

func TestThrottle_DifferentJobsIndependent(t *testing.T) {
	p := alertlog.NewThrottlePolicy(5 * time.Minute)
	p.Record("job-a")
	if !p.Allow("job-b") {
		t.Fatal("expected job-b to be allowed independently of job-a")
	}
}

func TestThrottle_AllowedAfterInterval(t *testing.T) {
	p := alertlog.NewThrottlePolicy(10 * time.Millisecond)
	p.Record("job-a")
	time.Sleep(20 * time.Millisecond)
	if !p.Allow("job-a") {
		t.Fatal("expected Allow to return true after interval has elapsed")
	}
}

func TestThrottle_Reset_ClearsRecord(t *testing.T) {
	p := alertlog.NewThrottlePolicy(5 * time.Minute)
	p.Record("job-a")
	if p.Allow("job-a") {
		t.Fatal("expected suppression before reset")
	}
	p.Reset("job-a")
	if !p.Allow("job-a") {
		t.Fatal("expected Allow to return true after Reset")
	}
}

func TestThrottle_Reset_UnknownJobIsNoop(t *testing.T) {
	p := alertlog.NewThrottlePolicy(5 * time.Minute)
	// Should not panic for unknown job.
	p.Reset("nonexistent")
	if !p.Allow("nonexistent") {
		t.Fatal("expected Allow to return true for unknown job after Reset")
	}
}
