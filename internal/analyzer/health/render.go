package health

import (
	"fmt"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer/core"
	"github.com/fatih/color"
)

// TerminalRenderer provides stylized terminal output for health metrics
type TerminalRenderer struct {
	CriticalColor *color.Color
	WarningColor  *color.Color
	HealthyColor   *color.Color
	ExcellentColor *color.Color
}

// NewTerminalRenderer initializes a customized terminal renderer
func NewTerminalRenderer() *TerminalRenderer {
	return &TerminalRenderer{
		CriticalColor: color.New(color.FgRed, color.Bold),
		WarningColor:  color.New(color.FgYellow, color.Bold),
		HealthyColor:   color.New(color.FgGreen, color.Bold),
		ExcellentColor: color.New(color.FgHiCyan, color.Bold),
	}
}

// RenderCategory formats a category string with its corresponding color
func (t *TerminalRenderer) RenderCategory(category core.ScoreCategory) string {
	switch category {
	case core.Critical:
		return t.CriticalColor.Sprint("CRITICAL")
	case core.Warning:
		return t.WarningColor.Sprint("WARNING")
	case core.Healthy:
		return t.HealthyColor.Sprint("HEALTHY")
	case core.Excellent:
		return t.ExcellentColor.Sprint("EXCELLENT")
	default:
		return string(category)
	}
}

// RenderScoreCard prints a beautiful score card for a specific health module
func (t *TerminalRenderer) RenderScoreCard(title string, score float64, category core.ScoreCategory) {
	fmt.Printf("\n==== %s ====\n", title)

	// Format the score
	scoreStr := fmt.Sprintf("%.1f / 100", score)

	// Create a visual pulse bar
	barLength := 20
	filled := int((score / 100.0) * float64(barLength))
	if filled > barLength {
		filled = barLength
	}
	empty := barLength - filled

	pulseBar := fmt.Sprintf("[%s%s]", strings.Repeat("=", filled), strings.Repeat("-", empty))

	// Colorize the pulse bar
	var coloredPulse string
	switch category {
	case core.Critical:
		coloredPulse = t.CriticalColor.Sprint(pulseBar)
	case core.Warning:
		coloredPulse = t.WarningColor.Sprint(pulseBar)
	case core.Healthy:
		coloredPulse = t.HealthyColor.Sprint(pulseBar)
	case core.Excellent:
		coloredPulse = t.ExcellentColor.Sprint(pulseBar)
	default:
		coloredPulse = pulseBar
	}

	fmt.Printf("Health: %s  %s  (Status: %s)\n", scoreStr, coloredPulse, t.RenderCategory(category))
	fmt.Println(strings.Repeat("=", len(title)+10))
}
