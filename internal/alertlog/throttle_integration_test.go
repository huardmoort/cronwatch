package alertlog_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/cronwatch/cronwatch/internal/alertlog"
)

func TestThrottle_ConcurrentAccess(t *testing.T) {
	p := alertlog.NewThrottlePolicy(50 * time.Millisecond)
	var allowed atomic.Int32
	var wg sync.WaitGroup

	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if p.Allow("shared-job") {
				allowed.Add(1)
				p.Record("shared-job")
			}
		}()
	}
	wg.Wait()

	// At least one goroutine must have been allowed.
	if allowed.Load() == 0 {
		t.Fatal("expected at least one goroutine to be allowed")
	}
}

func TestThrottle_MultipleJobsConcurrent(t *testing.T) {
	p := alertlog.NewThrottlePolicy(5 * time.Minute)
	jobs := []string{"job-1", "job-2", "job-3", "job-4"}
	results := make(map[string]bool)
	var mu sync.Mutex
	var wg sync.WaitGroup

	for _, job := range jobs {
		wg.Add(1)
		go func(j string) {
			defer wg.Done()
			allowed := p.Allow(j)
			if allowed {
				p.Record(j)
			}
			mu.Lock()
			results[j] = allowed
			mu.Unlock()
		}(job)
	}
	wg.Wait()

	for _, job := range jobs {
		if !results[job] {
			t.Errorf("expected job %q to be allowed on first call", job)
		}
	}
}
