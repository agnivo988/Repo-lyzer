package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// Client handles GitHub API requests
type Client struct {
	http  *http.Client
	token string
}

// User represents a GitHub user
type User struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

// NewClient creates a new GitHub API client
func NewClient() *Client {
	return &Client{
		http:  &http.Client{Timeout: 30 * time.Second},
		token: os.Getenv("GITHUB_TOKEN"),
	}
}

// HasToken returns true if a GitHub token is configured
func (c *Client) HasToken() bool {
	return c.token != ""
}

// get performs a GET request to the GitHub API and decodes the JSON response.
// It handles authentication and provides detailed error messages for rate limiting.
func (c *Client) get(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("network error: %w", err)
	}
	defer resp.Body.Close()

	// Handle rate limiting with detailed message
	if resp.StatusCode == http.StatusForbidden || resp.StatusCode == 429 {
		remaining := resp.Header.Get("X-RateLimit-Remaining")
		resetTime := resp.Header.Get("X-RateLimit-Reset")

		if remaining == "0" || resp.StatusCode == 429 {
			resetUnix, _ := strconv.ParseInt(resetTime, 10, 64)
			resetAt := time.Unix(resetUnix, 0)
			waitTime := time.Until(resetAt)

			if c.token == "" {
				return fmt.Errorf("🔴 Rate limit exceeded! Resets in %s\n"+
					"Tip: Set GITHUB_TOKEN env variable for 5000 requests/hour (vs 60 unauthenticated)",
					formatDuration(waitTime))
			}
			return fmt.Errorf("🔴 Rate limit exceeded! Resets in %s", formatDuration(waitTime))
		}
	}

	if resp.StatusCode == http.StatusNotFound {
		return fmt.Errorf("repository not found (check spelling or permissions)")
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return fmt.Errorf("authentication failed (check your GITHUB_TOKEN)")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API error: %s", resp.Status)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}

// formatDuration formats a duration in a human-readable way
func formatDuration(d time.Duration) string {
	if d < 0 {
		return "now"
	}
	if d < time.Minute {
		return fmt.Sprintf("%d seconds", int(d.Seconds()))
	}
	return fmt.Sprintf("%d min %d sec", int(d.Minutes()), int(d.Seconds())%60)
}

// GetUser fetches the authenticated user
func (c *Client) GetUser() (*User, error) {
	var u User
	err := c.get("https://api.github.com/user", &u)
	return &u, err
}

// GetFileContent fetches the content of a file from a repository
// Returns the base64 encoded content
func (c *Client) GetFileContent(owner, repo, path string) (string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, path)

	var result struct {
		Content  string `json:"content"`
		Encoding string `json:"encoding"`
	}

	if err := c.get(url, &result); err != nil {
		return "", err
	}

	// GitHub returns content with newlines, remove them for proper base64 decoding
	content := result.Content
	content = strings.ReplaceAll(content, "\n", "")

	return content, nil
}
