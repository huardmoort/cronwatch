package notify_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/notify"
)

func TestSend_Success(t *testing.T) {
	var received map[string]string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&received); err != nil {
			t.Errorf("decode body: %v", err)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	n := notify.New(srv.URL)
	evt := notify.Event{
		JobName:    "backup",
		Kind:       notify.KindFailure,
		Message:    "exit code 1",
		OccurredAt: time.Now(),
	}

	if err := n.Send(evt); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if received["job"] != "backup" {
		t.Errorf("expected job=backup, got %q", received["job"])
	}
	if received["kind"] != string(notify.KindFailure) {
		t.Errorf("expected kind=failure, got %q", received["kind"])
	}
}

func TestSend_NonSuccessStatus(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	n := notify.New(srv.URL)
	err := n.Send(notify.Event{
		JobName:    "cleanup",
		Kind:       notify.KindMissed,
		Message:    "missed window",
		OccurredAt: time.Now(),
	})

	if err == nil {
		t.Fatal("expected error for non-2xx status, got nil")
	}
}

func TestSend_EmptyWebhook(t *testing.T) {
	n := notify.New("")
	err := n.Send(notify.Event{
		JobName:    "noop",
		Kind:       notify.KindFailure,
		Message:    "irrelevant",
		OccurredAt: time.Now(),
	})
	if err != nil {
		t.Fatalf("expected no error for empty webhook, got %v", err)
	}
}
