package analyzer

import (
	"testing"
	"time" // This fixes the "undefined: time" error

	"github.com/agnivo988/Repo-lyzer/internal/github" // This fixes the "undefined: github" error
)

func TestCalculatePRVelocity(t *testing.T) {
	now := time.Now()
	// Create a PR that took exactly 24 hours to merge
	createdAt := now.Add(-24 * time.Hour)
	prs := []github.PullRequest{
		{
			CreatedAt: createdAt,
			MergedAt:  &now,
		},
	}

	velocity := CalculatePRVelocity(prs)
	// Replace your old if block with this:
	expected := 24.0
	tolerance := 0.0001
	if velocity < expected-tolerance || velocity > expected+tolerance {
		t.Errorf("Expected ~24.0 hours, got %f", velocity)
	}
}
