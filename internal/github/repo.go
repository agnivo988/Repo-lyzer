package github

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type Repo struct {
	Name          string    `json:"name"`
	Owner         User      `json:"owner"`
	FullName      string    `json:"full_name"`
	Stars         int       `json:"stargazers_count"`
	Forks         int       `json:"forks_count"`
	OpenIssues    int       `json:"open_issues_count"`
	Description   string    `json:"description"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	PushedAt      time.Time `json:"pushed_at"`
	WatchersCount int       `json:"watchers_count"`
	Language      string    `json:"language"`
	Fork          bool      `json:"fork"`
	Archived      bool      `json:"archived"`
	Private       bool      `json:"private"`
	DefaultBranch string    `json:"default_branch"`
	HTMLURL       string    `json:"html_url"`
	CloneURL      string    `json:"clone_url"`
}

func (c *Client) GetRepo(owner, repo string) (*Repo, error) {
	cacheKey := "repo:" + owner + "/" + repo
	if cached, found := c.cache.Get(cacheKey); found {
		r := cached.(Repo)
		return &r, nil
	}

	v, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		if cached, found := c.cache.Get(cacheKey); found {
			r := cached.(Repo)
			return &r, nil
		}

		var r Repo
		if err := c.get("https://api.github.com/repos/"+owner+"/"+repo, &r); err != nil {
			return nil, err
		}

		c.cache.Set(cacheKey, r, gocache.DefaultExpiration)
		return &r, nil
	})
	if err != nil {
		return nil, err
	}
	return v.(*Repo), nil
}
