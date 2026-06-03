package health

import (
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer/core"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// MaintainerAnalyzer uses the core weighted scoring engine for maintainer/stale metrics
type MaintainerAnalyzer struct {
	Engine *core.WeightedScoreEngine
}

// NewMaintainerAnalyzer initializes a maintainer analyzer with default thresholds
func NewMaintainerAnalyzer() *MaintainerAnalyzer {
	thresholds := core.Thresholds{
		Warning:   40,
		Healthy:   75,
		Excellent: 90,
	}
	return &MaintainerAnalyzer{
		Engine: core.NewWeightedScoreEngine(100.0, thresholds),
	}
}

// AnalyzeMaintainerActivity computes scores for maintainer activity and repository freshness
func (m *MaintainerAnalyzer) AnalyzeMaintainerActivity(
	repo *github.Repo,
	commits []github.Commit,
	issues []github.Issue,
	pulls []github.PullRequest,
) (float64, core.ScoreCategory) {
	metrics := []core.Metric{}

	// 1. Maintainer Commit Frequency (Using the repository owner or top committers)
	// For simplicity, we define "recent" as within the last 90 days
	recentCommits := 0
	maintainerCommits := 0
	now := time.Now()

	for _, commit := range commits {
		if now.Sub(commit.Commit.Author.Date) <= 90*24*time.Hour {
			recentCommits++
			// Heuristic: If committer matches repo owner, it's a maintainer commit
			if commit.Author != nil && commit.Author.Login == repo.Owner.Login {
				maintainerCommits++
			}
		}
	}

	maintainerScore := 0.0
	if recentCommits > 0 {
		// Calculate percentage of recent commits by the maintainer
		ratio := float64(maintainerCommits) / float64(recentCommits)
		if ratio > 0.1 {
			maintainerScore = 100.0 // Healthy if maintainer is active
		} else {
			maintainerScore = ratio * 10.0 * 100.0
		}
	} else if len(commits) > 0 {
		maintainerScore = 20.0 // Some activity but not recent
	}
	metrics = append(metrics, core.Metric{
		Name:        "Maintainer Commit Frequency",
		Score:       maintainerScore,
		Weight:      2.0,
		Description: "Measures how actively the maintainer commits code recently",
	})

	// 2. Release Freshness / Repository Freshness
	freshnessScore := 0.0
	if !repo.PushedAt.IsZero() {
		daysSincePush := now.Sub(repo.PushedAt).Hours() / 24
		switch {
		case daysSincePush <= 30:
			freshnessScore = 100.0
		case daysSincePush <= 90:
			freshnessScore = 80.0
		case daysSincePush <= 180:
			freshnessScore = 50.0
		case daysSincePush <= 365:
			freshnessScore = 20.0
		default:
			freshnessScore = 0.0
		}
	}
	metrics = append(metrics, core.Metric{
		Name:        "Release Freshness",
		Score:       freshnessScore,
		Weight:      1.5,
		Description: "Assesses how recently the repository was updated or released",
	})

	// 3. Stale Repository Detection (Issue Inactivity)
	issueActivityScore := 100.0
	staleIssuesCount := 0
	if len(issues) > 0 {
		for _, issue := range issues {
			if issue.State == "open" && now.Sub(issue.UpdatedAt).Hours()/24 > 180 {
				staleIssuesCount++
			}
		}
		staleRatio := float64(staleIssuesCount) / float64(len(issues))
		if staleRatio > 0.5 {
			issueActivityScore = 30.0
		} else if staleRatio > 0.2 {
			issueActivityScore = 60.0
		} else {
			issueActivityScore = 100.0
		}
	}
	metrics = append(metrics, core.Metric{
		Name:        "Issue Inactivity",
		Score:       issueActivityScore,
		Weight:      1.0,
		Description: "Detects proportion of stale or inactive issues",
	})

	// 4. Stale Repository Detection (PR Inactivity)
	prActivityScore := 100.0
	stalePRsCount := 0
	if len(pulls) > 0 {
		for _, pr := range pulls {
			if pr.State == "open" && now.Sub(pr.UpdatedAt).Hours()/24 > 90 {
				stalePRsCount++
			}
		}
		staleRatio := float64(stalePRsCount) / float64(len(pulls))
		if staleRatio > 0.3 {
			prActivityScore = 20.0
		} else if staleRatio > 0.1 {
			prActivityScore = 60.0
		} else {
			prActivityScore = 100.0
		}
	}
	metrics = append(metrics, core.Metric{
		Name:        "PR Inactivity",
		Score:       prActivityScore,
		Weight:      1.5,
		Description: "Detects proportion of abandoned or stale pull requests",
	})

	return m.Engine.CalculateScore(metrics)
}
