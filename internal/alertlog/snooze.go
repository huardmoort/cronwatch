package alertlog

import (
	"sync"
	"time"
)

// SnoozePolicy temporarily suppresses alerts for a job until a deadline.
// Unlike MutePolicy (which is operator-driven and persistent), a snooze
// is a short-lived, programmatic suppression that expires automatically.
type SnoozePolicy struct {
	mu      sync.Mutex
	records map[string]time.Time // job -> snooze-until
}

// NewSnoozePolicy returns an initialised SnoozePolicy.
func NewSnoozePolicy() *SnoozePolicy {
	return &SnoozePolicy{
		records: make(map[string]time.Time),
	}
}

// Snooze suppresses alerts for job for the given duration.
// Calling Snooze again before expiry extends the deadline.
func (s *SnoozePolicy) Snooze(job string, d time.Duration) {
	if d <= 0 {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records[job] = time.Now().Add(d)
}

// IsSnoozed reports whether alerts for job are currently suppressed.
func (s *SnoozePolicy) IsSnoozed(job string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	until, ok := s.records[job]
	if !ok {
		return false
	}
	if time.Now().After(until) {
		delete(s.records, job)
		return false
	}
	return true
}

// Cancel removes any active snooze for job immediately.
func (s *SnoozePolicy) Cancel(job string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.records, job)
}

// ActiveJobs returns the names of all currently snoozed jobs.
func (s *SnoozePolicy) ActiveJobs() []string {
	s.mu.Lock()
	defer s.mu.Unlock()
	now := time.Now()
	out := make([]string, 0, len(s.records))
	for job, until := range s.records {
		if now.Before(until) {
			out = append(out, job)
		}
	}
	return out
}
