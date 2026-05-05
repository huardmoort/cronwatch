package watcher_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/config"
	"github.com/user/cronwatch/internal/notify"
	"github.com/user/cronwatch/internal/store"
	"github.com/user/cronwatch/internal/watcher"
)

func TestRunLoop_CallsCheckAllRepeatedly(t *testing.T) {
	var alertCount int32

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&alertCount, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	cfg := &config.Config{
		Jobs: []config.Job{
			{Name: "tick-job", Schedule: "* * * * *", MaxDelay: "1s"},
		},
		Notify: config.Notify{Webhook: ts.URL},
	}

	st, _ := store.New(tempStore(t))
	// Record a run far in the past so the job appears missed.
	_ = st.Record("tick-job", time.Now().Add(-10*time.Minute))

	n := notify.New(cfg)
	w := watcher.New(cfg, st, n)

	ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
	defer cancel()

	// Use a very short interval so multiple ticks fire within the timeout.
	watcher.RunLoop(ctx, w, 80*time.Millisecond)

	// Expect at least 2 CheckAll calls (immediate + >=1 tick), each producing an alert.
	if got := atomic.LoadInt32(&alertCount); got < 2 {
		t.Errorf("expected at least 2 alerts, got %d", got)
	}
}

func TestRunLoop_StopsOnContextCancel(t *testing.T) {
	cfg := &config.Config{
		Jobs: []config.Job{},
		Notify: config.Notify{},
	}
	st, _ := store.New(tempStore(t))
	n := notify.New(cfg)
	w := watcher.New(cfg, st, n)

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel immediately

	done := make(chan struct{})
	go func() {
		watcher.RunLoop(ctx, w, time.Second)
		close(done)
	}()

	select {
	case <-done:
		// success
	case <-time.After(500 * time.Millisecond):
		t.Fatal("RunLoop did not stop after context cancellation")
	}
}
