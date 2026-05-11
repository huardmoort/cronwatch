package alertlog_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/alertlog"
)

type stubDigestSender struct {
	calls atomic.Int32
}

func (s *stubDigestSender) Send(d alertlog.Digest) error {
	s.calls.Add(1)
	return nil
}

func TestDigestScheduler_FiresOnInterval(t *testing.T) {
	path := tempPath(t)
	log, err := alertlog.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	sender := &stubDigestSender{}
	scheduler := alertlog.NewDigestScheduler(log, sender, 30*time.Millisecond, 1*time.Hour)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Millisecond)
	defer cancel()

	scheduler.Start(ctx)

	// Allow at least 2 ticks
	time.Sleep(90 * time.Millisecond)
	cancel()
	time.Sleep(20 * time.Millisecond)

	if got := sender.calls.Load(); got < 2 {
		t.Errorf("expected at least 2 digest sends, got %d", got)
	}
}

func TestDigestScheduler_StopsOnContextCancel(t *testing.T) {
	path := tempPath(t)
	log, err := alertlog.New(path)
	if err != nil {
		t.Fatalf("New: %v", err)
	}

	sender := &stubDigestSender{}
	scheduler := alertlog.NewDigestScheduler(log, sender, 20*time.Millisecond, 1*time.Hour)

	ctx, cancel := context.WithCancel(context.Background())
	scheduler.Start(ctx)
	time.Sleep(30 * time.Millisecond)
	cancel()
	time.Sleep(10 * time.Millisecond)

	before := sender.calls.Load()
	time.Sleep(50 * time.Millisecond)
	after := sender.calls.Load()

	if after != before {
		t.Errorf("scheduler kept firing after cancel: before=%d after=%d", before, after)
	}
}
