package analyzer

import (
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// helpers for building test fixtures

func repoWith(pushedDaysAgo int, openIssues int, description string) *github.Repo {
	pushed := time.Now().AddDate(0, 0, -pushedDaysAgo)
	return &github.Repo{
		PushedAt:    pushed,
		OpenIssues:  openIssues,
		Description: description,
	}
}

func makeContributors(n int) []github.Contributor {
	c := make([]github.Contributor, n)
	for i := range c {
		c[i] = github.Contributor{Login: "user", Commits: 10}
	}
	return c
}

func makeRiskCommits(n int) []github.Commit {
	c := make([]github.Commit, n)
	for i := range c {
		c[i] = github.Commit{}
	}
	return c
}

// readmeTree returns a minimal file tree that contains a README.
func readmeTree() []github.TreeEntry {
	return []github.TreeEntry{{Path: "README.md", Type: "blob"}}
}

func emptyTree() []github.TreeEntry { return nil }

// ─── riskLevel ────────────────────────────────────────────────────────────────

func TestRiskLevel(t *testing.T) {
	cases := []struct {
		score int
		want  string
	}{
		{0, "Low"},
		{30, "Low"},
		{31, "Moderate"},
		{60, "Moderate"},
		{61, "High"},
		{100, "High"},
	}
	for _, tc := range cases {
		got := riskLevel(tc.score)
		if got != tc.want {
			t.Errorf("riskLevel(%d) = %q, want %q", tc.score, got, tc.want)
		}
	}
}

// ─── AnalyzeRisk — clean repo ─────────────────────────────────────────────────

func TestAnalyzeRisk_CleanRepo(t *testing.T) {
	repo := repoWith(10, 5, "a project") // pushed 10 days ago, few issues
	contributors := makeContributors(5)
	commits := makeRiskCommits(20)
	tree := readmeTree()

	report := AnalyzeRisk(repo, contributors, commits, tree)

	if report.Score != 0 {
		t.Errorf("expected score 0 for clean repo, got %d; signals: %v", report.Score, report.Signals)
	}
	if report.Level != "Low" {
		t.Errorf("expected Low risk, got %q", report.Level)
	}
	if len(report.Warnings) != 0 {
		t.Errorf("expected no warnings, got %v", report.Warnings)
	}
}

// ─── AnalyzeRisk — individual rules ──────────────────────────────────────────

func TestAnalyzeRisk_LowBusFactor(t *testing.T) {
	// Single contributor dominates → bus factor ≤ 1 → +25
	contributors := []github.Contributor{
		{Login: "solo", Commits: 999},
	}
	repo := repoWith(10, 5, "desc")
	report := AnalyzeRisk(repo, contributors, makeRiskCommits(20), readmeTree())

	if report.Score < 25 {
		t.Errorf("expected ≥25 for low bus factor, got %d", report.Score)
	}
	if !containsSignal(report, "Low bus factor") {
		t.Errorf("expected 'Low bus factor' signal, got %v", report.Signals)
	}
}

func TestAnalyzeRisk_LowContributorCount(t *testing.T) {
	// Only 2 contributors → +15
	contributors := makeContributors(2)
	repo := repoWith(10, 5, "desc")
	report := AnalyzeRisk(repo, contributors, makeRiskCommits(20), readmeTree())

	if !containsSignal(report, "Low contributor count") {
		t.Errorf("expected 'Low contributor count' signal, got %v", report.Signals)
	}
}

func TestAnalyzeRisk_NoRecentCommits(t *testing.T) {
	repo := repoWith(10, 5, "desc")
	report := AnalyzeRisk(repo, makeContributors(5), makeRiskCommits(0), readmeTree())

	if !containsSignal(report, "No recent commits (90d)") {
		t.Errorf("expected 'No recent commits' signal, got %v", report.Signals)
	}
}

func TestAnalyzeRisk_HighIssueBacklog(t *testing.T) {
	repo := repoWith(10, 150, "desc") // 150 open issues
	report := AnalyzeRisk(repo, makeContributors(5), makeRiskCommits(20), readmeTree())

	if !containsSignal(report, "High open issue backlog") {
		t.Errorf("expected 'High open issue backlog' signal, got %v", report.Signals)
	}
}

func TestAnalyzeRisk_MissingREADME(t *testing.T) {
	repo := repoWith(10, 5, "desc")
	report := AnalyzeRisk(repo, makeContributors(5), makeRiskCommits(20), emptyTree())

	if !containsSignal(report, "Missing README") {
		t.Errorf("expected 'Missing README' signal, got %v", report.Signals)
	}
}

func TestAnalyzeRisk_StaleRepo(t *testing.T) {
	repo := repoWith(200, 5, "desc") // pushed 200 days ago → stale
	report := AnalyzeRisk(repo, makeContributors(5), makeRiskCommits(20), readmeTree())

	if !containsSignal(report, "Repository stale (>180d)") {
		t.Errorf("expected 'Repository stale' signal, got %v", report.Signals)
	}
}

// ─── AnalyzeRisk — score cap ──────────────────────────────────────────────────

func TestAnalyzeRisk_ScoreCapAt100(t *testing.T) {
	// All rules fire: single contributor, stale, no commits, >100 issues, no README
	repo := repoWith(200, 150, "")
	contributors := []github.Contributor{{Login: "solo", Commits: 999}}

	report := AnalyzeRisk(repo, contributors, makeRiskCommits(0), emptyTree())

	if report.Score > 100 {
		t.Errorf("score should be capped at 100, got %d", report.Score)
	}
	if report.Score < 61 {
		t.Errorf("score should be High risk when all rules fire, got %d (level: %s)", report.Score, report.Level)
	}
	if report.Level != "High" {
		t.Errorf("expected High risk level, got %q", report.Level)
	}
}

// ─── helpers ──────────────────────────────────────────────────────────────────

func containsSignal(r RiskReport, label string) bool {
	for _, s := range r.Signals {
		if s.Label == label {
			return true
		}
	}
	return false
}
