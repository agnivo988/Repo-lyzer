package health

import (
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer/core"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestPRVelocityAnalyzer_AnalyzePRVelocity(t *testing.T) {
	analyzer := NewPRVelocityAnalyzer()
	now := time.Now()

	t.Run("No PRs", func(t *testing.T) {
		score, category := analyzer.AnalyzePRVelocity([]github.PullRequest{})
		if score != 0.0 {
			t.Errorf("expected 0.0 for no PRs, got %v", score)
		}
		if category != core.Warning {
			t.Errorf("expected Warning category for no PRs, got %v", category)
		}
	})

	t.Run("Healthy PRs", func(t *testing.T) {
		mergedAt1 := now.Add(-1 * time.Hour) // Fast merge
		mergedAt2 := now.Add(-1 * time.Hour) // Fast merge
		pulls := []github.PullRequest{
			{State: "closed", CreatedAt: now.Add(-2 * time.Hour), MergedAt: &mergedAt1},
			{State: "closed", CreatedAt: now.Add(-3 * time.Hour), MergedAt: &mergedAt2},
		}

		score, category := analyzer.AnalyzePRVelocity(pulls)
		if score != 100.0 {
			t.Errorf("expected 100.0 for fast merged PRs, got %.2f", score)
		}
		if category != core.Excellent {
			t.Errorf("expected Excellent category, got %v", category)
		}
	})

	t.Run("Stale and Unmerged PRs", func(t *testing.T) {
		pulls := []github.PullRequest{
			{State: "open", CreatedAt: now.Add(-60 * 24 * time.Hour)}, // Very stale open PR
			{State: "closed", CreatedAt: now.Add(-10 * time.Hour)},    // Closed without merge
		}

		score, category := analyzer.AnalyzePRVelocity(pulls)
		// Poor merge ratio and stale ratio should drop score below 50
		if score > 50.0 {
			t.Errorf("expected score < 50.0, got %.2f", score)
		}
		if category == core.Excellent || category == core.Healthy {
			t.Errorf("expected Warning or Critical category, got %v", category)
		}
	})
}
