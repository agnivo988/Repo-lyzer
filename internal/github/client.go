package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	
	gocache "github.com/patrickmn/go-cache"
	"golang.org/x/sync/singleflight"
)

// Client handles GitHub API requests
type Client struct {
	http  *http.Client
	token string
	ctx   context.Context
	cache *gocache.Cache
	sf    singleflight.Group
}

// User represents a GitHub user
type User struct {
	Login     string `json:"login"`
	Name      string `json:"name"`
	AvatarURL string `json:"avatar_url"`
}

// NewClient creates a new GitHub API client
func NewClient() *Client {
	token := os.Getenv("REPO_LYZER_GITHUB_TOKEN")
	if token == "" {
		token = os.Getenv("GITHUB_TOKEN")
	}
	return &Client{
		http:  &http.Client{Timeout: 30 * time.Second},
		token: token,
		ctx:   context.Background(),
		cache: gocache.New(
			5*time.Minute,  // default expiration
			10*time.Minute, // cleanup interval
		),
	}
}

// NewClientWithToken creates a client pre-configured with the given token.
// If token is empty, behaves identically to NewClient().
func NewClientWithToken(token string) *Client {
	c := NewClient()
	if token != "" {
		c.SetToken(token)
	}
	return c
}

// SetContext sets the context used to cancel in-flight HTTP requests.
func (c *Client) SetContext(ctx context.Context) {
	if ctx == nil {
		c.ctx = context.Background()
		return
	}
	c.ctx = ctx
}

// HasToken returns true if a GitHub token is configured
func (c *Client) HasToken() bool {
	return c.token != ""
}

// SetToken sets the GitHub token for authentication
func (c *Client) SetToken(token string) {
	c.token = token
	c.cache.Flush()
}

// get performs a GET request to the GitHub API and decodes the JSON response.
// It handles authentication and provides detailed error messages for rate limiting.
func (c *Client) get(url string, target interface{}) error {
	// Implement retry with exponential backoff and respect rate-limit headers.
	const maxRetries = 3

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		req, err := http.NewRequestWithContext(c.ctx, "GET", url, nil)
		if err != nil {
			return err
		}

		req.Header.Set("Accept", "application/vnd.github+json")
		if c.token != "" {
			req.Header.Set("Authorization", "Bearer "+c.token)
		}

		resp, err := c.http.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("network error: %w", err)
			// Retry on transient network errors
			if attempt == maxRetries {
				return lastErr
			}
			backoff := time.Duration(1<<attempt) * time.Second
			select {
			case <-time.After(backoff):
				continue
			case <-c.ctx.Done():
				return c.ctx.Err()
			}
		}

		// Ensure body closed on retry or error branches
		if resp == nil {
			lastErr = fmt.Errorf("empty response")
			if attempt == maxRetries {
				return lastErr
			}
			continue
		}

		// Handle rate limiting: wait until reset and retry
		if resp.StatusCode == http.StatusForbidden || resp.StatusCode == 429 {
			remaining := resp.Header.Get("X-RateLimit-Remaining")
			resetTime := resp.Header.Get("X-RateLimit-Reset")

			if remaining == "0" || resp.StatusCode == 429 {
				resetUnix, _ := strconv.ParseInt(resetTime, 10, 64)
				resetAt := time.Unix(resetUnix, 0)
				waitTime := time.Until(resetAt)
				if waitTime < 0 {
					waitTime = time.Second
				}

				resp.Body.Close()
				// If unauthenticated, return informative error after waiting once
				if c.token == "" {
					select {
					case <-time.After(waitTime + time.Second):
						return fmt.Errorf("rate limit exceeded and no token configured; consider setting GITHUB_TOKEN")
					case <-c.ctx.Done():
						return c.ctx.Err()
					}
				}

				// Authenticated: wait and retry
				select {
				case <-time.After(waitTime + time.Second):
					continue
				case <-c.ctx.Done():
					return c.ctx.Err()
				}
			}

			resp.Body.Close()
			return fmt.Errorf("access forbidden (Status 403): the request was rejected by GitHub API or requires extended permissions")
		}

		// Not found
		if resp.StatusCode == http.StatusNotFound {
			resp.Body.Close()
			return fmt.Errorf("repository not found or inaccessible — it may be private or you may not have permission")
		}

		// Unauthorized
		if resp.StatusCode == http.StatusUnauthorized {
			resp.Body.Close()
			return fmt.Errorf("authentication failed (check your GITHUB_TOKEN)")
		}

		// Retry on 5xx
		if resp.StatusCode >= 500 && resp.StatusCode <= 599 {
			resp.Body.Close()
			if attempt == maxRetries {
				return fmt.Errorf("GitHub server error: %s", resp.Status)
			}
			backoff := time.Duration(1<<attempt) * time.Second
			select {
			case <-time.After(backoff):
				continue
			case <-c.ctx.Done():
				return c.ctx.Err()
			}
		}

		// Non-OK responses handled above, so decode on success
		if resp.StatusCode != http.StatusOK {
			errMsg := fmt.Errorf("GitHub API error: %s", resp.Status)
			resp.Body.Close()
			return errMsg
		}

		// Success
		defer resp.Body.Close()
		return json.NewDecoder(resp.Body).Decode(target)
	}

	if lastErr != nil {
		return lastErr
	}
	return fmt.Errorf("request failed: %s", url)
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
	cacheKey := "user:me"
	if cached, found := c.cache.Get(cacheKey); found {
		u := cached.(User)
		return &u, nil
	}

	v, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		if cached, found := c.cache.Get(cacheKey); found {
			u := cached.(User)
			return &u, nil
		}

		var u User
		if err := c.get("https://api.github.com/user", &u); err != nil {
			return nil, err
		}

		c.cache.Set(cacheKey, u, gocache.DefaultExpiration)
		return &u, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*User), nil
}

// GetUserByLogin fetches a user by their login/username
func (c *Client) GetUserByLogin(login string) (*User, error) {
	cacheKey := "user:" + login
	if cached, found := c.cache.Get(cacheKey); found {
		u := cached.(User)
		return &u, nil
	}

	v, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		if cached, found := c.cache.Get(cacheKey); found {
			u := cached.(User)
			return &u, nil
		}

		url := fmt.Sprintf("https://api.github.com/users/%s", login)
		var u User
		if err := c.get(url, &u); err != nil {
			return nil, err
		}

		c.cache.Set(cacheKey, u, gocache.DefaultExpiration)
		return &u, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*User), nil
}

// GetFileContent fetches the content of a file from a repository
// Returns the base64 encoded content
func (c *Client) GetFileContent(owner, repo, path string) (string, error) {
	cacheKey := "file:" + owner + "/" + repo + ":" + path
	if cached, found := c.cache.Get(cacheKey); found {
		return cached.(string), nil
	}

	v, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		if cached, found := c.cache.Get(cacheKey); found {
			return cached.(string), nil
		}

		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/contents/%s", owner, repo, path)

		var result struct {
			Content  string `json:"content"`
			Encoding string `json:"encoding"`
		}

		if err := c.get(url, &result); err != nil {
			return "", err
		}

		// GitHub returns content with newlines, remove them for proper base64 decoding
		content := strings.ReplaceAll(result.Content, "\n", "")

		c.cache.Set(cacheKey, content, gocache.DefaultExpiration)
		return content, nil
	})
	if err != nil {
		return "", err
	}
	return v.(string), nil
}
