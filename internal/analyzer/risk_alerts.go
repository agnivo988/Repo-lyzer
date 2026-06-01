package analyzer

import (
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/contribution"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// RiskSignal represents a single triggered risk rule and its point contribution.
type RiskSignal struct {
	Label  string
	Points int
}

// RiskReport is the result of a full repository risk analysis.
// Score ranges from 0 (no risk) to 100 (maximum risk).
// Level classifies the score into a human-readable tier.
// Signals lists every rule that fired, for audit/debug purposes.
// Warnings contains the user-facing maintenance alerts.
type RiskReport struct {
	Score    int
	Level    string
	Signals  []RiskSignal
	Warnings []string
}

// riskLevel converts a numeric risk score to a display label.
func riskLevel(score int) string {
	switch {
	case score >= 61:
		return "High"
	case score >= 31:
		return "Moderate"
	default:
		return "Low"
	}
}

// AnalyzeRisk evaluates repository risk using lightweight, metadata-only signals.
// It reuses data already fetched by the analyze command — no extra API calls are made.
//
// Parameters:
//   - repo:         Repository metadata (stars, open issues, pushed date, etc.)
//   - contributors: Full contributor list for bus factor and diversity checks.
//   - commits:      Recent commits (callers should pass the 90-day window).
//   - fileTree:     Repository file tree, used to detect a missing README.
//
// Returns a RiskReport with a score (0–100), a risk level string, the triggered
// signals, and human-readable warning messages.
func AnalyzeRisk(
	repo *github.Repo,
	contributors []github.Contributor,
	commits []github.Commit,
	fileTree []github.TreeEntry,
) RiskReport {
	var signals []RiskSignal
	var warnings []string
	total := 0

	add := func(label string, points int, warning string) {
		total += points
		signals = append(signals, RiskSignal{Label: label, Points: points})
		warnings = append(warnings, warning)
	}

	// Rule 1 — Bus factor ≤ 1 (+25)
	// Reuses the existing BusFactor() function to stay consistent with the
	// rest of the analysis output (bus factor is already shown to users).
	busFactor, _ := BusFactor(contributors)
	if busFactor <= 1 {
		add(
			"Low bus factor",
			25,
			"Low maintainer diversity detected",
		)
	}

	// Rule 2 — Contributor count ≤ 2 (+15)
	if len(contributors) <= 2 {
		add(
			"Low contributor count",
			15,
			"Very few contributors — project may lack community support",
		)
	}

	// Rule 3 — No recent commits in the past 90 days (+20)
	// Callers pass a pre-fetched recent-commit slice; we simply check its length.
	if len(commits) == 0 {
		add(
			"No recent commits (90d)",
			20,
			"Low recent development activity",
		)
	}

	// Rule 4 — High unresolved issue backlog (+15)
	if repo.OpenIssues > 100 {
		add(
			"High open issue backlog",
			15,
			"High unresolved issue backlog",
		)
	}

	// Rule 5 — Missing README (+10)
	// Delegates to the contribution package helper that already handles README
	// detection — avoids duplicating filename-matching logic.
	if contribution.FindReadmePath(fileTree) == "" {
		add(
			"Missing README",
			10,
			"Repository has no README file",
		)
	}

	// Rule 6 — Repository stale: not pushed to in >180 days (+15)
	if !repo.PushedAt.IsZero() && time.Since(repo.PushedAt) > 180*24*time.Hour {
		add(
			"Repository stale (>180d)",
			15,
			"Repository appears stale — no pushes in over 6 months",
		)
	}

	// Cap score at 100 to keep the scale meaningful.
	if total > 100 {
		total = 100
	}

	return RiskReport{
		Score:    total,
		Level:    riskLevel(total),
		Signals:  signals,
		Warnings: warnings,
	}
}

// ── Backward-compatible TUI shims ────────────────────────────────────────────
// The interactive dashboard (internal/ui) was built against RiskAlertsResult
// and AnalyzeRiskAlerts. These shims keep the TUI compiling without touching
// ui/types.go, ui/app.go, or ui/dashboard.go.

// RiskAlertsResult is the legacy result type used by the TUI dashboard.
type RiskAlertsResult struct {
	Alerts []string
}

// AnalyzeRiskAlerts is the legacy entry point called by the TUI.
// Preserved as-is so the dashboard rendering code continues to work unchanged.
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
