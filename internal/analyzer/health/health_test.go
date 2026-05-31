package health

import (
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer/core"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestHealthAnalyzer_CalculateRepositoryHealth(t *testing.T) {
	analyzer := NewHealthAnalyzer()

	now := time.Now()
	repo := &github.Repo{
		Description: "A great repo",
		Stars:       100,
		OpenIssues:  5,
		PushedAt:    now.AddDate(0, 0, -5), // 5 days ago
	}

	commits := make([]github.Commit, 15) // 15 commits

	score, category := analyzer.CalculateRepositoryHealth(repo, commits)

	if score != 100.0 {
		t.Errorf("Expected perfect score, got %.2f", score)
	}
	if category != core.Excellent {
		t.Errorf("Expected Excellent category, got %s", category)
	}
}

func TestHealthAnalyzer_PoorHealth(t *testing.T) {
	analyzer := NewHealthAnalyzer()

	now := time.Now()
	repo := &github.Repo{
		Description: "", // Empty description
		Stars:       0,
		OpenIssues:  150,                   // Massive open issues
		PushedAt:    now.AddDate(-2, 0, 0), // 2 years ago
	}

	commits := []github.Commit{} // No commits

	score, category := analyzer.CalculateRepositoryHealth(repo, commits)

	// Since everything is poor, score should be 0
	if score != 0.0 {
		t.Errorf("Expected score 0, got %.2f", score)
	}
	if category != core.Critical {
		t.Errorf("Expected Critical category, got %s", category)
	}
}
