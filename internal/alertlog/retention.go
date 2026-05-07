package alertlog

import "time"

// RetentionPolicy describes how long alert log entries should be kept.
type RetentionPolicy struct {
	// MaxAge is the maximum age of an entry before it is eligible for pruning.
	// A zero value disables age-based pruning.
	MaxAge time.Duration

	// MaxEntries is the maximum number of entries to retain per job.
	// A zero value disables count-based pruning.
	MaxEntries int
}

// Apply prunes entries from the log according to the policy.
// It returns the number of entries removed.
func (p RetentionPolicy) Apply(l *Log) int {
	if p.MaxAge == 0 && p.MaxEntries == 0 {
		return 0
	}

	l.mu.Lock()
	defer l.mu.Unlock()

	before := len(l.entries)

	if p.MaxAge > 0 {
		cutoff := time.Now().Add(-p.MaxAge)
		filtered := l.entries[:0]
		for _, e := range l.entries {
			if e.At.After(cutoff) {
				filtered = append(filtered, e)
			}
		}
		l.entries = filtered
	}

	if p.MaxEntries > 0 {
		// Group entries by job and keep only the most recent MaxEntries per job.
		byJob := make(map[string][]Entry)
		for _, e := range l.entries {
			byJob[e.Job] = append(byJob[e.Job], e)
		}
		l.entries = l.entries[:0]
		for _, entries := range byJob {
			if len(entries) > p.MaxEntries {
				entries = entries[len(entries)-p.MaxEntries:]
			}
			l.entries = append(l.entries, entries...)
		}
	}

	after := len(l.entries)
	return before - after
}
