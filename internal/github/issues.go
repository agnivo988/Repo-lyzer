package github

import (
	"fmt"
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

// GetIssues fetches all issues for a repository with pagination
// state can be "open", "closed", or "all"
func (c *Client) GetIssues(owner, repo string, state string) ([]Issue, error) {
	cacheKey := "issues:" + owner + "/" + repo + ":" + state
	if cached, found := c.cache.Get(cacheKey); found {
		return copyIssues(cached.([]Issue)), nil
	}

	v, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		if cached, found := c.cache.Get(cacheKey); found {
			return copyIssues(cached.([]Issue)), nil
		}

		var allIssues []Issue

		page := 1
		perPage := 100

		for {
			url := fmt.Sprintf(
				"https://api.github.com/repos/%s/%s/issues?state=%s&per_page=%d&page=%d",
				owner, repo, state, perPage, page,
			)

			var issues []Issue
			if err := c.get(url, &issues); err != nil {
				return nil, err
			}

			// Stop when no more issues
			if len(issues) == 0 {
				break
			}

			allIssues = append(allIssues, issues...)

			// Stop when fewer than per_page (last page)
			if len(issues) < perPage {
				break
			}

			page++
		}

		c.cache.Set(cacheKey, allIssues, gocache.DefaultExpiration)
		return copyIssues(allIssues), nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]Issue), nil
}
