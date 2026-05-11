package alertlog_test

import (
	"sync"
	"testing"
	"time"

	"github.com/example/cronwatch/internal/alertlog"
)

// TestMute_ConcurrentAccess verifies MutePolicy is safe under concurrent use.
func TestMute_ConcurrentAccess(t *testing.T) {
	p := alertlog.NewMutePolicy()
	const goroutines = 20
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(i int) {
			defer wg.Done()
			job := "job-concurrent"
			if i%3 == 0 {
				p.Mute(job, 10*time.Second)
			} else if i%3 == 1 {
				p.IsMuted(job)
			} else {
				p.Unmute(job)
			}
		}(i)
	}

	wg.Wait() // must not race or panic
}

// TestMute_MultipleJobsConcurrent ensures independent jobs do not interfere.
func TestMute_MultipleJobsConcurrent(t *testing.T) {
	p := alertlog.NewMutePolicy()
	jobs := []string{"alpha", "beta", "gamma", "delta"}
	var wg sync.WaitGroup

	for _, job := range jobs {
		wg.Add(1)
		go func(j string) {
			defer wg.Done()
			p.Mute(j, 5*time.Minute)
		}(job)
	}
	wg.Wait()

	active := p.ActiveMutes()
	if len(active) != len(jobs) {
		t.Fatalf("expected %d active mutes, got %d", len(jobs), len(active))
	}
	for _, job := range jobs {
		if _, ok := active[job]; !ok {
			t.Errorf("expected %s in active mutes", job)
		}
	}
}
