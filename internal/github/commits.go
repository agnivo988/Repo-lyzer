package github

import (
	"fmt"
	"strconv"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

type Commit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Date time.Time `json:"date"`
		} `json:"author"`
	} `json:"commit"`
	Author *struct {
		Login string `json:"login"`
	} `json:"author"`
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
	cacheKey := "commits:" + owner + "/" + repo + ":" + strconv.Itoa(days)
	if cached, found := c.cache.Get(cacheKey); found {
		return copyCommits(cached.([]Commit)), nil
	}

	v, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		if cached, found := c.cache.Get(cacheKey); found {
			return copyCommits(cached.([]Commit)), nil
		}

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
			if err := c.get(url, &commits); err != nil {
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

		c.cache.Set(cacheKey, allCommits, gocache.DefaultExpiration)
		return copyCommits(allCommits), nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]Commit), nil
}

func copyCommitDetail(d CommitDetail) CommitDetail {
	out := d
	out.Files = make([]CommitFile, len(d.Files))
	copy(out.Files, d.Files)
	return out
}

func (c *Client) GetCommit(owner, repo, sha string) (*CommitDetail, error) {
	cacheKey := "commit:" + owner + "/" + repo + ":" + sha
	if cached, found := c.cache.Get(cacheKey); found {
		d := copyCommitDetail(cached.(CommitDetail))
		return &d, nil
	}

	v, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		if cached, found := c.cache.Get(cacheKey); found {
			d := copyCommitDetail(cached.(CommitDetail))
			return &d, nil
		}

		url := fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s", owner, repo, sha)
		var commit CommitDetail
		if err := c.get(url, &commit); err != nil {
			return nil, err
		}

		c.cache.Set(cacheKey, commit, gocache.DefaultExpiration)
		d := copyCommitDetail(commit)
		return &d, nil
	})
	if err != nil {
		return nil, err
	}
	d := copyCommitDetail(*v.(*CommitDetail))
	return &d, nil
}
