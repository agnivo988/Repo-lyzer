package analyzer

import (
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestCalculateHealth(t *testing.T) {
	tests := []struct {
		name     string
		repo     *github.Repo
		commits  []github.Commit
		minScore int
		maxScore int
	}{
		{
			name: "healthy repo with recent activity",
			repo: &github.Repo{
				Stars:       100,
				Forks:       20,
				OpenIssues:  5,
				Description: "A great project",
				PushedAt:    time.Now().Add(-24 * time.Hour),
			},
			commits:  makeCommits(50),
			minScore: 50,
			maxScore: 100,
		},
		{
			name: "inactive repo",
			repo: &github.Repo{
				Stars:       10,
				Forks:       2,
				OpenIssues:  0,
				Description: "",
				PushedAt:    time.Now().Add(-365 * 24 * time.Hour),
			},
			commits:  makeCommits(5),
			minScore: 50,
			maxScore: 70,
		},
		{
			name: "popular but stale repo",
			repo: &github.Repo{
				Stars:       1000,
				Forks:       200,
				OpenIssues:  50,
				Description: "Popular project",
				PushedAt:    time.Now().Add(-180 * 24 * time.Hour),
			},
			commits:  makeCommits(10),
			minScore: 30,
			maxScore: 80,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			score := CalculateHealth(tt.repo, tt.commits)
			if score < tt.minScore || score > tt.maxScore {
				t.Errorf("CalculateHealth() = %d, want between %d and %d", score, tt.minScore, tt.maxScore)
			}
		})
	}
}

func TestCalculateHealth_ScoreBounds(t *testing.T) {
	// Test that score is always between 0 and 100
	repo := &github.Repo{
		Stars:      999999,
		Forks:      999999,
		OpenIssues: 0,
		PushedAt:   time.Now(),
	}
	commits := makeCommits(1000)

	score := CalculateHealth(repo, commits)
	if score < 0 || score > 100 {
		t.Errorf("Score %d is out of bounds [0, 100]", score)
	}
}

// Helper function to create test commits
func makeCommits(count int) []github.Commit {
	commits := make([]github.Commit, count)
	for i := 0; i < count; i++ {
		commits[i] = github.Commit{
			SHA: "abc123",
			Commit: struct {
				Author struct {
					Date time.Time `json:"date"`
				} `json:"author"`
			}{
				Author: struct {
					Date time.Time `json:"date"`
				}{
					Date: time.Now().Add(-time.Duration(i) * 24 * time.Hour),
				},
			},
		}
	}
	return commits
}
