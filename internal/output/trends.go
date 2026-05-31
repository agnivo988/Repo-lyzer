// Package output provides formatting and output functions for Repo-lyzer.
// This file implements trend analysis visualization.
package output

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/charmbracelet/lipgloss"
)

// Trend output styles
var (
	trendHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#00E5FF")).
				Padding(1, 0)

	trendLabelStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#888888"))

	trendValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF"))

	improvingStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF87"))

	decliningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F5F"))

	stableStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB000"))
)

// PrintTrendMetrics prints comprehensive trend analysis results
func PrintTrendMetrics(metrics *analyzer.TrendMetrics, detailedFlag bool) {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println(trendHeaderStyle.Render("📈 REPOSITORY TREND ANALYSIS"))
	fmt.Println(strings.Repeat("=", 60))

	// Repository info
	fmt.Printf("\nRepository: %s/%s\n", metrics.Owner, metrics.Repo)
	fmt.Printf("Analysis Period: Last %d months\n\n", metrics.AnalysisPeriod)

	// Overall Trend
	printOverallTrend(metrics)

	// Commit Trends
	printCommitTrends(metrics, detailedFlag)

	// Contributor Trends
	printContributorTrends(metrics)

	// Issue Resolution Trends
	printIssueTrends(metrics)

	// PR Trends
	printPRTrends(metrics)

	// Health Score Prediction
	printHealthPrediction(metrics)

	// Monthly Sparkline
	if len(metrics.MonthlyData) > 0 {
		printMonthlySparkline(metrics)
	}

	fmt.Println()
}

// printOverallTrend prints the overall trend assessment
func printOverallTrend(metrics *analyzer.TrendMetrics) {
	fmt.Println("OVERALL ASSESSMENT")
	fmt.Println(strings.Repeat("-", 40))

	trendStyle := getTrendStyle(metrics.OverallTrend)
	fmt.Printf("   Trend: %s\n\n", trendStyle.Render(string(metrics.OverallTrend)))
	fmt.Printf("   %s\n\n", metrics.Summary)
}

// printCommitTrends prints commit frequency trends
func printCommitTrends(metrics *analyzer.TrendMetrics, detailed bool) {
	fmt.Println("COMMIT TRENDS")
	fmt.Println(strings.Repeat("-", 40))

	trendStyle := getTrendStyle(metrics.CommitTrend)
	sparkline := GenerateSparkline(metrics.CommitTrendValues, 20)

	fmt.Printf("   Trend: %s\n", trendStyle.Render(string(metrics.CommitTrend)))
	fmt.Printf("   Change Rate: %.1f%%\n", metrics.CommitChangeRate)
	fmt.Printf("   Avg/Month: %.1f commits\n", metrics.AvgCommitsPerMonth)
	fmt.Printf("   Sparkline: %s\n\n", sparkline)
	if detailed {
		fmt.Println("Monthly Breakdown:")

		for _, month := range metrics.MonthlyData {
			fmt.Printf("- %s : %d commits\n", month.Month.Format("Jan 06"), month.Commits)
		}

		fmt.Println()
	}
}

// printContributorTrends prints contributor growth trends
func printContributorTrends(metrics *analyzer.TrendMetrics) {
	fmt.Println("CONTRIBUTOR TRENDS")
	fmt.Println(strings.Repeat("-", 40))

	trendStyle := getTrendStyle(metrics.ContributorTrend)
	sparkline := GenerateSparkline(metrics.ContributorTrendValues, 20)

	fmt.Printf("   Trend: %s\n", trendStyle.Render(string(metrics.ContributorTrend)))
	fmt.Printf("   Current Contributors: %d\n", metrics.CurrentContributors)
	fmt.Printf("   Change Rate: %.1f%%\n", metrics.ContributorChangeRate)

	if metrics.NewContributors > 0 {
		fmt.Printf("   New Contributors: +%d\n", metrics.NewContributors)
	}
	if metrics.LostContributors > 0 {
		fmt.Printf("   Lost Contributors: -%d\n", metrics.LostContributors)
	}

	fmt.Printf("   Sparkline: %s\n\n", sparkline)
}

// printIssueTrends prints issue resolution trends
func printIssueTrends(metrics *analyzer.TrendMetrics) {
	fmt.Println("ISSUE RESOLUTION TRENDS")
	fmt.Println(strings.Repeat("-", 40))

	trendStyle := getTrendStyle(metrics.IssueResolutionTrend)

	fmt.Printf("   Trend: %s\n", trendStyle.Render(string(metrics.IssueResolutionTrend)))
	fmt.Printf("   Resolution Rate: %.1f%%\n", metrics.ResolutionRate)
	fmt.Printf("   Avg Resolution Time: %s\n\n", formatTrendDuration(metrics.AvgResolutionTime))
}

// printPRTrends prints PR merge trends
func printPRTrends(metrics *analyzer.TrendMetrics) {
	fmt.Println("PULL REQUEST TRENDS")
	fmt.Println(strings.Repeat("-", 40))

	trendStyle := getTrendStyle(metrics.PRTrend)
	sparkline := GenerateSparkline(metrics.PRTrendValues, 20)

	fmt.Printf("   Trend: %s\n", trendStyle.Render(string(metrics.PRTrend)))
	fmt.Printf("   Merge Rate: %.1f%%\n", metrics.PRMergeRate)
	fmt.Printf("   Sparkline: %s\n\n", sparkline)
}

// printHealthPrediction prints health score prediction
func printHealthPrediction(metrics *analyzer.TrendMetrics) {
	fmt.Println("HEALTH SCORE PREDICTION")
	fmt.Println(strings.Repeat("-", 40))

	currentStyle := getHealthStyle(metrics.CurrentHealthScore)
	predictedStyle := getHealthStyle(metrics.PredictedHealthScore)
	trendStyle := getTrendStyle(metrics.HealthScoreTrend)

	fmt.Printf("   Current Score: %s\n", currentStyle.Render(fmt.Sprintf("%d/100", metrics.CurrentHealthScore)))
	fmt.Printf("   Predicted Score: %s\n", predictedStyle.Render(fmt.Sprintf("%d/100", metrics.PredictedHealthScore)))
	fmt.Printf("   Trend: %s\n\n", trendStyle.Render(string(metrics.HealthScoreTrend)))
}

// printMonthlySparkline prints monthly data as a sparkline
func printMonthlySparkline(metrics *analyzer.TrendMetrics) {
	fmt.Println("MONTHLY ACTIVITY")
	fmt.Println(strings.Repeat("-", 40))

	// Create a simple bar chart for monthly commits
	maxCommits := 0
	for _, m := range metrics.MonthlyData {
		if m.Commits > maxCommits {
			maxCommits = m.Commits
		}
	}

	if maxCommits > 0 {
		for _, m := range metrics.MonthlyData {
			monthStr := m.Month.Format("Jan 06")
			barLen := int(float64(m.Commits) / float64(maxCommits) * 30)
			bar := strings.Repeat("█", barLen)
			fmt.Printf("   %s | %s %d\n", monthStr, bar, m.Commits)
		}
	}
	fmt.Println()
}

// PrintTrendSummary prints a brief trend summary for dashboards
func PrintTrendSummary(metrics *analyzer.TrendMetrics) {
	trendStyle := getTrendStyle(metrics.OverallTrend)

	fmt.Printf("\nTrend: %s | ", trendStyle.Render(string(metrics.OverallTrend)))
	fmt.Printf("Commits: %s | ", getTrendStyle(metrics.CommitTrend).Render(abbrevTrend(string(metrics.CommitTrend))))
	fmt.Printf("Contributors: %s | ", getTrendStyle(metrics.ContributorTrend).Render(abbrevTrend(string(metrics.ContributorTrend))))
	fmt.Printf("Health: %d/100 -> %d/100\n", metrics.CurrentHealthScore, metrics.PredictedHealthScore)
}

// PrintTrendIndicator prints a single trend indicator with optional value
func PrintTrendIndicator(trend analyzer.TrendIndicator, value string) {
	trendStyle := getTrendStyle(trend)
	if value != "" {
		fmt.Printf("%s %s", trendStyle.Render(string(trend)), value)
	} else {
		fmt.Print(trendStyle.Render(string(trend)))
	}
}

// GenerateSparkline creates an ASCII sparkline from a slice of integers
func GenerateSparkline(data []int, width int) string {
	if len(data) == 0 {
		return ""
	}

	// Normalize data to fit in width
	minVal := math.MaxInt
	maxVal := math.MinInt
	for _, v := range data {
		if v < minVal {
			minVal = v
		}
		if v > maxVal {
			maxVal = v
		}
	}

	rangeVal := maxVal - minVal
	if rangeVal == 0 {
		rangeVal = 1
	}

	var result strings.Builder
	for i, v := range data {
		normalized := float64(v-minVal) / float64(rangeVal)
		barLen := int(normalized * float64(width-1))
		if barLen < 1 {
			barLen = 1
		}

		// Use different characters for up/down/stable based on context
		char := "▄"
		if i > 0 {
			if data[i] > data[i-1] {
				char = "▄"
			} else if data[i] < data[i-1] {
				char = "▀"
			} else {
				char = "▬"
			}
		}

		result.WriteString(strings.Repeat(char, barLen))
	}

	return result.String()
}

// getTrendStyle returns the appropriate style for a trend indicator
func getTrendStyle(trend analyzer.TrendIndicator) lipgloss.Style {
	switch trend {
	case analyzer.TrendImproving:
		return improvingStyle
	case analyzer.TrendDeclining:
		return decliningStyle
	default:
		return stableStyle
	}
}

// getHealthStyle returns the appropriate style for a health score
func getHealthStyle(score int) lipgloss.Style {
	if score >= 80 {
		return improvingStyle
	} else if score >= 60 {
		return stableStyle
	} else {
		return decliningStyle
	}
}

// abbrevTrend returns an abbreviated trend indicator
func abbrevTrend(trend string) string {
	switch trend {
	case "↗️ Improving":
		return "↑"
	case "↘️ Declining":
		return "↓"
	default:
		return "→"
	}
}

// formatTrendDuration formats a duration to a human-readable string
func formatTrendDuration(d time.Duration) string {
	days := int(d.Hours() / 24)
	hours := int(d.Hours())
	minutes := int(d.Minutes())

	switch {
	case days > 30:
		return fmt.Sprintf("~%d months", days/30)
	case days > 0:
		return fmt.Sprintf("%d days", days)
	case hours > 0:
		return fmt.Sprintf("%d hours", hours)
	case minutes > 0:
		return fmt.Sprintf("%d minutes", minutes)
	default:
		return "< 1 minute"
	}
}

// PrintTrendCompact prints trend data in compact JSON format
func PrintTrendCompact(metrics *analyzer.TrendMetrics) {
	fmt.Printf(`{"owner":"%s","repo":"%s","period":%d,"overall_trend":"%s","commit_trend":"%s","contributor_trend":"%s","current_health":%d,"predicted_health":%d}`,
		metrics.Owner,
		metrics.Repo,
		metrics.AnalysisPeriod,
		metrics.OverallTrend,
		metrics.CommitTrend,
		metrics.ContributorTrend,
		metrics.CurrentHealthScore,
		metrics.PredictedHealthScore)
}
