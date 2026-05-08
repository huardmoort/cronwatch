// Package notify provides alerting functionality for cronwatch.
// It supports sending notifications when cron jobs fail or are missed.
package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Event represents a notification event for a cron job.
type Event struct {
	JobName   string
	Kind      Kind
	Message   string
	OccurredAt time.Time
}

// Kind describes the type of notification event.
type Kind string

const (
	KindFailure Kind = "failure"
	KindMissed  Kind = "missed"
)

// Notifier sends notifications for cron job events.
type Notifier struct {
	webhookURL string
	client     *http.Client
}

// New creates a Notifier that posts alerts to the given webhook URL.
func New(webhookURL string) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
		client:     &http.Client{Timeout: 10 * time.Second},
	}
}

// Send dispatches an Event notification to the configured webhook.
func (n *Notifier) Send(evt Event) error {
	if n.webhookURL == "" {
		return nil
	}

	payload := map[string]string{
		"job":        evt.JobName,
		"kind":       string(evt.Kind),
		"message":    evt.Message,
		"occurred_at": evt.OccurredAt.Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notify: marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notify: post webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Read a limited portion of the response body to include in the error
		// for easier debugging, without risking unbounded memory consumption.
		respBody, _ := io.ReadAll(io.LimitReader(resp.Body, 256))
		if len(respBody) > 0 {
			return fmt.Errorf("notify: webhook returned status %d: %s", resp.StatusCode, respBody)
		}
		return fmt.Errorf("notify: webhook returned status %d", resp.StatusCode)
	}

	return nil
}
