package github

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type IssueLabel struct {
	Name  string `json:"name"`
	Color string `json:"color"`
}

type Issue struct {
	Number      int          `json:"number"`
	Title       string       `json:"title"`
	State       string       `json:"state"`
	CreatedAt   time.Time    `json:"created_at"`
	UpdatedAt   time.Time    `json:"updated_at"`
	ClosedAt    *time.Time   `json:"closed_at"`
	Comments    int          `json:"comments"`
	PullRequest *struct{}    `json:"pull_request,omitempty"`
	Labels      []IssueLabel `json:"labels"`
	User        User         `json:"user"`
}

func (c *Client) GetIssues(owner, repo string, state string) ([]Issue, error) {
	cacheKey := "issues:" + owner + "/" + repo + ":" + state
	if cached, found := c.cache.Get(cacheKey); found {
		return copyIssues(cached.([]Issue)), nil
	}

	v, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		if cached, found := c.cache.Get(cacheKey); found {
			return copyIssues(cached.([]Issue)), nil
		}

		var issues []Issue
		url := "https://api.github.com/repos/" + owner + "/" + repo + "/issues?state=" + state + "&per_page=100"
		if err := c.get(url, &issues); err != nil {
			return nil, err
		}

		c.cache.Set(cacheKey, issues, gocache.DefaultExpiration)
		return copyIssues(issues), nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]Issue), nil
}
