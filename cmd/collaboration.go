// Package cmd provides command-line interface commands for the Repo-lyzer application.
// It includes the collaboration command for analyzing team collaboration patterns.
package cmd

import (
	"fmt"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/agnivo988/Repo-lyzer/internal/output"
	"github.com/agnivo988/Repo-lyzer/internal/progress"
	"github.com/spf13/cobra"
)

// Avoid excessive GitHub API usage (1 review request per PR)
const maxPRsForReviewFetch = 50

// collaborationCmd defines the "collaboration" command for the CLI.
// It analyzes how contributors collaborate within a repository.
var collaborationCmd = &cobra.Command{
	Use:   "collaboration owner/repo",
	Short: "Analyze team collaboration patterns in a repository",
	Long: `Analyze team collaboration patterns including:
  • Code review participation rates
  • Cross-contributor commit patterns
  • Knowledge silos detection
  • Collaboration network metrics
  • Review response time statistics`,
	Example: `
  # Analyze collaboration patterns
  repo-lyzer collaboration golang/go

  # Analyze with network visualization
  repo-lyzer collaboration facebook/react`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runCollaboration(args[0])
	},
}

// runCollaboration performs the collaboration analysis
func runCollaboration(repoArg string) error {
	// Validate the repository URL format
	owner, repo, err := validateRepoURL(repoArg)
	if err != nil {
		return fmt.Errorf("invalid repository URL: %w", err)
	}

	// Record start time for analysis timing
	startTime := time.Now()

	// Initialize GitHub client
	client := github.NewClient()

	// Create overall progress tracker
	// Steps: contributors, PRs, reviews, analysis = 4 steps
	overallProgress := progress.NewOverallProgress(4)

	// Fetch contributors
	overallProgress.StartStep("👥 Fetching contributor information")

	contributors, err := client.GetContributors(owner, repo)
	if err != nil {
		overallProgress.Finish()
		return fmt.Errorf("failed to get contributors: %w", err)
	}

	overallProgress.CompleteStep(
		fmt.Sprintf("Contributors fetched (%d)", len(contributors)),
	)

	// Fetch pull requests
	overallProgress.StartStep("📝 Fetching pull requests")

	prs, err := client.GetPullRequests(owner, repo, "all")
	if err != nil {
		overallProgress.Finish()
		return fmt.Errorf("failed to get pull requests: %w", err)
	}

	overallProgress.CompleteStep(
		fmt.Sprintf("Pull requests fetched (%d)", len(prs)),
	)

	// Fetch reviews for pull requests
	overallProgress.StartStep("👀 Fetching PR reviews")

	reviews := make(map[int][]github.Review)

	fetchLimit := len(prs)
	if fetchLimit > maxPRsForReviewFetch {
		fetchLimit = maxPRsForReviewFetch

		fmt.Printf(
			"⚠️  Limiting review fetch to %d PRs to avoid GitHub API rate limiting\n",
			maxPRsForReviewFetch,
		)
	}

	reviewErrors := 0

	for _, pr := range prs[:fetchLimit] {
		prReviews, err := client.GetPullRequestReviews(
			owner,
			repo,
			pr.Number,
		)

		if err != nil {
			reviewErrors++
			continue
		}

		if len(prReviews) > 0 {
			reviews[pr.Number] = prReviews
		}
	}

	if reviewErrors > 0 {
		fmt.Printf(
			"⚠️  Failed to fetch reviews for %d PRs\n",
			reviewErrors,
		)
	}

	overallProgress.CompleteStep(
		fmt.Sprintf(
			"Reviews fetched for %d PRs (%d errors)",
			len(reviews),
			reviewErrors,
		),
	)

	// Analyze collaboration
	overallProgress.StartStep("📊 Analyzing collaboration patterns")

	metrics := analyzer.AnalyzeCollaboration(
		contributors,
		prs,
		reviews,
	)

	overallProgress.CompleteStep("Collaboration analysis complete")

	// Mark analysis as complete
	overallProgress.Finish()

	// Output results
	output.PrintCollaborationInsights(metrics)
	output.PrintCollaborationNetwork(metrics)

	// Track analysis duration
	duration := time.Since(startTime)

	fmt.Printf("\n⏱️  Analysis completed in %v\n", duration)

	return nil
}

func init() {
	rootCmd.AddCommand(collaborationCmd)
}
