// Package api provides the HTTP API for cronwatch status and health endpoints.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/cronwatch/internal/alertlog"
	"github.com/cronwatch/internal/metrics"
)

// Server holds dependencies for the HTTP API.
type Server struct {
	mux      *http.ServeMux
	metrics  *metrics.Collector
	alertLog *alertlog.Log
}

// New creates a new API Server with the given metrics collector and alert log.
func New(mc *metrics.Collector, al *alertlog.Log) *Server {
	s := &Server{
		mux:      http.NewServeMux(),
		metrics:  mc,
		alertLog: al,
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/healthz", s.handleHealthz())
	s.mux.HandleFunc("/status", s.handleStatus())
	s.mux.HandleFunc("/metrics", s.handleMetrics())
	if s.alertLog != nil {
		s.mux.HandleFunc("/alerts", s.alertLog.Handler())
	}
}

// ServeHTTP implements http.Handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *Server) handleHealthz() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}
}

func (s *Server) handleStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}
		snap := s.metrics.Snapshot()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(snap)
	}
}

func (s *Server) handleMetrics() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		snap := s.metrics.Snapshot()
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(snap)
	}
}
