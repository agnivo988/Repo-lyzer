package github

type Issue struct {
	State             string `json:"state"`
	PullRequestURL    string `json:"pull_request"`
	IsPullRequest     bool   `json:"is_pull_request"` // Custom field, may not be in API response
	HTMLURL           string `json:"html_url"`
}

func (c *Client) GetIssues(owner, repo string, state string) ([]Issue, error) {
	var issues []Issue
	url := "https://api.github.com/repos/" + owner + "/" + repo + "/issues?state=" + state
	err := c.get(url, &issues)
	return issues, err
}

func (c *Client) GetPullRequests(owner, repo string, state string) ([]Issue, error) {
	var prs []Issue
	url := "https://api.github.com/repos/" + owner + "/" + repo + "/pulls?state=" + state
	err := c.get(url, &prs)
	return prs, err
}
