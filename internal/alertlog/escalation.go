package alertlog

import (
	"sync"
	"time"
)

// EscalationPolicy tracks repeated alert occurrences and escalates
// notifications when a job exceeds a failure threshold within a window.
type EscalationPolicy struct {
	mu        sync.Mutex
	threshold int
	window    time.Duration
	counts    map[string][]time.Time
}

// NewEscalationPolicy creates an EscalationPolicy that escalates when a job
// triggers more than threshold alerts within window.
func NewEscalationPolicy(threshold int, window time.Duration) *EscalationPolicy {
	return &EscalationPolicy{
		threshold: threshold,
		window:    window,
		counts:    make(map[string][]time.Time),
	}
}

// Record registers an alert occurrence for job and reports whether the
// escalation threshold has been reached.
func (e *EscalationPolicy) Record(job string, at time.Time) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.threshold == 0 || e.window == 0 {
		return false
	}

	cutoff := at.Add(-e.window)
	times := e.counts[job]

	// prune entries outside the window
	filtered := times[:0]
	for _, t := range times {
		if t.After(cutoff) {
			filtered = append(filtered, t)
		}
	}
	filtered = append(filtered, at)
	e.counts[job] = filtered

	return len(filtered) >= e.threshold
}

// Reset clears the occurrence history for a specific job.
func (e *EscalationPolicy) Reset(job string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	delete(e.counts, job)
}

// Count returns the number of occurrences within the window for job.
func (e *EscalationPolicy) Count(job string, now time.Time) int {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.window == 0 {
		return 0
	}
	cutoff := now.Add(-e.window)
	var n int
	for _, t := range e.counts[job] {
		if t.After(cutoff) {
			n++
		}
	}
	return n
}
