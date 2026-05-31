package contribution

import (
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestCalculate(t *testing.T) {
	// Setup mock inputs
	now := time.Now()

	// 1. Minimum case: nothing set
	scoreMin := Calculate(
		false,
		"",
		[]github.Issue{},
		[]github.Commit{},
		[]github.Contributor{},
	)
	// Expect score = 1.0 (from low stale issues ratio since there are 0 issues)
	if scoreMin.Score != 1.0 {
		t.Errorf("Expected minimum score to be 1.0, got %f", scoreMin.Score)
	}
	if scoreMin.Level != "Needs Improvement 🔴" {
		t.Errorf("Expected level to be 'Needs Improvement 🔴', got '%s'", scoreMin.Level)
	}

	// 2. Maximum case: everything set
	// Mock issue with "good first issue" label
	issues := []github.Issue{
		{
			State:     "open",
			UpdatedAt: now,
			Labels: []github.IssueLabel{
				{Name: "good first issue"},
			},
		},
	}
	// Mock commits: recent one
	var mockCommit github.Commit
	mockCommit.Commit.Author.Date = now.Add(-2 * time.Hour)
	mockCommit.Author = &struct {
		Login string `json:"login"`
	}{Login: "dev1"}
	commits := []github.Commit{mockCommit}
	contributors := []github.Contributor{
		{Login: "dev1", Commits: 100},
	}

	scoreMax := Calculate(
		true,
		"Here is installation and setup info.",
		issues,
		commits,
		contributors,
	)

	// Expect score:
	// Contributing (2.0)
	// Readme (2.0)
	// Good first issue (1.5)
	// Recent commit (1.5)
	// Active maintainer (2.0)
	// Low stale issues (1.0)
	// Total = 10.0
	if scoreMax.Score != 10.0 {
		t.Errorf("Expected maximum score to be 10.0, got %f", scoreMax.Score)
	}
	if scoreMax.Level != "Contributor Friendly 🟢" {
		t.Errorf("Expected level to be 'Contributor Friendly 🟢', got '%s'", scoreMax.Level)
	}

	// 3. Stale open issues check
	staleIssues := []github.Issue{
		{
			State:     "open",
			UpdatedAt: now.AddDate(0, 0, -90), // 90 days old -> stale
			Labels:    []github.IssueLabel{},
		},
	}
	// Total open = 1, stale = 1, ratio = 100% (>= 30%), so no point (0 pt)
	scoreStale := Calculate(
		false,
		"",
		staleIssues,
		[]github.Commit{},
		[]github.Contributor{},
	)
	// Expect score = 0.0 (contributing=0, readme=0, good_first=0, recent_commit=0, active_maintainer=0, stale_issues=0)
	if scoreStale.Score != 0.0 {
		t.Errorf("Expected stale score to be 0.0, got %f", scoreStale.Score)
	}
}
