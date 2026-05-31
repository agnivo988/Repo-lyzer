package risk

import (
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer/core"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// MaintainerRiskAnalyzer detects risks associated with contributor concentration and overload
type MaintainerRiskAnalyzer struct {
	Engine *core.WeightedScoreEngine
}

// NewMaintainerRiskAnalyzer initializes the maintainer risk analyzer
func NewMaintainerRiskAnalyzer() *MaintainerRiskAnalyzer {
	thresholds := core.Thresholds{
		Warning:   40,
		Healthy:   0,
		Excellent: 80,
	}
	return &MaintainerRiskAnalyzer{
		Engine: core.NewWeightedScoreEngine(100.0, thresholds),
	}
}

// AnalyzeMaintainerRisk evaluates maintainer overload and contributor concentration
func (m *MaintainerRiskAnalyzer) AnalyzeMaintainerRisk(
	repo *github.Repo,
	contributors []github.Contributor,
	issues []github.Issue,
) (float64, RiskCategory) {
	metrics := []core.Metric{}

	// 1. Maintainer Overload Score
	// Check if the repo has many open issues but extremely few contributors
	overloadRiskScore := 0.0
	if repo.OpenIssues > 50 && len(contributors) <= 2 {
		overloadRiskScore = 90.0 // Very high risk of overload
	} else if repo.OpenIssues > 20 && len(contributors) <= 3 {
		overloadRiskScore = 60.0 // Warning
	} else if repo.OpenIssues > 0 && len(contributors) >= 10 {
		overloadRiskScore = 0.0 // Stable
	}

	metrics = append(metrics, core.Metric{
		Name:        "Maintainer Overload Risk",
		Score:       overloadRiskScore,
		Weight:      2.0,
		Description: "Assesses the ratio of open burden (issues/PRs) compared to the active contributor base",
	})

	// 2. Contributor Concentration
	concentrationRiskScore := 0.0
	if len(contributors) > 0 {
		totalCommits := 0
		for _, c := range contributors {
			totalCommits += c.Commits
		}
		topContributorRatio := float64(contributors[0].Commits) / float64(totalCommits)

		switch {
		case topContributorRatio > 0.85:
			concentrationRiskScore = 100.0 // Critical risk, single point of failure
		case topContributorRatio > 0.60:
			concentrationRiskScore = 60.0 // Warning
		default:
			concentrationRiskScore = 0.0 // Stable
		}
	} else {
		concentrationRiskScore = 50.0 // Unknown risk
	}

	metrics = append(metrics, core.Metric{
		Name:        "Contributor Concentration",
		Score:       concentrationRiskScore,
		Weight:      1.5,
		Description: "Risk of a single contributor holding too much institutional knowledge",
	})

	// 3. Inactive Maintainer Detection
	inactiveRiskScore := 0.0
	now := time.Now()
	if len(issues) > 0 {
		staleIssues := 0
		for _, issue := range issues {
			if issue.State == "open" && now.Sub(issue.UpdatedAt).Hours()/24 > 180 {
				staleIssues++
			}
		}
		staleRatio := float64(staleIssues) / float64(len(issues))
		if staleRatio > 0.5 {
			inactiveRiskScore = 80.0
		} else if staleRatio > 0.2 {
			inactiveRiskScore = 40.0
		}
	}
	metrics = append(metrics, core.Metric{
		Name:        "Inactive Maintainer Risk",
		Score:       inactiveRiskScore,
		Weight:      1.0,
		Description: "Detects if maintainers have abandoned triage duties",
	})

	score, _ := m.Engine.CalculateScore(metrics)

	var category RiskCategory
	if score >= 80.0 {
		category = RiskCritical
	} else if score >= 40.0 {
		category = RiskWarning
	} else {
		category = RiskStable
	}

	return score, category
}
