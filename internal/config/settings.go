// Package config provides application settings and configuration management.
// It handles persistence of user preferences including theme, export options,
// and GitHub token configuration.
//
// Settings are stored in: ~/.repo-lyzer/settings.json
package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// ExportFormat represents available export formats
type ExportFormat string

const (
	ExportJSON     ExportFormat = "json"
	ExportMarkdown ExportFormat = "markdown"
	ExportCSV      ExportFormat = "csv"
	ExportHTML     ExportFormat = "html"
)

// AllExportFormats returns all available export formats
func AllExportFormats() []ExportFormat {
	return []ExportFormat{ExportJSON, ExportMarkdown, ExportCSV, ExportHTML}
}

// AppSettings holds all user-configurable application settings
type AppSettings struct {
	// Theme settings
	ThemeName string `json:"theme_name"`

	// Export settings
	DefaultExportFormat ExportFormat `json:"default_export_format"`
	ExportDirectory     string       `json:"export_directory"`

	// GitHub settings
	GitHubToken string `json:"github_token"`

	// Analysis settings
	DefaultAnalysisType string `json:"default_analysis_type"` // "quick", "detailed", "custom"
}

// DefaultSettings returns the default application settings
func DefaultSettings() *AppSettings {
	home, _ := os.UserHomeDir()
	return &AppSettings{
		ThemeName:           "Catppuccin Mocha",
		DefaultExportFormat: ExportJSON,
		ExportDirectory:     filepath.Join(home, "Downloads"),
		GitHubToken:         "",
		DefaultAnalysisType: "quick",
	}
}

// getSettingsDir returns the settings directory path
func getSettingsDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".repo-lyzer"), nil
}

// getSettingsPath returns the full path to the settings file
func getSettingsPath() (string, error) {
	dir, err := getSettingsDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "settings.json"), nil
}

// LoadSettings loads settings from disk, or returns defaults if not found
func LoadSettings() (*AppSettings, error) {
	settingsPath, err := getSettingsPath()
	if err != nil {
		return DefaultSettings(), err
	}

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		// Return defaults if file doesn't exist
		if os.IsNotExist(err) {
			return DefaultSettings(), nil
		}
		return DefaultSettings(), err
	}

	settings := DefaultSettings()
	if err := json.Unmarshal(data, settings); err != nil {
		return DefaultSettings(), err
	}

	return settings, nil
}

// SaveSettings saves settings to disk
func (s *AppSettings) SaveSettings() error {
	dir, err := getSettingsDir()
	if err != nil {
		return err
	}

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	settingsPath := filepath.Join(dir, "settings.json")
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(settingsPath, data, 0644)
}

// ResetToDefaults resets all settings to default values and saves
func ResetToDefaults() (*AppSettings, error) {
	settings := DefaultSettings()
	err := settings.SaveSettings()
	return settings, err
}

// SetTheme updates the theme name and saves
func (s *AppSettings) SetTheme(themeName string) error {
	s.ThemeName = themeName
	return s.SaveSettings()
}

// SetExportFormat updates the default export format and saves
func (s *AppSettings) SetExportFormat(format ExportFormat) error {
	s.DefaultExportFormat = format
	return s.SaveSettings()
}

// SetExportDirectory updates the export directory and saves
func (s *AppSettings) SetExportDirectory(dir string) error {
	s.ExportDirectory = dir
	return s.SaveSettings()
}

// SetGitHubToken updates the GitHub token and saves
func (s *AppSettings) SetGitHubToken(token string) error {
	s.GitHubToken = token
	return s.SaveSettings()
}

// ClearGitHubToken removes the GitHub token and saves
func (s *AppSettings) ClearGitHubToken() error {
	s.GitHubToken = ""
	return s.SaveSettings()
}

// HasGitHubToken returns true if a GitHub token is configured
func (s *AppSettings) HasGitHubToken() bool {
	return s.GitHubToken != ""
}

// GetMaskedToken returns the token with most characters masked for display
func (s *AppSettings) GetMaskedToken() string {
	if s.GitHubToken == "" {
		return ""
	}
	if len(s.GitHubToken) <= 8 {
		return strings.Repeat("*", len(s.GitHubToken))
	}
	// Show first 4 and last 4 characters
	return s.GitHubToken[:4] + strings.Repeat("*", len(s.GitHubToken)-8) + s.GitHubToken[len(s.GitHubToken)-4:]
}

// CycleExportFormat cycles to the next export format
func (s *AppSettings) CycleExportFormat() ExportFormat {
	formats := AllExportFormats()
	for i, f := range formats {
		if f == s.DefaultExportFormat {
			nextIndex := (i + 1) % len(formats)
			s.DefaultExportFormat = formats[nextIndex]
			s.SaveSettings()
			return s.DefaultExportFormat
		}
	}
	// Default to first format if current not found
	s.DefaultExportFormat = formats[0]
	s.SaveSettings()
	return s.DefaultExportFormat
}

// FormatDisplayName returns a user-friendly name for the export format
func (f ExportFormat) DisplayName() string {
	switch f {
	case ExportJSON:
		return "JSON"
	case ExportMarkdown:
		return "Markdown"
	case ExportCSV:
		return "CSV"
	case ExportHTML:
		return "HTML"
	default:
		return string(f)
	}
}
