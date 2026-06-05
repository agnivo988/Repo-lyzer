// Package cmd provides command-line interface commands for the Repo-lyzer application.
// It includes the trends command for analyzing repository trends and forecasting.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/agnivo988/Repo-lyzer/internal/output"
	"github.com/agnivo988/Repo-lyzer/internal/predictive"
	"github.com/agnivo988/Repo-lyzer/internal/progress"
	"github.com/agnivo988/Repo-lyzer/internal/temporal"
	"github.com/spf13/cobra"
)

var buildTimelineFromGitHub = temporal.BuildTimelineFromGitHub
var forecastHealthFromTimeline = predictive.ForecastHealthFromTimeline
var newPredictor = predictive.NewPredictor

type forecastPredictionOutput struct {
	Timestamp  time.Time `json:"timestamp"`
	Value      float64   `json:"value"`
	LowerBound float64   `json:"lower_bound"`
	UpperBound float64   `json:"upper_bound"`
	Confidence float64   `json:"confidence"`
	Method     string    `json:"method"`
}

type forecastTrendOutput struct {
	CurrentHealth int                        `json:"current_health"`
	ForecastModel string                     `json:"forecast_model"`
	Predictions   []forecastPredictionOutput `json:"predictions"`
}

// trendsCmd defines the "trends" command for the CLI.
// It analyzes historical trends and predicts future repository trajectory.
var trendsCmd = &cobra.Command{
	Use:   "trends owner/repo",
	Short: "Analyze repository trends and forecast future trajectory",
	Long: `Analyze historical trends and predict future repository health:
  • Commit frequency trends over time
  • Contributor growth and decline rates
  • Issue resolution velocity
  • Pull request merge patterns
  • Health score prediction using linear regression
  • Trend indicators (Improving, Declining, Stable)`,
	Example: `
  # Analyze 6-month trends (default)
  repo-lyzer trends golang/go

  # Analyze 12-month trends
  repo-lyzer trends facebook/react --months=12

  # Detailed output with monthly breakdown
  repo-lyzer trends kubernetes/kubernetes --months=6 --detailed

  # Compact JSON output
	repo-lyzer trends dashkite/dolores --months=6 --json

	# Forecast health using the timeline builder
	repo-lyzer trends kubernetes/kubernetes --months=12 --forecast --model=linear`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runTrends(args[0], cmd)
	},
}

// runTrends performs the trend analysis
func runTrends(repoArg string, cmd *cobra.Command) error {
	// Validate the repository URL format
	owner, repo, err := validateRepoURL(repoArg)
	if err != nil {
		return fmt.Errorf("invalid repository URL: %w", err)
	}

	// Get flags
	monthsFlag, _ := cmd.Flags().GetInt("months")
	detailedFlag, _ := cmd.Flags().GetBool("detailed")
	jsonFlag, _ := cmd.Flags().GetBool("json")
	forecastFlag, _ := cmd.Flags().GetBool("forecast")
	demoFlag, _ := cmd.Flags().GetBool("demo")
	modelFlag, _ := cmd.Flags().GetString("model")

	// Use months flag or default
	months := monthsFlag
	if months < 1 {
		months = 6 // Default to 6 months
	}
	if months > 24 {
		months = 24 // Max 24 months
	}

	// Track detailed flag (for future use)
	_ = detailedFlag

	if forecastFlag {
		return runTrendsForecast(owner, repo, months, modelFlag, jsonFlag, demoFlag)
	}

	// Record start time for analysis timing
	startTime := time.Now()

	// Initialize GitHub client
	client := github.NewClient()

	// Create overall progress tracker
	// Steps: repo info, commits, contributors, issues, PRs, analysis = 6 steps
	overallProgress := progress.NewOverallProgress(6)

	// Fetch repository information
	overallProgress.StartStep("Fetching repository information")
	_, err = client.GetRepo(owner, repo)
	if err != nil {
		overallProgress.Finish()
		return fmt.Errorf("failed to get repository: %w", err)
	}
	overallProgress.CompleteStep("Repository information fetched")

	// Fetch commits for the analysis period
	overallProgress.StartStep(fmt.Sprintf("Fetching commits (%dd)", daysFromMonths(months)))
	commits, err := client.GetCommits(owner, repo, daysFromMonths(months))
	if err != nil {
		overallProgress.Finish()
		return fmt.Errorf("failed to get commits: %w", err)
	}
	overallProgress.CompleteStep(fmt.Sprintf("Commits fetched (%d)", len(commits)))

	// Fetch contributors
	overallProgress.StartStep("Fetching contributor information")
	contributors, err := client.GetContributors(owner, repo)
	if err != nil {
		overallProgress.Finish()
		return fmt.Errorf("failed to get contributors: %w", err)
	}
	overallProgress.CompleteStep(fmt.Sprintf("Contributors fetched (%d)", len(contributors)))

	// Fetch issues
	overallProgress.StartStep("Fetching issues")
	issues, err := client.GetIssues(owner, repo, "all")
	if err != nil {
		overallProgress.Finish()
		return fmt.Errorf("failed to get issues: %w", err)
	}
	overallProgress.CompleteStep(fmt.Sprintf("Issues fetched (%d)", len(issues)))

	// Fetch pull requests
	overallProgress.StartStep("Fetching pull requests")
	prs, err := client.GetPullRequests(owner, repo, "all")
	if err != nil {
		overallProgress.Finish()
		return fmt.Errorf("failed to get pull requests: %w", err)
	}
	overallProgress.CompleteStep(fmt.Sprintf("Pull requests fetched (%d)", len(prs)))

	// Analyze trends
	overallProgress.StartStep("Analyzing trends")
	metrics := analyzer.AnalyzeTrends(owner, repo, commits, contributors, issues, prs, months)
	overallProgress.CompleteStep("Trend analysis complete")

	// Mark analysis as complete
	overallProgress.Finish()

	// Output results
	if jsonFlag {
		output.PrintTrendCompact(metrics)
	} else {
		output.PrintTrendMetrics(metrics, detailedFlag)

		// Track analysis duration
		duration := time.Since(startTime)
		fmt.Printf("\nAnalysis completed in %v\n", duration)
	}

	return nil
}

func runTrendsForecast(owner, repo string, months int, modelName string, jsonFlag bool, demoFlag bool) error {
	if modelName == "" {
		modelName = "linear"
	}
	if modelName != "linear" {
		return fmt.Errorf("unsupported forecast model %q", modelName)
	}

	var timeline *temporal.Timeline
	if demoFlag {
		timeline = createDemoTimeline(owner, repo, months)
	} else {
		client := github.NewClient()
		var err error
		timeline, err = buildTimelineFromGitHub(client, owner, repo, months)
		if err != nil {
			return fmt.Errorf("failed to build timeline: %w", err)
		}
	}

	predictor := newPredictor()
	predictor.ForecastHorizon = months
	forecast, err := forecastHealthFromTimeline(predictor, timeline, months)
	if err != nil {
		return fmt.Errorf("forecast health: %w", err)
	}

	currentHealth := 0
	if latest := timeline.LatestSnapshot(); latest != nil {
		currentHealth = latest.Metrics.AverageHealth
	}

	outputPayload := forecastTrendOutput{
		CurrentHealth: currentHealth,
		ForecastModel: modelName,
		Predictions:   make([]forecastPredictionOutput, 0, len(forecast.Predictions)),
	}

	for _, prediction := range forecast.Predictions {
		outputPayload.Predictions = append(outputPayload.Predictions, forecastPredictionOutput{
			Timestamp:  prediction.Timestamp,
			Value:      prediction.Value,
			LowerBound: prediction.LowerBound,
			UpperBound: prediction.UpperBound,
			Confidence: prediction.Confidence,
			Method:     prediction.Method,
		})
	}

	// If JSON flag requested, emit compact JSON payload. Otherwise, print a human-friendly report.
	if jsonFlag {
		encoder := json.NewEncoder(os.Stdout)
		encoder.SetIndent("", "  ")
		return encoder.Encode(outputPayload)
	}

	// Human-friendly textual output similar to the example in the issue description.
	fmt.Println("Repository Trends")
	fmt.Println("────────────────────")
	fmt.Println()

	// Health score per month
	fmt.Println("Health Score")
	for _, s := range timeline.Snapshots {
		if s == nil {
			continue
		}
		fmt.Printf("%s: %d\n", s.Timestamp.Format("Jan"), s.Metrics.AverageHealth)
	}
	fmt.Println()

	// Overall trend label from forecast
	trendLabel := "Stable"
	if forecast.Trend != "" {
		// capitalize
		trendLabel = strings.Title(forecast.Trend)
	}
	fmt.Printf("Trend: %s\n\n", trendLabel)

	// Contributor growth
	fmt.Println("Contributor Growth")
	for _, s := range timeline.Snapshots {
		if s == nil {
			continue
		}
		fmt.Printf("%s: %d\n", s.Timestamp.Format("Jan"), s.Metrics.ContributorCount)
	}

	// Contributor trend: simple first vs last
	contribTrend := "Stable"
	if len(timeline.Snapshots) >= 2 {
		first := timeline.Snapshots[0].Metrics.ContributorCount
		last := timeline.LatestSnapshot().Metrics.ContributorCount
		if last > first {
			contribTrend = "Increasing"
		} else if last < first {
			contribTrend = "Decreasing"
		}
	}
	fmt.Printf("\nTrend: %s\n\n", contribTrend)

	// Predicted health score for ~30 and ~90 days (1 and 3 months)
	fmt.Println("Predicted Health Score")
	if len(outputPayload.Predictions) > 0 {
		// 30 days ~ first month
		one := outputPayload.Predictions[0].Value
		fmt.Printf("30 Days: %.0f\n", one)
		// 90 days ~ third month if available
		if len(outputPayload.Predictions) >= 3 {
			three := outputPayload.Predictions[2].Value
			fmt.Printf("90 Days: %.0f\n", three)
		} else if len(outputPayload.Predictions) > 0 {
			last := outputPayload.Predictions[len(outputPayload.Predictions)-1].Value
			fmt.Printf("90 Days: %.0f\n", last)
		}
	} else {
		fmt.Println("No predictions available")
	}

	return nil
}

// daysFromMonths converts months to approximate days
func daysFromMonths(months int) int {
	return months * 30
}

func init() {
	rootCmd.AddCommand(trendsCmd)
	trendsCmd.Flags().IntP("months", "m", 6, "Number of months to analyze (1-24)")
	trendsCmd.Flags().BoolP("detailed", "d", false, "Show detailed monthly breakdown")
	trendsCmd.Flags().BoolP("json", "j", false, "Output in compact JSON format")
	trendsCmd.Flags().Bool("forecast", false, "Forecast repository health using timeline data")
	trendsCmd.Flags().Bool("demo", false, "Use demo data (no network) for reproducible output and screenshots")
	trendsCmd.Flags().String("model", "linear", "Forecast model to use")
}

// createDemoTimeline returns a deterministic timeline for demo/screenshot purposes.
func createDemoTimeline(owner, repo string, months int) *temporal.Timeline {
	// Create a 4-month demo timeline matching the user's example values.
	t := temporal.NewTimeline(owner, repo)
	// Fixed months for demo: Jan..Apr of 2024
	dates := []time.Time{
		time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC),
		time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
	}
	health := []int{72, 75, 79, 83}
	contribs := []int{12, 15, 18, 22}

	for i := 0; i < len(dates); i++ {
		s := temporal.NewSnapshot(dates[i], nil)
		s.Metrics.AverageHealth = health[i]
		s.Metrics.ContributorCount = contribs[i]
		s.Contributors = make([]string, contribs[i])
		// fill contributor names
		for j := 0; j < contribs[i]; j++ {
			s.Contributors[j] = fmt.Sprintf("contrib-%d", j+1)
		}
		_ = t.AddSnapshot(s)
	}

	return t
}
