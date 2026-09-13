package analyzer

import (
    "testing"
    "time"

    "github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestCalculateHealthDetailed_Sanity(t *testing.T) {
    repo := &github.Repo{
        Stars:      120,
        OpenIssues: 12,
        PushedAt:   time.Now().AddDate(0, -1, 0),
    }

    commits := make([]github.Commit, 60)
    contributors := []github.Contributor{
        {Login: "alice", Commits: 40},
        {Login: "bob", Commits: 20},
    }
    // minimal merged PR (use non-nil time)
    t1 := time.Now()
    prs := []github.PullRequest{{Number: 1, MergedAt: &t1}}
    issues := []github.Issue{{State: "open"}}

    score := CalculateHealthDetailed(repo, commits, contributors, prs, issues)
    if score < 0 || score > 100 {
        t.Fatalf("unexpected score %d, want 0..100", score)
    }
}
