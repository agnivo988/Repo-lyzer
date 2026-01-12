package github

func (c *Client) GetLanguages(owner, repo string) (map[string]int, error) {
	var langs map[string]int
	err := c.get("https://api.github.com/repos/"+owner+"/"+repo+"/languages", &langs)
	return langs, err
}
