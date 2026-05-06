// Package alertlog records a persistent history of alerts fired by cronwatch.
//
// Each time the watcher detects a missed or failed cron job and dispatches a
// notification, an Entry is appended to the log. Entries are stored as a
// JSON array on disk so they survive daemon restarts.
//
// Typical usage:
//
//	log, err := alertlog.New("/var/lib/cronwatch/alerts.json")
//	if err != nil {
//	    // handle
//	}
//	if err := log.Record("backup", "missed run"); err != nil {
//	    // handle
//	}
//	for _, e := range log.Entries() {
//	    fmt.Println(e.Job, e.Reason, e.FiredAt)
//	}
package alertlog
