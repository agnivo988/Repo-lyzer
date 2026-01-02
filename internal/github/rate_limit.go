package github

import (
	"time"
)
type RateLimit struct {
	Resources struct {
		Core struct {
			Limit     int `json:"limit"`
			Remaining int `json:"remaining"`
			Reset     int `json:"reset"`
		} `json:"core"`
	} `json:"resources"`
}
func (c *Client) GetRateLimit() (*RateLimit, error) {
	var rateLimit RateLimit
	err := c.get("https://api.github.com/rate_limit", &rateLimit)
	if err != nil {
		return nil, err
	}
	return &rateLimit, nil
}
func (r *RateLimit) ResetTime() time.Time {
	return time.Unix(int64(r.Resources.Core.Reset), 0)
}
