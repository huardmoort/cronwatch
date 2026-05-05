// Package reporter provides status reporting for cronwatch.
//
// It queries the store for last-run timestamps, evaluates whether
// jobs are healthy or missed, and formats the results into a
// human-readable table.
//
// Basic usage:
//
//	s, _ := store.New("/var/lib/cronwatch/state.json")
//	r := reporter.New(s, os.Stdout)
//	statuses, err := r.Collect(jobMap, time.Now())
//	if err != nil { ... }
//	r.Print(statuses)
package reporter
