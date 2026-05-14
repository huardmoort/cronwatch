package alertlog

import (
	"testing"
	"time"
)

func TestSnooze_IsSnoozed_BeforeExpiry(t *testing.T) {
	s := NewSnoozePolicy()
	s.Snooze("backup", 10*time.Minute)
	if !s.IsSnoozed("backup") {
		t.Fatal("expected job to be snoozed")
	}
}

func TestSnooze_NotSnoozed_WhenNeverSet(t *testing.T) {
	s := NewSnoozePolicy()
	if s.IsSnoozed("backup") {
		t.Fatal("expected job not to be snoozed")
	}
}

func TestSnooze_Expired_ReturnsFalse(t *testing.T) {
	s := NewSnoozePolicy()
	s.Snooze("backup", 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)
	if s.IsSnoozed("backup") {
		t.Fatal("expected snooze to have expired")
	}
}

func TestSnooze_Cancel_ClearsSnooze(t *testing.T) {
	s := NewSnoozePolicy()
	s.Snooze("backup", 10*time.Minute)
	s.Cancel("backup")
	if s.IsSnoozed("backup") {
		t.Fatal("expected snooze to be cancelled")
	}
}

func TestSnooze_ExtendsDuration(t *testing.T) {
	s := NewSnoozePolicy()
	s.Snooze("backup", 1*time.Millisecond)
	s.Snooze("backup", 10*time.Minute) // extend
	time.Sleep(5 * time.Millisecond)
	if !s.IsSnoozed("backup") {
		t.Fatal("expected snooze to still be active after extension")
	}
}

func TestSnooze_ZeroDuration_IsNoop(t *testing.T) {
	s := NewSnoozePolicy()
	s.Snooze("backup", 0)
	if s.IsSnoozed("backup") {
		t.Fatal("expected zero-duration snooze to be a no-op")
	}
}

func TestSnooze_ActiveJobs_ReturnsOnlyActive(t *testing.T) {
	s := NewSnoozePolicy()
	s.Snooze("job-a", 10*time.Minute)
	s.Snooze("job-b", 1*time.Millisecond)
	time.Sleep(5 * time.Millisecond)

	active := s.ActiveJobs()
	if len(active) != 1 || active[0] != "job-a" {
		t.Fatalf("expected only job-a to be active, got %v", active)
	}
}

func TestSnooze_DifferentJobs_Independent(t *testing.T) {
	s := NewSnoozePolicy()
	s.Snooze("job-a", 10*time.Minute)
	if s.IsSnoozed("job-b") {
		t.Fatal("snoozed job-a should not affect job-b")
	}
}
