package analyzer

import (
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestAnalyzeMaintainerDashboard(t *testing.T) {
	// Mock repo
	repo := &github.Repo{
		FullName: "test/repo",
	}

	// Mock PRs
	now := time.Now()
	prs := []github.PullRequest{
		{
			Number:    1,
			Title:     "Stuck PR",
			State:     "open",
			CreatedAt: now.AddDate(0, 0, -10), // 10 days ago
			User:      github.User{Login: "dev1"},
		},
		{
			Number:    2,
			Title:     "Draft PR",
			State:     "open",
			Draft:     true,
			CreatedAt: now.AddDate(0, 0, -2), // 2 days ago
			User:      github.User{Login: "dev2"},
		},
		{
			Number:    3,
			Title:     "Active fresh PR",
			State:     "open",
			CreatedAt: now.AddDate(0, 0, -1), // 1 day ago
			User:      github.User{Login: "dev3"},
		},
	}

	// Mock Issues
	issues := []github.Issue{
		{
			Number:    101,
			Title:     "Unlabeled Issue",
			State:     "open",
			CreatedAt: now.AddDate(0, 0, -5),
			UpdatedAt: now.AddDate(0, 0, -5),
			Labels:    []github.IssueLabel{}, // No labels
			User:      github.User{Login: "user1"},
		},
		{
			Number:    102,
			Title:     "Stale Labeled Issue",
			State:     "open",
			CreatedAt: now.AddDate(0, 0, -40),
			UpdatedAt: now.AddDate(0, 0, -35), // Stale
			Labels:    []github.IssueLabel{{Name: "bug"}},
			User:      github.User{Login: "user2"},
		},
		{
			Number:      103,
			Title:       "Pull Request Issue",
			State:       "open",
			PullRequest: &struct{}{}, // Should be ignored as it's a PR
			User:        github.User{Login: "user3"},
		},
	}

	// Mock Contributors
	contributors := []github.Contributor{
		{Login: "contrib1", Commits: 100},
		{Login: "contrib2", Commits: 50},
		{Login: "contrib3", Commits: 1}, // Under limit of 2, should be ignored
	}

	// Run Analysis
	analysis := AnalyzeMaintainerDashboard(
		repo,
		prs,
		issues,
		75,    // low health (<80)
		2,     // low bus factor (<=2)
		false, // no lock file
		contributors,
	)

	// Verify PRs
	if len(analysis.PRsStuck) != 2 {
		t.Errorf("Expected 2 stuck PRs, got %d", len(analysis.PRsStuck))
	}
	if analysis.PRsStuck[0].Number != 1 || analysis.PRsStuck[1].Number != 2 {
		t.Errorf("Stuck PR numbers mismatched")
	}

	// Verify Issues
	if len(analysis.IssueCandidates) != 2 {
		t.Errorf("Expected 2 issue candidates, got %d", len(analysis.IssueCandidates))
	}

	// Verify suggested actions priorities
	hasHighHealth := false
	hasHighBusFactor := false
	hasMediumLockfile := false
	for _, action := range analysis.SuggestedActions {
		if action.Title == "Improve repository health" && action.Priority == "High" {
			hasHighHealth = true
		}
		if action.Title == "Onboard more core maintainers" && action.Priority == "High" {
			hasHighBusFactor = true
		}
		if action.Title == "Add package lock file" && action.Priority == "Medium" {
			hasMediumLockfile = true
		}
	}

	if !hasHighHealth {
		t.Error("Expected high priority action for improving health score")
	}
	if !hasHighBusFactor {
		t.Error("Expected high priority action for low bus factor")
	}
	if !hasMediumLockfile {
		t.Error("Expected medium priority action for missing lock file")
	}

	// Verify contributors deserving appreciation
	if len(analysis.ContributorsToAppreciate) != 2 {
		t.Errorf("Expected 2 contributors to appreciate, got %d", len(analysis.ContributorsToAppreciate))
	}
	if analysis.ContributorsToAppreciate[0].Login != "contrib1" || analysis.ContributorsToAppreciate[0].Contribution != "Top Contributor" {
		t.Errorf("First appreciated contributor mismatch")
	}
}
