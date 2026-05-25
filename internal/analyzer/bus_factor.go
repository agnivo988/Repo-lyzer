// Package analyzer provides functions for analyzing GitHub repository data.
// It includes calculations for repository health, maturity, bus factor, and other metrics.
package analyzer

import (
	"sort"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// BusFactor calculates the bus factor of a repository based on contributor commit distribution.
// The bus factor indicates how risky it is if key contributors leave the project.
// It returns a score from 1-3 and a risk level string.
// Parameters:
//   - contributors: Slice of repository contributors with their commit counts
//
// Returns:
//   - int: Risk score (1=High Risk, 2=Medium Risk, 3=Low Risk)
//   - string: Risk level description
//
// Example:
//
//	contributors := []github.Contributor{
//	    {Commits: 100},
//	    {Commits: 50},
//	    {Commits: 25},
//	}
//	score, risk := BusFactor(contributors)
//	// score: 2, risk: "Medium Risk"
// sortContributorsByCommitsDesc returns a copy sorted by commit count (highest first).
func sortContributorsByCommitsDesc(contributors []github.Contributor) []github.Contributor {
	sorted := append([]github.Contributor(nil), contributors...)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Commits == sorted[j].Commits {
			return sorted[i].Login < sorted[j].Login
		}
		return sorted[i].Commits > sorted[j].Commits
	})
	return sorted
}

func BusFactor(contributors []github.Contributor) (int, string) {
	if len(contributors) == 0 {
		return 0, "Unknown"
	}

	sorted := sortContributorsByCommitsDesc(contributors)

	total := 0
	for _, c := range sorted {
		total += c.Commits
	}
	if total == 0 {
		return 0, "Unknown"
	}

	top := sorted[0].Commits
	ratio := float64(top) / float64(total)

	switch {
	case ratio > 0.7:
		return 1, "High Risk"
	case ratio > 0.4:
		return 2, "Medium Risk"
	default:
		return 3, "Low Risk"
	}
}
