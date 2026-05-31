package output

import (
	"fmt"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/predictive"
	"github.com/fatih/color"
)

// RenderForecastDashboard renders the AI Trend Forecasting predictions
func RenderForecastDashboard(repo string, health *predictive.ForecastResult, issues *predictive.ForecastResult, risks []predictive.ContributorRiskForecast) {
	fmt.Printf("\n")
	color.New(color.FgHiCyan, color.Bold).Printf("📈 Repo-lyzer Trend Forecasting Engine: %s\n", repo)
	fmt.Printf("================================================================================\n\n")

	// 1. Repository Growth/Decline (Commit Velocity)
	if health != nil {
		color.New(color.FgHiBlue, color.Bold).Println("📊 Repository Trajectory (Commit Velocity)")
		fmt.Printf("  Trend: %s\n", colorizeTrend(health.Trend))
		fmt.Printf("  Risk : %s\n", colorizeRisk(health.RiskLevel))
		fmt.Printf("  Recs : %v\n\n", health.Recommendations)
	}

	// 2. Issue Accumulation
	if issues != nil {
		color.New(color.FgHiYellow, color.Bold).Println("🐛 Issue Accumulation Forecast (Next 30 Days)")
		fmt.Printf("  Trend: %s\n", colorizeTrend(issues.Trend))

		fmt.Printf("\n  %-10s %-15s\n", "Days Out", "Predicted Open")
		fmt.Printf("  %s\n", strings.Repeat("-", 26))
		for i, p := range issues.Predictions {
			fmt.Printf("  %-10d %-15.1f\n", i+1, p.Value)
		}
		fmt.Printf("\n")
	}

	// 3. Maintainer Burnout Risk
	color.New(color.FgHiMagenta, color.Bold).Println("🔥 Maintainer Burnout Risk")
	if len(risks) == 0 {
		fmt.Println("  No core maintainers found with high burnout risk.")
	} else {
		fmt.Printf("\n  %-20s %-15s %-15s %-15s\n", "Maintainer", "Burnout Risk", "Attrition Risk", "Trajectory")
		fmt.Printf("  %s\n", strings.Repeat("-", 70))
		for _, r := range risks {
			burnoutStr := fmt.Sprintf("%.2f", r.BurnoutRisk)
			attritionStr := fmt.Sprintf("%.2f", r.AttritionRisk)

			if r.BurnoutRisk > 0.7 {
				burnoutStr = color.RedString(burnoutStr)
			} else if r.BurnoutRisk > 0.4 {
				burnoutStr = color.YellowString(burnoutStr)
			} else {
				burnoutStr = color.GreenString(burnoutStr)
			}

			trajStr := r.Trajectory
			if trajStr == "worsening" {
				trajStr = color.RedString(trajStr)
			} else {
				trajStr = color.GreenString(trajStr)
			}

			fmt.Printf("  %-20s %-15s %-15s %-15s\n", r.ContributorID, burnoutStr, attritionStr, trajStr)
		}
	}
	fmt.Printf("\n")
}

func colorizeTrend(trend string) string {
	switch trend {
	case "improving":
		return color.GreenString("improving ↑")
	case "stable":
		return color.YellowString("stable →")
	case "degrading":
		return color.RedString("degrading ↓")
	}
	return trend
}

func colorizeRisk(risk string) string {
	switch risk {
	case "low":
		return color.GreenString(risk)
	case "medium":
		return color.YellowString(risk)
	case "high":
		return color.RedString(risk)
	}
	return risk
}
