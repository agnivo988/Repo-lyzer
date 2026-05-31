package ui

import (
	"strings"
	"testing"
)

func TestBuildBlobURL_EscapesUntrustedPathSegments(t *testing.T) {
	t.Parallel()

	got, err := buildBlobURL("github.com", "owner", "repo", "main", "/dir/a&b?.go")
	if err != nil {
		t.Fatalf("buildBlobURL returned error: %v", err)
	}

	want := "https://github.com/owner/repo/blob/main/dir/a&b%3F.go"
	if got != want {
		t.Fatalf("buildBlobURL() = %q, want %q", got, want)
	}
}

func TestBuildVSCodeBlobURL_EscapesPathAndUsesVSCodeHost(t *testing.T) {
	t.Parallel()

	got, err := buildVSCodeBlobURL("owner", "repo", "main", "/folder/name with spaces.ts")
	if err != nil {
		t.Fatalf("buildVSCodeBlobURL returned error: %v", err)
	}

	want := "https://vscode.dev/github/owner/repo/blob/main/folder/name%20with%20spaces.ts"
	if got != want {
		t.Fatalf("buildVSCodeBlobURL() = %q, want %q", got, want)
	}
}

func TestBuildBlobURL_CustomBranch(t *testing.T) {
	t.Parallel()

	got, err := buildBlobURL("github.com", "owner", "repo", "master", "/dir/file.go")
	if err != nil {
		t.Fatalf("buildBlobURL returned error: %v", err)
	}

	want := "https://github.com/owner/repo/blob/master/dir/file.go"
	if got != want {
		t.Fatalf("buildBlobURL() = %q, want %q", got, want)
	}
}

func TestBuildBlobURL_RejectsInvalidInputs(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		host     string
		owner    string
		repo     string
		branch   string
		filePath string
	}{
		{name: "empty owner", host: "github.com", owner: "", repo: "repo", branch: "main", filePath: "/main.go"},
		{name: "empty repo", host: "github.com", owner: "owner", repo: "", branch: "main", filePath: "/main.go"},
		{name: "empty file path", host: "github.com", owner: "owner", repo: "repo", branch: "main", filePath: ""},
		{name: "slash-only file path", host: "github.com", owner: "owner", repo: "repo", branch: "main", filePath: "/"},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := buildBlobURL(tt.host, tt.owner, tt.repo, tt.branch, tt.filePath)
			if err == nil {
				t.Fatalf("expected error for case %q", tt.name)
			}
		})
	}
}

func TestSanitizeRepoPath_RemovesExtraSeparatorsAndEscapes(t *testing.T) {
	t.Parallel()

	got, err := sanitizeRepoPath("//a///b c//d&x.go")
	if err != nil {
		t.Fatalf("sanitizeRepoPath returned error: %v", err)
	}

	if got != "a/b c/d&x.go" {
		t.Fatalf("sanitizeRepoPath() = %q, want %q", got, "a/b c/d&x.go")
	}

	if strings.Contains(got, "//") {
		t.Fatalf("sanitizeRepoPath() retained duplicate separators: %q", got)
	}
}
