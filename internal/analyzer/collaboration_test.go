package analyzer

import (
	"testing"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestAnalyzeCrossContributorPatterns(t *testing.T) {
	tests := []struct {
		name     string
		prs      []github.PullRequest
		wantKeys []string
	}{
		{
			name: "single digit PR number",
			prs: []github.PullRequest{
				{Number: 5, User: github.User{Login: "user1"}},
			},
			wantKeys: []string{"PR #5"},
		},
		{
			name: "multi digit PR number",
			prs: []github.PullRequest{
				{Number: 12, User: github.User{Login: "user1"}},
				{Number: 105, User: github.User{Login: "user2"}},
			},
			wantKeys: []string{"PR #12", "PR #105"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			patterns := analyzeCrossContributorPatterns(tt.prs)
			
			if len(patterns.TopCollaboratedFiles) != len(tt.wantKeys) {
				t.Errorf("analyzeCrossContributorPatterns() returned %d files, want %d", len(patterns.TopCollaboratedFiles), len(tt.wantKeys))
			}

			for _, wantKey := range tt.wantKeys {
				found := false
				for _, file := range patterns.TopCollaboratedFiles {
					if file.Filename == wantKey {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("analyzeCrossContributorPatterns() = %v, want key %s", patterns.TopCollaboratedFiles, wantKey)
				}
			}
		})
	}
}
