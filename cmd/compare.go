// Package cmd provides command-line interface commands for the Repo-lyzer application.
// It includes commands for analyzing repositories, comparing repositories, and running the interactive menu.
package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
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
		r1 := strings.Split(args[0], "/")
		r2 := strings.Split(args[1], "/")

		if len(r1) != 2 || len(r2) != 2 {
			return fmt.Errorf("repositories must be in owner/repo format")
		}

		client := github.NewClient()

		repo1, err := client.GetRepo(r1[0], r1[1])
		if err != nil {
			return err
		}

		_, _ = client.GetLanguages(r1[0], r1[1])
		commits1, _ := client.GetCommits(r1[0], r1[1], 14)
		contributors1, err := client.GetContributorsWithAvatars(r1[0], r1[1], 15)
		if err != nil {
			fmt.Printf("Error fetching contributors for %s/%s: %v\n", r1[0], r1[1], err)
			return err
		}
		_, _ = client.GetFileTree(r1[0], r1[1], repo1.DefaultBranch)
		bus1, risk1 := analyzer.BusFactor(contributors1)

		maturityScore1, maturityLevel1 :=
			analyzer.RepoMaturityScore(repo1, len(commits1), len(contributors1), false)

		// ---------- Fetch Repo 2 ----------
		repo2, err := client.GetRepo(r2[0], r2[1])
		if err != nil {
			return err
		}

		_, _ = client.GetLanguages(r2[0], r2[1])
		commits2, _ := client.GetCommits(r2[0], r2[1], 14)
		contributors2, err := client.GetContributorsWithAvatars(r2[0], r2[1], 15)
		if err != nil {
			fmt.Printf("Error fetching contributors for %s/%s: %v\n", r2[0], r2[1], err)
			return err
		}
		_, _ = client.GetFileTree(r2[0], r2[1], repo2.DefaultBranch)
		bus2, risk2 := analyzer.BusFactor(contributors2)

		maturityScore2, maturityLevel2 :=
			analyzer.RepoMaturityScore(repo2, len(commits2), len(contributors2), false)

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
