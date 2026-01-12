// Package cmd provides command-line interface commands for the Repo-lyzer application.
// It includes commands for analyzing repositories, comparing repositories, and running the interactive menu.
package cmd

import (
	"fmt"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/agnivo988/Repo-lyzer/internal/output"
	"github.com/spf13/cobra"
)

// RunAnalyze executes the analyze command for a given GitHub repository.
// It takes the owner and repository name, performs comprehensive analysis including
// repository info, languages, commits, contributors, and generates various reports.
// Parameters:
//   - owner: GitHub username or organization name
//   - repo: Repository name
// Returns an error if the analysis fails.
func RunAnalyze(owner, repo string) error {
	args := []string{owner + "/" + repo}
	analyzeCmd.SetArgs(args)
	return analyzeCmd.Execute()
}


// analyzeCmd defines the "analyze" command for the CLI.
// It analyzes a single GitHub repository and prints various metrics and reports.
// Usage example:
//   repo-lyzer analyze octocat/Hello-World
// This will fetch repository data, calculate health scores, bus factor, maturity,
// and display comprehensive analysis results including languages, commit activity,
// contributor information, and a recruiter summary.
var analyzeCmd = &cobra.Command{
	Use:   "analyze owner/repo",
	Short: "Analyze a GitHub repository",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse the repository argument into owner and repo parts
		parts := strings.Split(args[0], "/")
		if len(parts) != 2 {
			return fmt.Errorf("repository must be in owner/repo format")
		}

		// Initialize GitHub client
		client := github.NewClient()

		// Fetch repository information
		repo, err := client.GetRepo(parts[0], parts[1])
		if err != nil {
			return err
		}

		// Fetch programming languages used in the repository
		langs, err := client.GetLanguages(parts[0], parts[1])
		if err != nil {
			return fmt.Errorf("failed to get languages: %w", err)
		}

		// Fetch commits from the last 365 days
		commits, err := client.GetCommits(parts[0], parts[1], 365)
		if err != nil {
			return fmt.Errorf("failed to get commits: %w", err)
		}

		// Calculate repository health score
		score := analyzer.CalculateHealth(repo, commits)

		// Analyze commit activity per day
		activity := analyzer.CommitsPerDay(commits)

		// Fetch contributors
		contributors, err := client.GetContributors(parts[0], parts[1])
		if err != nil {
			return err
		}

		// Calculate bus factor and risk level
		busFactor, busRisk := analyzer.BusFactor(contributors)

		// Calculate repository maturity score and level
		maturityScore, maturityLevel :=
			analyzer.RepoMaturityScore(
				repo,
				len(commits),
				len(contributors),
				false, // Assuming no releases check for simplicity
			)

		// Build recruiter summary
		summary := analyzer.BuildRecruiterSummary(
			repo.FullName,
			repo.Forks,
			repo.Stars,
			len(commits),
			len(contributors),
			maturityScore,
			maturityLevel,
			busFactor,
			busRisk,
		)

		// Output the analysis results
		output.PrintRepo(repo)
		output.PrintLanguages(langs)
		output.PrintCommitActivity(activity, 14)
		output.PrintHealth(score)
		output.PrintGitHubAPIStatus(client)
		output.PrintRecruiterSummary(summary)

		return nil
	},
}
