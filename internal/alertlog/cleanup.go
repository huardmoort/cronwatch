package alertlog

import "time"

// Cleanup removes alert log entries older than the given retention duration.
// It rewrites the persisted log with only the entries that fall within the
// retention window. Returns the number of entries pruned.
func (l *Log) Cleanup(retention time.Duration) (int, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if retention <= 0 {
		return 0, nil
	}

	cutoff := time.Now().Add(-retention)

	kept := l.entries[:0]
	pruned := 0
	for _, e := range l.entries {
		if e.At.After(cutoff) {
			kept = append(kept, e)
		} else {
			pruned++
		}
	}

	if pruned == 0 {
		return 0, nil
	}

	l.entries = kept
	if err := l.persist(); err != nil {
		return 0, err
	}

	return pruned, nil
}

// OlderThan returns all entries recorded before the given time.
func (l *Log) OlderThan(t time.Time) []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()

	var out []Entry
	for _, e := range l.entries {
		if e.At.Before(t) {
			out = append(out, e)
		}
	}
	return out
}
