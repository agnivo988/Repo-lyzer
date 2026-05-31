package health

import (
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer/core"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// HealthAnalyzer uses the core weighted scoring engine
type HealthAnalyzer struct {
	Engine *core.WeightedScoreEngine
}

// NewHealthAnalyzer initializes a health analyzer with default thresholds
func NewHealthAnalyzer() *HealthAnalyzer {
	thresholds := core.Thresholds{
		Warning:   40,
		Healthy:   70,
		Excellent: 90,
	}
	return &HealthAnalyzer{
		Engine: core.NewWeightedScoreEngine(100.0, thresholds),
	}
}

// CalculateRepositoryHealth computes the health score based on basic metrics
func (h *HealthAnalyzer) CalculateRepositoryHealth(repo *github.Repo, commits []github.Commit) (float64, core.ScoreCategory) {
	metrics := []core.Metric{}

	// Description Metric
	descScore := 0.0
	if repo.Description != "" {
		descScore = 100.0
	}
	metrics = append(metrics, core.Metric{
		Name:   "Description Present",
		Score:  descScore,
		Weight: 1.0,
	})

	// Stars Metric
	starsScore := 0.0
	if repo.Stars > 50 {
		starsScore = 100.0
	} else {
		starsScore = float64(repo.Stars) / 50.0 * 100.0
	}
	metrics = append(metrics, core.Metric{
		Name:   "Stars Count",
		Score:  starsScore,
		Weight: 1.0,
	})

	// Commits Metric
	commitsScore := 0.0
	if len(commits) > 10 {
		commitsScore = 100.0
	} else {
		commitsScore = float64(len(commits)) / 10.0 * 100.0
	}
	metrics = append(metrics, core.Metric{
		Name:   "Recent Commits Activity",
		Score:  commitsScore,
		Weight: 2.0, // Higher weight for activity
	})

	// Open Issues Metric
	issuesScore := 100.0
	if repo.OpenIssues >= 20 {
		issuesScore = 50.0
		if repo.OpenIssues > 100 {
			issuesScore = 0.0
		}
	}
	metrics = append(metrics, core.Metric{
		Name:   "Open Issues Volume",
		Score:  issuesScore,
		Weight: 1.0,
	})

	// Recency Metric
	recencyScore := 0.0
	if !repo.PushedAt.IsZero() {
		since := time.Since(repo.PushedAt)
		switch {
		case since <= 30*24*time.Hour:
			recencyScore = 100.0
		case since <= 90*24*time.Hour:
			recencyScore = 80.0
		case since <= 180*24*time.Hour:
			recencyScore = 50.0
		case since <= 365*24*time.Hour:
			recencyScore = 30.0
		default:
			recencyScore = 0.0
		}
	}
	metrics = append(metrics, core.Metric{
		Name:   "Last Push Recency",
		Score:  recencyScore,
		Weight: 1.5,
	})

	return h.Engine.CalculateScore(metrics)
}
