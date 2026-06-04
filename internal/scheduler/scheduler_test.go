package scheduler

import (
	"net/http"
	"net/http/httptest"
	"strings"
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

	// Check if the error is a timeout error
	if !strings.Contains(err.Error(), "Client.Timeout exceeded") && !strings.Contains(err.Error(), "timeout") {
		t.Errorf("Expected timeout error message, but got: %v", err)
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
	err := s.sendToWebhook("http://invalid-url-that-does-not-exist.local", "test_report.json", []byte(`{"status":"ok"}`))

	if err == nil {
		t.Fatal("Expected error for invalid URL, but got nil")
	}
}
