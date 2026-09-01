package opportunity

import (
    "context"
    "time"
    "github.com/google/go-github/v55/github"
)

type GitHubService struct {
    client *github.Client
}

func NewGitHubService() *GitHubService {
    return &GitHubService{client: github.NewClient(nil)}
}

func (s *GitHubService) ListIssues(ctx context.Context, owner, repo string, state string) ([]*github.Issue, error) {
    var all []*github.Issue
    opt := &github.IssueListByRepoOptions{State: state, ListOptions: github.ListOptions{PerPage: 100}}
    for {
        issues, resp, err := s.client.Issues.ListByRepo(ctx, owner, repo, opt)
        if err != nil {
            return nil, err
        }
        all = append(all, issues...)
        if resp.NextPage == 0 {
            break
        }
        opt.Page = resp.NextPage
    }
    return all, nil
}

func (s *GitHubService) ListPullRequests(ctx context.Context, owner, repo string) ([]*github.PullRequest, error) {
    var all []*github.PullRequest
    opt := &github.PullRequestListOptions{State: "open", ListOptions: github.ListOptions{PerPage: 100}}
    for {
        prs, resp, err := s.client.PullRequests.List(ctx, owner, repo, opt)
        if err != nil {
            return nil, err
        }
        all = append(all, prs...)
        if resp.NextPage == 0 {
            break
        }
        opt.Page = resp.NextPage
    }
    return all, nil
}

func (s *GitHubService) ListContributors(ctx context.Context, owner, repo string) ([]*github.ContributorStats, error) {
    // GitHub provides contributor stats endpoint; may be delayed, handle 202 response.
    for {
        stats, resp, err := s.client.Repositories.ListContributorsStats(ctx, owner, repo)
        if err != nil {
            return nil, err
        }
        if resp.StatusCode == 202 {
            // GitHub is computing statistics, wait briefly.
            time.Sleep(2 * time.Second)
            continue
        }
        return stats, nil
    }
}
