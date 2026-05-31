package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync/atomic"
	"testing"
	"time"
)

// Test that transient 5xx errors are retried and eventual success is returned
func TestGet_RetriesOnServerError(t *testing.T) {
	var calls int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n <= 2 {
			// first two calls: server error
			http.Error(w, "temporary server error", http.StatusInternalServerError)
			return
		}

		// third call: success with user JSON
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"login": "alice", "name": "Alice"})
	}))
	defer srv.Close()

	c := NewClient()
	// use server's client to ensure same transport
	c.http = srv.Client()
	c.SetContext(context.Background())

	var u User
	if err := c.get(srv.URL+"/user", &u); err != nil {
		t.Fatalf("expected success after retries, got err: %v", err)
	}
	if u.Login != "alice" {
		t.Fatalf("expected login alice, got %q", u.Login)
	}
	if atomic.LoadInt32(&calls) < 3 {
		t.Fatalf("expected at least 3 calls, got %d", calls)
	}
}

// Test that when GitHub rate-limits (429) the client waits until reset and retries
func TestGet_WaitsOnRateLimit(t *testing.T) {
	var calls int32

	// first call: respond with 429 and X-RateLimit-Reset in ~1s
	// second call: return success
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n == 1 {
			resetAt := time.Now().Add(1 * time.Second).Unix()
			w.Header().Set("X-RateLimit-Remaining", "0")
			w.Header().Set("X-RateLimit-Reset", strconv.FormatInt(resetAt, 10))
			http.Error(w, "rate limited", http.StatusTooManyRequests)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"login": "ratelimited", "name": "Rate"})
	}))
	defer srv.Close()

	c := NewClient()
	c.http = srv.Client()
	// ensure client is treated as authenticated so it will wait and retry
	c.SetToken("fake-token")
	// give context enough time for wait+1s buffer used in client
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	c.SetContext(ctx)

	var u User
	start := time.Now()
	if err := c.get(srv.URL+"/user", &u); err != nil {
		t.Fatalf("expected success after rate-limit wait, got err: %v", err)
	}
	elapsed := time.Since(start)
	if elapsed < 1*time.Second {
		t.Fatalf("expected to wait at least ~1s for rate-limit reset, waited %v", elapsed)
	}
	if u.Login != "ratelimited" {
		t.Fatalf("expected login ratelimited, got %q", u.Login)
	}
	if atomic.LoadInt32(&calls) < 2 {
		t.Fatalf("expected at least 2 calls, got %d", calls)
	}
}
