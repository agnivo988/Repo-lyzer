package github

import (
	"fmt"
	"strconv"

	gocache "github.com/patrickmn/go-cache"
)

// Contributor represents a GitHub contributor
type Contributor struct {
	Login     string `json:"login"`
	Commits   int    `json:"contributions"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// GetContributors fetches ALL contributors (paginated)
func (c *Client) GetContributors(owner, repo string) ([]Contributor, error) {
	cacheKey := "contributors:" + owner + "/" + repo
	if cached, found := c.cache.Get(cacheKey); found {
		return copyContributors(cached.([]Contributor)), nil
	}

	v, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		if cached, found := c.cache.Get(cacheKey); found {
			return copyContributors(cached.([]Contributor)), nil
		}

		var allContributors []Contributor

		page := 1
		perPage := 100

		for {
			url := fmt.Sprintf(
				"https://api.github.com/repos/%s/%s/contributors?per_page=%d&page=%d",
				owner, repo, perPage, page,
			)

			var contributors []Contributor
			if err := c.get(url, &contributors); err != nil {
				return nil, err
			}

			// Stop when no more contributors
			if len(contributors) == 0 {
				break
			}

			allContributors = append(allContributors, contributors...)
			page++
		}

		c.cache.Set(cacheKey, allContributors, gocache.DefaultExpiration)
		return copyContributors(allContributors), nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]Contributor), nil
}

// GetContributorsWithAvatars fetches contributors and avatar URLs for top N contributors
func (c *Client) GetContributorsWithAvatars(owner, repo string, topN int) ([]Contributor, error) {
	cacheKey := "contributors_avatars:" + owner + "/" + repo + ":" + strconv.Itoa(topN)
	if cached, found := c.cache.Get(cacheKey); found {
		return copyContributors(cached.([]Contributor)), nil
	}

	v, err, _ := c.sf.Do(cacheKey, func() (interface{}, error) {
		if cached, found := c.cache.Get(cacheKey); found {
			return copyContributors(cached.([]Contributor)), nil
		}

		contributors, err := c.GetContributors(owner, repo)
		if err != nil {
			return nil, err
		}

		copied := copyContributors(contributors)

		// Fetch avatars for top contributors
		maxAvatars := topN
		if len(copied) < maxAvatars {
			maxAvatars = len(copied)
		}

		for i := 0; i < maxAvatars; i++ {
			user, err := c.GetUserByLogin(copied[i].Login)
			if err != nil {
				// Log error but continue
				continue
			}
			copied[i].AvatarURL = user.AvatarURL
		}

		c.cache.Set(cacheKey, copied, gocache.DefaultExpiration)
		return copyContributors(copied), nil
	})
	if err != nil {
		return nil, err
	}
	return v.([]Contributor), nil
}
