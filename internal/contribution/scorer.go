package contribution

import (
	"fmt"
	"strings"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// ContributionScore represents the computed contribution friendliness metrics.
type ContributionScore struct {
	Score      float64  `json:"score"`
	Level      string   `json:"level"`
	Strengths  []string `json:"strengths"`
	Weaknesses []string `json:"weaknesses"`
}

var readmeKeywords = []string{
	"installation",
	"setup",
	"getting started",
	"usage",
}

func hasSetupSection(readmeContent string) bool {
	contentLower := strings.ToLower(readmeContent)
	for _, kw := range readmeKeywords {
		if strings.Contains(contentLower, kw) {
			return true
		}
	}
	return false
}

// Calculate computes the contribution score based on various repository metrics.
func Calculate(
	hasContributing bool,
	readmeContent string,
	issues []github.Issue,
	commits []github.Commit,
	contributors []github.Contributor,
) ContributionScore {
	score := 0.0
	var strengths, weaknesses []string

	// 1. CONTRIBUTING.md — 2.0 pts
	if hasContributing {
		score += 2.0
		strengths = append(strengths, "CONTRIBUTING.md exists")
	} else {
		weaknesses = append(weaknesses, "No CONTRIBUTING.md found")
	}

	// 2. README keywords — 2.0 pts
	if hasSetupSection(readmeContent) {
		score += 2.0
		strengths = append(strengths, "README has setup/installation guide")
	} else {
		weaknesses = append(weaknesses, "README missing setup/installation section")
	}

	// 3. Good first issues — 1.5 pts
	hasGoodFirstIssue := false
	goodFirstIssueLabels := map[string]bool{
		"good first issue":     true,
		"good-first-issue":     true,
		"beginner-friendly":    true,
		"easy":                 true,
		"first-timers-only":    true,
		"first timers only":    true,
		"easy-fix":             true,
		"easy fix":             true,
		"starter bug":          true,
		"starter-bug":          true,
		"contribution welcome": true,
	}

	totalOpenIssues := 0
	staleOpenIssues := 0
	cutoffStale := time.Now().AddDate(0, 0, -60) // 60 days ago

	for _, issue := range issues {
		// Skip Pull Requests (GitHub issues API includes PRs)
		if issue.PullRequest != nil {
			continue
		}

		if issue.State == "open" {
			totalOpenIssues++
			if issue.UpdatedAt.Before(cutoffStale) {
				staleOpenIssues++
			}

			for _, label := range issue.Labels {
				labelNameLower := strings.ToLower(label.Name)
				if goodFirstIssueLabels[labelNameLower] {
					hasGoodFirstIssue = true
				}
			}
		}
	}

	if hasGoodFirstIssue {
		score += 1.5
		strengths = append(strengths, "Good first issues are available")
	} else {
		weaknesses = append(weaknesses, "No good first issues found")
	}

	// 4. Recent commits (last 14 days) — 1.5 pts
	cutoffRecentCommits := time.Now().AddDate(0, 0, -14)
	hasRecentCommits := false
	for _, commit := range commits {
		if commit.Commit.Author.Date.After(cutoffRecentCommits) {
			hasRecentCommits = true
			break
		}
	}

	if hasRecentCommits {
		score += 1.5
		strengths = append(strengths, "Active recent development (commits in last 14 days)")
	} else {
		weaknesses = append(weaknesses, "No commits in the last 14 days")
	}

	// 5. Active maintainers (commit in last 30 days) — 2.0 pts
	cutoffActiveMaintainers := time.Now().AddDate(0, 0, -30)
	hasActiveMaintainers := false

	// Identify core maintainers (e.g. top 3 contributors)
	maintainers := make(map[string]bool)
	maxMaintainers := 3
	if len(contributors) < maxMaintainers {
		maxMaintainers = len(contributors)
	}
	for i := 0; i < maxMaintainers; i++ {
		maintainers[strings.ToLower(contributors[i].Login)] = true
	}

	for _, commit := range commits {
		if commit.Commit.Author.Date.After(cutoffActiveMaintainers) {
			if commit.Author != nil && maintainers[strings.ToLower(commit.Author.Login)] {
				hasActiveMaintainers = true
				break
			}
		}
	}

	if hasActiveMaintainers {
		score += 2.0
		strengths = append(strengths, "Active maintainers (commit in last 30 days)")
	} else {
		weaknesses = append(weaknesses, "No maintainer commits in the last 30 days")
	}

	// 6. Low stale issue ratio — 1.0 pt
	staleRatio := 0.0
	if totalOpenIssues > 0 {
		staleRatio = float64(staleOpenIssues) / float64(totalOpenIssues)
	}

	if totalOpenIssues == 0 || staleRatio < 0.30 {
		score += 1.0
		strengths = append(strengths, "Low ratio of stale open issues")
	} else {
		weaknesses = append(weaknesses, fmt.Sprintf("High ratio of stale open issues (%.1f%%)", staleRatio*100))
	}

	// Calculate Level
	level := ""
	if score >= 8.0 {
		level = "Contributor Friendly 🟢"
	} else if score >= 5.0 {
		level = "Moderately Friendly 🟡"
	} else {
		level = "Needs Improvement 🔴"
	}

	return ContributionScore{
		Score:      score,
		Level:      level,
		Strengths:  strengths,
		Weaknesses: weaknesses,
	}
}
