package reporter_test

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/reporter"
	"github.com/example/cronwatch/internal/store"
)

func tempStore(t *testing.T) *store.Store {
	t.Helper()
	f, err := os.CreateTemp("", "cronwatch-reporter-*.json")
	if err != nil {
		t.Fatal(err)
	}
	_ = f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	s, err := store.New(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestCollect_NoLastRun(t *testing.T) {
	s := tempStore(t)
	r := reporter.New(s, nil)
	jobs := map[string]string{"backup": "0 2 * * *"}
	statuses, err := r.Collect(jobs, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(statuses) != 1 {
		t.Fatalf("expected 1 status, got %d", len(statuses))
	}
	if statuses[0].LastRun != nil {
		t.Error("expected nil LastRun for unseen job")
	}
	if statuses[0].Missed {
		t.Error("expected Missed=false when no last run recorded")
	}
}

func TestCollect_MissedJob(t *testing.T) {
	s := tempStore(t)
	past := time.Now().Add(-48 * time.Hour)
	if err := s.Record("backup", past); err != nil {
		t.Fatal(err)
	}
	r := reporter.New(s, nil)
	jobs := map[string]string{"backup": "0 2 * * *"}
	statuses, err := r.Collect(jobs, time.Now())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !statuses[0].Missed {
		t.Error("expected Missed=true for overdue job")
	}
}

func TestCollect_InvalidCron(t *testing.T) {
	s := tempStore(t)
	r := reporter.New(s, nil)
	jobs := map[string]string{"bad": "not-a-cron"}
	_, err := r.Collect(jobs, time.Now())
	if err == nil {
		t.Error("expected error for invalid cron expression")
	}
}

func TestPrint_ContainsHeaders(t *testing.T) {
	s := tempStore(t)
	r := reporter.New(s, nil)
	jobs := map[string]string{"myjob": "*/5 * * * *"}
	statuses, _ := r.Collect(jobs, time.Now())

	var buf bytes.Buffer
	rWithBuf := reporter.New(s, &buf)
	rWithBuf.Print(statuses)

	out := buf.String()
	for _, hdr := range []string{"JOB", "LAST RUN", "NEXT RUN", "STATUS"} {
		if !strings.Contains(out, hdr) {
			t.Errorf("output missing header %q", hdr)
		}
	}
	if !strings.Contains(out, "myjob") {
		t.Error("output missing job name")
	}
	_ = json.Marshal // suppress unused import if any
}
