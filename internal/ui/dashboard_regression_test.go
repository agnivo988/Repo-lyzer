package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	tea "github.com/charmbracelet/bubbletea"
)

func TestDashboardQualityView_KeepsTopTabsVisibleWhenContentOverflows(t *testing.T) {
	model := NewDashboardModel()
	model.currentView = viewQualityDashboard
	model.width = 80
	model.height = 12

	repo := &github.Repo{
		FullName:      "owner/repo",
		Description:   "Regression test repository",
		DefaultBranch: "main",
		HTMLURL:       "https://github.com/owner/repo",
		CreatedAt:     time.Now(),
		PushedAt:      time.Now(),
	}

	model.SetData(AnalysisResult{
		Repo: repo,
		QualityDashboard: &analyzer.QualityDashboard{
			OverallScore: 69,
			RiskLevel:    "Medium",
			QualityGrade: "C",
			KeyMetrics: analyzer.DashboardMetrics{
				HealthScore:      90,
				SecurityScore:    100,
				MaturityLevel:    "Prototype",
				BusFactor:        2,
				ActivityLevel:    "Low",
				ContributorCount: 17,
			},
			ProblemHotspots: []analyzer.ProblemHotspot{
				{Area: "Bus Factor", Severity: "High", Description: "Very low contributor diversity"},
				{Area: "Security", Severity: "Medium", Description: "Security checks need hardening"},
				{Area: "Activity", Severity: "Medium", Description: "Low activity in the last 90 days"},
				{Area: "Testing", Severity: "Low", Description: "Missing coverage on critical paths"},
				{Area: "Docs", Severity: "Low", Description: "Insufficient onboarding docs"},
			},
			Recommendations: []string{
				"👥 Encourage more contributors to reduce bus factor risk",
				"📚 Improve documentation to enable easier onboarding",
				"🧪 Add regression coverage for UI layout behavior",
				"🔒 Improve security checks in CI",
				"📈 Increase regular maintenance cadence",
			},
		},
	})

	view := model.View()

	if !strings.Contains(view, "Overview") || !strings.Contains(view, "Quality") {
		t.Fatalf("top tabs not visible in quality view output:\n%s", view)
	}
}

func TestDashboardOverviewViewShowsRepositoryTrends(t *testing.T) {
	model := NewDashboardModel()
	model.currentView = viewOverview
	model.width = 80
	model.height = 18

	repo := &github.Repo{
		FullName:      "owner/repo",
		Description:   "Regression test repository",
		DefaultBranch: "main",
		HTMLURL:       "https://github.com/owner/repo",
		CreatedAt:     time.Now(),
		PushedAt:      time.Now(),
	}

	commits := make([]github.Commit, 0)
	addCommit := func(date time.Time, login string) {
		commit := github.Commit{}
		commit.Commit.Author.Date = date
		commit.Author = &struct {
			Login string `json:"login"`
		}{Login: login}
		commits = append(commits, commit)
	}

	addCommit(time.Date(2026, 1, 5, 0, 0, 0, 0, time.UTC), "alice")
	addCommit(time.Date(2026, 2, 5, 0, 0, 0, 0, time.UTC), "alice")
	addCommit(time.Date(2026, 2, 10, 0, 0, 0, 0, time.UTC), "bob")
	addCommit(time.Date(2026, 3, 5, 0, 0, 0, 0, time.UTC), "alice")
	addCommit(time.Date(2026, 3, 10, 0, 0, 0, 0, time.UTC), "bob")
	addCommit(time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC), "carol")

	model.SetData(AnalysisResult{
		Repo:    repo,
		Commits: commits,
		Contributors: []github.Contributor{
			{Login: "alice", Commits: 3},
			{Login: "bob", Commits: 2},
			{Login: "carol", Commits: 1},
		},
		HealthScore:   78,
		BusFactor:     2,
		BusRisk:       "Medium",
		MaturityLevel: "Prototype",
	})

	view := model.View()
	for _, expected := range []string{"Repository Trends", "Health:", "Contributors:", "Forecast:"} {
		if !strings.Contains(view, expected) {
			t.Fatalf("overview view missing %q:\n%s", expected, view)
		}
	}
}

func TestDashboardNumericTabsFollowVisibleOrder(t *testing.T) {
	model := NewDashboardModel()

	checks := []struct {
		key  string
		want dashboardView
	}{
		{key: "1", want: viewOverview},
		{key: "2", want: viewQualityDashboard},
		{key: "3", want: viewRepo},
		{key: "4", want: viewLanguages},
		{key: "5", want: viewActivity},
		{key: "6", want: viewTrends},
		{key: "7", want: viewContributors},
		{key: "8", want: viewContributorInsights},
	}

	for _, check := range checks {
		updated, _ := model.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(check.key)})
		model = updated.(DashboardModel)
		if model.currentView != check.want {
			t.Fatalf("key %q set view %v, want %v", check.key, model.currentView, check.want)
		}
	}
}
