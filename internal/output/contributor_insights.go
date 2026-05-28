package output

import (
    "fmt"

    "github.com/agnivo988/Repo-lyzer/internal/analyzer"
    "github.com/agnivo988/Repo-lyzer/internal/github"
    "github.com/charmbracelet/lipgloss"
)

// PrintContributorInsights prints a concise, beginner-friendly contributor summary
func PrintContributorInsights(contribs []github.Contributor) {
    insights := analyzer.AnalyzeContributors(contribs)

    header := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00E5FF"))
    fmt.Println()
    fmt.Println(header.Render("👥 Contributor Insights"))
    fmt.Printf("  Total contributors : %d\n", insights.TotalContributors)
    fmt.Printf("  Active contributors: %d\n", insights.ActiveContributors)
    if insights.TopContributor != nil {
        fmt.Printf("  Top contributor    : %s (%d commits, %.1f%%)\n", insights.TopContributor.Login, insights.TopContributor.Commits, insights.TopContributor.Percentage)
    }
    fmt.Printf("  Diversity score    : %.1f/100\n", insights.DiversityScore)
    fmt.Printf("  Concentration risk : %s\n", insights.ConcentrationRisk)
    // Simple recommendation
    if len(insights.Recommendations) > 0 {
        fmt.Printf("  Recommendation     : %s\n", insights.Recommendations[0])
    }
    fmt.Println()
}
