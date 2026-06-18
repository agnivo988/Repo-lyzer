package github

import (
	"fmt"

	gocache "github.com/patrickmn/go-cache"
)

func (c *Client) HasReleases(owner, repo string) (bool, error) {
	cacheKey := "releases:" + owner + "/" + repo
	if cached, found := c.cache.Get(cacheKey); found {
		return cached.(bool), nil
	}

	v, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		if cached, found := c.cache.Get(cacheKey); found {
			return cached.(bool), nil
		}

		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/releases?per_page=1&page=1",
			owner, repo,
		)

		var releases []struct{}
		if err := c.get(url, &releases); err != nil {
			return false, err
		}

		hasReleases := len(releases) > 0
		c.cache.Set(cacheKey, hasReleases, gocache.DefaultExpiration)
		return hasReleases, nil
	})
	if err != nil {
		return false, err
	}
	return v.(bool), nil
}
