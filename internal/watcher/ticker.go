package watcher

import (
	"context"
	"log"
	"time"
)

// RunLoop starts a blocking loop that calls CheckAll on the given Watcher
// at the specified interval until ctx is cancelled.
func RunLoop(ctx context.Context, w *Watcher, interval time.Duration) {
	if interval <= 0 {
		interval = time.Minute
	}

	log.Printf("[cronwatch] starting watcher loop (interval=%s)", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Run once immediately before waiting for the first tick.
	if err := w.CheckAll(); err != nil {
		log.Printf("[cronwatch] CheckAll error: %v", err)
	}

	for {
		select {
		case <-ctx.Done():
			log.Println("[cronwatch] watcher loop stopped")
			return
		case <-ticker.C:
			if err := w.CheckAll(); err != nil {
				log.Printf("[cronwatch] CheckAll error: %v", err)
			}
		}
	}
}
