package cmd

import (
	"fmt"
	"sync"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/agnivo988/Repo-lyzer/internal/output"
	"github.com/agnivo988/Repo-lyzer/internal/progress"
	"github.com/spf13/cobra"
)

var analyzePRCmd = &cobra.Command{
	Use:   "analyze-pr owner/repo",
	Short: "Analyze Pull Request metrics for a GitHub repository",
	Long: `Analyze pull request metrics including:
  • Average time to merge
  • Review participation (% of PRs with 2+ reviewers)
  • PR size distribution
  • Abandoned PR ratio
  • First-time contributor friendliness

Note: Each PR requires 2 API calls (details + reviews). With authentication,
you have 5,000 requests/hour. Default limit is 100 PRs (200 requests).
Use --limit 0 for no limit, but be cautious of rate limits.

Examples:
  repo-lyzer analyze-pr golang/go
  repo-lyzer analyze-pr microsoft/vscode --state closed --limit 50
  repo-lyzer analyze-pr octocat/Hello-World --json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		state, _ := cmd.Flags().GetString("state")
		limit, _ := cmd.Flags().GetInt("limit")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		// Validate the repository URL format
		owner, repo, err := validateRepoURL(args[0])
		if err != nil {
			return fmt.Errorf("invalid repository URL: %w", err)
		}

		// Initialize GitHub client
		client := github.NewClient()

		// Create overall progress tracker
		overallProgress := progress.NewOverallProgress(3)

		// Fetch pull requests
		if !jsonOutput {
			overallProgress.StartStep("🔍 Fetching pull requests")
		}

		var prs []github.PullRequest
		if limit > 0 {
			prs, err = client.GetPullRequestsWithLimit(owner, repo, state, limit)
		} else {
			prs, err = client.GetPullRequests(owner, repo, state)
		}

		if err != nil {
			if !jsonOutput {
				overallProgress.Finish()
			}
			return fmt.Errorf("failed to fetch pull requests: %w", err)
		}

		if len(prs) == 0 {
			if !jsonOutput {
				overallProgress.Finish()
				fmt.Printf("No pull requests found for %s/%s with state '%s'\n", owner, repo, state)
			}
			return nil
		}

		if !jsonOutput {
			overallProgress.CompleteStep(fmt.Sprintf("Found %d pull requests", len(prs)))
		}

		if !jsonOutput {
			overallProgress.StartStep("🔄 Fetching PR details and reviews")
		}

		// Fetch PR details and reviews concurrently with worker pool
		type prResult struct {
			pr      *github.PullRequest
			reviews []github.Review
			index   int
			err     error
		}

		workers := 10
		semaphore := make(chan struct{}, workers)
		results := make(chan prResult, len(prs))

		var wg sync.WaitGroup

		// Launch goroutines for each PR
		for i, pr := range prs {
			wg.Add(1)

			go func(prNumber, index int) {
				defer wg.Done()

				semaphore <- struct{}{}
				defer func() { <-semaphore }()

				// Fetch detailed PR info
				prDetails, err := client.GetPullRequestDetails(owner, repo, prNumber)
				if err != nil {
					results <- prResult{
						index: index,
						err:   fmt.Errorf("PR #%d details: %w", prNumber, err),
					}
					return
				}

				// Fetch reviews
				prReviews, err := client.GetPullRequestReviews(owner, repo, prNumber)
				if err != nil {
					results <- prResult{
						index: index,
						err:   fmt.Errorf("PR #%d reviews: %w", prNumber, err),
					}
					return
				}

				results <- prResult{
					pr:      prDetails,
					reviews: prReviews,
					index:   index,
				}
			}(pr.Number, i)
		}

		// Close results channel after all goroutines complete
		go func() {
			wg.Wait()
			close(results)
		}()

		// Collect results
		updatedPRs := make([]*github.PullRequest, len(prs))
		reviews := make(map[int][]github.Review)

		errorCount := 0
		i := 0

		for result := range results {
			if !jsonOutput {
				i++
				percentage := i * 100 / len(prs)

				overallProgress.UpdateStep(
					fmt.Sprintf(
						"🔄 Fetching PR details and reviews [%d/%d - %d%%]",
						i,
						len(prs),
						percentage,
					),
				)
			}

			if result.err != nil {
				errorCount++
				continue
			}

			if result.pr == nil {
				errorCount++
				continue
			}

			updatedPRs[result.index] = result.pr
			reviews[result.pr.Number] = result.reviews
		}

		// Filter out nil entries (failed fetches)
		var finalPRs []github.PullRequest
		for _, pr := range updatedPRs {
			if pr != nil {
				finalPRs = append(finalPRs, *pr)
			}
		}

		if !jsonOutput {
			overallProgress.CompleteStep(
				fmt.Sprintf("PR details fetched (%d errors)", errorCount),
			)
		}

		if len(finalPRs) == 0 {
			if !jsonOutput {
				overallProgress.Finish()
			}
			return fmt.Errorf("no PRs could be fetched successfully")
		}

		// Use finalPRs instead of prs
		prs = finalPRs

		// Analyze pull requests
		if !jsonOutput {
			overallProgress.StartStep("📊 Analyzing pull request metrics")
		}

		analytics := analyzer.AnalyzePullRequests(prs, reviews)

		if !jsonOutput {
			overallProgress.CompleteStep("Pull request analysis complete")
			overallProgress.Finish()
		}

		// Output results
		if jsonOutput {
			jsonStr, err := output.FormatPRAnalyticsJSON(analytics)
			if err != nil {
				return fmt.Errorf("failed to format JSON: %w", err)
			}
			fmt.Println(jsonStr)
		} else {
			output.PrintPRAnalytics(analytics)
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzePRCmd)

	analyzePRCmd.Flags().String(
		"state",
		"all",
		"Filter PRs by state: open, closed, or all",
	)

	analyzePRCmd.Flags().Int(
		"limit",
		100,
		"Limit number of PRs to analyze (0 = no limit, use with caution)",
	)

	analyzePRCmd.Flags().Bool(
		"json",
		false,
		"Output results as JSON",
	)
}
