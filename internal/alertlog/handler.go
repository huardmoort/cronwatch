package alertlog

import (
	"encoding/json"
	"net/http"
)

// Handler returns an http.HandlerFunc that serves the alert log as JSON.
// It accepts an optional ?job= query parameter to filter entries by job name.
func (l *Log) Handler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		jobFilter := r.URL.Query().Get("job")

		l.mu.Lock()
		entries := make([]Entry, 0, len(l.entries))
		for _, e := range l.entries {
			if jobFilter == "" || e.Job == jobFilter {
				entries = append(entries, e)
			}
		}
		l.mu.Unlock()

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(entries); err != nil {
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
		}
	}
}
