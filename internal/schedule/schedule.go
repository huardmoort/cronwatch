package schedule

import (
	"fmt"
	"time"

	"github.com/robfig/cron/v3"
)

// NextRun returns the next scheduled run time for a cron expression.
func NextRun(expr string, from time.Time) (time.Time, error) {
	parser := cron.NewParser(
		cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)
	sched, err := parser.Parse(expr)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid cron expression %q: %w", expr, err)
	}
	return sched.Next(from), nil
}

// IsMissed reports whether a job with the given expression and last run time
// has missed its scheduled window as of now, given a grace period.
func IsMissed(expr string, lastRun time.Time, grace time.Duration) (bool, error) {
	now := time.Now()
	next, err := NextRun(expr, lastRun)
	if err != nil {
		return false, err
	}
	// If the next expected run is in the past beyond the grace period, it's missed.
	deadline := next.Add(grace)
	return now.After(deadline), nil
}

// DurationUntilNext returns the duration from now until the next scheduled run.
func DurationUntilNext(expr string) (time.Duration, error) {
	next, err := NextRun(expr, time.Now())
	if err != nil {
		return 0, err
	}
	return time.Until(next), nil
}
