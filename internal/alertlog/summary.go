package alertlog

import "time"

// Summary holds aggregated statistics for alerts across all jobs or a single job.
type Summary struct {
	TotalAlerts  int            `json:"total_alerts"`
	ByJob        map[string]int `json:"by_job"`
	OldestAlert  *time.Time     `json:"oldest_alert,omitempty"`
	NewestAlert  *time.Time     `json:"newest_alert,omitempty"`
}

// Summarize computes a Summary from all entries in the log.
// If job is non-empty, only entries for that job are included.
func (l *Log) Summarize(job string) Summary {
	l.mu.Lock()
	defer l.mu.Unlock()

	s := Summary{
		ByJob: make(map[string]int),
	}

	for _, e := range l.entries {
		if job != "" && e.Job != job {
			continue
		}

		s.TotalAlerts++
		s.ByJob[e.Job]++

		t := e.At
		if s.OldestAlert == nil || t.Before(*s.OldestAlert) {
			copy := t
			s.OldestAlert = &copy
		}
		if s.NewestAlert == nil || t.After(*s.NewestAlert) {
			copy := t
			s.NewestAlert = &copy
		}
	}

	return s
}
