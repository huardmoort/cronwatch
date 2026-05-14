package alertlog_test

import (
	"sync"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/alertlog"
)

func TestSnooze_ConcurrentAccess(t *testing.T) {
	s := alertlog.NewSnoozePolicy()
	const workers = 20
	var wg sync.WaitGroup
	wg.Add(workers)
	for i := 0; i < workers; i++ {
		go func() {
			defer wg.Done()
			s.Snooze("job", 50*time.Millisecond)
			_ = s.IsSnoozed("job")
			s.Cancel("job")
		}()
	}
	wg.Wait()
}

func TestSnooze_MultipleJobsConcurrent(t *testing.T) {
	s := alertlog.NewSnoozePolicy()
	jobs := []string{"alpha", "beta", "gamma", "delta"}
	var wg sync.WaitGroup
	for _, job := range jobs {
		wg.Add(1)
		go func(j string) {
			defer wg.Done()
			s.Snooze(j, 10*time.Minute)
			if !s.IsSnoozed(j) {
				t.Errorf("expected %s to be snoozed", j)
			}
		}(job)
	}
	wg.Wait()

	active := s.ActiveJobs()
	if len(active) != len(jobs) {
		t.Fatalf("expected %d active jobs, got %d", len(jobs), len(active))
	}
}
