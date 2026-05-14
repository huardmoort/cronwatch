package alertlog_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/alertlog"
)

func newSnoozeServer(t *testing.T) (*alertlog.SnoozePolicy, *httptest.Server) {
	t.Helper()
	p := alertlog.NewSnoozePolicy()
	srv := httptest.NewServer(alertlog.NewSnoozeHandler(p))
	t.Cleanup(srv.Close)
	return p, srv
}

func TestSnoozeHandler_GET_ReturnsActiveSnoozed(t *testing.T) {
	p, srv := newSnoozeServer(t)
	p.Snooze("backup", 10*time.Minute)

	resp, err := http.Get(srv.URL + "/snooze")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	var body map[string][]string
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatal(err)
	}
	if len(body["snoozed"]) != 1 || body["snoozed"][0] != "backup" {
		t.Fatalf("unexpected snoozed jobs: %v", body["snoozed"])
	}
}

func TestSnoozeHandler_POST_SnoozeJob(t *testing.T) {
	p, srv := newSnoozeServer(t)

	resp, err := http.Post(srv.URL+"/snooze?job=nightly&minutes=30", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
	if !p.IsSnoozed("nightly") {
		t.Fatal("expected nightly to be snoozed")
	}
}

func TestSnoozeHandler_POST_MissingJob(t *testing.T) {
	_, srv := newSnoozeServer(t)
	resp, err := http.Post(srv.URL+"/snooze?minutes=5", "", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestSnoozeHandler_DELETE_CancelsSnooze(t *testing.T) {
	p, srv := newSnoozeServer(t)
	p.Snooze("cleanup", 10*time.Minute)

	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/snooze?job=cleanup", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
	if p.IsSnoozed("cleanup") {
		t.Fatal("expected snooze to be cancelled")
	}
}

func TestSnoozeHandler_MethodNotAllowed(t *testing.T) {
	_, srv := newSnoozeServer(t)
	req, _ := http.NewRequest(http.MethodPut, srv.URL+"/snooze", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
}
