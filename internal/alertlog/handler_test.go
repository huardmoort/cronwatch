package alertlog_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/cronwatch/internal/alertlog"
)

func newTestLog(t *testing.T) (*alertlog.Log, string) {
	t.Helper()
	f, err := os.CreateTemp("", "alertlog-handler-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	l, err := alertlog.New(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	return l, f.Name()
}

func TestHandler_ReturnsAllEntries(t *testing.T) {
	l, _ := newTestLog(t)
	_ = l.Record("backup", "missed deadline")
	_ = l.Record("cleanup", "missed deadline")

	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	rr := httptest.NewRecorder()
	l.Handler()(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var entries []map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
}

func TestHandler_FiltersByJob(t *testing.T) {
	l, _ := newTestLog(t)
	_ = l.Record("backup", "missed")
	_ = l.Record("cleanup", "missed")
	_ = l.Record("backup", "missed again")

	req := httptest.NewRequest(http.MethodGet, "/alerts?job=backup", nil)
	rr := httptest.NewRecorder()
	l.Handler()(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var entries []map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatal(err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 backup entries, got %d", len(entries))
	}
	for _, e := range entries {
		if e["job"] != "backup" {
			t.Errorf("unexpected job in filtered result: %v", e["job"])
		}
	}
}

func TestHandler_MethodNotAllowed(t *testing.T) {
	l, _ := newTestLog(t)
	req := httptest.NewRequest(http.MethodPost, "/alerts", nil)
	rr := httptest.NewRecorder()
	l.Handler()(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}

func TestHandler_EmptyLog(t *testing.T) {
	l, _ := newTestLog(t)
	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	rr := httptest.NewRecorder()
	l.Handler()(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var entries []interface{}
	if err := json.NewDecoder(rr.Body).Decode(&entries); err != nil {
		t.Fatal(err)
	}
	if len(entries) != 0 {
		t.Fatalf("expected empty list, got %d entries", len(entries))
	}
}

func TestHandler_ContentType(t *testing.T) {
	l, _ := newTestLog(t)
	_ = l.Record("myjob", "reason")
	_ = time.Now() // suppress unused import
	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	rr := httptest.NewRecorder()
	l.Handler()(rr, req)
	ct := rr.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}
