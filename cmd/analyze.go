// Package cmd provides command-line interface commands for the Repo-lyzer application.
// It includes commands for analyzing repositories, comparing repositories, and running the interactive menu.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

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

// validateRepoURL validates the repository URL format and provides clear error messages
func validateRepoURL(repoArg string) (owner, repo string, err error) {
	if repoArg == "" {
		return "", "", fmt.Errorf("repository URL cannot be empty")
	}

	if strings.Contains(repoArg, " ") {
		return "", "", fmt.Errorf("repository URL cannot contain spaces")
	}

	parts := strings.Split(repoArg, "/")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("repository must be in 'owner/repo' format (found %d parts separated by '/')", len(parts))
	}

	owner, repo = parts[0], parts[1]

	if owner == "" {
		return "", "", fmt.Errorf("owner name cannot be empty")
	}

	if repo == "" {
		return "", "", fmt.Errorf("repository name cannot be empty")
	}

	// Basic validation for GitHub username/repo name patterns
	if len(owner) > 39 {
		return "", "", fmt.Errorf("owner name is too long (maximum 39 characters)")
	}

	if len(owner) < 1 {
		return "", "", fmt.Errorf("owner name is too short (minimum 1 character)")
	}

	if strings.HasPrefix(owner, "-") || strings.HasSuffix(owner, "-") {
		return "", "", fmt.Errorf("owner name cannot start or end with a hyphen")
	}

	if strings.Contains(owner, "--") {
		return "", "", fmt.Errorf("owner name cannot contain consecutive hyphens")
	}

	// Check for valid characters (alphanumeric, hyphens)
	for _, char := range owner {
		if !((char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-') {
			return "", "", fmt.Errorf("owner name contains invalid character '%c' (only alphanumeric characters and hyphens allowed)", char)
		}
	}

	if len(repo) > 100 {
		return "", "", fmt.Errorf("repository name is too long (maximum 100 characters)")
	}

	if len(repo) < 1 {
		return "", "", fmt.Errorf("repository name is too short (minimum 1 character)")
	}

	// Repository names can contain more characters than usernames
	for _, char := range repo {
		if char == ' ' || char == '\t' || char == '\n' || char == '\r' {
			return "", "", fmt.Errorf("repository name cannot contain whitespace")
		}
	}

	return owner, repo, nil
}

// runDryRun performs a dry run of the analysis, validating the repository URL
// and displaying what metrics would be calculated without making API calls.
func runDryRun(repoArg string) error {
	fmt.Printf("🔍 Dry Run Mode - Validating repository: %s\n\n", repoArg)

	// Use the same validation as the full run
	owner, repo, err := validateRepoURL(repoArg)
	if err != nil {
		return fmt.Errorf("invalid repository URL: %w", err)
	}

	fmt.Printf("✅ Repository URL format is valid: %s/%s\n", owner, repo)
	fmt.Println("📊 The following metrics would be calculated:")
	fmt.Println("  • Repository information (stars, forks, description, etc.)")
	fmt.Println("  • Programming languages used")
	fmt.Println("  • Commit activity over the last 365 days")
	fmt.Println("  • Repository health score")
	fmt.Println("  • Daily commit activity (last 14 days)")
	fmt.Println("  • Contributor information")
	fmt.Println("  • Bus factor and risk assessment")
	fmt.Println("  • Repository maturity score and level")
	fmt.Println("  • Recruiter summary with key insights")
	fmt.Println()
	fmt.Println("💡 This dry run does not consume API rate limits or perform actual computations.")
	fmt.Println("   Run without --dry-run to execute the full analysis.")

	return nil
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
		dryRun, _ := cmd.Flags().GetBool("dry-run")

		if dryRun {
			return runDryRun(args[0])
		}

		// Validate the repository URL format
		owner, repo, err := validateRepoURL(args[0])
		if err != nil {
			return fmt.Errorf("invalid repository URL: %w", err)
		}

		// Record start time for analysis timing
		startTime := time.Now()

		// Initialize GitHub client
		client := github.NewClient()

		// Fetch repository information
		repoInfo, err := client.GetRepo(owner, repo)
		if err != nil {
			// Check if it's a private repo error and no token is set
			if strings.Contains(err.Error(), "repository not found") && !client.HasToken() {
				fmt.Print("This appears to be a private repository. Please enter your GitHub access token: ")
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					token := strings.TrimSpace(scanner.Text())
					if token != "" {
						client.SetToken(token)
						// Retry fetching the repo with the token
						repoInfo, err = client.GetRepo(owner, repo)
						if err != nil {
							return fmt.Errorf("failed to access repository even with token: %w", err)
						}
					} else {
						return fmt.Errorf("no token provided, cannot access private repository")
					}
				} else {
					return fmt.Errorf("failed to read token input")
				}
			} else {
				return err
			}
		}

		// Fetch programming languages used in the repository
		langs, err := client.GetLanguages(owner, repo)
		if err != nil {
			return fmt.Errorf("failed to get languages: %w", err)
		}

		// Fetch commits from the last 365 days
		commits, err := client.GetCommits(owner, repo, 365)
		if err != nil {
			return fmt.Errorf("failed to get commits: %w", err)
		}

		// Calculate repository health score
		score := analyzer.CalculateHealth(repoInfo, commits)

		// Analyze commit activity per day
		activity := analyzer.CommitsPerDay(commits)

		// Fetch contributors
		contributors, err := client.GetContributorsWithAvatars(owner, repo, 15)
		if err != nil {
			return err
		}

		// Calculate bus factor and risk level
		busFactor, busRisk := analyzer.BusFactor(contributors)

		// Calculate repository maturity score and level
		maturityScore, maturityLevel :=
			analyzer.RepoMaturityScore(
				repoInfo,
				len(commits),
				len(contributors),
				false, // Assuming no releases check for simplicity
			)

		// Build recruiter summary
		summary := analyzer.BuildRecruiterSummary(
			repoInfo.FullName,
			repoInfo.Forks,
			repoInfo.Stars,
			len(commits),
			len(contributors),
			maturityScore,
			maturityLevel,
			busFactor,
			busRisk,
		)

		// Output the analysis results
		output.PrintRepo(repoInfo)
		output.PrintLanguages(langs)
		output.PrintCommitActivity(activity, 14)
		output.PrintHealth(score)
		output.PrintGitHubAPIStatus(client)
		output.PrintRecruiterSummary(summary)

		// Display analysis time
		duration := time.Since(startTime)
		fmt.Printf("\n⏱️  Analysis completed in %.2f seconds\n", duration.Seconds())

		return nil
	},
}

func init() {
	rootCmd.AddCommand(analyzeCmd)
	analyzeCmd.Flags().Bool("dry-run", false, "Validate repository URL and show what metrics would be calculated without making API calls")
}
