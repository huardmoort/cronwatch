package alertlog

import (
	"strings"
	"testing"
	"time"
)

func makeLogWithEntries(t *testing.T, entries []Entry) *Log {
	t.Helper()
	p := tempPath(t)
	l, err := New(p)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	l.mu.Lock()
	l.entries = entries
	l.mu.Unlock()
	return l
}

func TestDigest_EmptyLog(t *testing.T) {
	l := makeLogWithEntries(t, nil)
	got := Digest(l, time.Hour)
	if len(got) != 0 {
		t.Fatalf("expected empty digest, got %d entries", len(got))
	}
}

func TestDigest_CountsWithinWindow(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Job: "backup", At: now.Add(-10 * time.Minute), Message: "missed"},
		{Job: "backup", At: now.Add(-20 * time.Minute), Message: "missed"},
		{Job: "backup", At: now.Add(-2 * time.Hour), Message: "old"},
	}
	l := makeLogWithEntries(t, entries)
	got := Digest(l, time.Hour)
	if len(got) != 1 {
		t.Fatalf("expected 1 digest entry, got %d", len(got))
	}
	if got[0].Count != 2 {
		t.Errorf("expected count 2, got %d", got[0].Count)
	}
	if got[0].Job != "backup" {
		t.Errorf("unexpected job name: %s", got[0].Job)
	}
}

func TestDigest_MultipleJobs(t *testing.T) {
	now := time.Now()
	entries := []Entry{
		{Job: "alpha", At: now.Add(-5 * time.Minute), Message: "m"},
		{Job: "beta", At: now.Add(-5 * time.Minute), Message: "m"},
		{Job: "beta", At: now.Add(-15 * time.Minute), Message: "m"},
	}
	l := makeLogWithEntries(t, entries)
	got := Digest(l, time.Hour)
	if len(got) != 2 {
		t.Fatalf("expected 2 digest entries, got %d", len(got))
	}
	counts := map[string]int{}
	for _, e := range got {
		counts[e.Job] = e.Count
	}
	if counts["alpha"] != 1 {
		t.Errorf("alpha: expected 1, got %d", counts["alpha"])
	}
	if counts["beta"] != 2 {
		t.Errorf("beta: expected 2, got %d", counts["beta"])
	}
}

func TestFormatDigest_NoAlerts(t *testing.T) {
	out := FormatDigest(nil, time.Hour)
	if !strings.Contains(out, "No alerts") {
		t.Errorf("expected 'No alerts' message, got: %s", out)
	}
}

func TestFormatDigest_ContainsJobName(t *testing.T) {
	now := time.Now()
	entries := []DigestEntry{
		{Job: "myjob", Count: 3, FirstAlert: now.Add(-30 * time.Minute), LastAlert: now},
	}
	out := FormatDigest(entries, time.Hour)
	if !strings.Contains(out, "myjob") {
		t.Errorf("expected job name in output, got: %s", out)
	}
	if !strings.Contains(out, "3") {
		t.Errorf("expected count in output, got: %s", out)
	}
}
