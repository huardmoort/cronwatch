// Package heartbeat provides an HTTP handler that cron jobs can ping
// to record a successful execution in the store.
package heartbeat

import (
	"fmt"
	"net/http"
	"time"

	"github.com/example/cronwatch/internal/store"
)

// Recorder persists heartbeat pings for named jobs.
type Recorder struct {
	store *store.Store
}

// New returns a Recorder backed by the given store.
func New(s *store.Store) *Recorder {
	return &Recorder{store: s}
}

// Handler returns an http.HandlerFunc that records a heartbeat for the job
// name supplied as the "job" query parameter.
//
// Example request:
//
//	GET /ping?job=backup-db
func (r *Recorder) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		job := req.URL.Query().Get("job")
		if job == "" {
			http.Error(w, "missing 'job' query parameter", http.StatusBadRequest)
			return
		}

		if err := r.store.Record(job, time.Now()); err != nil {
			http.Error(w, fmt.Sprintf("failed to record heartbeat: %v", err), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "ok: recorded heartbeat for %q\n", job)
	}
}
