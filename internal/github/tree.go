package github

import gocache "github.com/patrickmn/go-cache"

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
		if err := c.get("https://api.github.com/repos/"+owner+"/"+repo+"/git/trees/"+branch+"?recursive=1", &t); err != nil {
			return nil, err
		}

		c.cache.Set(cacheKey, t.Tree, gocache.DefaultExpiration)
		return copyTreeEntries(t.Tree), nil
	})
	if err != nil {
		return nil, err
	}
	src := v.([]TreeEntry)
	out := make([]TreeEntry, len(src))
	copy(out, src)
	return out, nil
}
