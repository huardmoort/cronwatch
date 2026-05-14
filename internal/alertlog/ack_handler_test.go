package alertlog_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/alertlog"
)

func newAckServer(t *testing.T) (*alertlog.AckPolicy, *httptest.Server) {
	t.Helper()
	p := alertlog.NewAckPolicy()
	return p, httptest.NewServer(alertlog.NewAckHandler(p))
}

func TestAckHandler_GET_NotAcked(t *testing.T) {
	_, srv := newAckServer(t)
	defer srv.Close()
	resp, err := http.Get(srv.URL + "/ack?job=backup")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	var body map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&body) //nolint:errcheck
	if body["acked"].(bool) {
		t.Fatal("expected not acked")
	}
}

func TestAckHandler_POST_AckJob(t *testing.T) {
	p, srv := newAckServer(t)
	defer srv.Close()
	payload, _ := json.Marshal(map[string]string{"job": "backup", "duration": "10m"})
	resp, err := http.Post(srv.URL+"/ack", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
	if !p.IsAcked("backup") {
		t.Fatal("expected job to be acked")
	}
}

func TestAckHandler_DELETE_ClearsAck(t *testing.T) {
	p, srv := newAckServer(t)
	defer srv.Close()
	p.Acknowledge("backup", 10*time.Minute)
	req, _ := http.NewRequest(http.MethodDelete, srv.URL+"/ack?job=backup", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent {
		t.Fatalf("expected 204, got %d", resp.StatusCode)
	}
	if p.IsAcked("backup") {
		t.Fatal("expected ack to be cleared")
	}
}

func TestAckHandler_POST_MissingJob(t *testing.T) {
	_, srv := newAckServer(t)
	defer srv.Close()
	payload, _ := json.Marshal(map[string]string{"duration": "5m"})
	resp, err := http.Post(srv.URL+"/ack", "application/json", bytes.NewReader(payload))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", resp.StatusCode)
	}
}

func TestAckHandler_MethodNotAllowed(t *testing.T) {
	_, srv := newAckServer(t)
	defer srv.Close()
	req, _ := http.NewRequest(http.MethodPatch, srv.URL+"/ack", nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected 405, got %d", resp.StatusCode)
	}
}
