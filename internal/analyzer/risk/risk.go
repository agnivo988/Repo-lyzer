package risk

import (
	"github.com/agnivo988/Repo-lyzer/internal/analyzer/core"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// RiskAnalyzer orchestrates risk metrics across the repository
type RiskAnalyzer struct {
	Engine *core.WeightedScoreEngine
}

// NewRiskAnalyzer initializes a risk analyzer with inverted thresholds
// For risk, a HIGH score (e.g. 100) means HIGH RISK (Critical),
// and a LOW score (e.g. 0) means LOW RISK (Healthy).
// We configure the core engine thresholds accordingly.
func NewRiskAnalyzer() *RiskAnalyzer {
	// Let's define the thresholds for risk (0-100 scale where higher is riskier):
	// >= 80 : Critical Risk
	// >= 50 : Warning Risk
	// >= 0  : Healthy (Stable)

	// Since our core engine returns 'Excellent' for highest scores, we need to map categories correctly.
	// We'll create custom wrapper functions below to translate the core engine output to Risk terminology.
	thresholds := core.Thresholds{
		Warning:   50,
		Healthy:   0,  // Technically not used directly since we remap
		Excellent: 80, // We will map Excellent to Critical for risk
	}
	return &RiskAnalyzer{
		Engine: core.NewWeightedScoreEngine(100.0, thresholds),
	}
}

// RiskCategory represents the inverted risk levels
type RiskCategory string

const (
	RiskCritical RiskCategory = "CRITICAL"
	RiskWarning  RiskCategory = "WARNING"
	RiskStable   RiskCategory = "STABLE"
)

// TranslateScoreToRisk maps a 0-100 risk score to a RiskCategory
// where higher score = higher risk.
func (r *RiskAnalyzer) TranslateScoreToRisk(score float64) RiskCategory {
	if score >= 80.0 {
		return RiskCritical
	} else if score >= 40.0 {
		return RiskWarning
	}
	return RiskStable
}

// CalculateOverallRisk computes a rudimentary base risk score (to be expanded by sub-analyzers)
func (r *RiskAnalyzer) CalculateOverallRisk(repo *github.Repo, commits []github.Commit) (float64, RiskCategory) {
	metrics := []core.Metric{}

	// Basic placeholder risk: If a repo is archived, it's immediately high risk
	archivedScore := 0.0
	if repo.Archived {
		archivedScore = 100.0 // Max risk
	}
	metrics = append(metrics, core.Metric{
		Name:        "Archived Status",
		Score:       archivedScore,
		Weight:      5.0,
		Description: "Flag if the repository is officially archived by the owner",
	})

	score, _ := r.Engine.CalculateScore(metrics)
	return score, r.TranslateScoreToRisk(score)
}
