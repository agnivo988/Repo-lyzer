package github

import "time"

func copyCommits(src []Commit) []Commit {
	out := make([]Commit, len(src))
	for i, c := range src {
		out[i] = c
		if c.Author != nil {
			authorCopy := *c.Author
			out[i].Author = &authorCopy
		}
	}
	return out
}

func copyContributors(src []Contributor) []Contributor {
	out := make([]Contributor, len(src))
	copy(out, src)
	return out
}

func copyIssues(src []Issue) []Issue {
	out := make([]Issue, len(src))
	for i, issue := range src {
		out[i] = issue
		if issue.ClosedAt != nil {
			closedCopy := *issue.ClosedAt
			out[i].ClosedAt = &closedCopy
		}
		if issue.PullRequest != nil {
			prCopy := *issue.PullRequest
			out[i].PullRequest = &prCopy
		}
	}
	return out
}

func copyLanguagesMap(src map[string]int) map[string]int {
	out := make(map[string]int, len(src))
	for k, v := range src {
		out[k] = v
	}
	return out
}

func copyPullRequests(src []PullRequest) []PullRequest {
	out := make([]PullRequest, len(src))
	for i, pr := range src {
		out[i] = pr
		if pr.MergedAt != nil {
			mergedCopy := *pr.MergedAt
			out[i].MergedAt = &mergedCopy
		}
		if pr.ClosedAt != nil {
			closedCopy := *pr.ClosedAt
			out[i].ClosedAt = &closedCopy
		}
	}
	return out
}

func copyReviews(src []Review) []Review {
	out := make([]Review, len(src))
	copy(out, src)
	return out
}

func copyTreeEntries(src []TreeEntry) []TreeEntry {
	out := make([]TreeEntry, len(src))
	copy(out, src)
	return out
}

func copyTime(t *time.Time) *time.Time {
	if t == nil {
		return nil
	}
	copy := *t
	return &copy
}
