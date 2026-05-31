package risk

import (
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer/core"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// ReleaseRiskAnalyzer assesses the risks associated with missing dependencies locks and stale releases
type ReleaseRiskAnalyzer struct {
	Engine *core.WeightedScoreEngine
}

// NewReleaseRiskAnalyzer initializes the release risk analyzer
func NewReleaseRiskAnalyzer() *ReleaseRiskAnalyzer {
	thresholds := core.Thresholds{
		Warning:   40,
		Healthy:   0,
		Excellent: 80,
	}
	return &ReleaseRiskAnalyzer{
		Engine: core.NewWeightedScoreEngine(100.0, thresholds),
	}
}

// AnalyzeReleaseAndDependencyRisk evaluates release cadence and lockfile hygiene
func (r *ReleaseRiskAnalyzer) AnalyzeReleaseAndDependencyRisk(releases []github.Release, tree []github.TreeEntry) (float64, RiskCategory) {
	metrics := []core.Metric{}

	// 1. Dependency Lockfile Presence (Hygiene Risk)
	// If a project doesn't pin its dependencies, it's highly susceptible to supply chain attacks
	hasLockfile := false
	for _, entry := range tree {
		if entry.Type == "blob" {
			switch entry.Path {
			case "go.sum", "package-lock.json", "yarn.lock", "pnpm-lock.yaml", "Cargo.lock", "poetry.lock", "Gemfile.lock", "composer.lock":
				hasLockfile = true
				break
			}
		}
	}

	depRiskScore := 100.0 // Assume critical risk if no lockfile
	if len(tree) == 0 {
		// Tree wasn't loaded or empty repo, neutral
		depRiskScore = 0.0
	} else if hasLockfile {
		depRiskScore = 0.0 // Stable, lockfile exists
	}

	metrics = append(metrics, core.Metric{
		Name:        "Missing Dependency Lockfile",
		Score:       depRiskScore,
		Weight:      3.0,
		Description: "High risk if project lacks pinned dependency lockfiles (e.g. package-lock.json, go.sum)",
	})

	// 2. Release Cadence Risk
	releaseRiskScore := 0.0
	if len(releases) == 0 {
		releaseRiskScore = 80.0 // High risk if the repository has absolutely zero releases
	} else {
		// Calculate the time since the latest release
		latestRelease := releases[0]
		for _, rel := range releases {
			if rel.PublishedAt.After(latestRelease.PublishedAt) {
				latestRelease = rel
			}
		}

		daysSinceRelease := time.Since(latestRelease.PublishedAt).Hours() / 24
		switch {
		case daysSinceRelease > 365:
			releaseRiskScore = 90.0 // Stale release, high risk
		case daysSinceRelease > 180:
			releaseRiskScore = 50.0 // Warning
		default:
			releaseRiskScore = 0.0 // Stable
		}
	}

	metrics = append(metrics, core.Metric{
		Name:        "Stale Release Cadence",
		Score:       releaseRiskScore,
		Weight:      2.0,
		Description: "Assesses risk of abandoned or highly stale releases in production",
	})

	score, _ := r.Engine.CalculateScore(metrics)

	// Translate using the same risk mapping logic
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
