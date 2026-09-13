package opportunity

import (
    "fmt"
    "github.com/google/go-github/v55/github"
    "context"
    "time"
)

type Engine struct {
    client *github.Client
    // cache could be added later
}

func NewEngine() *Engine {
    // Use unauthenticated client; token support can be added later
    client := github.NewClient(nil)
    return &Engine{client: client}
}

// Report holds the analysis results (placeholder for now)
type Report struct {
    Owner string
    Repo  string
    Score int
}

func (e *Engine) Run(ctx context.Context, owner, repo string) (*Report, error) {
    // Use GitHubService to fetch issues and contributors
    svc := NewGitHubService()
    // Fetch open issues
    issues, err := svc.ListIssues(ctx, owner, repo, "open")
    if err != nil {
        return nil, fmt.Errorf("failed to list issues: %w", err)
    }
    // Fetch contributors
    contributors, err := svc.ListContributors(ctx, owner, repo)
    if err != nil {
        return nil, fmt.Errorf("failed to list contributors: %w", err)
    }
    // Simple scoring: 2 points per good-first-issue, 1 point per help-wanted, plus number of contributors
    goodFirst := 0
    helpWanted := 0
    for _, iss := range issues {
        cls := Classify(iss)
        if cls == ClassificationGoodFirstIssue {
            goodFirst++
        } else if cls == ClassificationHelpWanted {
            helpWanted++
        }
    }
    score := goodFirst*2 + helpWanted + len(contributors)
    fmt.Printf("Running opportunity analysis for %s/%s...\n", owner, repo)
    return &Report{Owner: owner, Repo: repo, Score: score}, nil
}
