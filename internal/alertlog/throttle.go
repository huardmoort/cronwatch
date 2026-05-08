package alertlog

import (
	"sync"
	"time"
)

// ThrottlePolicy controls how frequently digest notifications are sent
// for a given job to avoid alert fatigue.
type ThrottlePolicy struct {
	mu       sync.Mutex
	lastSent map[string]time.Time
	Interval time.Duration
}

// NewThrottlePolicy creates a ThrottlePolicy with the given minimum interval
// between repeated digest alerts for the same job.
func NewThrottlePolicy(interval time.Duration) *ThrottlePolicy {
	return &ThrottlePolicy{
		lastSent: make(map[string]time.Time),
		Interval: interval,
	}
}

// Allow returns true if enough time has elapsed since the last digest was
// sent for the given job, or if no digest has been sent yet.
func (t *ThrottlePolicy) Allow(job string) bool {
	if t.Interval == 0 {
		return true
	}
	t.mu.Lock()
	defer t.mu.Unlock()
	last, ok := t.lastSent[job]
	if !ok || time.Since(last) >= t.Interval {
		return true
	}
	return false
}

// Record marks the current time as the last sent time for the given job.
func (t *ThrottlePolicy) Record(job string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.lastSent[job] = time.Now()
}

// Reset clears the last-sent record for the given job, allowing the next
// call to Allow to return true immediately.
func (t *ThrottlePolicy) Reset(job string) {
	t.mu.Lock()
	defer t.mu.Unlock()
	delete(t.lastSent, job)
}
