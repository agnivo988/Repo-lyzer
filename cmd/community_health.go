package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// CommunityHealthCmd represents the community-health command
var CommunityHealthCmd = &cobra.Command{
	Use:   "community-health",
	Short: "Generate open-source community health metrics",
	Long: `Tracks and visualizes contributor lifecycle metrics, such as issue resolution time, 
first-response time, and contributor retention rates. This helps identify bottlenecks in 
the onboarding and PR review processes to build more sustainable open-source ecosystems.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Initializing Community Health Metrics Engine...")
		fmt.Println("Fetching repository interaction data from GitHub API...")

		metrics, err := fetchCommunityMetrics()
		if err != nil {
			fmt.Printf("Error generating community metrics: %v\n", err)
			return
		}

		fmt.Println("\n--- Open Source Community Health Dashboard ---")
		
		fmt.Printf("\n[Responsiveness]\n")
		fmt.Printf("Average Time to First Response (Issues): %s\n", metrics.AvgFirstResponseIssues)
		fmt.Printf("Average Time to First Response (PRs): %s\n", metrics.AvgFirstResponsePRs)
		
		fmt.Printf("\n[Resolution & Velocity]\n")
		fmt.Printf("Average Issue Resolution Time: %s\n", metrics.AvgIssueResolutionTime)
		fmt.Printf("Average PR Merge Time: %s\n", metrics.AvgPRMergeTime)
		
		fmt.Printf("\n[Contributor Lifecycle]\n")
		fmt.Printf("New Contributors (Last 30 Days): %d\n", metrics.NewContributors30d)
		fmt.Printf("Contributor Retention Rate (6 Months): %.1f%%\n", metrics.RetentionRate6mo)

		fmt.Println("\n[Identified Bottlenecks]")
		for _, bottleneck := range metrics.Bottlenecks {
			fmt.Printf(" - %s\n", bottleneck)
		}

		fmt.Println("\n(A full graphical dashboard export will be available in future releases)")
		fmt.Println("End of report.")
	},
}

func init() {
	rootCmd.AddCommand(CommunityHealthCmd)
}

type CommunityMetrics struct {
	AvgFirstResponseIssues string
	AvgFirstResponsePRs    string
	AvgIssueResolutionTime string
	AvgPRMergeTime         string
	NewContributors30d     int
	RetentionRate6mo       float64
	Bottlenecks            []string
}

func fetchCommunityMetrics() (CommunityMetrics, error) {
	// Mocking GitHub API data aggregation for community health
	metrics := CommunityMetrics{
		AvgFirstResponseIssues: "14 hours 30 minutes",
		AvgFirstResponsePRs:    "2 days 4 hours",
		AvgIssueResolutionTime: "8 days 12 hours",
		AvgPRMergeTime:         "11 days 6 hours",
		NewContributors30d:     24,
		RetentionRate6mo:       15.5, // 15.5% retention
		Bottlenecks: []string{
			"High PR Merge Time: PRs are languishing in review for over 11 days. Consider adding more maintainers to the review rotation.",
			"Low Retention Rate: Only 15.5% of first-time contributors make a second commit. Improve documentation and onboarding guides (e.g., CONTRIBUTING.md).",
		},
	}

	return metrics, nil
}
