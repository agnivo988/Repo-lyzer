package analyzer

import (
	"testing"
)

func TestBuildRecruiterSummary_Metrics(t *testing.T) {
	tests := []struct {
		name               string
		openIssues         int
		prMergeRate        float64
		expectedIssueHealth string
		expectedPRHealth    string
	}{
		{
			name:                "Healthy issues and active PRs",
			openIssues:          15,
			prMergeRate:         75.0,
			expectedIssueHealth: "Healthy",
			expectedPRHealth:    "Active",
		},
		{
			name:                "Backlogged issues and slow PRs",
			openIssues:          30,
			prMergeRate:         45.0,
			expectedIssueHealth: "Backlogged",
			expectedPRHealth:    "Slow",
		},
		{
			name:                "Critical issues and stalled PRs",
			openIssues:          55,
			prMergeRate:         15.0,
			expectedIssueHealth: "Critical",
			expectedPRHealth:    "Stalled",
		},
		{
			name:                "Boundary test: 20 issues",
			openIssues:          20,
			prMergeRate:         60.0,
			expectedIssueHealth: "Backlogged",
			expectedPRHealth:    "Slow",
		},
		{
			name:                "Boundary test: 50 issues",
			openIssues:          50,
			prMergeRate:         30.0,
			expectedIssueHealth: "Critical",
			expectedPRHealth:    "Slow",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			summary := BuildRecruiterSummary(
				"test/repo",
				100, 50, // stars, forks
				200, 10, // commits, contributors
				80, "High", // maturity
				5, "Low", // bus factor
				tt.openIssues,
				tt.prMergeRate,
			)

			if summary.IssueHealth != tt.expectedIssueHealth {
				t.Errorf("IssueHealth = %v, want %v", summary.IssueHealth, tt.expectedIssueHealth)
			}
			if summary.PRHealth != tt.expectedPRHealth {
				t.Errorf("PRHealth = %v, want %v", summary.PRHealth, tt.expectedPRHealth)
			}
		})
	}
}
