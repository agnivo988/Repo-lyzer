package cmd

import "testing"

func TestValidateCompareSavePath(t *testing.T) {
	tests := []struct {
		name  string
		path  string
		valid bool
	}{
		{name: "filename", path: "report.html", valid: true},
		{name: "nested relative path", path: "reports/compare.json", valid: true},
		{name: "parent traversal", path: "../report.html", valid: false},
		{name: "nested parent traversal", path: "reports/../../report.html", valid: false},
		{name: "absolute path", path: "/tmp/report.html", valid: false},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := validateCompareSavePath(test.path) == nil
			if got != test.valid {
				t.Fatalf("validateCompareSavePath(%q) valid = %v, want %v", test.path, got, test.valid)
			}
		})
	}
}
