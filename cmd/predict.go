package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/agnivo988/Repo-lyzer/internal/output"
	"github.com/agnivo988/Repo-lyzer/internal/predictive"
	"github.com/agnivo988/Repo-lyzer/internal/progress"
	"github.com/spf13/cobra"
)

var predictContributor string

var predictCmd = &cobra.Command{
	Use:   "predict owner/repo",
	Short: "Run the AI Contribution Engine predictions",
	Long: `Predicts the best first issues, calculates pull request merge probabilities, 
and optionally determines contributor fit using heuristic AI logic.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		repo := args[0]
		parts := strings.Split(repo, "/")
		if len(parts) != 2 {
			fmt.Println("Error: repository must be in owner/repo format")
			os.Exit(1)
		}
		owner, name := parts[0], parts[1]

		client := github.NewClient()

		overallProgress := progress.NewOverallProgress(4)

		overallProgress.StartStep("Fetching issues...")
		issues, err := client.GetIssues(owner, name, "open")
		if err != nil {
			overallProgress.Finish()
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		overallProgress.CompleteStep("Issues fetched")

		overallProgress.StartStep("Fetching pull requests...")
		prs, err := client.GetPullRequests(owner, name, "open")
		if err != nil {
			overallProgress.Finish()
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		overallProgress.CompleteStep("Pull requests fetched")

		overallProgress.StartStep("Fetching commits...")
		commits, err := client.GetCommits(owner, name, 365)
		if err != nil {
			overallProgress.Finish()
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		overallProgress.CompleteStep("Commits fetched")

		overallProgress.StartStep("Running AI predictions...")

		// 1. Difficulty Estimation
		bestIssues := predictive.PredictBestFirstIssues(issues, 5)

		// 2. Merge Probability
		prProbs := predictive.PredictTopPRs(prs, 5)

		// 3. Contributor Fit Analysis
		var fit predictive.ContributorFit
		if predictContributor != "" && len(bestIssues) > 0 {
			var targetIssue github.Issue
			for _, iss := range issues {
				if iss.Number == bestIssues[0].IssueNumber {
					targetIssue = iss
					break
				}
			}
			fit = predictive.AnalyzeContributorFit(targetIssue, predictContributor, commits)
		}

		overallProgress.Finish()

		// Output Rendering
		output.RenderPredictiveDashboard(repo, bestIssues, prProbs, fit)
	},
}

func init() {
	predictCmd.Flags().StringVarP(&predictContributor, "contributor", "c", "", "GitHub username to run Contributor Fit Analysis against the easiest issue")
	rootCmd.AddCommand(predictCmd)
}
