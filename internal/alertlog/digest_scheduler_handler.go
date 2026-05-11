package alertlog

import (
	"encoding/json"
	"net/http"
	"time"
)

// DigestHandlerConfig holds configuration for the digest HTTP handler.
type DigestHandlerConfig struct {
	Log    *Log
	Window time.Duration
}

// NewDigestHandler returns an HTTP handler that serves an on-demand digest
// of alert activity within the configured window.
func NewDigestHandler(cfg DigestHandlerConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		window := cfg.Window
		if window <= 0 {
			window = 24 * time.Hour
		}

		d := Digest(cfg.Log, window)

		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(d); err != nil {
			http.Error(w, "encode error", http.StatusInternalServerError)
		}
	}
}
