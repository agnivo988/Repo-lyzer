package risk

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

// RiskRenderer provides specialized terminal output for risk metrics
type RiskRenderer struct {
	CriticalColor *color.Color
	WarningColor  *color.Color
	StableColor   *color.Color
}

// NewRiskRenderer initializes a customized terminal renderer for risk
func NewRiskRenderer() *RiskRenderer {
	return &RiskRenderer{
		CriticalColor: color.New(color.FgHiWhite, color.BgRed, color.Bold),
		WarningColor:  color.New(color.FgHiBlack, color.BgYellow, color.Bold),
		StableColor:   color.New(color.FgGreen, color.Bold),
	}
}

// RenderCategory formats a risk category string with its corresponding color
func (r *RiskRenderer) RenderCategory(category RiskCategory) string {
	switch category {
	case RiskCritical:
		return r.CriticalColor.Sprint(" [ CRITICAL RISK ] ")
	case RiskWarning:
		return r.WarningColor.Sprint(" [ WARNING RISK ] ")
	case RiskStable:
		return r.StableColor.Sprint(" [ STABLE ] ")
	default:
		return string(category)
	}
}

// RenderRiskCard prints an enterprise-grade score card for risk
// Since higher score = higher risk, the pulse bar grows with risk!
func (r *RiskRenderer) RenderRiskCard(title string, score float64, category RiskCategory) {
	fmt.Printf("\n==== %s ====\n", title)

	scoreStr := fmt.Sprintf("%.1f / 100", score)

	// Pulse bar: higher score = more danger filled
	barLength := 20
	filled := int((score / 100.0) * float64(barLength))
	if filled > barLength {
		filled = barLength
	}
	empty := barLength - filled

	// Use 'X' for danger, '-' for safe
	pulseBar := fmt.Sprintf("[%s%s]", strings.Repeat("X", filled), strings.Repeat("-", empty))

	var coloredPulse string
	switch category {
	case RiskCritical:
		coloredPulse = color.New(color.FgRed, color.Bold).Sprint(pulseBar)
	case RiskWarning:
		coloredPulse = color.New(color.FgYellow).Sprint(pulseBar)
	case RiskStable:
		coloredPulse = color.New(color.FgGreen).Sprint(pulseBar)
	}

	fmt.Printf("Risk Level: %s  %s  %s\n", scoreStr, coloredPulse, r.RenderCategory(category))
	fmt.Println(strings.Repeat("=", len(title)+10))
}
