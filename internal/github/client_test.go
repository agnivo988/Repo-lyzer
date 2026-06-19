package github

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
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

// Test that a 403 with Retry-After (secondary rate limit) is retried after waiting
func TestGet_RetriesOnSecondaryRateLimit(t *testing.T) {
	var calls int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n == 1 {
			// Simulate GitHub secondary rate limit: 403 + Retry-After, remaining is NOT 0
			w.Header().Set("X-RateLimit-Remaining", "50")
			w.Header().Set("Retry-After", "1")
			http.Error(w, "secondary rate limit", http.StatusForbidden)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"login": "secondary", "name": "Secondary"})
	}))
	defer srv.Close()

	c := NewClient()
	c.http = srv.Client()
	c.SetToken("fake-token")
	ctx, cancel := context.WithTimeout(context.Background(), 6*time.Second)
	defer cancel()
	c.SetContext(ctx)

	var u User
	start := time.Now()
	if err := c.get(srv.URL+"/user", &u); err != nil {
		t.Fatalf("expected success after secondary rate-limit wait, got err: %v", err)
	}
	elapsed := time.Since(start)
	if elapsed < 1*time.Second {
		t.Fatalf("expected to wait at least 1s for Retry-After, waited %v", elapsed)
	}
	if u.Login != "secondary" {
		t.Fatalf("expected login secondary, got %q", u.Login)
	}
	if atomic.LoadInt32(&calls) < 2 {
		t.Fatalf("expected at least 2 calls, got %d", calls)
	}
}

// Test that a genuine 403 (no Retry-After, no rate-limit headers) returns access forbidden error
func TestGet_ForbiddenWithoutRetryAfterReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Pure 403 — no Retry-After, remaining is non-zero
		w.Header().Set("X-RateLimit-Remaining", "50")
		http.Error(w, "forbidden", http.StatusForbidden)
	}))
	defer srv.Close()

	c := NewClient()
	c.http = srv.Client()
	c.SetToken("fake-token")

	var u User
	err := c.get(srv.URL+"/user", &u)
	if err == nil {
		t.Fatal("expected error for genuine 403, got nil")
	}
	if !strings.Contains(err.Error(), "access forbidden") {
		t.Fatalf("expected 'access forbidden' in error, got: %v", err)
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
