package watcher_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cronwatch/internal/config"
	"github.com/cronwatch/internal/notify"
	"github.com/cronwatch/internal/store"
	"github.com/cronwatch/internal/watcher"
)

func tempStore(t *testing.T) *store.Store {
	t.Helper()
	f, err := os.CreateTemp("", "watcher-store-*.json")
	if err != nil {
		t.Fatalf("create temp file: %v", err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	st, err := store.New(f.Name())
	if err != nil {
		t.Fatalf("store.New: %v", err)
	}
	return st
}

func TestCheckAll_MissedJob_SendsAlert(t *testing.T) {
	alerted := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		alerted = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := &config.Config{
		WebhookURL: ts.URL,
		Jobs: []config.Job{
			{Name: "test-job", Schedule: "* * * * *"},
		},
	}
	st := tempStore(t)
	// Record a last run 10 minutes ago so the job appears missed.
	st.Record("test-job", time.Now().Add(-10*time.Minute))

	n := notify.New(cfg.WebhookURL)
	w := watcher.New(cfg, st, n)
	w.CheckAllAt(time.Now())

	if !alerted {
		t.Error("expected alert to be sent for missed job")
	}
}

func TestCheckAll_NoLastRun_NoAlert(t *testing.T) {
	alerted := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		alerted = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := &config.Config{
		WebhookURL: ts.URL,
		Jobs: []config.Job{
			{Name: "unseen-job", Schedule: "* * * * *"},
		},
	}
	st := tempStore(t)
	n := notify.New(cfg.WebhookURL)
	w := watcher.New(cfg, st, n)
	w.CheckAllAt(time.Now())

	if alerted {
		t.Error("expected no alert when job has never run")
	}
}

func TestCheckAll_RecentRun_NoAlert(t *testing.T) {
	alerted := false
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		alerted = true
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := &config.Config{
		WebhookURL: ts.URL,
		Jobs: []config.Job{
			{Name: "fresh-job", Schedule: "@hourly"},
		},
	}
	st := tempStore(t)
	st.Record("fresh-job", time.Now().Add(-30*time.Minute))

	n := notify.New(cfg.WebhookURL)
	w := watcher.New(cfg, st, n)
	w.CheckAllAt(time.Now())

	if alerted {
		t.Error("expected no alert for recently-run job")
	}
}
