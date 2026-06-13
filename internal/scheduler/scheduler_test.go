package scheduler

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

// TestSendToWebhookTimesOutWhenReceiverStalls verifies webhook delivery fails
// fast when a receiver accepts the request but never responds.
func TestSendToWebhookTimesOutWhenReceiverStalls(t *testing.T) {
	originalClient := webhookHTTPClient
	webhookHTTPClient = &http.Client{Timeout: 25 * time.Millisecond}
	t.Cleanup(func() {
		webhookHTTPClient = originalClient
	})

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	start := time.Now()
	err := (&Scheduler{}).sendToWebhook(server.URL, "report.json", []byte(`{"score":90}`))
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("expected stalled webhook receiver to return a timeout error")
	}
	if !strings.Contains(err.Error(), "failed to send webhook") {
		t.Fatalf("expected webhook send error, got %v", err)
	}
	if elapsed > time.Second {
		t.Fatalf("webhook call was not bounded by the configured timeout; elapsed %s", elapsed)
	}
}
