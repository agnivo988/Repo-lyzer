package temporal

import (
	"fmt"
	"sort"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

var timelineNow = time.Now

var fetchRepo = func(client *github.Client, owner, repo string) (*github.Repo, error) {
	return client.GetRepo(owner, repo)
}

var fetchCommits = func(client *github.Client, owner, repo string, days int) ([]github.Commit, error) {
	return client.GetCommits(owner, repo, days)
}

var fetchIssues = func(client *github.Client, owner, repo, state string) ([]github.Issue, error) {
	return client.GetIssues(owner, repo, state)
}

var fetchPullRequests = func(client *github.Client, owner, repo, state string) ([]github.PullRequest, error) {
	return client.GetPullRequests(owner, repo, state)
}

// BuildTimelineFromGitHub builds a monthly repository timeline from GitHub activity.
func BuildTimelineFromGitHub(client *github.Client, owner string, repo string, months int) (*Timeline, error) {
	if client == nil {
		return nil, fmt.Errorf("client is nil")
	}
	if owner == "" || repo == "" {
		return nil, fmt.Errorf("owner and repo are required")
	}
	if months <= 0 {
		return nil, fmt.Errorf("months must be greater than zero, got %d", months)
	}

	repoInfo, err := fetchRepo(client, owner, repo)
	if err != nil {
		return nil, fmt.Errorf("fetch repo: %w", err)
	}

	lookbackDays := months * 31
	commits, err := fetchCommits(client, owner, repo, lookbackDays)
	if err != nil {
		return nil, fmt.Errorf("fetch commits: %w", err)
	}

	issues, err := fetchIssues(client, owner, repo, "all")
	if err != nil {
		return nil, fmt.Errorf("fetch issues: %w", err)
	}

	prs, err := fetchPullRequests(client, owner, repo, "all")
	if err != nil {
		return nil, fmt.Errorf("fetch pull requests: %w", err)
	}

	now := startOfMonthUTC(timelineNow().UTC())
	snapshots := buildMonthlySnapshots(repoInfo, commits, issues, prs, months, now)

	timeline := NewTimeline(owner, repo)
	for _, snapshot := range snapshots {
		if err := timeline.AddSnapshot(snapshot); err != nil {
			return nil, err
		}
	}

	return timeline, nil
}

func buildMonthlySnapshots(repoInfo *github.Repo, commits []github.Commit, issues []github.Issue, prs []github.PullRequest, months int, now time.Time) []*Snapshot {
	snapshots := make([]*Snapshot, 0, months)
	if months <= 0 {
		return snapshots
	}

	lookbackStart := startOfMonthUTC(now).AddDate(0, -(months-1), 0)

	for index := 0; index < months; index++ {
		monthStart := lookbackStart.AddDate(0, index, 0)
		monthEnd := monthStart.AddDate(0, 1, 0)

		cumulativeCommits := commitsUpTo(commits, monthEnd)
		monthlyActiveContributors := activeContributorsInMonth(commits, monthStart, monthEnd)
		cumulativeContributors := contributorsUpTo(commits, monthEnd)
		openIssues := countOpenIssuesAtMonthEnd(issues, monthEnd)
		openPRs := countOpenPullRequestsAtMonthEnd(prs, monthEnd)

		metrics := NewRepositoryMetrics()
		metrics.CommitCount = len(cumulativeCommits)
		metrics.ContributorCount = len(cumulativeContributors)
		metrics.ActiveContributors = len(monthlyActiveContributors)
		metrics.IssuesOpen = openIssues
		metrics.PullRequestsOpen = openPRs
		metrics.AverageCommitFrequency = averageCommitFrequency(len(cumulativeCommits), lookbackStart, monthEnd)
		metrics.AverageHealth = calculateMonthlyHealth(repoInfo, cumulativeCommits, openIssues, monthEnd)

		snapshot := NewSnapshot(monthStart, nil)
		snapshot.Metrics = metrics
		snapshot.Contributors = sortedStrings(monthlyActiveContributors)
		snapshots = append(snapshots, snapshot)
	}

	return snapshots
}

func startOfMonthUTC(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), 1, 0, 0, 0, 0, time.UTC)
}

func monthKey(t time.Time) string {
	return startOfMonthUTC(t).Format("2006-01")
}

func commitsUpTo(commits []github.Commit, monthEnd time.Time) []github.Commit {
	result := make([]github.Commit, 0, len(commits))
	for _, commit := range commits {
		if commit.Commit.Author.Date.Before(monthEnd) {
			result = append(result, commit)
		}
	}
	return result
}

func contributorsUpTo(commits []github.Commit, monthEnd time.Time) map[string]struct{} {
	result := make(map[string]struct{})
	for _, commit := range commits {
		if commit.Commit.Author.Date.After(monthEnd) {
			continue
		}
		if id := commitContributorID(commit); id != "" {
			result[id] = struct{}{}
		}
	}
	return result
}

func activeContributorsInMonth(commits []github.Commit, monthStart, monthEnd time.Time) map[string]struct{} {
	result := make(map[string]struct{})
	for _, commit := range commits {
		commitDate := commit.Commit.Author.Date
		if commitDate.Before(monthStart) || !commitDate.Before(monthEnd) {
			continue
		}
		if id := commitContributorID(commit); id != "" {
			result[id] = struct{}{}
		}
	}
	return result
}

func countOpenIssuesAtMonthEnd(issues []github.Issue, monthEnd time.Time) int {
	count := 0
	for _, issue := range issues {
		if issue.PullRequest != nil {
			continue
		}
		if !issue.CreatedAt.Before(monthEnd) {
			continue
		}
		if issue.ClosedAt != nil && issue.ClosedAt.Before(monthEnd) {
			continue
		}
		if issue.State == "closed" && issue.ClosedAt == nil {
			continue
		}
		count++
	}
	return count
}

func countOpenPullRequestsAtMonthEnd(prs []github.PullRequest, monthEnd time.Time) int {
	count := 0
	for _, pr := range prs {
		if !pr.CreatedAt.Before(monthEnd) {
			continue
		}
		if pr.MergedAt != nil && pr.MergedAt.Before(monthEnd) {
			continue
		}
		if pr.ClosedAt != nil && pr.ClosedAt.Before(monthEnd) {
			continue
		}
		count++
	}
	return count
}

func averageCommitFrequency(commitCount int, start time.Time, monthEnd time.Time) float64 {
	if commitCount == 0 {
		return 0
	}
	days := monthEnd.Sub(start).Hours() / 24
	if days <= 0 {
		days = 1
	}
	return float64(commitCount) / days
}

func calculateMonthlyHealth(repoInfo *github.Repo, commits []github.Commit, openIssues int, monthEnd time.Time) int {
	if repoInfo == nil {
		return 0
	}
	if len(commits) == 0 && openIssues == 0 {
		return 0
	}

	repoCopy := *repoInfo
	repoCopy.OpenIssues = openIssues

	latestCommit := latestCommitTime(commits)
	if !latestCommit.IsZero() {
		repoCopy.PushedAt = latestCommit
	} else if repoCopy.PushedAt.After(monthEnd) {
		repoCopy.PushedAt = monthEnd
	}

	return analyzer.CalculateHealth(&repoCopy, commits)
}

func latestCommitTime(commits []github.Commit) time.Time {
	var latest time.Time
	for _, commit := range commits {
		if commit.Commit.Author.Date.After(latest) {
			latest = commit.Commit.Author.Date
		}
	}
	return latest
}

func commitContributorID(commit github.Commit) string {
	if commit.Author != nil && commit.Author.Login != "" {
		return commit.Author.Login
	}
	return ""
}

func sortedStrings(values map[string]struct{}) []string {
	if len(values) == 0 {
		return []string{}
	}
	result := make([]string, 0, len(values))
	for value := range values {
		result = append(result, value)
	}
	sort.Strings(result)
	return result
}
