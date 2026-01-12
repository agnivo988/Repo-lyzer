package analyzer

import (
	"testing"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestAnalyzeContributors_Empty(t *testing.T) {
	insights := AnalyzeContributors([]github.Contributor{})

	if insights.TotalContributors != 0 {
		t.Errorf("TotalContributors = %d, want 0", insights.TotalContributors)
	}
	if insights.ConcentrationRisk != "Unknown" {
		t.Errorf("ConcentrationRisk = %s, want Unknown", insights.ConcentrationRisk)
	}
	if insights.TeamSize != "None" {
		t.Errorf("TeamSize = %s, want None", insights.TeamSize)
	}
}

func TestAnalyzeContributors_SingleContributor(t *testing.T) {
	contributors := []github.Contributor{
		{Login: "solo-dev", Commits: 100},
	}

	insights := AnalyzeContributors(contributors)

	if insights.TotalContributors != 1 {
		t.Errorf("TotalContributors = %d, want 1", insights.TotalContributors)
	}
	if insights.TeamSize != "Solo" {
		t.Errorf("TeamSize = %s, want Solo", insights.TeamSize)
	}
	if insights.ConcentrationRisk != "Critical" {
		t.Errorf("ConcentrationRisk = %s, want Critical", insights.ConcentrationRisk)
	}
	if insights.TopContributor == nil {
		t.Fatal("TopContributor should not be nil")
	}
	if insights.TopContributor.Login != "solo-dev" {
		t.Errorf("TopContributor.Login = %s, want solo-dev", insights.TopContributor.Login)
	}
	if insights.TopContributor.Percentage != 100 {
		t.Errorf("TopContributor.Percentage = %.1f, want 100", insights.TopContributor.Percentage)
	}
}

func TestAnalyzeContributors_BalancedTeam(t *testing.T) {
	contributors := []github.Contributor{
		{Login: "dev1", Commits: 25},
		{Login: "dev2", Commits: 25},
		{Login: "dev3", Commits: 25},
		{Login: "dev4", Commits: 25},
	}

	insights := AnalyzeContributors(contributors)

	if insights.TotalContributors != 4 {
		t.Errorf("TotalContributors = %d, want 4", insights.TotalContributors)
	}
	// Team size is based on active contributors (>1% commits)
	// All 4 have 25% each, so all are active -> Medium team
	if insights.TeamSize != "Medium" {
		t.Errorf("TeamSize = %s, want Medium", insights.TeamSize)
	}
	if insights.ConcentrationRisk != "Low" {
		t.Errorf("ConcentrationRisk = %s, want Low", insights.ConcentrationRisk)
	}
	// Diversity should be high for balanced team
	if insights.DiversityScore < 70 {
		t.Errorf("DiversityScore = %.1f, want >= 70 for balanced team", insights.DiversityScore)
	}
}

func TestAnalyzeContributors_UnbalancedTeam(t *testing.T) {
	contributors := []github.Contributor{
		{Login: "main-dev", Commits: 900},
		{Login: "helper1", Commits: 50},
		{Login: "helper2", Commits: 30},
		{Login: "helper3", Commits: 20},
	}

	insights := AnalyzeContributors(contributors)

	if insights.ConcentrationRisk != "Critical" {
		t.Errorf("ConcentrationRisk = %s, want Critical (top contributor has 90%%)", insights.ConcentrationRisk)
	}
	// Diversity should be low for unbalanced team
	if insights.DiversityScore > 30 {
		t.Errorf("DiversityScore = %.1f, want < 30 for unbalanced team", insights.DiversityScore)
	}
}

func TestAnalyzeContributors_LargeTeam(t *testing.T) {
	contributors := make([]github.Contributor, 15)
	for i := 0; i < 15; i++ {
		contributors[i] = github.Contributor{
			Login:   "dev",
			Commits: 100 - i*5, // Decreasing commits
		}
	}

	insights := AnalyzeContributors(contributors)

	if insights.TotalContributors != 15 {
		t.Errorf("TotalContributors = %d, want 15", insights.TotalContributors)
	}
	if insights.TeamSize != "Large" {
		t.Errorf("TeamSize = %s, want Large", insights.TeamSize)
	}
}

func TestContributorDetail_Classification(t *testing.T) {
	tests := []struct {
		commits    int
		percentage float64
		wantType   string
	}{
		{500, 50, "Core"},
		{200, 15, "Regular"},
		{50, 3, "Occasional"},
		{5, 0.5, "New"},
	}

	for _, tt := range tests {
		got := classifyContributor(tt.commits, tt.percentage)
		if got != tt.wantType {
			t.Errorf("classifyContributor(%d, %.1f) = %s, want %s",
				tt.commits, tt.percentage, got, tt.wantType)
		}
	}
}

func TestCommitDistribution(t *testing.T) {
	contributors := []github.Contributor{
		{Login: "top", Commits: 50},
		{Login: "mid1", Commits: 20},
		{Login: "mid2", Commits: 15},
		{Login: "low1", Commits: 10},
		{Login: "low2", Commits: 5},
	}

	insights := AnalyzeContributors(contributors)
	dist := insights.CommitDistribution

	// Top 1% (1 person) should have 50% of commits
	if dist.Top1Percent != 50 {
		t.Errorf("Top1Percent = %.1f, want 50", dist.Top1Percent)
	}

	// Top 50% should have most commits
	if dist.Top50Percent < 70 {
		t.Errorf("Top50Percent = %.1f, want >= 70", dist.Top50Percent)
	}

	// Gini should be between 0 and 1
	if dist.GiniCoefficient < 0 || dist.GiniCoefficient > 1 {
		t.Errorf("GiniCoefficient = %.2f, want between 0 and 1", dist.GiniCoefficient)
	}
}

func TestRecommendations(t *testing.T) {
	// Critical concentration should generate warnings
	contributors := []github.Contributor{
		{Login: "solo", Commits: 100},
	}

	insights := AnalyzeContributors(contributors)

	if len(insights.Recommendations) == 0 {
		t.Error("Expected recommendations for solo contributor")
	}

	// Check that recommendations contain relevant advice
	hasWarning := false
	for _, rec := range insights.Recommendations {
		if len(rec) > 0 {
			hasWarning = true
			break
		}
	}
	if !hasWarning {
		t.Error("Expected warning recommendations for critical concentration")
	}
}

func TestDiversityScore_Bounds(t *testing.T) {
	// Test with various team compositions
	testCases := []struct {
		name         string
		contributors []github.Contributor
	}{
		{
			name: "single contributor",
			contributors: []github.Contributor{
				{Login: "solo", Commits: 100},
			},
		},
		{
			name: "perfectly balanced",
			contributors: []github.Contributor{
				{Login: "a", Commits: 25},
				{Login: "b", Commits: 25},
				{Login: "c", Commits: 25},
				{Login: "d", Commits: 25},
			},
		},
		{
			name: "highly unbalanced",
			contributors: []github.Contributor{
				{Login: "main", Commits: 99},
				{Login: "other", Commits: 1},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			insights := AnalyzeContributors(tc.contributors)
			if insights.DiversityScore < 0 || insights.DiversityScore > 100 {
				t.Errorf("DiversityScore = %.1f, want between 0 and 100", insights.DiversityScore)
			}
		})
	}
}
