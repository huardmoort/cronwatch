package watcher

import (
	"log"
	"time"

	"github.com/cronwatch/internal/config"
	"github.com/cronwatch/internal/notify"
	"github.com/cronwatch/internal/schedule"
	"github.com/cronwatch/internal/store"
)

// Watcher monitors cron jobs for missed runs and alerts on failures.
type Watcher struct {
	cfg      *config.Config
	store    *store.Store
	notifier *notify.Notifier
	stop     chan struct{}
}

// New creates a new Watcher instance.
func New(cfg *config.Config, st *store.Store, notifier *notify.Notifier) *Watcher {
	return &Watcher{
		cfg:      cfg,
		store:    st,
		notifier: notifier,
		stop:     make(chan struct{}),
	}
}

// Start begins the watch loop, checking jobs at the configured interval.
func (w *Watcher) Start(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	log.Printf("watcher: started, check interval=%s", interval)
	for {
		select {
		case <-ticker.C:
			w.checkAll(time.Now())
		case <-w.stop:
			log.Println("watcher: stopped")
			return
		}
	}
}

// Stop signals the watch loop to exit.
func (w *Watcher) Stop() {
	close(w.stop)
}

// checkAll iterates over all configured jobs and checks for missed runs.
func (w *Watcher) checkAll(now time.Time) {
	for _, job := range w.cfg.Jobs {
		last, ok := w.store.LastRun(job.Name)
		if !ok {
			log.Printf("watcher: no run recorded for job %q, skipping", job.Name)
			continue
		}
		missed, err := schedule.IsMissed(job.Schedule, last, now)
		if err != nil {
			log.Printf("watcher: invalid schedule for job %q: %v", job.Name, err)
			continue
		}
		if missed {
			msg := "job " + job.Name + " missed its scheduled run"
			log.Printf("watcher: ALERT — %s", msg)
			if err := w.notifier.Send(msg); err != nil {
				log.Printf("watcher: failed to send alert for job %q: %v", job.Name, err)
			}
		}
	}
}
