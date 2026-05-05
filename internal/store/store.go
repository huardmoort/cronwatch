package store

import (
	"encoding/json"
	"os"
	"sync"
	"time"
)

// JobRun records the result of a single cron job execution.
type JobRun struct {
	JobName   string    `json:"job_name"`
	StartedAt time.Time `json:"started_at"`
	Success   bool      `json:"success"`
	Note      string    `json:"note,omitempty"`
}

// Store persists job run history to a JSON file.
type Store struct {
	mu      sync.RWMutex
	path    string
	records []JobRun
}

// New loads an existing store from path, or creates a new empty one.
func New(path string) (*Store, error) {
	s := &Store{path: path}
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return s, nil
	}
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &s.records); err != nil {
		return nil, err
	}
	return s, nil
}

// Record appends a job run entry and flushes to disk.
func (s *Store) Record(run JobRun) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.records = append(s.records, run)
	return s.flush()
}

// LastRun returns the most recent JobRun for the given job name, and whether one was found.
func (s *Store) LastRun(jobName string) (JobRun, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	for i := len(s.records) - 1; i >= 0; i-- {
		if s.records[i].JobName == jobName {
			return s.records[i], true
		}
	}
	return JobRun{}, false
}

// flush writes records to disk; caller must hold the write lock.
func (s *Store) flush() error {
	data, err := json.MarshalIndent(s.records, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path, data, 0o644)
}
