package alertlog_test

import (
	"sync"
	"testing"
	"time"

	"github.com/user/cronwatch/internal/alertlog"
)

func TestEscalation_ConcurrentAccess(t *testing.T) {
	policy := alertlog.NewEscalationPolicy(50, time.Minute)
	now := time.Now()
	var wg sync.WaitGroup

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			policy.Record("concurrent-job", now.Add(time.Duration(i)*time.Millisecond))
		}(i)
	}
	wg.Wait()

	count := policy.Count("concurrent-job", now.Add(time.Minute))
	if count != 100 {
		t.Fatalf("expected 100 recorded entries, got %d", count)
	}
}

func TestEscalation_MultipleJobs_Concurrent(t *testing.T) {
	policy := alertlog.NewEscalationPolicy(3, time.Minute)
	now := time.Now()
	jobs := []string{"alpha", "beta", "gamma"}
	var wg sync.WaitGroup

	for _, job := range jobs {
		for i := 0; i < 5; i++ {
			wg.Add(1)
			go func(j string, idx int) {
				defer wg.Done()
				policy.Record(j, now.Add(time.Duration(idx)*time.Millisecond))
			}(job, i)
		}
	}
	wg.Wait()

	for _, job := range jobs {
		count := policy.Count(job, now.Add(time.Minute))
		if count != 5 {
			t.Errorf("job %s: expected 5 entries, got %d", job, count)
		}
	}
}
