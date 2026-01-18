package analyzer

type RiskAlertsResult struct {
	Alerts []string
}

func AnalyzeRiskAlerts(
	busFactor int,
	healthScore int,
	commitsLast90Days int,
	hasCriticalVulns bool,
) *RiskAlertsResult {

	var alerts []string

	if busFactor <= 1 {
		alerts = append(alerts, "Low bus factor (single contributor dependency)")
	}

	if commitsLast90Days == 0 {
		alerts = append(alerts, "No commit activity in the last 90 days")
	}

	if healthScore < 40 {
		alerts = append(alerts, "Very low repository health score")
	}

	if hasCriticalVulns {
		alerts = append(alerts, "Critical dependency vulnerabilities detected")
	}

	if len(alerts) == 0 {
		return nil
	}

	return &RiskAlertsResult{Alerts: alerts}
}
