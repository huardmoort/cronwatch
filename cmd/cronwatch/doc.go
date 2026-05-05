// Package main is the entry point for the cronwatch daemon.
//
// Usage:
//
//	cronwatch [flags]
//
// Flags:
//
//	-config string
//		Path to the YAML configuration file (default "cronwatch.yaml").
//
//	-report
//		Print a one-shot status report of all monitored jobs and exit.
//		Does not start the daemon or heartbeat server.
//
// In normal (daemon) mode cronwatch:
//  1. Loads configuration and opens the persistent job-run store.
//  2. Starts an HTTP server that accepts heartbeat pings from cron jobs.
//  3. Periodically checks each job for missed runs and sends webhook
//     alerts via the configured notifier.
//  4. Shuts down gracefully on SIGINT or SIGTERM.
package main
