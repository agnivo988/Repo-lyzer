package github

import "fmt"

// Contributor represents a GitHub contributor
type Contributor struct {
	Login     string `json:"login"`
	Commits   int    `json:"contributions"`
	AvatarURL string `json:"avatar_url,omitempty"`
}

// GetContributors fetches ALL contributors (paginated)
func (c *Client) GetContributors(owner, repo string) ([]Contributor, error) {
	var allContributors []Contributor

	page := 1
	perPage := 100

	for {
		url := fmt.Sprintf(
			"https://api.github.com/repos/%s/%s/contributors?per_page=%d&page=%d",
			owner, repo, perPage, page,
		)

		var contributors []Contributor
		err := c.get(url, &contributors)
		if err != nil {
			return nil, err
		}

		// Stop when no more contributors
		if len(contributors) == 0 {
			break
		}

		allContributors = append(allContributors, contributors...)
		page++
	}

	return allContributors, nil
}

// GetContributorsWithAvatars fetches contributors and avatar URLs for top N contributors
func (c *Client) GetContributorsWithAvatars(owner, repo string, topN int) ([]Contributor, error) {
	contributors, err := c.GetContributors(owner, repo)
	if err != nil {
		return nil, err
	}

	// Fetch avatars for top contributors
	maxAvatars := topN
	if len(contributors) < maxAvatars {
		maxAvatars = len(contributors)
	}

	for i := 0; i < maxAvatars; i++ {
		user, err := c.GetUserByLogin(contributors[i].Login)
		if err != nil {
			// Log error but continue
			continue
		}
		contributors[i].AvatarURL = user.AvatarURL
	}

	return contributors, nil
}
