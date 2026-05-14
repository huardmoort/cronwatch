package alertlog

import (
	"sync"
	"time"
)

// AckPolicy tracks acknowledgements for alerts, suppressing further
// notifications for a job until the acknowledgement expires or is cleared.
type AckPolicy struct {
	mu   sync.Mutex
	acks map[string]time.Time
}

// NewAckPolicy returns an initialised AckPolicy.
func NewAckPolicy() *AckPolicy {
	return &AckPolicy{
		acks: make(map[string]time.Time),
	}
}

// Acknowledge records an acknowledgement for job that expires after d.
// If d is zero the acknowledgement is permanent until explicitly cleared.
func (a *AckPolicy) Acknowledge(job string, d time.Duration) {
	a.mu.Lock()
	defer a.mu.Unlock()
	var expiry time.Time
	if d > 0 {
		expiry = time.Now().Add(d)
	}
	a.acks[job] = expiry
}

// IsAcked reports whether job currently has an active acknowledgement.
func (a *AckPolicy) IsAcked(job string) bool {
	a.mu.Lock()
	defer a.mu.Unlock()
	expiry, ok := a.acks[job]
	if !ok {
		return false
	}
	if expiry.IsZero() {
		return true
	}
	if time.Now().After(expiry) {
		delete(a.acks, job)
		return false
	}
	return true
}

// Clear removes any active acknowledgement for job.
func (a *AckPolicy) Clear(job string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.acks, job)
}

// AckedJobs returns the names of all currently acknowledged jobs.
func (a *AckPolicy) AckedJobs() []string {
	a.mu.Lock()
	defer a.mu.Unlock()
	out := make([]string, 0, len(a.acks))
	for job, expiry := range a.acks {
		if expiry.IsZero() || time.Now().Before(expiry) {
			out = append(out, job)
		} else {
			delete(a.acks, job)
		}
	}
	return out
}
