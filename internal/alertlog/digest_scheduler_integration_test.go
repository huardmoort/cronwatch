package alertlog_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/alertlog"
)

func TestDigestScheduler_Integration_SendsOverHTTP(t *testing.T) {
	var received atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		received.Add(1)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	path := tempPath(t)
	log, err := alertlog.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	// Record a couple of entries so the digest has something to report
	if err := log.Record("job-a", "missed run"); err != nil {
		t.Fatalf("Record: %v", err)
	}

	sender := alertlog.NewWebhookSender(server.URL)
	scheduler := alertlog.NewDigestScheduler(log, sender, 30*time.Millisecond, 1*time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	scheduler.Start(ctx)
	<-ctx.Done()
	time.Sleep(20 * time.Millisecond)

	if got := received.Load(); got < 1 {
		t.Errorf("expected at least 1 HTTP call, got %d", got)
	}
}
