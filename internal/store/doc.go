// Package store provides a lightweight, file-backed persistence layer for
// cronwatch job run history.
//
// Each JobRun entry captures the job name, start time, success flag, and an
// optional human-readable note (e.g. exit code or error message). Records are
// appended to a JSON file on disk so that cronwatch can survive restarts and
// still detect missed or failing jobs by consulting historical data.
//
// Typical usage:
//
//	s, err := store.New("/var/lib/cronwatch/runs.json")
//	if err != nil { ... }
//
//	// after a job completes:
//	_ = s.Record(store.JobRun{
//		JobName:   "nightly-backup",
//		StartedAt: time.Now(),
//		Success:   true,
//	})
//
//	// to check when a job last ran:
//	if run, ok := s.LastRun("nightly-backup"); ok {
//		fmt.Println("last run:", run.StartedAt)
//	}
package store
