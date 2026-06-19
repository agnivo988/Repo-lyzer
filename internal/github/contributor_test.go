package github

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
)

// makeContributors returns a slice of n dummy Contributor values.
func makeContributors(n int) []Contributor {
	out := make([]Contributor, n)
	for i := range out {
		out[i] = Contributor{Login: "user", Commits: 1}
	}
	return out
}

// paginationLoop runs the same pagination logic as GetContributors so tests
// can exercise the exact loop condition without needing to redirect api.github.com.
func paginationLoop(c *Client, baseURL string) ([]Contributor, int32, error) {
	var calls int32
	perPage := 100
	var all []Contributor

	for page := 1; ; page++ {
		url := baseURL + "?per_page=100&page=" + itoa(page)
		atomic.AddInt32(&calls, 1)
		var batch []Contributor
		if err := c.get(url, &batch); err != nil {
			return nil, calls, err
		}
		if len(batch) == 0 {
			break
		}
		all = append(all, batch...)
		if len(batch) < perPage { // the fix: early exit on partial page
			break
		}
	}
	return all, calls, nil
}

// Test: 150 contributors (full page of 100 + partial of 50).
// Fixed code: 2 API calls. Old buggy code: 3 calls (also fetched empty page 3).
func TestGetContributors_PartialLastPageSavesOneCall(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		w.Header().Set("Content-Type", "application/json")
		if page == "1" {
			json.NewEncoder(w).Encode(makeContributors(100))
		} else {
			json.NewEncoder(w).Encode(makeContributors(50))
		}
	}))
	defer srv.Close()

	c := NewClient()
	c.http = srv.Client()
	c.SetToken("fake-token")

	all, calls, err := paginationLoop(c, srv.URL+"/contributors")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 2 {
		t.Fatalf("expected 2 API calls for 150 contributors, got %d", calls)
	}
	if len(all) != 150 {
		t.Fatalf("expected 150 contributors, got %d", len(all))
	}
}

// Test: 200 contributors (two full pages of 100 + empty page 3).
// Neither old nor new code can avoid the 3rd call here — both need the empty
// page to know there is nothing left. Verifies correct count is returned.
func TestGetContributors_ExactMultipleOfPerPage(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		page := r.URL.Query().Get("page")
		w.Header().Set("Content-Type", "application/json")
		switch page {
		case "1", "2":
			json.NewEncoder(w).Encode(makeContributors(100))
		default:
			json.NewEncoder(w).Encode([]Contributor{})
		}
	}))
	defer srv.Close()

	c := NewClient()
	c.http = srv.Client()
	c.SetToken("fake-token")

	all, calls, err := paginationLoop(c, srv.URL+"/contributors")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 API calls for 200 contributors (2 full + 1 empty), got %d", calls)
	}
	if len(all) != 200 {
		t.Fatalf("expected 200 contributors, got %d", len(all))
	}
}

// Test: 0 contributors — empty page on first call, loop stops immediately.
func TestGetContributors_EmptyRepo(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode([]Contributor{})
	}))
	defer srv.Close()

	c := NewClient()
	c.http = srv.Client()
	c.SetToken("fake-token")

	all, calls, err := paginationLoop(c, srv.URL+"/contributors")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 API call for empty repo, got %d", calls)
	}
	if len(all) != 0 {
		t.Fatalf("expected 0 contributors, got %d", len(all))
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	b := make([]byte, 0, 10)
	for n > 0 {
		b = append([]byte{byte('0' + n%10)}, b...)
		n /= 10
	}
	return string(b)
}
