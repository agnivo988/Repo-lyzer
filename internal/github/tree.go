package github

import "fmt"

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
	var t TreeResponse
	// recursive=1 to get full tree
	err := c.get("https://api.github.com/repos/"+owner+"/"+repo+"/git/trees/"+branch+"?recursive=1", &t)
	if err != nil {
		return nil, err
	}
	if t.Truncated {
		return t.Tree, fmt.Errorf("file tree truncated by GitHub: results may be incomplete for large repositories")
	}
	return t.Tree, nil
}
