package analyzer

import "Repo-lyzer/internal/github"

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

	if score > 100 {
		score = 100
	}
	return score
}
