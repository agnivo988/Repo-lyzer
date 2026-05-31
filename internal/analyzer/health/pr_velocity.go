package health

import (
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer/core"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// PRVelocityAnalyzer analyzes pull request merge velocity and responsiveness
type PRVelocityAnalyzer struct {
	Engine *core.WeightedScoreEngine
}

// NewPRVelocityAnalyzer initializes a PR velocity analyzer
func NewPRVelocityAnalyzer() *PRVelocityAnalyzer {
	thresholds := core.Thresholds{
		Warning:   40,
		Healthy:   70,
		Excellent: 85,
	}
	return &PRVelocityAnalyzer{
		Engine: core.NewWeightedScoreEngine(100.0, thresholds),
	}
}

// AnalyzePRVelocity processes pull requests to determine merge velocity and health
func (p *PRVelocityAnalyzer) AnalyzePRVelocity(pulls []github.PullRequest) (float64, core.ScoreCategory) {
	metrics := []core.Metric{}
	if len(pulls) == 0 {
		return 0, core.Warning // Cannot analyze without PRs, neutral/warning
	}

	var totalMergeDuration time.Duration
	var mergedCount int
	var staleCount int
	var closedWithoutMergeCount int
	var openCount int

	now := time.Now()

	for _, pr := range pulls {
		switch pr.State {
		case "open":
			openCount++
			if now.Sub(pr.CreatedAt).Hours()/24 > 30 {
				staleCount++
			}
		case "closed":
			if pr.MergedAt != nil {
				mergedCount++
				totalMergeDuration += pr.MergedAt.Sub(pr.CreatedAt)
			} else {
				closedWithoutMergeCount++
			}
		}
	}

	// 1. Average PR Merge Duration
	mergeDurationScore := 100.0
	if mergedCount > 0 {
		avgMergeHours := (totalMergeDuration.Hours() / float64(mergedCount))
		switch {
		case avgMergeHours <= 48: // Merged within 2 days
			mergeDurationScore = 100.0
		case avgMergeHours <= 168: // Merged within 1 week
			mergeDurationScore = 80.0
		case avgMergeHours <= 336: // Merged within 2 weeks
			mergeDurationScore = 50.0
		default: // Merged > 2 weeks
			mergeDurationScore = 30.0
		}
	} else if openCount > 0 {
		mergeDurationScore = 0.0 // PRs exist but none merged
	}
	metrics = append(metrics, core.Metric{
		Name:        "Average PR Merge Duration",
		Score:       mergeDurationScore,
		Weight:      2.0,
		Description: "Assesses how fast pull requests are merged",
	})

	// 2. Stale PR Ratio
	stalePRScore := 100.0
	if openCount > 0 {
		staleRatio := float64(staleCount) / float64(openCount)
		stalePRScore = (1.0 - staleRatio) * 100.0
	}
	metrics = append(metrics, core.Metric{
		Name:        "Stale PR Ratio",
		Score:       stalePRScore,
		Weight:      1.5,
		Description: "Ratio of pull requests that are open and inactive for a long time",
	})

	// 3. PR Close Rate (Merged vs Closed Without Merge)
	closeRateScore := 100.0
	totalResolved := mergedCount + closedWithoutMergeCount
	if totalResolved > 0 {
		mergeRatio := float64(mergedCount) / float64(totalResolved)
		if mergeRatio > 0.8 {
			closeRateScore = 100.0
		} else if mergeRatio > 0.5 {
			closeRateScore = 70.0
		} else {
			closeRateScore = 40.0
		}
	} else if openCount > 0 {
		closeRateScore = 50.0 // Only open PRs exist
	}
	metrics = append(metrics, core.Metric{
		Name:        "PR Resolution Ratio",
		Score:       closeRateScore,
		Weight:      1.0,
		Description: "Evaluates the ratio of successfully merged PRs vs abandoned/closed PRs",
	})

	return p.Engine.CalculateScore(metrics)
}
