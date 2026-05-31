package output

import (
	"fmt"

	"github.com/agnivo988/Repo-lyzer/internal/contribution"
	"github.com/charmbracelet/lipgloss"
)

// PrintContributionScore prints the contribution friendliness score to the terminal.
func PrintContributionScore(score contribution.ContributionScore) {
	color := "#ff5555" // Red
	if score.Score >= 8.0 {
		color = "#50fa7b" // Green
	} else if score.Score >= 5.0 {
		color = "#ffb86c" // Orange/Yellow
	}

	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(color))

	fmt.Println(style.Render(
		fmt.Sprintf("🤝 Contribution Score: %.1f/10 (%s)", score.Score, score.Level),
	))

	if len(score.Strengths) > 0 {
		fmt.Println("  ✅ Strengths:")
		for _, s := range score.Strengths {
			fmt.Printf("    • %s\n", s)
		}
	}
	if len(score.Weaknesses) > 0 {
		fmt.Println("  ❌ Areas for Improvement:")
		for _, w := range score.Weaknesses {
			fmt.Printf("    • %s\n", w)
		}
	}
	fmt.Println()
}
