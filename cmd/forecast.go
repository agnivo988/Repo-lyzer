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

var forecastCmd = &cobra.Command{
	Use:   "forecast owner/repo",
	Short: "Run the AI Trend Forecasting predictions",
	Long: `Forecasts repository growth, issue accumulation, and maintainer burnout risks 
using historical data linear regression.`,
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

		overallProgress.StartStep("Fetching recent issues...")
		issues, err := client.GetIssues(owner, name, "all")
		if err != nil {
			overallProgress.Finish()
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		overallProgress.CompleteStep("Issues fetched")

		overallProgress.StartStep("Fetching historical commits...")
		commits, err := client.GetCommits(owner, name, 365) // 1 year of history
		if err != nil {
			overallProgress.Finish()
			fmt.Printf("Error: %v\n", err)
			os.Exit(1)
		}
		overallProgress.CompleteStep("Commits fetched")

		overallProgress.StartStep("Running Linear Regression models...")
		predictor := predictive.NewPredictor()

		healthResult, _ := predictor.ForecastHealth(commits, 6)          // Forecast 6 months out
		issueResult, _ := predictor.ForecastIssueAccumulation(issues, 3) // 30, 60, 90? Actually period is 1,2,3 for days or something. Wait, the implementation uses days (1, 2, 3). Let's use 3.
		overallProgress.CompleteStep("Regression models trained")

		overallProgress.StartStep("Detecting burnout risks...")
		risks, _ := predictor.ForecastContributorRisk(commits)
		overallProgress.CompleteStep("Burnout analysis complete")

		overallProgress.Finish()

		// Output Rendering
		output.RenderForecastDashboard(repo, healthResult, issueResult, risks)
	},
}

func init() {
	rootCmd.AddCommand(forecastCmd)
}
