// Package cmd provides command-line interface commands for the Repo-lyzer application.
package cmd

import (
	"fmt"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/spf13/cobra"
)

// summaryCmd defines the "summary" command for the CLI.
// It provides a quick 6-line summary of a GitHub repository.
// Usage example:
//
//	repo-lyzer summary octocat/Hello-World
var summaryCmd = &cobra.Command{
	Use:   "summary owner/repo",
	Short: "Display a quick 6-line repository summary",
	Long: `Display a quick 6-line summary of a GitHub repository for fast evaluation.

Perfect for recruiters and developers who need a rapid assessment of a repository 
without opening the full interactive dashboard.

Summary includes:
  • Commits in last 30 days
  • Top programming language
  • Total contributors
  • Repository health score (0-100)
  • Last commit timestamp

Examples:
  # Quick check on a repository
  repo-lyzer summary golang/go

  # Evaluate multiple repositories quickly
  repo-lyzer summary facebook/react
  repo-lyzer summary vuejs/vue
  repo-lyzer summary angular/angular`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate the repository URL format
		owner, repo, err := validateRepoURL(args[0])
		if err != nil {
			return fmt.Errorf("invalid repository URL: %w", err)
		}

		// Initialize GitHub client
		client := github.NewClient()

		// Fetch repository information
		repoInfo, err := client.GetRepo(owner, repo)
		if err != nil {
			return fmt.Errorf("failed to get repository: %w", err)
		}

		// Fetch commits from the last 30 days
		commits30d, err := client.GetCommits(owner, repo, 30)
		if err != nil {
			return fmt.Errorf("failed to get commits: %w", err)
		}

		// Fetch programming languages
		langs, err := client.GetLanguages(owner, repo)
		if err != nil {
			return fmt.Errorf("failed to get languages: %w", err)
		}

		// Fetch contributors
		contributors, err := client.GetContributors(owner, repo)
		if err != nil {
			return fmt.Errorf("failed to get contributors: %w", err)
		}

		// Fetch commits for health calculation (last 365 days)
		commitsYear, err := client.GetCommits(owner, repo, 365)
		if err != nil {
			return fmt.Errorf("failed to get yearly commits: %w", err)
		}

		// Calculate health score
		healthScore := analyzer.CalculateHealth(repoInfo, commitsYear)

		// Get top language
		topLang := getTopLanguage(langs)
		if topLang == "" && repoInfo.Language != "" {
			topLang = repoInfo.Language
		}
		if topLang == "" {
			topLang = "Unknown"
		}

		// Format last commit date
		lastCommit := formatTimeAgo(repoInfo.PushedAt)

		// Print the 5-line summary
		fmt.Printf("📊 Repository Summary: %s\n", repoInfo.FullName)
		fmt.Printf("   Commits (30d): %d\n", len(commits30d))
		fmt.Printf("   Top Language: %s\n", topLang)
		fmt.Printf("   Contributors: %d\n", len(contributors))
		fmt.Printf("   Health Score: %d/100\n", healthScore)
		fmt.Printf("   Last Commit: %s\n", lastCommit)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(summaryCmd)
}
