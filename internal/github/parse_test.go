package github

import (
	"testing"
)

func TestParseGitHubURL(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedOwner string
		expectedRepo  string
		expectError   bool
	}{
		{
			name:          "standard owner/repo",
			input:         "golang/go",
			expectedOwner: "golang",
			expectedRepo:  "go",
			expectError:   false,
		},
		{
			name:          "https URL",
			input:         "https://github.com/golang/go",
			expectedOwner: "golang",
			expectedRepo:  "go",
			expectError:   false,
		},
		{
			name:          "https URL with trailing slash",
			input:         "https://github.com/golang/go/",
			expectedOwner: "golang",
			expectedRepo:  "go",
			expectError:   false,
		},
		{
			name:          "https URL with .git",
			input:         "https://github.com/golang/go.git",
			expectedOwner: "golang",
			expectedRepo:  "go",
			expectError:   false,
		},
		{
			name:          "ssh URL",
			input:         "git@github.com:golang/go.git",
			expectedOwner: "golang",
			expectedRepo:  "go",
			expectError:   false,
		},
		{
			name:          "URL with subpage",
			input:         "https://github.com/golang/go/tree/master/src",
			expectedOwner: "golang",
			expectedRepo:  "go",
			expectError:   false,
		},
		{
			name:          "URL with query params",
			input:         "https://github.com/golang/go?tab=readme-ov-file",
			expectedOwner: "golang",
			expectedRepo:  "go",
			expectError:   false,
		},
		{
			name:          "URL with hash anchor",
			input:         "https://github.com/golang/go#readme",
			expectedOwner: "golang",
			expectedRepo:  "go",
			expectError:   false,
		},
		{
			name:        "empty input",
			input:       "",
			expectError: true,
		},
		{
			name:        "spaces only",
			input:       "   ",
			expectError: true,
		},
		{
			name:        "invalid owner format - consecutive hyphens",
			input:       "go--lang/go",
			expectError: true,
		},
		{
			name:        "invalid owner format - start with hyphen",
			input:       "-golang/go",
			expectError: true,
		},
		{
			name:        "invalid owner format - end with hyphen",
			input:       "golang-/go",
			expectError: true,
		},
		{
			name:        "invalid characters in owner",
			input:       "go.lang/go",
			expectError: true,
		},
		{
			name:        "owner too long",
			input:       "abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz/go",
			expectError: true,
		},
		{
			name:        "whitespace in repository name",
			input:       "golang/go lang",
			expectError: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			owner, repo, err := ParseGitHubURL(tc.input)
			if tc.expectError {
				if err == nil {
					t.Errorf("expected error for input %q, but got nil", tc.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for input %q: %v", tc.input)
				}
				if owner != tc.expectedOwner {
					t.Errorf("expected owner %q, got %q for input %q", tc.expectedOwner, owner, tc.input)
				}
				if repo != tc.expectedRepo {
					t.Errorf("expected repo %q, got %q for input %q", tc.expectedRepo, repo, tc.input)
				}
			}
		})
	}
}
