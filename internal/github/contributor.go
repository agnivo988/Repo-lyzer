package github

import "fmt"

// Contributor represents a GitHub contributor
type Contributor struct {
	Login   string `json:"login"`
	Commits int    `json:"contributions"`
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
