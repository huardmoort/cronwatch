package alertlog

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestWebhookSender_Send_Success(t *testing.T) {
	var received WebhookPayload
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Fatalf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	sender := NewWebhookSender(ts.URL, nil)
	d := Digest{
		GeneratedAt: time.Now().UTC(),
		WindowHours: 24,
		Jobs: []DigestJobStat{{Job: "backup", Count: 3}},
	}

	if err := sender.Send(d); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if received.WindowHours != 24 {
		t.Errorf("window_hours: got %d, want 24", received.WindowHours)
	}
	if len(received.Jobs) != 1 || received.Jobs[0].Job != "backup" {
		t.Errorf("unexpected jobs: %+v", received.Jobs)
	}
}

func TestWebhookSender_Send_EmptyEndpoint(t *testing.T) {
	sender := NewWebhookSender("", nil)
	err := sender.Send(Digest{})
	if err != nil {
		t.Fatalf("expected no error for empty endpoint, got: %v", err)
	}
}

func TestWebhookSender_Send_NonSuccessStatus(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	sender := NewWebhookSender(ts.URL, nil)
	err := sender.Send(Digest{Jobs: []DigestJobStat{}})
	if err == nil {
		t.Fatal("expected error for 500 response")
	}
}
