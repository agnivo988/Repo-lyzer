package github

import (
	"fmt"
	"time"

	gocache "github.com/patrickmn/go-cache"
)

// Release represents a GitHub release
type Release struct {
	ID          int       `json:"id"`
	TagName     string    `json:"tag_name"`
	Name        string    `json:"name"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
}

// GetReleases fetches releases for a repository
func (c *Client) GetReleases(owner, repo string) ([]Release, error) {
	cacheKey := "releases:" + owner + "/" + repo
	if cached, found := c.cache.Get(cacheKey); found {
		return copyReleases(cached.([]Release)), nil
	}

	v, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		if cached, found := c.cache.Get(cacheKey); found {
			return copyReleases(cached.([]Release)), nil
		}

		var allReleases []Release

		page := 1
		perPage := 100

		for {
			url := fmt.Sprintf(
				"https://api.github.com/repos/%s/%s/releases?per_page=%d&page=%d",
				owner, repo, perPage, page,
			)

			var releases []Release
			if err := c.get(url, &releases); err != nil {
				return nil, err
			}

			if len(releases) == 0 {
				break
			}

			allReleases = append(allReleases, releases...)

			if len(releases) < perPage {
				break
			}

			page++
		}

		c.cache.Set(cacheKey, allReleases, gocache.DefaultExpiration)
		return copyReleases(allReleases), nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]Release), nil
}

func copyReleases(r []Release) []Release {
	out := make([]Release, len(r))
	copy(out, r)
	return out
}
