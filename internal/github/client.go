package github

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type Client struct {
	http *http.Client
	token string
}

func NewClient() *Client {
	return &Client{
		http:  &http.Client{},
		token: os.Getenv("GITHUB_TOKEN"),
	}
}

func (c *Client) get(url string, target interface{}) error {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/vnd.github+json")

	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf(
	    "GitHub API error: %s (tip: set GITHUB_TOKEN env variable)",
	       resp.Status,
)
	}

	return json.NewDecoder(resp.Body).Decode(target)
}
