// Package api exposes a small HTTP interface for cronwatch.
//
// Endpoints:
//
//	 GET /healthz  – liveness probe; always returns 200 with a JSON body.
//	 GET /status   – returns the current status of all monitored cron jobs
//	                 as a JSON array of JobStatus objects.
//	 GET /metrics  – returns a JSON snapshot of internal runtime counters
//	                 (alerts sent, heartbeats received, uptime, etc.).
//
// The server is created with New and satisfies http.Handler, so it can be
// passed directly to http.ListenAndServe or wrapped with middleware.
package api
