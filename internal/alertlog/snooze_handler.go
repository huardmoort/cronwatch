package alertlog

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
)

type snoozeHandler struct {
	policy *SnoozePolicy
}

// NewSnoozeHandler returns an http.Handler that exposes snooze management.
//
// POST /snooze?job=<name>&minutes=<n>  — snooze a job for n minutes
// DELETE /snooze?job=<name>            — cancel snooze for a job
// GET /snooze                          — list currently snoozed jobs
func NewSnoozeHandler(p *SnoozePolicy) http.Handler {
	return &snoozeHandler{policy: p}
}

func (h *snoozeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.listSnoozed(w)
	case http.MethodPost:
		h.addSnooze(w, r)
	case http.MethodDelete:
		h.cancelSnooze(w, r)
	default:
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *snoozeHandler) listSnoozed(w http.ResponseWriter) {
	jobs := h.policy.ActiveJobs()
	if jobs == nil {
		jobs = []string{}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(map[string][]string{"snoozed": jobs})
}

func (h *snoozeHandler) addSnooze(w http.ResponseWriter, r *http.Request) {
	job := r.URL.Query().Get("job")
	if job == "" {
		http.Error(w, "missing job parameter", http.StatusBadRequest)
		return
	}
	minStr := r.URL.Query().Get("minutes")
	min, err := strconv.Atoi(minStr)
	if err != nil || min <= 0 {
		http.Error(w, "minutes must be a positive integer", http.StatusBadRequest)
		return
	}
	h.policy.Snooze(job, time.Duration(min)*time.Minute)
	w.WriteHeader(http.StatusNoContent)
}

func (h *snoozeHandler) cancelSnooze(w http.ResponseWriter, r *http.Request) {
	job := r.URL.Query().Get("job")
	if job == "" {
		http.Error(w, "missing job parameter", http.StatusBadRequest)
		return
	}
	h.policy.Cancel(job)
	w.WriteHeader(http.StatusNoContent)
}
