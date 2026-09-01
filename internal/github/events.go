package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Event represents a generic GitHub activity event
type Event struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Actor     User            `json:"actor"`
	Repo      Repo            `json:"repo"`
	Payload   json.RawMessage `json:"payload"`
	CreatedAt time.Time       `json:"created_at"`
}

// GetRepositoryEvents fetches the latest events for a repository using ETag for safe polling
func (c *Client) GetRepositoryEvents(owner, repo, etag string) ([]Event, string, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/events?per_page=100", owner, repo)

	req, err := http.NewRequestWithContext(c.ctx, "GET", url, nil)
	if err != nil {
		return nil, etag, err
	}

	req.Header.Set("Accept", "application/vnd.github+json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
	if etag != "" {
		req.Header.Set("If-None-Match", etag)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, etag, fmt.Errorf("network error fetching events: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotModified {
		// No new events, return empty list and same ETag
		return []Event{}, etag, nil
	}

	if resp.StatusCode != http.StatusOK {
		return nil, etag, fmt.Errorf("github API error: %s", resp.Status)
	}

	newETag := resp.Header.Get("ETag")
	var events []Event
	if err := json.NewDecoder(resp.Body).Decode(&events); err != nil {
		return nil, etag, err
	}

	return events, newETag, nil
}
