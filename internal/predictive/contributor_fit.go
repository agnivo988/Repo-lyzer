package predictive

import (
	"math"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// ContributorFit represents how well a contributor matches an issue
type ContributorFit struct {
	IssueNumber int      `json:"issue_number"`
	Contributor string   `json:"contributor"`
	Score       float64  `json:"score"` // 0 to 100 (100 = Perfect Fit)
	Level       string   `json:"level"` // "Excellent", "Good", "Moderate", "Poor"
	Reasons     []string `json:"reasons"`
}

// AnalyzeContributorFit determines if a given contributor is a good fit for an issue
func AnalyzeContributorFit(issue github.Issue, contributor string, pastCommits []github.Commit) ContributorFit {
	score := 50.0 // Base score
	var reasons []string
	contributorLower := strings.ToLower(contributor)

	// 1. Check Historical Context (Have they contributed here before?)
	commitCount := 0
	for _, commit := range pastCommits {
		if commit.Author != nil && strings.ToLower(commit.Author.Login) == contributorLower {
			commitCount++
		}
	}

	if commitCount == 0 {
		score -= 10.0
		reasons = append(reasons, "First-time contributor to this repository")
	} else if commitCount > 10 {
		score += 30.0
		reasons = append(reasons, "Highly experienced contributor in this repository")
	} else {
		score += 15.0
		reasons = append(reasons, "Has previous context in this repository")
	}

	// 2. Issue Difficulty Match
	difficulty := EstimateDifficulty(issue)
	if difficulty.Level == "Trivial" || difficulty.Level == "Easy" {
		if commitCount == 0 {
			score += 25.0
			reasons = append(reasons, "Perfect starting difficulty for a new contributor")
		} else {
			score += 10.0
			reasons = append(reasons, "Easy issue, quick win for existing contributor")
		}
	} else if difficulty.Level == "Expert" || difficulty.Level == "Hard" {
		if commitCount == 0 {
			score -= 30.0
			reasons = append(reasons, "Extremely difficult issue for a first-time contributor")
		} else {
			score += 10.0
			reasons = append(reasons, "Appropriate challenge for experienced contributor")
		}
	}

	// Bound the score
	score = math.Max(0.0, math.Min(100.0, score))

	// Categorize Level
	level := "Moderate"
	if score >= 80 {
		level = "Excellent"
	} else if score >= 60 {
		level = "Good"
	} else if score >= 40 {
		level = "Moderate"
	} else {
		level = "Poor"
	}

	return ContributorFit{
		IssueNumber: issue.Number,
		Contributor: contributor,
		Score:       score,
		Level:       level,
		Reasons:     reasons,
	}
}
