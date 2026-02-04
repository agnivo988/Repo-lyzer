package output

import (
	"io/ioutil"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// EnhancedPDFConfig contains configuration for enhanced PDF reports
type EnhancedPDFConfig struct {
	Company  CompanyConfig  `yaml:"company"`
	Branding BrandingConfig `yaml:"branding"`
	Report   ReportConfig   `yaml:"report"`
	Sections SectionsConfig `yaml:"sections"`
}

// CompanyConfig contains company information for branding
type CompanyConfig struct {
	Name string `yaml:"name"`
	Logo string `yaml:"logo"` // Path to logo file
}

// BrandingConfig contains visual branding settings
type BrandingConfig struct {
	PrimaryColor   string `yaml:"primary_color"`
	SecondaryColor string `yaml:"secondary_color"`
	Footer         string `yaml:"footer"`
}

// ReportConfig contains report generation settings
type ReportConfig struct {
	ShowCoverPage bool   `yaml:"show_cover_page"`
	ShowTOC       bool   `yaml:"show_toc"`
	ShowCharts    bool   `yaml:"show_charts"`
	ColorScheme   string `yaml:"color_scheme"` // "default", "corporate", "dark"
}

// SectionsConfig controls which sections to include
type SectionsConfig struct {
	ExecutiveSummary   bool `yaml:"executive_summary"`
	RepositoryOverview bool `yaml:"repository_overview"`
	CodeQuality        bool `yaml:"code_quality"`
	Security           bool `yaml:"security"`
	Contributors       bool `yaml:"contributors"`
	Recommendations    bool `yaml:"recommendations"`
}

// DefaultPDFConfig returns default configuration
func DefaultPDFConfig() EnhancedPDFConfig {
	return EnhancedPDFConfig{
		Company: CompanyConfig{
			Name: "",
			Logo: "",
		},
		Branding: BrandingConfig{
			PrimaryColor:   "#1E40AF",
			SecondaryColor: "#10B981",
			Footer:         "",
		},
		Report: ReportConfig{
			ShowCoverPage: true,
			ShowTOC:       true,
			ShowCharts:    true,
			ColorScheme:   "default",
		},
		Sections: SectionsConfig{
			ExecutiveSummary:   true,
			RepositoryOverview: true,
			CodeQuality:        true,
			Security:           true,
			Contributors:       true,
			Recommendations:    true,
		},
	}
}

// LoadPDFConfig loads configuration from a YAML file
func LoadPDFConfig(path string) (*EnhancedPDFConfig, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config EnhancedPDFConfig
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}

// LoadPDFConfigOrDefault tries to load config from standard locations or returns default
func LoadPDFConfigOrDefault() *EnhancedPDFConfig {
	// Try loading from .repo-lyzer/pdf-config.yml
	configPath := filepath.Join(".repo-lyzer", "pdf-config.yml")
	if config, err := LoadPDFConfig(configPath); err == nil {
		return config
	}

	// Return default config
	defaultConfig := DefaultPDFConfig()
	return &defaultConfig
}
