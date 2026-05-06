package ratelimit_test

import (
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/ratelimit"
)

// TestConcurrentAllow verifies that the Limiter is safe under concurrent
// access from multiple goroutines alerting on the same job.
func TestConcurrentAllow(t *testing.T) {
	l := ratelimit.New(10 * time.Minute)

	var allowed atomic.Int64
	var wg sync.WaitGroup

	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if l.Allow("shared-job") {
				allowed.Add(1)
			}
		}()
	}
	wg.Wait()

	// Only one goroutine should have been allowed through.
	if got := allowed.Load(); got != 1 {
		t.Fatalf("expected exactly 1 allowed call, got %d", got)
	}
}

// TestMultipleJobsConcurrent ensures independent jobs are each allowed once.
func TestMultipleJobsConcurrent(t *testing.T) {
	l := ratelimit.New(10 * time.Minute)
	jobs := []string{"job-1", "job-2", "job-3", "job-4", "job-5"}

	var wg sync.WaitGroup
	var allowed atomic.Int64

	for _, name := range jobs {
		name := name
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if l.Allow(name) {
					allowed.Add(1)
				}
			}()
		}
	}
	wg.Wait()

	if got := allowed.Load(); got != int64(len(jobs)) {
		t.Fatalf("expected %d allowed calls (one per job), got %d", len(jobs), got)
	}
}
