package analyzer

import "testing"

func TestGetSeverity(t *testing.T) {
	tests := []struct {
		name     string
		vuln     osvVuln
		expected string
	}{
		{
			name: "CVSS V3 Critical",
			vuln: osvVuln{
				Severity: []struct {
					Type  string `json:"type"`
					Score string `json:"score"`
				}{
					{Type: "CVSS_V3", Score: "9.5"},
				},
			},
			expected: "CRITICAL",
		},
		{
			name: "CVSS V3 High",
			vuln: osvVuln{
				Severity: []struct {
					Type  string `json:"type"`
					Score string `json:"score"`
				}{
					{Type: "CVSS_V3", Score: "8.5"},
				},
			},
			expected: "HIGH",
		},
		{
			name: "CVSS V2 Fallback High",
			vuln: osvVuln{
				Severity: []struct {
					Type  string `json:"type"`
					Score string `json:"score"`
				}{
					{Type: "CVSS_V2", Score: "7.5"},
				},
			},
			expected: "HIGH",
		},
		{
			name: "CVSS V2 Fallback Medium",
			vuln: osvVuln{
				Severity: []struct {
					Type  string `json:"type"`
					Score string `json:"score"`
				}{
					{Type: "CVSS_V2", Score: "5.0"},
				},
			},
			expected: "MEDIUM",
		},
		{
			name: "CVSS V2 Fallback Low",
			vuln: osvVuln{
				Severity: []struct {
					Type  string `json:"type"`
					Score string `json:"score"`
				}{
					{Type: "CVSS_V2", Score: "2.0"},
				},
			},
			expected: "LOW",
		},
		{
			name: "CVSS V3 Precedence over V2",
			vuln: osvVuln{
				Severity: []struct {
					Type  string `json:"type"`
					Score string `json:"score"`
				}{
					{Type: "CVSS_V2", Score: "9.5"},
					{Type: "CVSS_V3", Score: "5.0"},
				},
			},
			expected: "MEDIUM",
		},
		{
			name: "No Score Unknown",
			vuln: osvVuln{
				Severity: []struct {
					Type  string `json:"type"`
					Score string `json:"score"`
				}{},
			},
			expected: "UNKNOWN",
		},
		{
			name: "Text-only record Unknown",
			vuln: osvVuln{
				Severity: []struct {
					Type  string `json:"type"`
					Score string `json:"score"`
				}{
					{Type: "TEXT", Score: "HIGH"},
				},
			},
			expected: "UNKNOWN",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getSeverity(tt.vuln)
			if result != tt.expected {
				t.Errorf("getSeverity() = %v, want %v", result, tt.expected)
			}
		})
	}
}
