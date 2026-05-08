package alertlog_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/alertlog"
)

func TestDigestWebhook_FullLifecycle(t *testing.T) {
	var callCount int32
	var lastPayload struct {
		WindowHours int                      `json:"window_hours"`
		Jobs        []alertlog.DigestJobStat `json:"jobs"`
	}

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&callCount, 1)
		json.NewDecoder(r.Body).Decode(&lastPayload) //nolint:errcheck
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	path := filepath.Join(t.TempDir(), "alerts.json")
	log, err := alertlog.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	now := time.Now().UTC()
	for i := 0; i < 4; i++ {
		log.Record("nightly-backup", "missed", now.Add(-time.Duration(i)*time.Hour))
	}
	log.Record("weekly-report", "failed", now.Add(-2*time.Hour))

	// Persist so Digest reads from file.
	if err := os.WriteFile(path, mustMarshal(t, log), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}

	log2, _ := alertlog.New(path)
	d := alertlog.Digest(log2.Entries(), 48, now)

	sender := alertlog.NewWebhookSender(ts.URL, nil)
	if err := sender.Send(d); err != nil {
		t.Fatalf("Send: %v", err)
	}

	if atomic.LoadInt32(&callCount) != 1 {
		t.Errorf("expected 1 webhook call, got %d", callCount)
	}
	if lastPayload.WindowHours != 48 {
		t.Errorf("window_hours: got %d, want 48", lastPayload.WindowHours)
	}
	if len(lastPayload.Jobs) < 2 {
		t.Errorf("expected at least 2 job entries, got %d", len(lastPayload.Jobs))
	}
}

func mustMarshal(t *testing.T, v interface{}) []byte {
	t.Helper()
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("marshal: %v", err)
	}
	return b
}
