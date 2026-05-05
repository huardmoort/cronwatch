// Package reporter provides structured status reporting for cron job health.
package reporter

import (
	"fmt"
	"io"
	"os"
	"text/tabwriter"
	"time"

	"github.com/example/cronwatch/internal/schedule"
	"github.com/example/cronwatch/internal/store"
)

// JobStatus holds the computed health status for a single job.
type JobStatus struct {
	Name     string
	CronExpr string
	LastRun  *time.Time
	NextRun  time.Time
	Missed   bool
}

// Reporter generates status reports for monitored jobs.
type Reporter struct {
	store  *store.Store
	writer io.Writer
}

// New creates a Reporter writing to the given writer.
// If w is nil, os.Stdout is used.
func New(s *store.Store, w io.Writer) *Reporter {
	if w == nil {
		w = os.Stdout
	}
	return &Reporter{store: s, writer: w}
}

// Collect builds a JobStatus slice for the given job names and cron expressions.
func (r *Reporter) Collect(jobs map[string]string, now time.Time) ([]JobStatus, error) {
	statuses := make([]JobStatus, 0, len(jobs))
	for name, expr := range jobs {
		last := r.store.LastRun(name)
		next, err := schedule.NextRun(expr, now)
		if err != nil {
			return nil, fmt.Errorf("job %q: %w", name, err)
		}
		missed := false
		if last != nil {
			missed = schedule.IsMissed(expr, *last, now)
		}
		statuses = append(statuses, JobStatus{
			Name:     name,
			CronExpr: expr,
			LastRun:  last,
			NextRun:  next,
			Missed:   missed,
		})
	}
	return statuses, nil
}

// Print writes a human-readable table of job statuses to the reporter's writer.
func (r *Reporter) Print(statuses []JobStatus) {
	tw := tabwriter.NewWriter(r.writer, 0, 0, 2, ' ', 0)
	fmt.Fprintln(tw, "JOB\tLAST RUN\tNEXT RUN\tSTATUS")
	for _, s := range statuses {
		lastRun := "never"
		if s.LastRun != nil {
			lastRun = s.LastRun.Format(time.RFC3339)
		}
		status := "OK"
		if s.Missed {
			status = "MISSED"
		}
		fmt.Fprintf(tw, "%s\t%s\t%s\t%s\n",
			s.Name, lastRun, s.NextRun.Format(time.RFC3339), status)
	}
	tw.Flush()
}
