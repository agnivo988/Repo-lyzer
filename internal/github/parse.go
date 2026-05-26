package github

import (
	"fmt"
	"strings"
)

// ParseGitHubURL parses a GitHub repository URL or a string in 'owner/repo' format
// and validates it. It supports formats like:
// - owner/repo
// - https://github.com/owner/repo
// - git@github.com:owner/repo.git
// - https://github.com/owner/repo/tree/main
// - and other variations (with trailing slashes, query parameters, etc.)
func ParseGitHubURL(input string) (owner, repo string, err error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return "", "", fmt.Errorf("repository identifier/URL cannot be empty")
	}

	// Remove ssh prefix if it exists
	if strings.HasPrefix(input, "git@github.com:") {
		input = strings.TrimPrefix(input, "git@github.com:")
	} else if strings.HasPrefix(input, "git@github.com/") {
		input = strings.TrimPrefix(input, "git@github.com/")
	}

	// Remove schema prefixes
	input = strings.TrimPrefix(input, "https://")
	input = strings.TrimPrefix(input, "http://")
	input = strings.TrimPrefix(input, "www.")

	// Remove github.com/ prefix
	input = strings.TrimPrefix(input, "github.com/")

	// Strip query parameters and anchors
	if idx := strings.Index(input, "?"); idx != -1 {
		input = input[:idx]
	}
	if idx := strings.Index(input, "#"); idx != -1 {
		input = input[:idx]
	}

	// Remove trailing slashes
	input = strings.TrimRight(input, "/")

	// Remove trailing .git
	input = strings.TrimSuffix(input, ".git")

	// Split by slashes to get owner and repository name
	parts := strings.Split(input, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid repository format: must contain both owner and repository name")
	}

	owner = parts[0]
	repo = parts[1]

	if owner == "" {
		return "", "", fmt.Errorf("owner name cannot be empty")
	}
	if repo == "" {
		return "", "", fmt.Errorf("repository name cannot be empty")
	}

	// Basic validation for GitHub owner name patterns
	if len(owner) > 39 {
		return "", "", fmt.Errorf("owner name '%s' is too long (maximum 39 characters)", owner)
	}
	if strings.HasPrefix(owner, "-") || strings.HasSuffix(owner, "-") {
		return "", "", fmt.Errorf("owner name '%s' cannot start or end with a hyphen", owner)
	}
	if strings.Contains(owner, "--") {
		return "", "", fmt.Errorf("owner name '%s' cannot contain consecutive hyphens", owner)
	}

	// Check for valid characters in owner (alphanumeric, hyphens)
	for _, char := range owner {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-') {
			return "", "", fmt.Errorf("owner name '%s' contains invalid character '%c' (only alphanumeric characters and hyphens allowed)", owner, char)
		}
	}

	// Basic validation for repository name patterns
	if len(repo) > 100 {
		return "", "", fmt.Errorf("repository name '%s' is too long (maximum 100 characters)", repo)
	}

	for _, char := range repo {
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
			return "", "", fmt.Errorf("repository name '%s' cannot contain whitespace", repo)
		}
	}

	return owner, repo, nil
}
