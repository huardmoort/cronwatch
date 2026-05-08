package alertlog

import (
	"bytes"
	"fmt"
	"net/http"
)

// postJSON sends a JSON-encoded body via HTTP POST to url using client.
// It returns an error when the response status code is not in the 2xx range.
func postJSON(client *http.Client, url string, body []byte) error {
	resp, err := client.Post(url, "application/json", bytes.NewReader(body)) //nolint:noctx
	if err != nil {
		return fmt.Errorf("alertlog: webhook post: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("alertlog: webhook returned status %d", resp.StatusCode)
	}
	return nil
}
