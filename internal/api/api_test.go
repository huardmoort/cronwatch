package api_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/cronwatch/internal/alertlog"
	"github.com/cronwatch/internal/api"
	"github.com/cronwatch/internal/metrics"
)

func newTestServer(t *testing.T) *api.Server {
	t.Helper()
	mc := metrics.New()
	f, err := os.CreateTemp("", "api-alertlog-*.json")
	if err != nil {
		t.Fatal(err)
	}
	f.Close()
	t.Cleanup(func() { os.Remove(f.Name()) })
	al, err := alertlog.New(f.Name())
	if err != nil {
		t.Fatal(err)
	}
	return api.New(mc, al)
}

func TestHealthz_ReturnsOK(t *testing.T) {
	srv := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var body map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if body["status"] != "ok" {
		t.Errorf("expected status ok, got %s", body["status"])
	}
}

func TestStatus_ReturnsJSON(t *testing.T) {
	srv := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/status", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Errorf("expected application/json, got %s", ct)
	}
}

func TestStatus_MethodNotAllowed(t *testing.T) {
	srv := newTestServer(t)
	req := httptest.NewRequest(http.MethodPost, "/status", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", rr.Code)
	}
}

func TestMetrics_ReturnsSnapshot(t *testing.T) {
	srv := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/metrics", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rr.Code)
	}
	var snap map[string]interface{}
	if err := json.NewDecoder(rr.Body).Decode(&snap); err != nil {
		t.Fatal(err)
	}
}

func TestAlerts_RouteRegistered(t *testing.T) {
	srv := newTestServer(t)
	req := httptest.NewRequest(http.MethodGet, "/alerts", nil)
	rr := httptest.NewRecorder()
	srv.ServeHTTP(rr, req)
	if rr.Code != http.StatusOK {
		t.Fatalf("expected 200 from /alerts, got %d", rr.Code)
	}
}
