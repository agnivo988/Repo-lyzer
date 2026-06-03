package analyzer

import (
	"reflect"
	"testing"
)

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

