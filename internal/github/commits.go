package github

import (
	"fmt"
	"time"
)

type Commit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Date time.Time `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

type CommitFile struct {
	Filename  string `json:"filename"`
	Additions int    `json:"additions"`
	Deletions int    `json:"deletions"`
	Changes   int    `json:"changes"`
	Status    string `json:"status"`
}

type CommitDetail struct {
	SHA   string       `json:"sha"`
	Files []CommitFile `json:"files"`
}

func (c *Client) GetCommits(owner, repo string, days int) ([]Commit, error) {
	var allCommits []Commit
	since := time.Now().AddDate(0, 0, -days).Format(time.RFC3339)

	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/commits?since=%s&per_page=%d&page=%d",
			owner, repo, since, perPage, page,
		)

		var commits []Commit
		err := c.get(url, &commits)
		if err != nil {
			return nil, err
		}

		// Stop when no more commits or fewer than per_page
		if len(commits) == 0 || len(commits) < perPage {
			allCommits = append(allCommits, commits...)
			break
		}

		allCommits = append(allCommits, commits...)
		page++
	}

	return allCommits, nil
}

func (c *Client) GetCommit(owner, repo, sha string) (*CommitDetail, error) {
	url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repo, sha)
	var commit CommitDetail
	err := c.get(url, &commit)
	if err != nil {
		return nil, err
	}
	return &commit, nil
}
