package analyzer

import (
	"time"

	"Repo-lyzer/internal/github"
)

func RepoMaturityScore(repo *github.Repo, commits int, contributors int, hasReleases bool) (int, string) {
	score := 0

	// Age
	ageYears := time.Since(repo.CreatedAt).Hours() / (24 * 365)
	if ageYears >= 1 {
		score += 20
	}

	// Activity
	if commits > 100 {
		score += 25
	}

	// Contributors
	if contributors > 1 {
		score += 20
	}

	// Releases
	if hasReleases {
		score += 20
	}

	// Issues sanity
	if repo.OpenIssues < 50 {
		score += 15
	}

	level := "Prototype"
	switch {
	case score >= 80:
		level = "Production-Ready"
	case score >= 60:
		level = "Stable"
	case score >= 40:
		level = "Growing"
	}

	return score, level
}
