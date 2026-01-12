package analyzer

import (
	"testing"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestBusFactor(t *testing.T) {
	tests := []struct {
		name         string
		contributors []github.Contributor
		wantFactor   int
		wantRisk     string
	}{
		{
			name:         "no contributors",
			contributors: []github.Contributor{},
			wantFactor:   0,
			wantRisk:     "Unknown",
		},
		{
			name: "single contributor (high risk)",
			contributors: []github.Contributor{
				{Login: "dev1", Commits: 100},
			},
			wantFactor: 1,
			wantRisk:   "High Risk",
		},
		{
			name: "two contributors with uneven distribution",
			contributors: []github.Contributor{
				{Login: "dev1", Commits: 90},
				{Login: "dev2", Commits: 10},
			},
			wantFactor: 1,
			wantRisk:   "High Risk",
		},
		{
			name: "balanced team (medium risk)",
			contributors: []github.Contributor{
				{Login: "dev1", Commits: 50},
				{Login: "dev2", Commits: 50},
			},
			wantFactor: 2,
			wantRisk:   "Medium Risk",
		},
		{
			name: "large diverse team (low risk)",
			contributors: []github.Contributor{
				{Login: "dev1", Commits: 20},
				{Login: "dev2", Commits: 20},
				{Login: "dev3", Commits: 20},
				{Login: "dev4", Commits: 20},
				{Login: "dev5", Commits: 20},
			},
			wantFactor: 3,
			wantRisk:   "Low Risk",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			factor, risk := BusFactor(tt.contributors)
			
			if factor != tt.wantFactor {
				t.Errorf("BusFactor() factor = %d, want %d", factor, tt.wantFactor)
			}
			if risk != tt.wantRisk {
				t.Errorf("BusFactor() risk = %s, want %s", risk, tt.wantRisk)
			}
		})
	}
}

func TestBusFactor_RiskLevels(t *testing.T) {
	// Test risk levels based on top contributor ratio
	testCases := []struct {
		name         string
		contributors []github.Contributor
		wantRisk     string
	}{
		{
			name:         "empty returns Unknown",
			contributors: []github.Contributor{},
			wantRisk:     "Unknown",
		},
		{
			name: "ratio > 0.7 is High Risk",
			contributors: []github.Contributor{
				{Login: "dev1", Commits: 80},
				{Login: "dev2", Commits: 20},
			},
			wantRisk: "High Risk",
		},
		{
			name: "ratio 0.4-0.7 is Medium Risk",
			contributors: []github.Contributor{
				{Login: "dev1", Commits: 50},
				{Login: "dev2", Commits: 50},
			},
			wantRisk: "Medium Risk",
		},
		{
			name: "ratio < 0.4 is Low Risk",
			contributors: []github.Contributor{
				{Login: "dev1", Commits: 30},
				{Login: "dev2", Commits: 30},
				{Login: "dev3", Commits: 40},
			},
			wantRisk: "Low Risk",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, risk := BusFactor(tc.contributors)
			if risk != tc.wantRisk {
				t.Errorf("BusFactor() risk = %s, want %s", risk, tc.wantRisk)
			}
		})
	}
}
