package analyzer

import (
	"fmt"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

type MaintainerAction struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Priority    string `json:"priority"` // "High", "Medium", "Low"
}

type IssueCandidate struct {
	Number   int    `json:"number"`
	Title    string `json:"title"`
	Comments int    `json:"comments"`
	DaysOpen int    `json:"days_open"`
	Reason   string `json:"reason"`
}

type PRStuck struct {
	Number   int    `json:"number"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	DaysOpen int    `json:"days_open"`
	Reason   string `json:"reason"`
}

type ContributorAppreciation struct {
	Login        string `json:"login"`
	Commits      int    `json:"commits"`
	Contribution string `json:"contribution"`
}

type MaintainerAnalysis struct {
	RepoFullName             string                    `json:"repo_full_name"`
	SuggestedActions         []MaintainerAction        `json:"suggested_actions"`
	IssueCandidates          []IssueCandidate          `json:"issue_candidates"`
	PRsStuck                 []PRStuck                 `json:"prs_stuck"`
	ContributorsToAppreciate []ContributorAppreciation `json:"contributors_to_appreciate"`
}

func AnalyzeMaintainerDashboard(
	repo *github.Repo,
	prs []github.PullRequest,
	issues []github.Issue,
	healthScore int,
	busFactor int,
	hasLockFile bool,
	contributors []github.Contributor,
) *MaintainerAnalysis {
	analysis := &MaintainerAnalysis{
		RepoFullName:             repo.FullName,
		SuggestedActions:         []MaintainerAction{},
		IssueCandidates:          []IssueCandidate{},
		PRsStuck:                 []PRStuck{},
		ContributorsToAppreciate: []ContributorAppreciation{},
	}

	// 1. Analyze PRs Stuck in Review
	for _, pr := range prs {
		if pr.State != "open" {
			continue
		}
		daysOpen := int(time.Since(pr.CreatedAt).Hours() / 24)
		if daysOpen > 7 {
			analysis.PRsStuck = append(analysis.PRsStuck, PRStuck{
				Number:   pr.Number,
				Title:    pr.Title,
				Author:   pr.User.Login,
				DaysOpen: daysOpen,
				Reason:   "Open for >7 days without merge",
			})
		} else if pr.Draft {
			analysis.PRsStuck = append(analysis.PRsStuck, PRStuck{
				Number:   pr.Number,
				Title:    pr.Title,
				Author:   pr.User.Login,
				DaysOpen: daysOpen,
				Reason:   "Draft PR (stuck)",
			})
		}
	}

	// 2. Analyze Issue Candidates to Close or Label
	unlabeledCount := 0
	staleCount := 0
	for _, issue := range issues {
		// Filter out pull requests from issues list (issues containing a pull_request link are PRs)
		if issue.PullRequest != nil {
			continue
		}
		if issue.State != "open" {
			continue
		}

		daysOpen := int(time.Since(issue.CreatedAt).Hours() / 24)
		daysSinceUpdate := int(time.Since(issue.UpdatedAt).Hours() / 24)

		if len(issue.Labels) == 0 {
			unlabeledCount++
			analysis.IssueCandidates = append(analysis.IssueCandidates, IssueCandidate{
				Number:   issue.Number,
				Title:    issue.Title,
				Comments: issue.Comments,
				DaysOpen: daysOpen,
				Reason:   "Unlabeled issue",
			})
		} else if daysSinceUpdate > 30 {
			staleCount++
			analysis.IssueCandidates = append(analysis.IssueCandidates, IssueCandidate{
				Number:   issue.Number,
				Title:    issue.Title,
				Comments: issue.Comments,
				DaysOpen: daysOpen,
				Reason:   "Inactive for >30 days",
			})
		} else if issue.Comments == 0 && daysOpen > 14 {
			analysis.IssueCandidates = append(analysis.IssueCandidates, IssueCandidate{
				Number:   issue.Number,
				Title:    issue.Title,
				Comments: issue.Comments,
				DaysOpen: daysOpen,
				Reason:   "No comments (2+ weeks)",
			})
		}
	}

	// 3. Generate Suggested Next Actions (High, Medium, Low Priorities)
	// High Priority
	if healthScore < 80 {
		analysis.SuggestedActions = append(analysis.SuggestedActions, MaintainerAction{
			Title:       "Improve repository health",
			Description: fmt.Sprintf("Current health score is below 80 (%d/100). Review risk alerts to improve quality.", healthScore),
			Priority:    "High",
		})
	}
	if busFactor <= 2 {
		analysis.SuggestedActions = append(analysis.SuggestedActions, MaintainerAction{
			Title:       "Onboard more core maintainers",
			Description: fmt.Sprintf("Bus factor is low (%d). Share knowledge and onboard active contributors to reduce concentration risk.", busFactor),
			Priority:    "High",
		})
	}
	if len(analysis.PRsStuck) > 0 {
		analysis.SuggestedActions = append(analysis.SuggestedActions, MaintainerAction{
			Title:       "Review stuck pull requests",
			Description: fmt.Sprintf("There are %d pull requests stuck in review. Review and merge them to avoid pipeline blockage.", len(analysis.PRsStuck)),
			Priority:    "High",
		})
	}

	// Medium Priority
	if staleCount > 0 {
		analysis.SuggestedActions = append(analysis.SuggestedActions, MaintainerAction{
			Title:       "Triage stale issues",
			Description: fmt.Sprintf("There are %d issues inactive for more than 30 days. Consider asking for status updates or closing them.", staleCount),
			Priority:    "Medium",
		})
	}
	if !hasLockFile {
		analysis.SuggestedActions = append(analysis.SuggestedActions, MaintainerAction{
			Title:       "Add package lock file",
			Description: "No dependency lock file was detected. Adding a lock file ensures predictable and reproducible builds.",
			Priority:    "Medium",
		})
	}

	// Low Priority
	if unlabeledCount > 0 {
		analysis.SuggestedActions = append(analysis.SuggestedActions, MaintainerAction{
			Title:       "Label unlabeled issues",
			Description: fmt.Sprintf("There are %d unlabeled issues. Adding appropriate labels (e.g., 'bug', 'good first issue') helps contributors.", unlabeledCount),
			Priority:    "Low",
		})
	}
	if len(analysis.SuggestedActions) == 0 {
		analysis.SuggestedActions = append(analysis.SuggestedActions, MaintainerAction{
			Title:       "Onboard new contributors",
			Description: "Repository is in excellent state! Keep onboarding and supporting first-time contributors.",
			Priority:    "Low",
		})
	}

	// 4. Contributors Deserving Appreciation
	maxAppreciate := 5
	if len(contributors) < maxAppreciate {
		maxAppreciate = len(contributors)
	}
	for i := 0; i < maxAppreciate; i++ {
		c := contributors[i]
		if c.Commits < 2 {
			continue
		}
		contribution := "Active Contributor"
		if i == 0 {
			contribution = "Top Contributor"
		} else if i < 3 {
			contribution = "Core Contributor"
		}
		analysis.ContributorsToAppreciate = append(analysis.ContributorsToAppreciate, ContributorAppreciation{
			Login:        c.Login,
			Commits:      c.Commits,
			Contribution: contribution,
		})
	}

	return analysis
}
