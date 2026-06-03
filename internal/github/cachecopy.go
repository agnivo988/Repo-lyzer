package github

func copyCommits(src []Commit) []Commit {
	out := make([]Commit, len(src))
	copy(out, src)
	return out
}

func copyContributors(src []Contributor) []Contributor {
	out := make([]Contributor, len(src))
	copy(out, src)
	return out
}

func copyIssues(src []Issue) []Issue {
	out := make([]Issue, len(src))
	copy(out, src)
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
	copy(out, src)
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
