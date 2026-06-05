package temporal

import (
	"fmt"
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestBuildTimelineFromGitHub_EmptyRepository(t *testing.T) {
	withStubbedTimelineData(t, func(base time.Time) {
		repo := &github.Repo{FullName: "owner/repo"}
		stubGitHubFetches(t, repo, nil, nil, nil)

		timeline, err := BuildTimelineFromGitHub(github.NewClient(), "owner", "repo", 3)
		if err != nil {
			t.Fatalf("BuildTimelineFromGitHub failed: %v", err)
		}
		if timeline == nil {
			t.Fatal("expected timeline, got nil")
		}
		if len(timeline.Snapshots) != 3 {
			t.Fatalf("expected 3 snapshots, got %d", len(timeline.Snapshots))
		}
		for i, snapshot := range timeline.Snapshots {
			if snapshot.Metrics.CommitCount != 0 || snapshot.Metrics.ContributorCount != 0 || snapshot.Metrics.ActiveContributors != 0 || snapshot.Metrics.IssuesOpen != 0 || snapshot.Metrics.PullRequestsOpen != 0 || snapshot.Metrics.AverageCommitFrequency != 0 || snapshot.Metrics.AverageHealth != 0 {
				t.Fatalf("snapshot %d expected zero metrics, got %+v", i, snapshot.Metrics)
			}
		}
	})
}

func TestBuildTimelineFromGitHub_SingleMonthActivity(t *testing.T) {
	withStubbedTimelineData(t, func(base time.Time) {
		commits := []github.Commit{
			makeCommit(base.AddDate(0, 0, -1), "alice", "c1"),
		}
		issues := []github.Issue{
			makeIssue(base.AddDate(0, 0, -2), "open", nil, false, 1),
		}
		prs := []github.PullRequest{
			makePullRequest(base.AddDate(0, 0, -3), "open", nil, nil, 1),
		}
		repo := &github.Repo{FullName: "owner/repo", Description: "demo", Stars: 10, OpenIssues: 1, PushedAt: base.AddDate(0, 0, -1)}
		stubGitHubFetches(t, repo, commits, issues, prs)

		timeline, err := BuildTimelineFromGitHub(github.NewClient(), "owner", "repo", 1)
		if err != nil {
			t.Fatalf("BuildTimelineFromGitHub failed: %v", err)
		}
		if len(timeline.Snapshots) != 1 {
			t.Fatalf("expected 1 snapshot, got %d", len(timeline.Snapshots))
		}
		snapshot := timeline.Snapshots[0]
		if snapshot.Metrics.CommitCount != 1 {
			t.Fatalf("expected 1 commit, got %d", snapshot.Metrics.CommitCount)
		}
		if snapshot.Metrics.ContributorCount != 1 {
			t.Fatalf("expected 1 contributor, got %d", snapshot.Metrics.ContributorCount)
		}
		if snapshot.Metrics.ActiveContributors != 1 {
			t.Fatalf("expected 1 active contributor, got %d", snapshot.Metrics.ActiveContributors)
		}
		if snapshot.Metrics.IssuesOpen != 1 {
			t.Fatalf("expected 1 open issue, got %d", snapshot.Metrics.IssuesOpen)
		}
		if snapshot.Metrics.PullRequestsOpen != 1 {
			t.Fatalf("expected 1 open PR, got %d", snapshot.Metrics.PullRequestsOpen)
		}
		if snapshot.Metrics.AverageHealth <= 0 {
			t.Fatalf("expected positive average health, got %d", snapshot.Metrics.AverageHealth)
		}
	})
}

func TestBuildTimelineFromGitHub_MultipleMonthActivity(t *testing.T) {
	withStubbedTimelineData(t, func(base time.Time) {
		commits := []github.Commit{
			makeCommit(base.AddDate(0, -2, -2), "alice", "c1"),
			makeCommit(base.AddDate(0, -1, -2), "bob", "c2"),
			makeCommit(base.AddDate(0, 0, -3), "alice", "c3"),
			makeCommit(base.AddDate(0, 0, -1), "carol", "c4"),
		}
		issues := []github.Issue{}
		prs := []github.PullRequest{}
		repo := &github.Repo{FullName: "owner/repo", Description: "demo", Stars: 20, OpenIssues: 0, PushedAt: base.AddDate(0, 0, -1)}
		stubGitHubFetches(t, repo, commits, issues, prs)

		timeline, err := BuildTimelineFromGitHub(github.NewClient(), "owner", "repo", 3)
		if err != nil {
			t.Fatalf("BuildTimelineFromGitHub failed: %v", err)
		}
		if len(timeline.Snapshots) != 3 {
			t.Fatalf("expected 3 snapshots, got %d", len(timeline.Snapshots))
		}
		if timeline.Snapshots[0].Metrics.CommitCount > timeline.Snapshots[1].Metrics.CommitCount || timeline.Snapshots[1].Metrics.CommitCount > timeline.Snapshots[2].Metrics.CommitCount {
			t.Fatalf("expected cumulative commit counts to be non-decreasing")
		}
		if timeline.Snapshots[2].Metrics.ContributorCount != 3 {
			t.Fatalf("expected 3 unique contributors by final month, got %d", timeline.Snapshots[2].Metrics.ContributorCount)
		}
	})
}

func TestBuildTimelineFromGitHub_ContributorAggregation(t *testing.T) {
	withStubbedTimelineData(t, func(base time.Time) {
		commits := []github.Commit{
			makeCommit(base.AddDate(0, -1, -10), "alice", "c1"),
			makeCommit(base.AddDate(0, 0, -12), "bob", "c2"),
			makeCommit(base.AddDate(0, 0, -4), "alice", "c3"),
		}
		stubGitHubFetches(t, &github.Repo{FullName: "owner/repo"}, commits, nil, nil)

		timeline, err := BuildTimelineFromGitHub(github.NewClient(), "owner", "repo", 2)
		if err != nil {
			t.Fatalf("BuildTimelineFromGitHub failed: %v", err)
		}
		if timeline.Snapshots[0].Metrics.ContributorCount != 1 {
			t.Fatalf("expected 1 contributor in first month, got %d", timeline.Snapshots[0].Metrics.ContributorCount)
		}
		if timeline.Snapshots[1].Metrics.ContributorCount != 2 {
			t.Fatalf("expected 2 contributors by second month, got %d", timeline.Snapshots[1].Metrics.ContributorCount)
		}
	})
}

func TestBuildTimelineFromGitHub_IgnoresAnonymousContributorIds(t *testing.T) {
	withStubbedTimelineData(t, func(base time.Time) {
		commits := []github.Commit{
			makeCommit(base.AddDate(0, -1, -10), "alice", "c1"),
			makeCommit(base.AddDate(0, 0, -12), "", "c2"),
			makeCommit(base.AddDate(0, 0, -4), "", "c3"),
		}
		stubGitHubFetches(t, &github.Repo{FullName: "owner/repo"}, commits, nil, nil)

		timeline, err := BuildTimelineFromGitHub(github.NewClient(), "owner", "repo", 2)
		if err != nil {
			t.Fatalf("BuildTimelineFromGitHub failed: %v", err)
		}
		if timeline.Snapshots[0].Metrics.ContributorCount != 1 {
			t.Fatalf("expected 1 contributor in first month, got %d", timeline.Snapshots[0].Metrics.ContributorCount)
		}
		if timeline.Snapshots[1].Metrics.ContributorCount != 1 {
			t.Fatalf("expected anonymous commits to be ignored, got %d contributors", timeline.Snapshots[1].Metrics.ContributorCount)
		}
	})
}

func TestBuildTimelineFromGitHub_OpenIssueCounting(t *testing.T) {
	withStubbedTimelineData(t, func(base time.Time) {
		issues := []github.Issue{
			makeIssue(time.Date(2024, 5, 10, 12, 0, 0, 0, time.UTC), "open", nil, false, 1),
			makeIssue(time.Date(2024, 6, 10, 12, 0, 0, 0, time.UTC), "open", nil, false, 2),
		}
		stubGitHubFetches(t, &github.Repo{FullName: "owner/repo"}, nil, issues, nil)

		timeline, err := BuildTimelineFromGitHub(github.NewClient(), "owner", "repo", 2)
		if err != nil {
			t.Fatalf("BuildTimelineFromGitHub failed: %v", err)
		}
		if timeline.Snapshots[0].Metrics.IssuesOpen != 1 {
			t.Fatalf("expected 1 open issue in first month, got %d", timeline.Snapshots[0].Metrics.IssuesOpen)
		}
		if timeline.Snapshots[1].Metrics.IssuesOpen != 2 {
			t.Fatalf("expected 2 open issues in second month, got %d", timeline.Snapshots[1].Metrics.IssuesOpen)
		}
	})
}

func TestBuildTimelineFromGitHub_OpenPRCounting(t *testing.T) {
	withStubbedTimelineData(t, func(base time.Time) {
		prs := []github.PullRequest{
			makePullRequest(time.Date(2024, 5, 10, 12, 0, 0, 0, time.UTC), "open", nil, nil, 1),
			makePullRequest(time.Date(2024, 6, 10, 12, 0, 0, 0, time.UTC), "open", nil, nil, 2),
		}
		stubGitHubFetches(t, &github.Repo{FullName: "owner/repo"}, nil, nil, prs)

		timeline, err := BuildTimelineFromGitHub(github.NewClient(), "owner", "repo", 2)
		if err != nil {
			t.Fatalf("BuildTimelineFromGitHub failed: %v", err)
		}
		if timeline.Snapshots[0].Metrics.PullRequestsOpen != 1 {
			t.Fatalf("expected 1 open PR in first month, got %d", timeline.Snapshots[0].Metrics.PullRequestsOpen)
		}
		if timeline.Snapshots[1].Metrics.PullRequestsOpen != 2 {
			t.Fatalf("expected 2 open PRs in second month, got %d", timeline.Snapshots[1].Metrics.PullRequestsOpen)
		}
	})
}

func TestBuildTimelineFromGitHub_SnapshotOrdering(t *testing.T) {
	withStubbedTimelineData(t, func(base time.Time) {
		commits := []github.Commit{
			makeCommit(base.AddDate(0, -2, -1), "alice", "c1"),
			makeCommit(base.AddDate(0, -1, -1), "alice", "c2"),
			makeCommit(base.AddDate(0, 0, -1), "alice", "c3"),
		}
		stubGitHubFetches(t, &github.Repo{FullName: "owner/repo"}, commits, nil, nil)

		timeline, err := BuildTimelineFromGitHub(github.NewClient(), "owner", "repo", 3)
		if err != nil {
			t.Fatalf("BuildTimelineFromGitHub failed: %v", err)
		}
		for i := 1; i < len(timeline.Snapshots); i++ {
			if !timeline.Snapshots[i-1].Timestamp.Before(timeline.Snapshots[i].Timestamp) {
				t.Fatalf("snapshots not ordered chronologically: %v then %v", timeline.Snapshots[i-1].Timestamp, timeline.Snapshots[i].Timestamp)
			}
		}
	})
}

func TestBuildTimelineFromGitHub_ExclusiveMonthEndBoundaries(t *testing.T) {
	withStubbedTimelineData(t, func(base time.Time) {
		monthEnd := time.Date(2024, 7, 1, 0, 0, 0, 0, time.UTC)
		commits := []github.Commit{
			makeCommit(time.Date(2024, 6, 30, 23, 59, 0, 0, time.UTC), "alice", "c0"),
			makeCommit(monthEnd, "alice", "c1"),
		}
		issues := []github.Issue{
			makeIssue(monthEnd, "open", nil, false, 1),
		}
		prs := []github.PullRequest{
			makePullRequest(monthEnd, "open", nil, nil, 1),
		}
		stubGitHubFetches(t, &github.Repo{FullName: "owner/repo"}, commits, issues, prs)

		timeline, err := BuildTimelineFromGitHub(github.NewClient(), "owner", "repo", 1)
		if err != nil {
			t.Fatalf("BuildTimelineFromGitHub failed: %v", err)
		}
		if len(timeline.Snapshots) != 1 {
			t.Fatalf("expected 1 snapshot, got %d", len(timeline.Snapshots))
		}

		snapshot := timeline.Snapshots[0]
		if snapshot.Metrics.CommitCount != 1 {
			t.Fatalf("expected only the pre-boundary commit to count, got %d", snapshot.Metrics.CommitCount)
		}
		if snapshot.Metrics.IssuesOpen != 0 {
			t.Fatalf("expected 0 open issues at month end boundary, got %d", snapshot.Metrics.IssuesOpen)
		}
		if snapshot.Metrics.PullRequestsOpen != 0 {
			t.Fatalf("expected 0 open PRs at month end boundary, got %d", snapshot.Metrics.PullRequestsOpen)
		}
	})
}

func withStubbedTimelineData(t *testing.T, fn func(base time.Time)) {
	t.Helper()
	originalNow := timelineNow
	timelineNow = func() time.Time {
		return time.Date(2024, 6, 15, 12, 0, 0, 0, time.UTC)
	}
	t.Cleanup(func() {
		timelineNow = originalNow
	})

	fn(timelineNow())
}

func stubGitHubFetches(t *testing.T, repo *github.Repo, commits []github.Commit, issues []github.Issue, prs []github.PullRequest) {
	t.Helper()
	originalRepo := fetchRepo
	originalCommits := fetchCommits
	originalIssues := fetchIssues
	originalPRs := fetchPullRequests

	fetchRepo = func(_ *github.Client, owner, repoName string) (*github.Repo, error) {
		if repo == nil {
			return &github.Repo{FullName: fmt.Sprintf("%s/%s", owner, repoName)}, nil
		}
		return repo, nil
	}
	fetchCommits = func(_ *github.Client, owner, repoName string, days int) ([]github.Commit, error) {
		return commits, nil
	}
	fetchIssues = func(_ *github.Client, owner, repoName, state string) ([]github.Issue, error) {
		return issues, nil
	}
	fetchPullRequests = func(_ *github.Client, owner, repoName, state string) ([]github.PullRequest, error) {
		return prs, nil
	}

	t.Cleanup(func() {
		fetchRepo = originalRepo
		fetchCommits = originalCommits
		fetchIssues = originalIssues
		fetchPullRequests = originalPRs
	})
}

func makeCommit(date time.Time, login, sha string) github.Commit {
	return github.Commit{
		SHA: sha,
		Commit: struct {
			Author struct {
				Date time.Time `json:"date"`
			} `json:"author"`
		}{
			Author: struct {
				Date time.Time `json:"date"`
			}{Date: date.UTC()},
		},
		Author: &struct {
			Login string `json:"login"`
		}{Login: login},
	}
}

func makeIssue(date time.Time, state string, closedAt *time.Time, isPR bool, number int) github.Issue {
	var pullRequest *struct{}
	if isPR {
		pullRequest = &struct{}{}
	}
	return github.Issue{
		Number:    number,
		State:     state,
		CreatedAt: date.UTC(),
		ClosedAt:  closedAt,
		PullRequest: pullRequest,
	}
}

func makePullRequest(date time.Time, state string, mergedAt *time.Time, closedAt *time.Time, number int) github.PullRequest {
	return github.PullRequest{
		Number:    number,
		State:     state,
		CreatedAt: date.UTC(),
		MergedAt:  mergedAt,
		ClosedAt:  closedAt,
	}
}
