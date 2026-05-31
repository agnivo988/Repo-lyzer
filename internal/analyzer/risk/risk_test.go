package risk

import (
	"testing"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestRiskAnalyzer_CalculateOverallRisk(t *testing.T) {
	analyzer := NewRiskAnalyzer()

	t.Run("Stable Repo", func(t *testing.T) {
		repo := &github.Repo{
			Archived: false,
		}
		score, category := analyzer.CalculateOverallRisk(repo, nil)

		if score > 10.0 {
			t.Errorf("expected low risk score, got %.2f", score)
		}
		if category != RiskStable {
			t.Errorf("expected RiskStable, got %s", category)
		}
	})

	t.Run("Archived Repo", func(t *testing.T) {
		repo := &github.Repo{
			Archived: true, // Should trigger maximum risk immediately
		}
		score, category := analyzer.CalculateOverallRisk(repo, nil)

		if score < 80.0 {
			t.Errorf("expected critical risk score for archived repo, got %.2f", score)
		}
		if category != RiskCritical {
			t.Errorf("expected RiskCritical, got %s", category)
		}
	})
}
