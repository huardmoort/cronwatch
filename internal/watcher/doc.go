// Package watcher implements the core monitoring loop for cronwatch.
//
// A Watcher is initialised with the application config, a persistent Store,
// and a Notifier.  Calling Start blocks and periodically calls checkAll,
// which iterates over every configured job, retrieves its last recorded run
// time from the store, and uses the schedule package to determine whether the
// job has missed its expected execution window.  When a missed run is
// detected an alert message is dispatched via the Notifier.
//
// Usage:
//
//	w := watcher.New(cfg, st, notifier)
//	go w.Start(time.Minute)
//	// … later …
//	w.Stop()
package watcher
