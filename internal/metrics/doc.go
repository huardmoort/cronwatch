// Package metrics exposes lightweight, thread-safe counters that track
// cronwatch daemon activity at runtime.
//
// # Overview
//
// A single [Collector] is created at startup and passed to components that
// need to record events. Callers increment individual counters via helper
// methods (IncChecks, IncAlerts, IncHeartbeats, IncMissed). At any point a
// [Snapshot] can be obtained for logging or reporting without blocking
// ongoing increments.
//
// # Usage
//
//	col := metrics.New()
//	col.IncChecks()
//	snap := col.Snapshot()
//	fmt.Println(snap.ChecksTotal)
package metrics
