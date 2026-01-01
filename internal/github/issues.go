package github

type Issue struct {
	State string `json:"state"`
}

func (c *Client) GetIssues(owner, repo string, state string) ([]Issue, error) {
	var issues []Issue
	url := "https://api.github.com/repos/" + owner + "/" + repo + "/issues?state=" + state
	err := c.get(url, &issues)
	return issues, err
}
