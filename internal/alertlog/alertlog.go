// Package alertlog provides a persistent log of alerts sent by cronwatch.
// It records when alerts were fired, for which job, and the reason, so that
// operators can review alert history via the report or API.
package alertlog

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// Entry represents a single alert that was fired.
type Entry struct {
	Job       string    `json:"job"`
	Reason    string    `json:"reason"`
	FiredAt   time.Time `json:"fired_at"`
}

// Log holds an in-memory list of alert entries backed by a JSON file.
type Log struct {
	mu      sync.Mutex
	path    string
	entries []Entry
}

// New loads an existing alert log from path, or starts empty if the file does
// not exist. Returns an error if the file exists but cannot be parsed.
func New(path string) (*Log, error) {
	l := &Log{path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return l, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &l.entries); err != nil {
		return nil, err
	}
	return l, nil
}

// Record appends a new alert entry and persists the log to disk.
func (l *Log) Record(job, reason string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.entries = append(l.entries, Entry{
		Job:     job,
		Reason:  reason,
		FiredAt: time.Now().UTC(),
	})
	return l.save()
}

// Entries returns a copy of all recorded alert entries.
func (l *Log) Entries() []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	out := make([]Entry, len(l.entries))
	copy(out, l.entries)
	return out
}

// EntriesForJob returns a copy of all recorded alert entries for the given job name.
func (l *Log) EntriesForJob(job string) []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []Entry
	for _, e := range l.entries {
		if e.Job == job {
			out = append(out, e)
		}
	}
	return out
}

// EntriesSince returns a copy of all recorded alert entries fired at or after
// the given time. Useful for querying recent alerts in reports or dashboards.
func (l *Log) EntriesSince(t time.Time) []Entry {
	l.mu.Lock()
	defer l.mu.Unlock()
	var out []Entry
	for _, e := range l.entries {
		if !e.FiredAt.Before(t) {
			out = append(out, e)
		}
	}
	return out
}

// save writes the current entries slice to disk. Caller must hold l.mu.
func (l *Log) save() error {
	data, err := json.MarshalIndent(l.entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(l.path, data, 0o644)
}
