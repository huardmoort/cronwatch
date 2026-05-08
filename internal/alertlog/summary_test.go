package alertlog

import (
	"testing"
	"time"
)

func TestSummarize_EmptyLog(t *testing.T) {
	l := &Log{entries: []Entry{}}
	s := l.Summarize("")
	if s.TotalAlerts != 0 {
		t.Fatalf("expected 0 alerts, got %d", s.TotalAlerts)
	}
	if s.OldestAlert != nil || s.NewestAlert != nil {
		t.Fatal("expected nil timestamps for empty log")
	}
}

func TestSummarize_AllJobs(t *testing.T) {
	now := time.Now()
	l := &Log{
		entries: []Entry{
			{Job: "backup", At: now.Add(-2 * time.Hour), Reason: "missed"},
			{Job: "backup", At: now.Add(-1 * time.Hour), Reason: "missed"},
			{Job: "cleanup", At: now, Reason: "missed"},
		},
	}

	s := l.Summarize("")
	if s.TotalAlerts != 3 {
		t.Fatalf("expected 3 total alerts, got %d", s.TotalAlerts)
	}
	if s.ByJob["backup"] != 2 {
		t.Errorf("expected 2 backup alerts, got %d", s.ByJob["backup"])
	}
	if s.ByJob["cleanup"] != 1 {
		t.Errorf("expected 1 cleanup alert, got %d", s.ByJob["cleanup"])
	}
	if s.OldestAlert == nil || !s.OldestAlert.Equal(now.Add(-2*time.Hour)) {
		t.Errorf("unexpected oldest alert: %v", s.OldestAlert)
	}
	if s.NewestAlert == nil || !s.NewestAlert.Equal(now) {
		t.Errorf("unexpected newest alert: %v", s.NewestAlert)
	}
}

func TestSummarize_FilterByJob(t *testing.T) {
	now := time.Now()
	l := &Log{
		entries: []Entry{
			{Job: "backup", At: now.Add(-1 * time.Hour), Reason: "missed"},
			{Job: "cleanup", At: now, Reason: "missed"},
		},
	}

	s := l.Summarize("backup")
	if s.TotalAlerts != 1 {
		t.Fatalf("expected 1 alert for backup, got %d", s.TotalAlerts)
	}
	if _, ok := s.ByJob["cleanup"]; ok {
		t.Error("cleanup job should not appear in filtered summary")
	}
}

func TestSummarize_SingleEntry_OldestEqualsNewest(t *testing.T) {
	now := time.Now()
	l := &Log{
		entries: []Entry{
			{Job: "sync", At: now, Reason: "missed"},
		},
	}

	s := l.Summarize("")
	if s.OldestAlert == nil || s.NewestAlert == nil {
		t.Fatal("expected non-nil timestamps")
	}
	if !s.OldestAlert.Equal(*s.NewestAlert) {
		t.Errorf("expected oldest == newest for single entry")
	}
}
