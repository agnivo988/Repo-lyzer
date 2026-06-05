package ui

import (
	"strings"
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestRepositoryTrendsViewRendersTrendSections(t *testing.T) {
	model := NewDashboardModel()
	model.currentView = viewTrends

	repo := &github.Repo{
		FullName:    "owner/repo",
		Description: "Repository used for trend rendering tests",
		Stars:       120,
		OpenIssues:  7,
		PushedAt:    time.Now(),
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
	addCommit(time.Date(2026, 3, 15, 0, 0, 0, 0, time.UTC), "carol")
	addCommit(time.Date(2026, 4, 5, 0, 0, 0, 0, time.UTC), "alice")
	addCommit(time.Date(2026, 4, 10, 0, 0, 0, 0, time.UTC), "bob")
	addCommit(time.Date(2026, 4, 15, 0, 0, 0, 0, time.UTC), "carol")
	addCommit(time.Date(2026, 4, 20, 0, 0, 0, 0, time.UTC), "dave")

	model.SetData(AnalysisResult{
		Repo:    repo,
		Commits: commits,
	})

	view := model.View()
	for _, expected := range []string{
		"Repository Trends",
		"Health Score",
		"Contributor Growth",
		"Predicted Health Score",
	} {
		if !strings.Contains(view, expected) {
			t.Fatalf("trends view missing %q:\n%s", expected, view)
		}
	}
}
