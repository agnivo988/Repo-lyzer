package analyzer

import (
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

type ContributorActivityResult struct {
	Last90Days  int
	Last180Days int
	Trend       string
	Insight     string
}

func AnalyzeContributorActivity(commits []github.Commit) ContributorActivityResult {
	now := time.Now()

	last90 := 0
	last180 := 0

	for _, c := range commits {
		commitTime := c.Commit.Author.Date
		daysAgo := now.Sub(commitTime).Hours() / 24

		if daysAgo <= 90 {
			last90++
		}
		if daysAgo <= 180 {
			last180++
		}
	}

	result := ContributorActivityResult{
		Last90Days:  last90,
		Last180Days: last180,
	}

	// Determine trend
	switch {
	case last90 < last180:
		result.Trend = "Declining"
		result.Insight = "Commit activity has decreased in recent months"
	case last90 > last180:
		result.Trend = "Growing"
		result.Insight = "Commit activity has increased recently"
	default:
		result.Trend = "Stable"
		result.Insight = "Commit activity is stable"
	}

	// Risk insight
	if last90 == 0 {
		result.Insight = "No commit activity in the last 90 days"
	} else if last90 < 5 {
		result.Insight = "Very low recent development activity"
	}

	return result
}
