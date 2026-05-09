package alertlog

import (
	"context"
	"log"
	"time"
)

// DigestScheduler periodically computes and dispatches alert digests via a
// configured sender. It is designed to run as a background goroutine and
// respects context cancellation for clean shutdown.
type DigestScheduler struct {
	log      *Log
	sender   DigestSender
	window   time.Duration
	interval time.Duration
}

// DigestSender is the interface satisfied by any type that can deliver a
// formatted digest string to an external destination (e.g. a webhook).
type DigestSender interface {
	Send(ctx context.Context, body string) error
}

// NewDigestScheduler creates a DigestScheduler that reads from log, formats
// digests over window, and dispatches them every interval via sender.
//
// A typical production setup uses a 24-hour window with a 24-hour interval so
// that operators receive one daily summary of all alert activity.
func NewDigestScheduler(l *Log, sender DigestSender, window, interval time.Duration) *DigestScheduler {
	return &DigestScheduler{
		log:      l,
		sender:   sender,
		window:   window,
		interval: interval,
	}
}

// Run blocks until ctx is cancelled, dispatching a digest on every tick of the
// configured interval. The first dispatch occurs after one full interval has
// elapsed, not immediately on start.
func (ds *DigestScheduler) Run(ctx context.Context) {
	if ds.interval <= 0 {
		log.Println("[digest_scheduler] interval is zero or negative — scheduler disabled")
		return
	}

	ticker := time.NewTicker(ds.interval)
	defer ticker.Stop()

	log.Printf("[digest_scheduler] started — window=%s interval=%s", ds.window, ds.interval)

	for {
		select {
		case <-ctx.Done():
			log.Println("[digest_scheduler] stopping")
			return
		case t := <-ticker.C:
			ds.dispatch(ctx, t)
		}
	}
}

// dispatch computes the current digest and sends it. Errors are logged but do
// not stop the scheduler; the next tick will attempt delivery again.
func (ds *DigestScheduler) dispatch(ctx context.Context, at time.Time) {
	entries := ds.log.Entries()
	d := Digest(entries, at, ds.window)
	body := FormatDigest(d)

	if err := ds.sender.Send(ctx, body); err != nil {
		log.Printf("[digest_scheduler] send error: %v", err)
		return
	}

	log.Printf("[digest_scheduler] digest dispatched at %s (jobs=%d total_alerts=%d)",
		at.Format(time.RFC3339), len(d.Jobs), d.TotalAlerts)
}
