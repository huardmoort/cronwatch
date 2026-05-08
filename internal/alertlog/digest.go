package alertlog

import (
	"fmt"
	"strings"
	"time"
)

// DigestEntry represents a summarised view of alerts for a single job
// over a given time window.
type DigestEntry struct {
	Job        string
	Count      int
	FirstAlert time.Time
	LastAlert  time.Time
}

// Digest builds a human-readable digest of alert activity over the
// provided duration window relative to now. Jobs with no alerts in
// the window are omitted.
func Digest(l *Log, window time.Duration) []DigestEntry {
	l.mu.RLock()
	defer l.mu.RUnlock()

	cutoff := time.Now().Add(-window)

	type bucket struct {
		count int
		first time.Time
		last  time.Time
	}

	buckets := make(map[string]*bucket)

	for _, e := range l.entries {
		if e.At.Before(cutoff) {
			continue
		}
		b, ok := buckets[e.Job]
		if !ok {
			buckets[e.Job] = &bucket{count: 1, first: e.At, last: e.At}
			continue
		}
		b.count++
		if e.At.Before(b.first) {
			b.first = e.At
		}
		if e.At.After(b.last) {
			b.last = e.At
		}
	}

	out := make([]DigestEntry, 0, len(buckets))
	for job, b := range buckets {
		out = append(out, DigestEntry{
			Job:        job,
			Count:      b.count,
			FirstAlert: b.first,
			LastAlert:  b.last,
		})
	}
	return out
}

// FormatDigest renders a slice of DigestEntry values as a plain-text
// summary suitable for logging or a webhook message body.
func FormatDigest(entries []DigestEntry, window time.Duration) string {
	if len(entries) == 0 {
		return fmt.Sprintf("No alerts in the last %s.", window)
	}
	var sb strings.Builder
	fmt.Fprintf(&sb, "Alert digest (last %s):\n", window)
	for _, e := range entries {
		fmt.Fprintf(&sb, "  %-30s  alerts: %d  first: %s  last: %s\n",
			e.Job,
			e.Count,
			e.FirstAlert.Format(time.RFC3339),
			e.LastAlert.Format(time.RFC3339),
		)
	}
	return sb.String()
}
