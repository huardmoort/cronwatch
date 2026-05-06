package ratelimit

import (
	"testing"
	"time"
)

func TestAllow_FirstCallPermitted(t *testing.T) {
	l := New(5 * time.Minute)
	if !l.Allow("job-a") {
		t.Fatal("expected first call to be allowed")
	}
}

func TestAllow_SecondCallSuppressed(t *testing.T) {
	l := New(5 * time.Minute)
	l.Allow("job-a")
	if l.Allow("job-a") {
		t.Fatal("expected second call within cooldown to be suppressed")
	}
}

func TestAllow_DifferentJobsIndependent(t *testing.T) {
	l := New(5 * time.Minute)
	l.Allow("job-a")
	if !l.Allow("job-b") {
		t.Fatal("expected different job to be allowed independently")
	}
}

func TestAllow_ZeroCooldownAlwaysPermits(t *testing.T) {
	l := New(0)
	for i := 0; i < 5; i++ {
		if !l.Allow("job-a") {
			t.Fatalf("expected call %d to be allowed with zero cooldown", i)
		}
	}
}

func TestReset_ClearsSuppressionRecord(t *testing.T) {
	l := New(5 * time.Minute)
	l.Allow("job-a")
	l.Reset("job-a")
	if !l.Allow("job-a") {
		t.Fatal("expected call to be allowed after reset")
	}
}

func TestResetAll_ClearsAllRecords(t *testing.T) {
	l := New(5 * time.Minute)
	l.Allow("job-a")
	l.Allow("job-b")
	l.ResetAll()
	if !l.Allow("job-a") || !l.Allow("job-b") {
		t.Fatal("expected all jobs to be allowed after ResetAll")
	}
}

func TestAllow_PermitsAfterCooldownExpires(t *testing.T) {
	cooldown := 20 * time.Millisecond
	l := New(cooldown)
	l.Allow("job-a")
	time.Sleep(cooldown + 5*time.Millisecond)
	if !l.Allow("job-a") {
		t.Fatal("expected call to be allowed after cooldown elapsed")
	}
}
