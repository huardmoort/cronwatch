package alertlog

// NewWebhookSender is re-exported for integration tests in the alertlog_test package.
var NewWebhookSender = newWebhookSender

func newWebhookSender(endpoint string, client *http.Client) *WebhookSender {
	return NewWebhookSender(endpoint, client)
}
