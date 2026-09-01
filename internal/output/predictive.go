package output

import (
	"fmt"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/predictive"
	"github.com/fatih/color"
)

// RenderPredictiveDashboard renders the AI Contribution Engine predictions
func RenderPredictiveDashboard(repo string, bestIssues []predictive.IssueDifficulty, prProbs []predictive.PRProbability, fit predictive.ContributorFit) {
	fmt.Printf("\n")
	color.New(color.FgHiCyan, color.Bold).Printf("🤖 Repo-lyzer Predictive Intelligence Engine: %s\n", repo)
	fmt.Printf("================================================================================\n\n")

	// 1. Best First Issues
	color.New(color.FgHiYellow, color.Bold).Println("🌟 Best First Issues (Difficulty Estimation)")
	if len(bestIssues) == 0 {
		fmt.Println("  No open issues found.")
	} else {
		fmt.Printf("%-10s %-40s %-15s %-6s %s\n", "Issue #", "Title", "Level", "Score", "Heuristics")
		fmt.Printf("%s\n", strings.Repeat("-", 100))

		for _, issue := range bestIssues {
			levelStr := issue.Level
			switch issue.Level {
			case "Trivial", "Easy":
				levelStr = color.GreenString("%-15s", issue.Level)
			case "Medium":
				levelStr = color.YellowString("%-15s", issue.Level)
			case "Hard", "Expert":
				levelStr = color.RedString("%-15s", issue.Level)
			}

			// ColorString padding can be tricky with ANSI codes, so format string first
			title := truncateString(issue.Title, 40)
			reason := truncateString(fmt.Sprintf("%v", issue.Reasons), 40)

			fmt.Printf("#%-9d %-40s %s %-6.1f %s\n",
				issue.IssueNumber,
				title,
				levelStr,
				issue.Score,
				reason,
			)
		}
	}
	fmt.Printf("\n")

	// 2. PR Merge Probability
	color.New(color.FgHiMagenta, color.Bold).Println("🔀 Pull Request Merge Probability")
	if len(prProbs) == 0 {
		fmt.Println("  No open pull requests found.")
	} else {
		fmt.Printf("%-10s %-40s %-15s %-6s %s\n", "PR #", "Title", "Probability", "Score", "Heuristics")
		fmt.Printf("%s\n", strings.Repeat("-", 100))
		for _, pr := range prProbs {
			levelStr := pr.Level
			switch pr.Level {
			case "High":
				levelStr = color.GreenString("%-15s", "High")
			case "Medium":
				levelStr = color.YellowString("%-15s", "Medium")
			case "Low", "Critical Risk":
				levelStr = color.RedString("%-15s", pr.Level)
			}

			title := truncateString(pr.Title, 40)
			reason := truncateString(fmt.Sprintf("%v", pr.Reasons), 40)

			fmt.Printf("#%-9d %-40s %s %-6.1f %s\n",
				pr.PRNumber,
				title,
				levelStr,
				pr.Score,
				reason,
			)
		}
	}
	fmt.Printf("\n")

	// 3. Contributor Fit
	if fit.Contributor != "" {
		color.New(color.FgHiGreen, color.Bold).Printf("👤 Contributor Fit Analysis: %s\n", fit.Contributor)
		fmt.Printf("  Target Issue: #%d\n", fit.IssueNumber)
		fmt.Printf("  Match Level : %s (%.1f/100)\n", colorizeFit(fit.Level), fit.Score)
		fmt.Printf("  Analysis    : %v\n", fit.Reasons)
		fmt.Printf("\n")
	}
}

func truncateString(s string, max int) string {
	if len(s) > max {
		return s[:max-3] + "..."
	}
	return s
}

func colorizeFit(level string) string {
	switch level {
	case "Excellent", "Good":
		return color.GreenString(level)
	case "Moderate":
		return color.YellowString(level)
	case "Poor":
		return color.RedString(level)
	}
	return level
}
