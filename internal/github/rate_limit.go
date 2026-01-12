package github

import (
	"fmt"
	"time"
)

// RateLimit represents GitHub API rate limit information
type RateLimit struct {
	Resources struct {
		Core struct {
			Limit     int `json:"limit"`
			Remaining int `json:"remaining"`
			Reset     int `json:"reset"`
			Used      int `json:"used"`
		} `json:"core"`
		Search struct {
			Limit     int `json:"limit"`
			Remaining int `json:"remaining"`
			Reset     int `json:"reset"`
		} `json:"search"`
	} `json:"resources"`
}

// GetRateLimit fetches current rate limit status from GitHub API
func (c *Client) GetRateLimit() (*RateLimit, error) {
	var rateLimit RateLimit
	err := c.get("https://api.github.com/rate_limit", &rateLimit)
	if err != nil {
		return nil, err
	}
	return &rateLimit, nil
}

// ResetTime returns the time when rate limit resets
func (r *RateLimit) ResetTime() time.Time {
	return time.Unix(int64(r.Resources.Core.Reset), 0)
}

// TimeUntilReset returns duration until rate limit resets
func (r *RateLimit) TimeUntilReset() time.Duration {
	return time.Until(r.ResetTime())
}

// IsLimited returns true if rate limit is exhausted
func (r *RateLimit) IsLimited() bool {
	return r.Resources.Core.Remaining == 0
}

// UsagePercent returns percentage of rate limit used
func (r *RateLimit) UsagePercent() float64 {
	if r.Resources.Core.Limit == 0 {
		return 0
	}
	return float64(r.Resources.Core.Limit-r.Resources.Core.Remaining) / float64(r.Resources.Core.Limit) * 100
}

// FormatResetTime returns human-readable time until reset
func (r *RateLimit) FormatResetTime() string {
	d := r.TimeUntilReset()
	if d < 0 {
		return "now"
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	return fmt.Sprintf("%dm %ds", int(d.Minutes()), int(d.Seconds())%60)
}

// GetRateLimitStatus returns a formatted status string
func (r *RateLimit) GetRateLimitStatus() string {
	if r.IsLimited() {
		return fmt.Sprintf("🔴 Rate Limited! Resets in %s", r.FormatResetTime())
	}
	if r.Resources.Core.Remaining < 10 {
		return fmt.Sprintf("🟡 Low: %d/%d remaining (resets in %s)",
			r.Resources.Core.Remaining, r.Resources.Core.Limit, r.FormatResetTime())
	}
	return fmt.Sprintf("🟢 OK: %d/%d remaining",
		r.Resources.Core.Remaining, r.Resources.Core.Limit)
}
