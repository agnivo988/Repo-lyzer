package analyzer

import (
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// TestAnalyzeContributorActivityGrowing verifies that increasing activity is correctly classified
func TestAnalyzeContributorActivityGrowing(t *testing.T) {
	now := time.Now()

	// Create commits: 20 in last 90 days, 5 in previous 90 days (90-180)
	commits := []github.Commit{}

	// Recent period (last 90 days) - 20 commits
	for i := 0; i < 20; i++ {
		c := github.Commit{}
		c.Commit.Author.Date = now.AddDate(0, 0, -(i % 85))
		commits = append(commits, c)
	}

	// Previous period (days 90-180) - 5 commits
	for i := 0; i < 5; i++ {
		c := github.Commit{}
		c.Commit.Author.Date = now.AddDate(0, 0, -(95 + i*10))
		commits = append(commits, c)
	}

	result := AnalyzeContributorActivity(commits)

	if result.Trend != "Growing" {
		t.Errorf("expected trend 'Growing', got '%s'", result.Trend)
	}
	if result.Last90Days != 20 {
		t.Errorf("expected Last90Days=20, got %d", result.Last90Days)
	}
	if result.Last180Days != 25 {
		t.Errorf("expected Last180Days=25, got %d", result.Last180Days)
	}
	if result.Insight != "Commit activity has increased recently" {
		t.Errorf("expected insight about activity increasing, got '%s'", result.Insight)
	}
}

// TestAnalyzeContributorActivityDeclining verifies that decreasing activity is correctly classified
func TestAnalyzeContributorActivityDeclining(t *testing.T) {
	now := time.Now()

	// Create commits: 5 in last 90 days, 20 in previous 90 days (90-180)
	commits := []github.Commit{}

	// Recent period (last 90 days) - 5 commits
	for i := 0; i < 5; i++ {
		c := github.Commit{}
		c.Commit.Author.Date = now.AddDate(0, 0, -(i % 85))
		commits = append(commits, c)
	}

	// Previous period (days 90-180) - 20 commits
	for i := 0; i < 20; i++ {
		c := github.Commit{}
		c.Commit.Author.Date = now.AddDate(0, 0, -(95 + i*4))
		commits = append(commits, c)
	}

	result := AnalyzeContributorActivity(commits)

	if result.Trend != "Declining" {
		t.Errorf("expected trend 'Declining', got '%s'", result.Trend)
	}
	if result.Last90Days != 5 {
		t.Errorf("expected Last90Days=5, got %d", result.Last90Days)
	}
	if result.Last180Days != 25 {
		t.Errorf("expected Last180Days=25, got %d", result.Last180Days)
	}
	if result.Insight != "Commit activity has decreased in recent months" {
		t.Errorf("expected insight about activity decreasing, got '%s'", result.Insight)
	}
}

// TestAnalyzeContributorActivityStable verifies that consistent activity is correctly classified
func TestAnalyzeContributorActivityStable(t *testing.T) {
	now := time.Now()

	// Create commits: 10 in last 90 days, 10 in previous 90 days (90-180)
	commits := []github.Commit{}

	// Recent period (last 90 days) - 10 commits
	for i := 0; i < 10; i++ {
		c := github.Commit{}
		c.Commit.Author.Date = now.AddDate(0, 0, -(i % 85))
		commits = append(commits, c)
	}

	// Previous period (days 90-180) - 10 commits
	for i := 0; i < 10; i++ {
		c := github.Commit{}
		c.Commit.Author.Date = now.AddDate(0, 0, -(95 + i*8))
		commits = append(commits, c)
	}

	result := AnalyzeContributorActivity(commits)

	if result.Trend != "Stable" {
		t.Errorf("expected trend 'Stable', got '%s'", result.Trend)
	}
	if result.Last90Days != 10 {
		t.Errorf("expected Last90Days=10, got %d", result.Last90Days)
	}
	if result.Last180Days != 20 {
		t.Errorf("expected Last180Days=20, got %d", result.Last180Days)
	}
	if result.Insight != "Commit activity is stable" {
		t.Errorf("expected insight about activity being stable, got '%s'", result.Insight)
	}
}

// TestAnalyzeContributorActivityZeroCommits verifies handling of repositories with no commits
func TestAnalyzeContributorActivityZeroCommits(t *testing.T) {
	commits := []github.Commit{}

	result := AnalyzeContributorActivity(commits)

	if result.Trend != "Stable" {
		t.Errorf("expected trend 'Stable' for zero commits, got '%s'", result.Trend)
	}
	if result.Last90Days != 0 {
		t.Errorf("expected Last90Days=0, got %d", result.Last90Days)
	}
	if result.Last180Days != 0 {
		t.Errorf("expected Last180Days=0, got %d", result.Last180Days)
	}
	if result.Insight != "No commit activity in the last 90 days" {
		t.Errorf("expected insight about no recent activity, got '%s'", result.Insight)
	}
}

// TestAnalyzeContributorActivityOnlyRecentCommits verifies handling of repositories with only recent commits
func TestAnalyzeContributorActivityOnlyRecentCommits(t *testing.T) {
	now := time.Now()

	// Create commits: 15 in last 90 days only
	commits := []github.Commit{}
	for i := 0; i < 15; i++ {
		c := github.Commit{}
		c.Commit.Author.Date = now.AddDate(0, 0, -(i % 85))
		commits = append(commits, c)
	}

	result := AnalyzeContributorActivity(commits)

	// 15 recent commits vs 0 old commits = growing
	if result.Trend != "Growing" {
		t.Errorf("expected trend 'Growing' for recent-only commits, got '%s'", result.Trend)
	}
	if result.Last90Days != 15 {
		t.Errorf("expected Last90Days=15, got %d", result.Last90Days)
	}
	if result.Last180Days != 15 {
		t.Errorf("expected Last180Days=15, got %d", result.Last180Days)
	}
}

// TestAnalyzeContributorActivityOnlyOldCommits verifies handling of repositories with only old commits
func TestAnalyzeContributorActivityOnlyOldCommits(t *testing.T) {
	now := time.Now()

	// Create commits: 12 in days 90-180 only
	commits := []github.Commit{}
	for i := 0; i < 12; i++ {
		c := github.Commit{}
		c.Commit.Author.Date = now.AddDate(0, 0, -(100 + i*6))
		commits = append(commits, c)
	}

	result := AnalyzeContributorActivity(commits)

	// 0 recent commits vs 12 old commits = declining
	if result.Trend != "Declining" {
		t.Errorf("expected trend 'Declining' for old-only commits, got '%s'", result.Trend)
	}
	if result.Last90Days != 0 {
		t.Errorf("expected Last90Days=0, got %d", result.Last90Days)
	}
	if result.Last180Days != 12 {
		t.Errorf("expected Last180Days=12, got %d", result.Last180Days)
	}
	if result.Insight != "No commit activity in the last 90 days" {
		t.Errorf("expected insight about no recent activity, got '%s'", result.Insight)
	}
}
