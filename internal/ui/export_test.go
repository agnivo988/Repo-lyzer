package ui

import (
	"strings"
	"testing"
)

func TestValidateExportFormat(t *testing.T) {
	tests := []struct {
		format   string
		expected bool
	}{
		{"json", true},
		{"JSON", true}, // case insensitive
		{"markdown", true},
		{"Markdown", true},
		{"csv", true},
		{"CSV", true},
		{"html", true},
		{"HTML", true},
		{"pdf", true},
		{"PDF", true},
		{"xml", false},
		{"txt", false},
		{"", false},
		{"invalid", false},
	}

	for _, test := range tests {
		t.Run(test.format, func(t *testing.T) {
			err := ValidateExportFormat(test.format)
			if test.expected && err != nil {
				t.Errorf("ValidateExportFormat(%s) = error %v, expected no error", test.format, err)
			}
			if !test.expected && err == nil {
				t.Errorf("ValidateExportFormat(%s) = no error, expected error", test.format)
			}
			if !test.expected && err != nil {
				// Check that error message contains the invalid format
				if !strings.Contains(err.Error(), test.format) {
					t.Errorf("Error message should contain the invalid format '%s'", test.format)
				}
				// Check that error message lists supported formats
				if !strings.Contains(err.Error(), "json") || !strings.Contains(err.Error(), "markdown") {
					t.Errorf("Error message should list supported formats")
				}
			}
		})
	}
}

func TestExportAnalysis(t *testing.T) {
	// Test that validation works for the ExportAnalysis function
	validFormats := []string{"json", "markdown", "csv", "html", "pdf"}
	
	for _, format := range validFormats {
		t.Run("valid_"+format, func(t *testing.T) {
			err := ValidateExportFormat(format)
			if err != nil {
				t.Errorf("ValidateExportFormat should accept %s", format)
			}
		})
	}

	invalidFormats := []string{"xml", "txt", "invalid"}
	for _, format := range invalidFormats {
		t.Run("invalid_"+format, func(t *testing.T) {
			err := ValidateExportFormat(format)
			if err == nil {
				t.Errorf("ValidateExportFormat should reject %s", format)
			}
		})
	}
}