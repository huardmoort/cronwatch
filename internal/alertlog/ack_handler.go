package alertlog

import (
	"encoding/json"
	"net/http"
	"time"
)

// NewAckHandler returns an http.Handler that exposes GET / POST / DELETE
// endpoints for managing alert acknowledgements.
//
//   GET    /ack?job=<name>  — check whether a job is currently acknowledged
//   POST   /ack             — acknowledge a job (body: {"job":"…","duration":"15m"})
//   DELETE /ack?job=<name>  — clear acknowledgement for a job
func NewAckHandler(p *AckPolicy) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			job := r.URL.Query().Get("job")
			if job == "" {
				http.Error(w, "missing job parameter", http.StatusBadRequest)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{ //nolint:errcheck
				"job":   job,
				"acked": p.IsAcked(job),
			})

		case http.MethodPost:
			var req struct {
				Job      string `json:"job"`
				Duration string `json:"duration"`
			}
			if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Job == "" {
				http.Error(w, "invalid request body", http.StatusBadRequest)
				return
			}
			var d time.Duration
			if req.Duration != "" {
				var err error
				if d, err = time.ParseDuration(req.Duration); err != nil {
					http.Error(w, "invalid duration", http.StatusBadRequest)
					return
				}
			}
			p.Acknowledge(req.Job, d)
			w.WriteHeader(http.StatusNoContent)

		case http.MethodDelete:
			job := r.URL.Query().Get("job")
			if job == "" {
				http.Error(w, "missing job parameter", http.StatusBadRequest)
				return
			}
			p.Clear(job)
			w.WriteHeader(http.StatusNoContent)

		default:
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
