package github

import "time"

type Commit struct {
	SHA    string `json:"sha"`
	Commit struct {
		Author struct {
			Date time.Time `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

func (c *Client) GetCommits(owner, repo string, days int) ([]Commit, error) {
	var commits []Commit
	since := time.Now().AddDate(0, 0, -days).Format(time.RFC3339)

	url := "https://api.github.com/repos/" + owner + "/" + repo + "/commits?since=" + since
	err := c.get(url, &commits)
	return commits, err
}
