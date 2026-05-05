// Package notify implements alerting for cronwatch.
//
// A Notifier sends structured JSON payloads to a configured webhook URL
// whenever a monitored cron job fails or is detected as missed.
//
// Usage:
//
//	n := notify.New("https://hooks.example.com/alert")
//	err := n.Send(notify.Event{
//		JobName:    "daily-backup",
//		Kind:       notify.KindFailure,
//		Message:    "process exited with code 2",
//		OccurredAt: time.Now(),
//	})
//
// If the webhook URL is empty, Send is a no-op, making the notifier safe
// to use even when alerting is not configured.
package notify
