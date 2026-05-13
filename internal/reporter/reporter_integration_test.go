package reporter_test

import (
	"bytes"
	"strings"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/reporter"
)

func TestPrint_MissedStatusLabel(t *testing.T) {
	s := tempStore(t)
	past := time.Now().Add(-72 * time.Hour)
	if err := s.Record("nightly", past); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	r := reporter.New(s, &buf)
	jobs := map[string]string{"nightly": "0 3 * * *"}
	statuses, err := r.Collect(jobs, time.Now())
	if err != nil {
		t.Fatalf("collect: %v", err)
	}
	r.Print(statuses)

	if !strings.Contains(buf.String(), "MISSED") {
		t.Error("expected MISSED label in output for overdue job")
	}
}

func TestPrint_OKStatusLabel(t *testing.T) {
	s := tempStore(t)
	recent := time.Now().Add(-1 * time.Minute)
	if err := s.Record("frequent", recent); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	r := reporter.New(s, &buf)
	jobs := map[string]string{"frequent": "*/5 * * * *"}
	statuses, err := r.Collect(jobs, time.Now())
	if err != nil {
		t.Fatalf("collect: %v", err)
	}
	r.Print(statuses)

	if !strings.Contains(buf.String(), "OK") {
		t.Error("expected OK label in output for healthy job")
	}
}

func TestPrint_JobNameAppearsInOutput(t *testing.T) {
	s := tempStore(t)
	recent := time.Now().Add(-1 * time.Minute)
	if err := s.Record("my-unique-job", recent); err != nil {
		t.Fatal(err)
	}

	var buf bytes.Buffer
	r := reporter.New(s, &buf)
	jobs := map[string]string{"my-unique-job": "*/5 * * * *"}
	statuses, err := r.Collect(jobs, time.Now())
	if err != nil {
		t.Fatalf("collect: %v", err)
	}
	r.Print(statuses)

	if !strings.Contains(buf.String(), "my-unique-job") {
		t.Error("expected job name to appear in output")
	}
}
