// Package ratelimit provides per-job alert rate limiting for cronwatch.
//
// When a cron job fails or is missed, the watcher may call the notifier
// on every check cycle. The Limiter in this package ensures that alerts
// for the same job are suppressed until a configurable cooldown window
// has elapsed, preventing notification floods.
//
// Usage:
//
//	limiter := ratelimit.New(30 * time.Minute)
//	if limiter.Allow(jobName) {
//		notifier.Send(ctx, jobName, message)
//	}
package ratelimit
