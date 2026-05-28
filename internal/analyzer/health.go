package analyzer

import (
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func CalculateHealth(repo *github.Repo, commits []github.Commit) int {
	score := 50

	if repo.Description != "" {
		score += 10
	}
	if repo.Stars > 50 {
		score += 10
	}
	if len(commits) > 10 {
		score += 20
	}
	if repo.OpenIssues < 20 {
		score += 10
	}

	if !repo.PushedAt.IsZero() {
		since := time.Since(repo.PushedAt)
		switch {
		case since <= 30*24*time.Hour:
			score += 10
		case since <= 90*24*time.Hour:
			score += 5
		case since > 365*24*time.Hour:
			score -= 10
		case since > 180*24*time.Hour:
			score -= 5
		}
	}

	if score > 100 {
		score = 100
	}
	if score < 0 {
		score = 0
	}
	return score
}

// RepositoryHealthDetails contains a breakdown of the computed health score
type RepositoryHealthDetails struct {
	Score               int
	StarsScore          int
	CommitsScore        int
	IssueScore          int
	PRScore             int
	ContributorScore    int
}

// CalculateHealthDetailed computes a health score using repository metadata,
// commits, contributors, pull requests and issues. It returns a 0-100 score.
func CalculateHealthDetailed(repo *github.Repo, commits []github.Commit, contributors []github.Contributor, prs []github.PullRequest, issues []github.Issue) int {
	// We weight components to keep scores interpretable:
	// stars: 15, commits: 25, issues: 15, prs: 20, contributors/activity: 25
	var details RepositoryHealthDetails

	// Stars (15)
	if repo != nil {
		switch {
		case repo.Stars >= 1000:
			details.StarsScore = 15
		case repo.Stars >= 500:
			details.StarsScore = 12
		case repo.Stars >= 100:
			details.StarsScore = 8
		case repo.Stars >= 50:
			details.StarsScore = 5
		default:
			// proportional up to 50 stars
			details.StarsScore = int(float64(repo.Stars) / 50.0 * 5.0)
		}
	}

	// Commits (25)
	commitCount := len(commits)
	switch {
	case commitCount > 1000:
		details.CommitsScore = 25
	case commitCount > 500:
		details.CommitsScore = 20
	case commitCount > 100:
		details.CommitsScore = 12
	case commitCount > 25:
		details.CommitsScore = 6
	default:
		details.CommitsScore = int(float64(commitCount) / 25.0 * 6.0)
	}

	// Issues (15) - prefer low open issues and good resolution
	if repo != nil {
		open := repo.OpenIssues
		if open < 5 {
			details.IssueScore = 15
		} else if open < 20 {
			details.IssueScore = 10
		} else if open < 50 {
			details.IssueScore = 5
		} else {
			details.IssueScore = 2
		}
	}

	// PRs (20) - use merge rate as quality indicator
	if len(prs) > 0 {
		merged := 0
		for _, p := range prs {
			if p.MergedAt != nil {
				merged++
			}
		}
		mergeRate := float64(merged) / float64(len(prs))
		details.PRScore = int(mergeRate*20.0 + 0.5)
	} else {
		// no PRs -> neutral small score
		details.PRScore = 6
	}

	// Contributors & activity (25) - combine contributor diversity and recent activity
	contribScore := 0
	if len(contributors) > 0 {
		insights := AnalyzeContributors(contributors)
		// Use diversity score (0-100) mapped to 0-20
		contribScore = int(insights.DiversityScore/100.0*20.0 + 0.5)
	}

	// Activity portion (last 90 days) (0-5)
	activity := AnalyzeContributorActivity(commits)
	activityScore := 0
	if activity.Last90Days >= 50 {
		activityScore = 5
	} else if activity.Last90Days >= 20 {
		activityScore = 3
	} else if activity.Last90Days >= 5 {
		activityScore = 2
	}

	details.ContributorScore = contribScore + activityScore

	// Sum weighted components
	total := details.StarsScore + details.CommitsScore + details.IssueScore + details.PRScore + details.ContributorScore

	// Clamp and return
	if total > 100 {
		total = 100
	}
	if total < 0 {
		total = 0
	}

	details.Score = total
	return details.Score
}
