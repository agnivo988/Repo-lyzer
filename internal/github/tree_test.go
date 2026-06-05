package github

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetFileTree_Truncated(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		resp := TreeResponse{
			Sha: "tree-sha",
			Tree: []TreeEntry{
				{Path: "file1.txt", Type: "blob"},
			},
			Truncated: true,
		}
		_ = json.NewEncoder(w).Encode(resp)
	}))
	defer srv.Close()

	c := NewClient()
	c.http = srv.Client()
	c.BaseURL = srv.URL

	tree, err := c.GetFileTree("owner", "repo", "main")
	if !errors.Is(err, ErrTreeTruncated) {
		t.Fatalf("expected ErrTreeTruncated, got %v", err)
	}

	if len(tree) != 1 {
		t.Errorf("expected 1 tree entry, got %d", len(tree))
	}
	if tree[0].Path != "file1.txt" {
		t.Errorf("expected file1.txt, got %s", tree[0].Path)
	}
}
