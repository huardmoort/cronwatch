// Package metrics provides lightweight in-process counters for cronwatch
// operational telemetry such as alerts sent, heartbeats received, and checks
// performed. Values are safe for concurrent use.
package metrics

import (
	"sync"
	"sync/atomic"
	"time"
)

// Snapshot holds a point-in-time copy of all counters.
type Snapshot struct {
	ChecksTotal    int64
	AlertsSent     int64
	HeartbeatsRecv int64
	MissedJobs     int64
	CapturedAt     time.Time
}

// Collector accumulates operational counters.
type Collector struct {
	checksTotal    atomic.Int64
	alertsSent     atomic.Int64
	heartbeatsRecv atomic.Int64
	missedJobs     atomic.Int64

	mu      sync.Mutex
	started time.Time
}

// New returns an initialised Collector with the start time set to now.
func New() *Collector {
	return &Collector{started: time.Now()}
}

// IncChecks records one completed check cycle.
func (c *Collector) IncChecks() { c.checksTotal.Add(1) }

// IncAlerts records one alert successfully dispatched.
func (c *Collector) IncAlerts() { c.alertsSent.Add(1) }

// IncHeartbeats records one heartbeat received from a job.
func (c *Collector) IncHeartbeats() { c.heartbeatsRecv.Add(1) }

// IncMissed records one missed-job detection.
func (c *Collector) IncMissed() { c.missedJobs.Add(1) }

// Snapshot returns a consistent copy of all counters.
func (c *Collector) Snapshot() Snapshot {
	c.mu.Lock()
	defer c.mu.Unlock()
	return Snapshot{
		ChecksTotal:    c.checksTotal.Load(),
		AlertsSent:     c.alertsSent.Load(),
		HeartbeatsRecv: c.heartbeatsRecv.Load(),
		MissedJobs:     c.missedJobs.Load(),
		CapturedAt:     time.Now(),
	}
}

// Uptime returns the duration since the Collector was created.
func (c *Collector) Uptime() time.Duration {
	return time.Since(c.started)
}
