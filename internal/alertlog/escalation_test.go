package alertlog

import (
	"testing"
	"time"
)

func TestEscalation_ZeroThreshold_NeverEscalates(t *testing.T) {
	policy := NewEscalationPolicy(0, time.Minute)
	now := time.Now()
	for i := 0; i < 10; i++ {
		if policy.Record("job1", now) {
			t.Fatal("expected no escalation with zero threshold")
		}
	}
}

func TestEscalation_ThresholdReached(t *testing.T) {
	policy := NewEscalationPolicy(3, time.Minute)
	now := time.Now()

	if policy.Record("job1", now) {
		t.Fatal("should not escalate on first alert")
	}
	if policy.Record("job1", now.Add(time.Second)) {
		t.Fatal("should not escalate on second alert")
	}
	if !policy.Record("job1", now.Add(2*time.Second)) {
		t.Fatal("expected escalation on third alert")
	}
}

func TestEscalation_WindowExpiry_ResetsCount(t *testing.T) {
	policy := NewEscalationPolicy(2, 30*time.Second)
	now := time.Now()

	policy.Record("job1", now.Add(-60*time.Second)) // outside window
	if policy.Record("job1", now) {
		t.Fatal("old entry should have been pruned; should not escalate")
	}
}

func TestEscalation_DifferentJobs_Independent(t *testing.T) {
	policy := NewEscalationPolicy(2, time.Minute)
	now := time.Now()

	policy.Record("jobA", now)
	policy.Record("jobB", now)

	if policy.Record("jobA", now.Add(time.Second)) {
		t.Fatal("jobA escalated unexpectedly")
	}
	// jobB second record should escalate independently
	if !policy.Record("jobB", now.Add(time.Second)) {
		t.Fatal("expected jobB to escalate")
	}
}

func TestEscalation_Reset_ClearsHistory(t *testing.T) {
	policy := NewEscalationPolicy(2, time.Minute)
	now := time.Now()

	policy.Record("job1", now)
	policy.Reset("job1")

	if policy.Record("job1", now.Add(time.Second)) {
		t.Fatal("expected no escalation after reset")
	}
}

func TestEscalation_Count_ReflectsWindow(t *testing.T) {
	policy := NewEscalationPolicy(5, time.Minute)
	now := time.Now()

	policy.Record("job1", now.Add(-90*time.Second)) // outside
	policy.Record("job1", now.Add(-10*time.Second)) // inside
	policy.Record("job1", now)                      // inside

	if got := policy.Count("job1", now); got != 2 {
		t.Fatalf("expected count 2, got %d", got)
	}
}
