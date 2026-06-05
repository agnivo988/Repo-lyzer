package github

import (
	"errors"
	"fmt"

	gocache "github.com/patrickmn/go-cache"
)

var ErrTreeTruncated = errors.New("repository file tree is truncated by the GitHub API (>100k files)")

type TreeEntry struct {
	Path string `json:"path"`
	Mode string `json:"mode"`
	Type string `json:"type"`
	Size int    `json:"size"`
	Sha  string `json:"sha"`
}

type TreeResponse struct {
	Sha       string      `json:"sha"`
	Url       string      `json:"url"`
	Tree      []TreeEntry `json:"tree"`
	Truncated bool        `json:"truncated"`
}

func (c *Client) GetFileTree(owner, repo, branch string) ([]TreeEntry, error) {
	cacheKey := "tree:" + owner + "/" + repo + ":" + branch
	if cached, found := c.cache.Get(cacheKey); found {
		return copyTreeEntries(cached.([]TreeEntry)), nil
	}

	v, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		if cached, found := c.cache.Get(cacheKey); found {
			return copyTreeEntries(cached.([]TreeEntry)), nil
		}

		var t TreeResponse
		// recursive=1 to get full tree
		url := fmt.Sprintf("%s/repos/%s/%s/git/trees/%s?recursive=1", c.BaseURL, owner, repo, branch)
		if err := c.get(url, &t); err != nil {
			return nil, err
		}

		var truncationErr error
		if t.Truncated {
			fmt.Println("⚠️ Warning: Repository file tree is truncated by the GitHub API (>100k files). Downstream analysis may be incomplete.")
			truncationErr = ErrTreeTruncated
		}

		c.cache.Set(cacheKey, t.Tree, gocache.DefaultExpiration)
		return copyTreeEntries(t.Tree), truncationErr
	})
	if err != nil && !errors.Is(err, ErrTreeTruncated) {
		return nil, err
	}
	src := v.([]TreeEntry)
	out := make([]TreeEntry, len(src))
	copy(out, src)
	return out, err
}
