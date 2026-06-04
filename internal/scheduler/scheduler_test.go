package scheduler

import (
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSendToWebhookTimeout(t *testing.T) {
	// Create a test server that delays its response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(12 * time.Second) // Delay longer than the 10s client timeout
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	s := &Scheduler{}
	err := s.sendToWebhook(server.URL, "test_report.json", []byte(`{"status":"ok"}`))

	if err == nil {
		t.Fatal("Expected timeout error, but got nil")
	}

	// Check if the error is a timeout error using net.Error assertion
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		// Valid timeout path
	} else {
		t.Errorf("Expected timeout error, but got: %v", err)
	}
}

func TestSendToWebhookSuccess(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	s := &Scheduler{}
	err := s.sendToWebhook(server.URL, "test_report.json", []byte(`{"status":"ok"}`))

	if err != nil {
		t.Fatalf("Expected no error, but got: %v", err)
	}
}

func TestSendToWebhookInvalidURL(t *testing.T) {
	s := &Scheduler{}
	// Use an explicitly malformed URL to guarantee immediate failure
	err := s.sendToWebhook("http://invalid url with spaces", "test_report.json", []byte(`{"status":"ok"}`))

	if err == nil {
		t.Fatal("Expected error for invalid URL, but got nil")
	}
}
