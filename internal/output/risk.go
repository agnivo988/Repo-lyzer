package output

import (
	"fmt"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/charmbracelet/lipgloss"
)

// PrintRiskReport renders the repository risk analysis to the terminal.
// It follows the same pattern as PrintHealth and PrintContributionScore:
// a styled header, a score line, and a bullet list of warnings.
func PrintRiskReport(r analyzer.RiskReport) {
	// Choose header colour based on risk level, mirroring PrintHealth's
	// red/amber/green convention that users already recognise.
	var levelColor string
	var levelIcon string
	switch r.Level {
	case "High":
		levelColor = "#FF5F5F"
		levelIcon = "🔴"
	case "Moderate":
		levelColor = "#FFB000"
		levelIcon = "🟡"
	default: // Low
		levelColor = "#00FF87"
		levelIcon = "🟢"
	}

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(levelColor))

	fmt.Println(headerStyle.Render("\n⚠  Repository Risk Report"))

	fmt.Printf("Risk Level: %s %s (%d/100)\n", levelIcon, r.Level, r.Score)

	if len(r.Warnings) == 0 {
		fmt.Println(SuccessStyle.Render("  ✅ No risk signals detected"))
		fmt.Println()
		return
	}

	fmt.Println(WarningStyle.Render("\nWarnings:"))
	for _, w := range r.Warnings {
		fmt.Printf("  • %s\n", w)
	}
	fmt.Println()
}
