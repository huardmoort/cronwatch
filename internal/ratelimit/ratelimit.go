// Package ratelimit provides a simple token-bucket rate limiter for
// suppressing repeated alert notifications for the same cron job.
package ratelimit

import (
	"sync"
	"time"
)

// Limiter tracks per-job alert suppression windows.
type Limiter struct {
	mu       sync.Mutex
	cooldown time.Duration
	last     map[string]time.Time
}

// New returns a Limiter that suppresses repeated alerts within the given
// cooldown window. A zero or negative cooldown disables suppression.
func New(cooldown time.Duration) *Limiter {
	return &Limiter{
		cooldown: cooldown,
		last:     make(map[string]time.Time),
	}
}

// Allow reports whether an alert for the named job should be sent.
// It returns true the first time a job is seen, and again only after
// the cooldown window has elapsed since the last allowed alert.
func (l *Limiter) Allow(job string) bool {
	if l.cooldown <= 0 {
		return true
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if t, ok := l.last[job]; ok && now.Sub(t) < l.cooldown {
		return false
	}
	l.last[job] = now
	return true
}

// Reset clears the suppression record for the named job, allowing the
// next alert to pass through immediately.
func (l *Limiter) Reset(job string) {
	l.mu.Lock()
	defer l.mu.Unlock()
	delete(l.last, job)
}

// ResetAll clears suppression records for all jobs.
func (l *Limiter) ResetAll() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.last = make(map[string]time.Time)
}
