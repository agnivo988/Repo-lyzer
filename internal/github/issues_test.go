package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetIssuesPagination(t *testing.T) {
	// Create dummy issues for page 1 (full page) and page 2 (partial page)
	page1 := make([]Issue, 100)
	for i := 0; i < 100; i++ {
		page1[i] = Issue{Number: i + 1}
	}
	page2 := []Issue{{Number: 101}, {Number: 102}}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		w.Header().Set("Content-Type", "application/json")
		
		if page == "1" {
			json.NewEncoder(w).Encode(page1)
		} else if page == "2" {
			json.NewEncoder(w).Encode(page2)
		} else {
			json.NewEncoder(w).Encode([]Issue{})
		}
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.http = server.Client() // Use server's client for the transport

	issues, err := client.GetIssues("owner", "repo", "open")
	if err != nil {
		t.Fatalf("GetIssues failed: %v", err)
	}

	expectedCount := 102
	if len(issues) != expectedCount {
		t.Errorf("Expected %d issues, got %d", expectedCount, len(issues))
	}

	if issues[0].Number != 1 {
		t.Errorf("Expected first issue number 1, got %d", issues[0].Number)
	}

	if issues[101].Number != 102 {
		t.Errorf("Expected last issue number 102, got %d", issues[101].Number)
	}
}

func TestGetIssuesEmpty(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Issue{})
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.http = server.Client()

	issues, err := client.GetIssues("owner", "repo", "open")
	if err != nil {
		t.Fatalf("GetIssues failed: %v", err)
	}

	if len(issues) != 0 {
		t.Errorf("Expected 0 issues, got %d", len(issues))
	}
}
