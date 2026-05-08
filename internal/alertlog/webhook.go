package alertlog

import (
	"encoding/json"
	"net/http"
	"time"
)

// WebhookPayload is the JSON body sent to a webhook endpoint on digest delivery.
type WebhookPayload struct {
	GeneratedAt time.Time       `json:"generated_at"`
	WindowHours int             `json:"window_hours"`
	Jobs        []DigestJobStat `json:"jobs"`
}

// WebhookSender delivers a digest payload to an HTTP endpoint.
type WebhookSender struct {
	client  *http.Client
	endpoint string
}

// NewWebhookSender returns a WebhookSender targeting the given URL.
func NewWebhookSender(endpoint string, client *http.Client) *WebhookSender {
	if client == nil {
		client = &http.Client{Timeout: 10 * time.Second}
	}
	return &WebhookSender{client: client, endpoint: endpoint}
}

// Send marshals the digest and POSTs it to the configured endpoint.
// It returns an error if the endpoint is empty, the request fails, or
// the server responds with a non-2xx status code.
func (w *WebhookSender) Send(d Digest) error {
	if w.endpoint == "" {
		return nil
	}

	payload := WebhookPayload{
		GeneratedAt: d.GeneratedAt,
		WindowHours: d.WindowHours,
		Jobs:        d.Jobs,
	}

	buf, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	return postJSON(w.client, w.endpoint, buf)
}
