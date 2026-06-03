package health

import (
	"sort"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer/core"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// BusFactorAnalyzer assesses risk related to contributor dependency
type BusFactorAnalyzer struct {
	Engine *core.WeightedScoreEngine
}

// NewBusFactorAnalyzer initializes the bus factor analyzer
func NewBusFactorAnalyzer() *BusFactorAnalyzer {
	thresholds := core.Thresholds{
		Warning:   40,
		Healthy:   70,
		Excellent: 85,
	}
	return &BusFactorAnalyzer{
		Engine: core.NewWeightedScoreEngine(100.0, thresholds),
	}
}

// AnalyzeContributorStability computes the bus factor and overall contributor stability
func (b *BusFactorAnalyzer) AnalyzeContributorStability(contributors []github.Contributor) (float64, core.ScoreCategory) {
	metrics := []core.Metric{}

	if len(contributors) == 0 {
		return 0, core.Critical
	}

	totalCommits := 0
	for _, c := range contributors {
		totalCommits += c.Commits
	}

	// 1. Bus Factor Score
	// Calculates the minimum number of contributors required to account for >50% of the commits
	// Sort contributors by descending commit count to correctly compute bus factor and top ratio
	sort.Slice(contributors, func(i, j int) bool {
		return contributors[i].Commits > contributors[j].Commits
	})

	commitsCount := 0
	busFactor := 0

	for _, c := range contributors {
		commitsCount += c.Commits
		busFactor++
		if float64(commitsCount) > float64(totalCommits)*0.5 {
			break
		}
	}

	busFactorScore := 0.0
	switch {
	case busFactor >= 5:
		busFactorScore = 100.0
	case busFactor >= 3:
		busFactorScore = 80.0
	case busFactor == 2:
		busFactorScore = 50.0
	case busFactor == 1:
		busFactorScore = 20.0
	}
	metrics = append(metrics, core.Metric{
		Name:        "Bus Factor",
		Score:       busFactorScore,
		Weight:      3.0,
		Description: "Assesses how heavily the project relies on a few core contributors",
	})

	// 2. Contributor Diversity
	// Score based on the total number of contributors
	diversityScore := 0.0
	switch {
	case len(contributors) >= 20:
		diversityScore = 100.0
	case len(contributors) >= 10:
		diversityScore = 80.0
	case len(contributors) >= 5:
		diversityScore = 50.0
	default:
		diversityScore = 20.0
	}
	metrics = append(metrics, core.Metric{
		Name:        "Contributor Diversity",
		Score:       diversityScore,
		Weight:      1.0,
		Description: "Evaluates the overall size and diversity of the contributor base",
	})

	// 3. Maintainer Dependency Concentration
	// Checks if the absolute top contributor dominates an unhealthy amount (e.g., >80% alone)
	concentrationScore := 100.0
	if len(contributors) > 0 {
		topRatio := float64(contributors[0].Commits) / float64(totalCommits)
		if topRatio > 0.8 {
			concentrationScore = 10.0 // Extremely risky
		} else if topRatio > 0.6 {
			concentrationScore = 40.0 // High risk
		} else if topRatio > 0.4 {
			concentrationScore = 70.0 // Moderate
		}
	}
	metrics = append(metrics, core.Metric{
		Name:        "Maintainer Concentration Risk",
		Score:       concentrationScore,
		Weight:      1.5,
		Description: "Penalizes projects dominated almost entirely by a single author",
	})

	return b.Engine.CalculateScore(metrics)
}
