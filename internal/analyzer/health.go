package analyzer

import (
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func CalculateHealth(repo *github.Repo, commits []github.Commit) int {
	score := 50

	if repo.Description != "" {
		score += 10
	}
	if repo.Stars > 50 {
		score += 10
	}
	if len(commits) > 10 {
		score += 20
	}
	if repo.OpenIssues < 20 {
		score += 10
	}

	if !repo.PushedAt.IsZero() {
		since := time.Since(repo.PushedAt)
		switch {
		case since <= 30*24*time.Hour:
			score += 10
		case since <= 90*24*time.Hour:
			score += 5
		case since > 365*24*time.Hour:
			score -= 10
		case since > 180*24*time.Hour:
			score -= 5
		}
	}

	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}
	return score
}
