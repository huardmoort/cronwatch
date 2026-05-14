package alertlog_test

import (
	"sync"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/alertlog"
)

func TestAck_ConcurrentAccess(t *testing.T) {
	p := alertlog.NewAckPolicy()
	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			p.Acknowledge("job", 10*time.Minute)
			_ = p.IsAcked("job")
			p.Clear("job")
		}()
	}
	wg.Wait()
}

func TestAck_MultipleJobs_Concurrent(t *testing.T) {
	p := alertlog.NewAckPolicy()
	jobs := []string{"alpha", "beta", "gamma", "delta"}
	var wg sync.WaitGroup
	for _, job := range jobs {
		wg.Add(1)
		go func(j string) {
			defer wg.Done()
			p.Acknowledge(j, 5*time.Minute)
			if !p.IsAcked(j) {
				t.Errorf("expected %s to be acked", j)
			}
		}(job)
	}
	wg.Wait()
	all := p.AckedJobs()
	if len(all) != len(jobs) {
		t.Fatalf("expected %d acked jobs, got %d", len(jobs), len(all))
	}
}
