package alertlog

import (
	"testing"
	"time"
)

func TestRetentionPolicy_ZeroIsNoop(t *testing.T) {
	l := &Log{}
	l.entries = []Entry{
		{Job: "job1", At: time.Now().Add(-48 * time.Hour)},
		{Job: "job2", At: time.Now().Add(-1 * time.Hour)},
	}
	p := RetentionPolicy{}
	removed := p.Apply(l)
	if removed != 0 {
		t.Fatalf("expected 0 removed, got %d", removed)
	}
	if len(l.entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(l.entries))
	}
}

func TestRetentionPolicy_MaxAge_RemovesOld(t *testing.T) {
	l := &Log{}
	now := time.Now()
	l.entries = []Entry{
		{Job: "job1", At: now.Add(-72 * time.Hour)},
		{Job: "job1", At: now.Add(-1 * time.Hour)},
		{Job: "job2", At: now.Add(-24 * time.Hour)},
	}
	p := RetentionPolicy{MaxAge: 48 * time.Hour}
	removed := p.Apply(l)
	if removed != 1 {
		t.Fatalf("expected 1 removed, got %d", removed)
	}
	if len(l.entries) != 2 {
		t.Fatalf("expected 2 entries remaining, got %d", len(l.entries))
	}
}

func TestRetentionPolicy_MaxEntries_KeepsMostRecent(t *testing.T) {
	l := &Log{}
	now := time.Now()
	l.entries = []Entry{
		{Job: "job1", At: now.Add(-3 * time.Hour)},
		{Job: "job1", At: now.Add(-2 * time.Hour)},
		{Job: "job1", At: now.Add(-1 * time.Hour)},
	}
	p := RetentionPolicy{MaxEntries: 2}
	removed := p.Apply(l)
	if removed != 1 {
		t.Fatalf("expected 1 removed, got %d", removed)
	}
	if len(l.entries) != 2 {
		t.Fatalf("expected 2 entries remaining, got %d", len(l.entries))
	}
}

func TestRetentionPolicy_Combined(t *testing.T) {
	l := &Log{}
	now := time.Now()
	l.entries = []Entry{
		{Job: "job1", At: now.Add(-96 * time.Hour)}, // too old
		{Job: "job1", At: now.Add(-10 * time.Hour)}, // kept but over limit
		{Job: "job1", At: now.Add(-5 * time.Hour)},  // kept
		{Job: "job1", At: now.Add(-1 * time.Hour)},  // kept
	}
	p := RetentionPolicy{MaxAge: 48 * time.Hour, MaxEntries: 2}
	removed := p.Apply(l)
	if removed != 2 {
		t.Fatalf("expected 2 removed, got %d", removed)
	}
	if len(l.entries) != 2 {
		t.Fatalf("expected 2 entries remaining, got %d", len(l.entries))
	}
}

func TestRetentionPolicy_MultipleJobs_Independent(t *testing.T) {
	l := &Log{}
	now := time.Now()
	l.entries = []Entry{
		{Job: "a", At: now.Add(-3 * time.Hour)},
		{Job: "a", At: now.Add(-2 * time.Hour)},
		{Job: "a", At: now.Add(-1 * time.Hour)},
		{Job: "b", At: now.Add(-1 * time.Hour)},
	}
	p := RetentionPolicy{MaxEntries: 2}
	removed := p.Apply(l)
	if removed != 1 {
		t.Fatalf("expected 1 removed, got %d", removed)
	}
	if len(l.entries) != 3 {
		t.Fatalf("expected 3 entries remaining, got %d", len(l.entries))
	}
}
