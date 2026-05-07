package alertlog_test

import (
	"os"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/alertlog"
)

func TestRetentionPolicy_PersistsAfterApply(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "alertlog-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	l, err := alertlog.New(f.Name())
	if err != nil {
		t.Fatal(err)
	}

	now := time.Now()
	if err := l.Record("job1", "missed run", now.Add(-72*time.Hour)); err != nil {
		t.Fatal(err)
	}
	if err := l.Record("job1", "missed run", now.Add(-1*time.Hour)); err != nil {
		t.Fatal(err)
	}

	p := alertlog.RetentionPolicy{MaxAge: 48 * time.Hour}
	removed := p.Apply(l)
	if removed != 1 {
		t.Fatalf("expected 1 removed, got %d", removed)
	}

	if err := l.Save(); err != nil {
		t.Fatal(err)
	}

	// Reload and verify only 1 entry persisted.
	l2, err := alertlog.New(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	entries := l2.All()
	if len(entries) != 1 {
		t.Fatalf("expected 1 persisted entry, got %d", len(entries))
	}
}

func TestRetentionPolicy_NoOp_LeavesFileUnchanged(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "alertlog-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()

	l, err := alertlog.New(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	if err := l.Record("job1", "missed", time.Now()); err != nil {
		t.Fatal(err)
	}

	p := alertlog.RetentionPolicy{}
	removed := p.Apply(l)
	if removed != 0 {
		t.Fatalf("expected 0 removed, got %d", removed)
	}
	if len(l.All()) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(l.All()))
	}
}
