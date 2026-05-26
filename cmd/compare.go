// Package cmd provides command-line interface commands for the Repo-lyzer application.
// It includes commands for analyzing repositories, comparing repositories, and running the interactive menu.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/agnivo988/Repo-lyzer/internal/progress"
)

// RunCompare executes the compare command for two GitHub repositories.
// It takes two repository identifiers in owner/repo format, analyzes both repositories,
// and displays a comparison table with metrics like stars, forks, commits, contributors,
// bus factor, and maturity scores.
// Parameters:
//   - r1: First repository in owner/repo format
//   - r2: Second repository in owner/repo format
//
// Returns an error if the comparison fails.
func RunCompare(r1, r2 string) error {
	compareCmd.SetArgs([]string{r1, r2})
	return compareCmd.Execute()
}

var compareCmd = &cobra.Command{
	Use:   "compare owner1/repo1 owner2/repo2",
	Short: "Compare two GitHub repositories side-by-side",
	Long: `Compare two GitHub repositories and display a side-by-side comparison
of their key metrics and health indicators.

Comparison includes:
  • Stars, Forks, and Open Issues
  • Commit activity (past year)
  • Contributor count and engagement
  • Bus Factor and risk assessment  
  • Repository maturity scores
  • Verdict on which repository is more mature/stable

Examples:
  # Compare popular frameworks
  repo-lyzer compare facebook/react vuejs/vue

  # Compare similar tools
  repo-lyzer compare golang/go rust-lang/rust

  # Compare forks
  repo-lyzer compare original/repo fork/repo`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {

		// Parse repo names
		owner1, repo1Name, err := validateRepoURL(args[0])
		if err != nil {
			return fmt.Errorf("invalid first repository URL: %w", err)
		}
		owner2, repo2Name, err := validateRepoURL(args[1])
		if err != nil {
			return fmt.Errorf("invalid second repository URL: %w", err)
		}

		client := github.NewClient()

		// Create progress spinner
		spinner := progress.NewSpinner()

		// Fetch first repository
		spinner.Start(fmt.Sprintf("🔍 Analyzing %s/%s...", owner1, repo1Name))
		repo1, err := client.GetRepo(owner1, repo1Name)
		if err != nil {
			spinner.Stop()
			if strings.Contains(err.Error(), "repository not found") && !client.HasToken() {
				fmt.Printf("The first repository %s/%s appears to be private. Please enter your GitHub access token: ", owner1, repo1Name)
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					token := strings.TrimSpace(scanner.Text())
					if token != "" {
						client.SetToken(token)
						spinner.Start(fmt.Sprintf("🔍 Analyzing %s/%s...", owner1, repo1Name))
						repo1, err = client.GetRepo(owner1, repo1Name)
						if err != nil {
							spinner.Stop()
							return fmt.Errorf("failed to access first repository even with token: %w", err)
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

		_, _ = client.GetLanguages(owner1, repo1Name)
		commits1, _ := client.GetCommits(owner1, repo1Name, 14)
		contributors1, err := client.GetContributorsWithAvatars(owner1, repo1Name, 15)
		if err != nil {
			spinner.Stop()
			fmt.Printf("Error fetching contributors for %s/%s: %v\n", owner1, repo1Name, err)
			return err
		}
		_, _ = client.GetFileTree(owner1, repo1Name, repo1.DefaultBranch)
		bus1, risk1 := analyzer.BusFactor(contributors1)

		maturityScore1, maturityLevel1 :=
			analyzer.RepoMaturityScore(repo1, len(commits1), len(contributors1), false)

		spinner.StopWithMessage(fmt.Sprintf("Analyzed %s/%s", owner1, repo1Name))

		// ---------- Fetch Repo 2 ----------
		spinner.Start(fmt.Sprintf("🔍 Analyzing %s/%s...", owner2, repo2Name))
		repo2, err := client.GetRepo(owner2, repo2Name)
		if err != nil {
			spinner.Stop()
			if strings.Contains(err.Error(), "repository not found") && !client.HasToken() {
				fmt.Printf("The second repository %s/%s appears to be private. Please enter your GitHub access token: ", owner2, repo2Name)
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					token := strings.TrimSpace(scanner.Text())
					if token != "" {
						client.SetToken(token)
						spinner.Start(fmt.Sprintf("🔍 Analyzing %s/%s...", owner2, repo2Name))
						repo2, err = client.GetRepo(owner2, repo2Name)
						if err != nil {
							spinner.Stop()
							return fmt.Errorf("failed to access second repository even with token: %w", err)
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

		_, _ = client.GetLanguages(owner2, repo2Name)
		commits2, _ := client.GetCommits(owner2, repo2Name, 14)
		contributors2, err := client.GetContributorsWithAvatars(owner2, repo2Name, 15)
		if err != nil {
			spinner.Stop()
			fmt.Printf("Error fetching contributors for %s/%s: %v\n", owner2, repo2Name, err)
			return err
		}
		_, _ = client.GetFileTree(owner2, repo2Name, repo2.DefaultBranch)
		bus2, risk2 := analyzer.BusFactor(contributors2)

		maturityScore2, maturityLevel2 :=
			analyzer.RepoMaturityScore(repo2, len(commits2), len(contributors2), false)

		spinner.StopWithMessage(fmt.Sprintf("Analyzed %s/%s", owner2, repo2Name))

		// ---------- Output Table ----------
		fmt.Println("\n📊 Repository Comparison")

		table := tablewriter.NewWriter(os.Stdout)
		table.Header([]string{"Metric", repo1.FullName, repo2.FullName})

		table.Append([]string{"⭐ Stars",
			fmt.Sprintf("%d", repo1.Stars),
			fmt.Sprintf("%d", repo2.Stars),
		})

		table.Append([]string{"🍴 Forks",
			fmt.Sprintf("%d", repo1.Forks),
			fmt.Sprintf("%d", repo2.Forks),
		})

		table.Append([]string{"📦 Commits (1y)",
			fmt.Sprintf("%d", len(commits1)),
			fmt.Sprintf("%d", len(commits2)),
		})

		table.Append([]string{"👥 Contributors",
			fmt.Sprintf("%d", len(contributors1)),
			fmt.Sprintf("%d", len(contributors2)),
		})

		table.Append([]string{"⚠️ Bus Factor",
			fmt.Sprintf("%d (%s)", bus1, risk1),
			fmt.Sprintf("%d (%s)", bus2, risk2),
		})

		table.Append([]string{"🏗️ Maturity",
			fmt.Sprintf("%s (%d)", maturityLevel1, maturityScore1),
			fmt.Sprintf("%s (%d)", maturityLevel2, maturityScore2),
		})

		// Check if repositories are identical
		if repo1.Stars == repo2.Stars &&
			repo1.Forks == repo2.Forks &&
			len(commits1) == len(commits2) &&
			len(contributors1) == len(contributors2) &&
			bus1 == bus2 &&
			maturityScore1 == maturityScore2 {

			fmt.Println("\n✅ No differences found between the two repositories.")
			fmt.Println("Both repositories have identical metrics.")
			return nil
		}

		table.Render()

		// ---------- Verdict ----------
		fmt.Println("\n Verdict")
		if maturityScore1 > maturityScore2 {
			fmt.Printf("➡️ %s appears more mature and stable.\n", repo1.FullName)
		} else if maturityScore2 > maturityScore1 {
			fmt.Printf("➡️ %s appears more mature and stable.\n", repo2.FullName)
		} else {
			fmt.Println("➡️ Both repositories are similarly mature.")
		}

		return nil
	},
}

func init() {
	rootCmd.AddCommand(compareCmd)
}

// countTreeStats counts files, directories, and total size from tree entries
func countTreeStats(tree []github.TreeEntry) (files, dirs, totalSize int) {
	for _, entry := range tree {
		if entry.Type == "blob" {
			files++
			totalSize += entry.Size
		} else if entry.Type == "tree" {
			dirs++
		}
	}
	return
}
