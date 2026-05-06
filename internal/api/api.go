// Package api provides a lightweight HTTP server exposing cronwatch status
// endpoints, including job health, metrics snapshots, and manual heartbeat
// injection for testing purposes.
package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/example/cronwatch/internal/metrics"
	"github.com/example/cronwatch/internal/reporter"
)

// Server holds dependencies for the HTTP API.
type Server struct {
	mux      *http.ServeMux
	reporter *reporter.Reporter
	metrics  *metrics.Collector
}

// New creates a new API Server and registers its routes.
func New(r *reporter.Reporter, m *metrics.Collector) *Server {
	s := &Server{
		mux:      http.NewServeMux(),
		reporter: r,
		metrics:  m,
	}
	s.mux.HandleFunc("/healthz", s.handleHealthz)
	s.mux.HandleFunc("/status", s.handleStatus)
	s.mux.HandleFunc("/metrics", s.handleMetrics)
	return s
}

// ServeHTTP implements http.Handler so Server can be used directly.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

// handleHealthz returns a simple liveness response.
func (s *Server) handleHealthz(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": "ok",
		"time":   time.Now().UTC().Format(time.RFC3339),
	})
}

// handleStatus returns the collected job statuses as JSON.
func (s *Server) handleStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	statuses, err := s.reporter.Collect()
	if err != nil {
		http.Error(w, "failed to collect statuses", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(statuses)
}

// handleMetrics returns a snapshot of runtime metrics as JSON.
func (s *Server) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	snap := s.metrics.Snapshot()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(snap)
}
