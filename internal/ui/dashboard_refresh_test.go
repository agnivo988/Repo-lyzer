package ui

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/cache"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestCountNewIssues(t *testing.T) {
	got := countNewIssues(
		[]github.Issue{{Number: 1}, {Number: 2}, {Number: 3}},
		[]github.Issue{{Number: 2}, {Number: 3}, {Number: 4}, {Number: 5}},
	)
	if got != 2 {
		t.Fatalf("expected 2 new issues, got %d", got)
	}
}

func TestCountNewPullRequests(t *testing.T) {
	got := countNewPullRequests(
		[]github.PullRequest{{Number: 10}, {Number: 11}},
		[]github.PullRequest{{Number: 11}, {Number: 12}, {Number: 13}},
	)
	if got != 2 {
		t.Fatalf("expected 2 new pull requests, got %d", got)
	}
}

func TestCountNewContributors(t *testing.T) {
	got := countNewContributors(
		[]github.Contributor{{Login: "alice"}, {Login: "bob"}},
		[]github.Contributor{{Login: "Alice"}, {Login: "carol"}, {Login: "dave"}},
	)
	if got != 2 {
		t.Fatalf("expected 2 new contributors, got %d", got)
	}
}

func TestBuildRefreshSummaryFromCacheEntry(t *testing.T) {
	lastSync := time.Date(2026, 5, 31, 10, 0, 0, 0, time.UTC)
	baseline := AnalysisResult{
		Issues:       []github.Issue{{Number: 1}, {Number: 2}},
		PRs:          []github.PullRequest{{Number: 10}},
		Contributors: []github.Contributor{{Login: "alice"}},
	}
	analysis, err := json.Marshal(baseline)
	if err != nil {
		t.Fatalf("failed to marshal baseline: %v", err)
	}

	// Simulate fetched current data (no network)
	fetchedIssues := []github.Issue{{Number: 1}, {Number: 2}, {Number: 3}}
	fetchedPRs := []github.PullRequest{{Number: 10}, {Number: 11}}
	fetchedContribs := []github.Contributor{{Login: "alice"}, {Login: "bob"}}

	entry := &cache.CacheEntry{CachedAt: lastSync, Analysis: analysis}

	summary := RefreshSummary{
		LastSync:        entry.CachedAt,
		NewIssues:       countNewIssues(baseline.Issues, fetchedIssues),
		NewPullRequests: countNewPullRequests(baseline.PRs, fetchedPRs),
		NewContributors: countNewContributors(baseline.Contributors, fetchedContribs),
	}

	if !summary.LastSync.Equal(lastSync) {
		t.Fatalf("expected last sync %v, got %v", lastSync, summary.LastSync)
	}
	if summary.NewIssues != 1 {
		t.Fatalf("expected 1 new issue, got %d", summary.NewIssues)
	}
	if summary.NewPullRequests != 1 {
		t.Fatalf("expected 1 new pull request, got %d", summary.NewPullRequests)
	}
	if summary.NewContributors != 1 {
		t.Fatalf("expected 1 new contributor, got %d", summary.NewContributors)
	}
}

func TestDashboardOverviewShowsRefreshSummary(t *testing.T) {
	model := NewDashboardModel()
	model.SetData(AnalysisResult{Repo: &github.Repo{FullName: "owner/repo"}})
	model.SetRefreshSummary(&RefreshSummary{
		LastSync:        time.Date(2026, 5, 31, 10, 0, 0, 0, time.UTC),
		NewIssues:       12,
		NewPullRequests: 3,
		NewContributors: 1,
	})

	view := model.overviewView()
	for _, want := range []string{"Last Sync:", "Fetched:", "12 new issues", "3 new pull requests", "1 new contributor"} {
		if !strings.Contains(view, want) {
			t.Fatalf("expected overview to contain %q, got:\n%s", want, view)
		}
	}
}

func TestRefreshSummaryMarksDashboardSynced(t *testing.T) {
	model := NewMainModel(nil, nil)
	model.state = stateDashboard
	model.dashboard.SetData(AnalysisResult{Repo: &github.Repo{FullName: "owner/repo"}})

	next, _ := model.Update(refreshSummaryMsg{Summary: RefreshSummary{
		LastSync:        time.Now(),
		NewIssues:       1,
		NewPullRequests: 2,
		NewContributors: 3,
	}})
	got := next.(MainModel)

	if got.dashboard.cacheStatus != "synced" {
		t.Fatalf("expected dashboard status synced, got %s", got.dashboard.cacheStatus)
	}
	if got.dashboard.refreshSummary == nil {
		t.Fatal("expected refresh summary to be set")
	}
}

func TestSyncAnalysisUpdatesDashboardDataEverywhere(t *testing.T) {
	model := NewMainModel(nil, nil)
	model.state = stateDashboard
	model.dashboard.SetData(AnalysisResult{
		Repo:   &github.Repo{FullName: "owner/repo", OpenIssues: 1},
		Issues: []github.Issue{{Number: 1}},
	})

	next, _ := model.Update(syncAnalysisMsg{
		Result: AnalysisResult{
			Repo:   &github.Repo{FullName: "owner/repo", OpenIssues: 2},
			Issues: []github.Issue{{Number: 1}, {Number: 2}},
		},
		Summary: RefreshSummary{
			LastSync:  time.Now(),
			NewIssues: 1,
		},
	})
	got := next.(MainModel)

	if got.dashboard.data.Repo.OpenIssues != 2 {
		t.Fatalf("expected synced repo details to show 2 open issues, got %d", got.dashboard.data.Repo.OpenIssues)
	}
	if len(got.dashboard.data.Issues) != 2 {
		t.Fatalf("expected synced issue list to contain 2 issues, got %d", len(got.dashboard.data.Issues))
	}
	if got.dashboard.cacheStatus != "synced" {
		t.Fatalf("expected dashboard status synced, got %s", got.dashboard.cacheStatus)
	}
}

func TestDashboardSyncErrorReplacesSyncingStatus(t *testing.T) {
	model := NewMainModel(nil, nil)
	model.state = stateDashboard
	model.dashboard.SetData(AnalysisResult{Repo: &github.Repo{FullName: "owner/repo"}})
	model.dashboard.statusMsg = "Syncing repository data..."

	next, _ := model.Update(fmt.Errorf("cache not available"))
	got := next.(MainModel)

	if !strings.Contains(got.dashboard.statusMsg, "Sync failed: cache not available") {
		t.Fatalf("expected sync error status, got %q", got.dashboard.statusMsg)
	}
}
