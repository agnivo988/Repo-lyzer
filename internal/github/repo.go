package github

import "time"

type Repo struct {
	Name        string `json:"name"`
	FullName    string `json:"full_name"`
	Stars       int    `json:"stargazers_count"`
	Forks       int    `json:"forks_count"`
	OpenIssues  int    `json:"open_issues_count"`
	Description string `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

func (c *Client) GetRepo(owner, repo string) (*Repo, error) {
	var r Repo
	err := c.get("https://api.github.com/repos/"+owner+"/"+repo, &r)
	return &r, err
}
