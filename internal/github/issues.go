package github

import (
	"fmt"
	"time"
)

type Issue struct {
	Number    int        `json:"number"`
	State     string     `json:"state"`
	CreatedAt time.Time  `json:"created_at"`
	ClosedAt  *time.Time `json:"closed_at"`
	PullRequest *struct{} `json:"pull_request"`
}

func (c *Client) GetIssues(owner, repo string, state string) ([]Issue, error) {
	var all []Issue
	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/issues?state=%s&per_page=%d&page=%d",
			owner, repo, state, perPage, page,
		)
		var batch []Issue
		if err := c.get(url, &batch); err != nil {
			return nil, err
		}
		for _, issue := range batch {
    if issue.PullRequest != nil {
        continue
    }
    all = append(all, issue)
}
		if len(batch) < perPage {
			break
		}
		page++
	}
	return all, nil
}