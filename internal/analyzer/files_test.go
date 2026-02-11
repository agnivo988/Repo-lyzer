package analyzer

import (
	"testing"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestParseExtensionFilter(t *testing.T) {
	filter := ParseExtensionFilter("go, .MD,txt")
	if !filter.IsActive() {
		t.Fatal("expected filter to be active")
	}

	cases := []struct {
		path string
		want bool
	}{
		{"main.go", true},
		{"README.md", true},
		{"notes.TXT", true},
		{"script.py", false},
		{"Dockerfile", false},
	}

	for _, tc := range cases {
		if got := filter.Matches(tc.path); got != tc.want {
			t.Errorf("Matches(%q) = %v, want %v", tc.path, got, tc.want)
		}
	}
}

func TestFilterTreeEntriesByExtensions(t *testing.T) {
	entries := []github.TreeEntry{
		{Path: "main.go", Type: "blob"},
		{Path: "README.md", Type: "blob"},
		{Path: "docs", Type: "tree"},
	}

	filter := ParseExtensionFilter("go")
	filtered := FilterTreeEntriesByExtensions(entries, filter)

	if len(filtered) != 1 || filtered[0].Path != "main.go" {
		t.Fatalf("filtered entries = %v, want only main.go", filtered)
	}
}

func TestLanguageSizesFromTree(t *testing.T) {
	entries := []github.TreeEntry{
		{Path: "main.go", Type: "blob", Size: 120},
		{Path: "README.md", Type: "blob", Size: 80},
		{Path: "script", Type: "blob", Size: 10},
		{Path: "docs", Type: "tree", Size: 1},
	}

	sizes := LanguageSizesFromTree(entries)
	if sizes["go"] != 120 {
		t.Errorf("go size = %d, want 120", sizes["go"])
	}
	if sizes["md"] != 80 {
		t.Errorf("md size = %d, want 80", sizes["md"])
	}
	if sizes["unknown"] != 10 {
		t.Errorf("unknown size = %d, want 10", sizes["unknown"])
	}
}
