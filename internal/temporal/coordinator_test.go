package temporal

import (
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/evolution"
	"github.com/agnivo988/Repo-lyzer/internal/predictive"
)

func TestFinalize_UsesComputedSummaryMetrics(t *testing.T) {
	coordinator := NewCoordinator("owner", "repo")
	coordinator.Timeline.AddSnapshot(NewSnapshot(time.Now(), nil))

	coordinator.RiskIndicators = []evolution.RiskIndicator{
		{
			Name:       "dependency drift",
			Severity:   "high",
			Trajectory: "worsening",
		},
	}
	coordinator.DriftIndicators = []evolution.DriftIndicator{
		{
			SubsystemID: "core",
			Severity:    "medium",
			Direction:   "increasing",
		},
	}
	coordinator.HealthForecast = &predictive.ForecastResult{
		Trend:     "improving",
		RiskLevel: "medium",
	}

	result := coordinator.Finalize()

	if result.HealthScore != 74 {
		t.Fatalf("HealthScore = %d, want 74", result.HealthScore)
	}
	if result.HealthTrend != "improving" {
		t.Fatalf("HealthTrend = %q, want improving", result.HealthTrend)
	}
	if result.OverallRiskLevel != "high" {
		t.Fatalf("OverallRiskLevel = %q, want high", result.OverallRiskLevel)
	}
}

func TestFinalize_WithNoForecastFallsBackToSignalCounts(t *testing.T) {
	coordinator := NewCoordinator("owner", "repo")
	coordinator.Timeline.AddSnapshot(NewSnapshot(time.Now(), nil))

	coordinator.RiskIndicators = []evolution.RiskIndicator{
		{
			Name:       "activity drop",
			Severity:   "medium",
			Trajectory: "improving",
		},
	}
	coordinator.DriftIndicators = []evolution.DriftIndicator{
		{
			SubsystemID: "ui",
			Severity:    "low",
			Direction:   "decreasing",
		},
	}

	result := coordinator.Finalize()

	if result.HealthTrend != "improving" {
		t.Fatalf("HealthTrend = %q, want improving", result.HealthTrend)
	}
	if result.HealthScore != 86 {
		t.Fatalf("HealthScore = %d, want 86", result.HealthScore)
	}
	if result.OverallRiskLevel != "medium" {
		t.Fatalf("OverallRiskLevel = %q, want medium", result.OverallRiskLevel)
	}
}

func TestFinalize_CollectsCriticalIssues(t *testing.T) {
	coordinator := NewCoordinator("owner", "repo")
	coordinator.Timeline.AddSnapshot(NewSnapshot(time.Now(), nil))

	coordinator.RiskIndicators = []evolution.RiskIndicator{
		{
			Name:     "critical issue",
			Severity: "critical",
		},
		{
			Name:     "low issue",
			Severity: "low",
		},
	}

	result := coordinator.Finalize()

	if len(result.CriticalIssues) != 1 {
		t.Fatalf("CriticalIssues len = %d, want 1", len(result.CriticalIssues))
	}
	if result.CriticalIssues[0] != "critical issue" {
		t.Fatalf("CriticalIssues[0] = %q, want %q", result.CriticalIssues[0], "critical issue")
	}
}
