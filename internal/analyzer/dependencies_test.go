package analyzer

import (
	"reflect"
	"testing"
)

// TestParseRequirementsTxt tests the parsing of Python requirements.txt files
// with various version specifier formats including spaced comparators.
func TestParseRequirementsTxt(t *testing.T) {
	tests := []struct {
		name     string
		content  string
		expected []Dependency
	}{
		{
			name:    "Simple package",
			content: "flask==2.0.0",
			expected: []Dependency{
				{Name: "flask", Version: "==2.0.0", Type: "production"},
			},
		},
		{
			name:    "Package with dot",
			content: "ruamel.yaml>=0.17.0",
			expected: []Dependency{
				{Name: "ruamel.yaml", Version: ">=0.17.0", Type: "production"},
			},
		},
		{
			name:    "Package with extras",
			content: "requests[security]==2.28.0",
			expected: []Dependency{
				{Name: "requests", Version: "==2.28.0", Type: "production"},
			},
		},
		{
			name:    "Package with environment markers",
			content: "dataclasses; python_version < \"3.7\"",
			expected: []Dependency{
				{Name: "dataclasses", Version: "*", Type: "production"},
			},
		},
		{
			name:    "Package with version and marker",
			content: "requests==2.28.0 ; python_version > \"3.6\"",
			expected: []Dependency{
				{Name: "requests", Version: "==2.28.0", Type: "production"},
			},
		},
		{
			name:    "Package with space after comparator (>= )",
			content: "numpy >= 1.0",
			expected: []Dependency{
				{Name: "numpy", Version: ">= 1.0", Type: "production"},
			},
		},
		{
			name:    "Package with space after comparator (== )",
			content: "flask == 2.0.0",
			expected: []Dependency{
				{Name: "flask", Version: "== 2.0.0", Type: "production"},
			},
		},
		{
			name:    "Package with space after comparator (!= )",
			content: "requests != 2.0.0",
			expected: []Dependency{
				{Name: "requests", Version: "!= 2.0.0", Type: "production"},
			},
		},
		{
			name:    "Package with space after comparator (<= )",
			content: "numpy <= 1.20",
			expected: []Dependency{
				{Name: "numpy", Version: "<= 1.20", Type: "production"},
			},
		},
		{
			name:    "Package with space after compatible release comparator (~= )",
			content: "requests ~= 2.28.0",
			expected: []Dependency{
				{Name: "requests", Version: "~= 2.28.0", Type: "production"},
			},
		},
		{
			name:    "Package with space after arbitrary equality comparator (=== )",
			content: "numpy === 1.20.0",
			expected: []Dependency{
				{Name: "numpy", Version: "=== 1.20.0", Type: "production"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, _ := parseRequirementsTxt([]byte(tt.content))
			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("parseRequirementsTxt() = %v, want %v", got, tt.expected)
			}
		})
	}
}

// TestParseCargoToml tests the parsing of Rust Cargo.toml files
// with both simple and inline table dependency formats.
func TestParseCargoToml(t *testing.T) {
	content := `
[dependencies]
serde = "1.0"
tokio = { version = "1.15", features = ["full"] }
rand = { version = '0.8.5' }
other-pkg = { path = "../other" }

[dev-dependencies]
tokio-test = "0.4"
tempfile = { version = "3.3" }
`

	expected := []Dependency{
		{Name: "serde", Version: "1.0", Type: "production"},
		{Name: "tokio", Version: "1.15", Type: "production"},
		{Name: "rand", Version: "0.8.5", Type: "production"},
		{Name: "other-pkg", Version: "*", Type: "production"},
		{Name: "tokio-test", Version: "0.4", Type: "dev"},
		{Name: "tempfile", Version: "3.3", Type: "dev"},
	}

	got, _ := parseCargoToml([]byte(content))
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("parseCargoToml() =\n%v\nwant:\n%v", got, expected)
	}
}

