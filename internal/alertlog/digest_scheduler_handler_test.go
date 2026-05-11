package alertlog_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/alertlog"
)

func TestDigestHandler_ReturnsJSON(t *testing.T) {
	path := tempPath(t)
	log, err := alertlog.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := log.Record("job-x", "missed"); err != nil {
		t.Fatalf("Record: %v", err)
	}

	h := alertlog.NewDigestHandler(alertlog.DigestHandlerConfig{
		Log:    log,
		Window: 1 * time.Hour,
	})

	req := httptest.NewRequest(http.MethodGet, "/digest", nil)
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var d map[string]interface{}
	if err := json.NewDecoder(rec.Body).Decode(&d); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if _, ok := d["total_alerts"]; !ok {
		t.Error("response missing total_alerts field")
	}
}

func TestDigestHandler_MethodNotAllowed(t *testing.T) {
	path := tempPath(t)
	log, err := alertlog.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	h := alertlog.NewDigestHandler(alertlog.DigestHandlerConfig{
		Log:    log,
		Window: 1 * time.Hour,
	})

	req := httptest.NewRequest(http.MethodPost, "/digest", nil)
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusMethodNotAllowed {
		t.Errorf("expected 405, got %d", rec.Code)
	}
}

func TestDigestHandler_DefaultWindow(t *testing.T) {
	path := tempPath(t)
	log, err := alertlog.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// Window of 0 should default to 24h without panicking
	h := alertlog.NewDigestHandler(alertlog.DigestHandlerConfig{
		Log:    log,
		Window: 0,
	})

	req := httptest.NewRequest(http.MethodGet, "/digest", nil)
	rec := httptest.NewRecorder()
	h(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("expected 200, got %d", rec.Code)
	}
}
