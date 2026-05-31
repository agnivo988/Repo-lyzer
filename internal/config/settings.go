// Package config provides application settings and configuration management.
// It handles persistence of user preferences including theme, export options,
// and GitHub token configuration.
//
// Settings are stored in: ~/.repo-lyzer/settings.json
package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/zalando/go-keyring"
	"github.com/robfig/cron/v3"
)

const (
	secureStoreService = "repo-lyzer"
	secureStoreUser    = "github_token"
)

// ExportFormat represents available export formats
type ExportFormat string

const (
	ExportJSON     ExportFormat = "json"
	ExportMarkdown ExportFormat = "markdown"
	ExportCSV      ExportFormat = "csv"
	ExportHTML     ExportFormat = "html"
	ExportPDF      ExportFormat = "pdf"
)

// AllExportFormats returns all available export formats
func AllExportFormats() []ExportFormat {
	return []ExportFormat{ExportJSON, ExportMarkdown, ExportCSV, ExportHTML, ExportPDF}
}

// AppSettings holds all user-configurable application settings
type AppSettings struct {
	// Theme settings
	ThemeName string `json:"theme_name"`

	// Export settings
	DefaultExportFormat ExportFormat `json:"default_export_format"`
	ExportDirectory     string       `json:"export_directory"`

	// GitHub settings
	GitHubToken string `json:"-"`

	// Analysis settings
	DefaultAnalysisType string `json:"default_analysis_type"` // "quick", "detailed", "custom"

	// Log settings
	LogLevel string `json:"log_level"`

	// Monitoring settings
	MonitoringEnabled      bool          `json:"monitoring_enabled"`
	DefaultMonitorInterval time.Duration `json:"default_monitor_interval"`
	NotificationEnabled    bool          `json:"notification_enabled"`

	// Scheduled jobs for automated report scheduling
	ScheduledJobs []ScheduledJob `json:"scheduled_jobs"`
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
		LogLevel:            "info",
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
	if envPath := os.Getenv("REPO_LYZER_CONFIG_PATH"); envPath != "" {
		return envPath, nil
	}
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
		settings := applyEnvOverrides(DefaultSettings())
		if settings.GitHubToken == "" {
			if token, secureErr := loadTokenFromSecureStore(); secureErr != nil {
				return settings, secureErr
			} else {
				settings.GitHubToken = token
			}
		}
		return settings, err
	}

	settings := DefaultSettings()

	data, err := os.ReadFile(settingsPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return applyEnvOverrides(settings), err
		}
		// File doesn't exist, just apply env overrides to defaults
	} else {
		if err := json.Unmarshal(data, settings); err != nil {
			return applyEnvOverrides(DefaultSettings()), err
		}
	}

	settings = applyEnvOverrides(settings)
	if settings.GitHubToken == "" {
		if token, secureErr := loadTokenFromSecureStore(); secureErr != nil {
			return settings, secureErr
		} else {
			settings.GitHubToken = token
		}
	}

	return settings, nil
}

// applyEnvOverrides applies environment variable overrides to the given settings
func applyEnvOverrides(settings *AppSettings) *AppSettings {
	if token := os.Getenv("REPO_LYZER_GITHUB_TOKEN"); token != "" {
		settings.GitHubToken = token
	}

	if interval := os.Getenv("REPO_LYZER_INTERVAL"); interval != "" {
		if d, err := time.ParseDuration(interval); err == nil {
			settings.DefaultMonitorInterval = d
		}
	}

	if logLevel := os.Getenv("REPO_LYZER_LOG_LEVEL"); logLevel != "" {
		settings.LogLevel = strings.ToLower(logLevel)
	}

	return settings
}

// SaveSettings saves settings to disk
func (s *AppSettings) SaveSettings() error {
	var dir, settingsPath string

	if envPath := os.Getenv("REPO_LYZER_CONFIG_PATH"); envPath != "" {
		settingsPath = envPath
		dir = filepath.Dir(envPath)
	} else {
		var err error
		dir, err = getSettingsDir()
		if err != nil {
			return err
		}
		settingsPath = filepath.Join(dir, "settings.json")
	}

	// Ensure directory exists
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(settingsPath, data, 0600); err != nil {
		return err
	}

	// Best effort: keep file private when permissions are supported.
	_ = os.Chmod(settingsPath, 0600)
	return nil
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
	if token != "" {
		if err := saveTokenToSecureStore(token); err != nil {
			return err
		}
	}
	return s.SaveSettings()
}

// ClearGitHubToken removes the GitHub token and saves
func (s *AppSettings) ClearGitHubToken() error {
	s.GitHubToken = ""
	if err := clearTokenFromSecureStore(); err != nil {
		return err
	}
	return s.SaveSettings()
}

func loadTokenFromSecureStore() (string, error) {
	token, err := keyring.Get(secureStoreService, secureStoreUser)
	if err == nil {
		return token, nil
	}

	if errors.Is(err, keyring.ErrNotFound) || errors.Is(err, keyring.ErrUnsupportedPlatform) {
		return "", nil
	}

	return "", err
}

func saveTokenToSecureStore(token string) error {
	err := keyring.Set(secureStoreService, secureStoreUser, token)
	if errors.Is(err, keyring.ErrUnsupportedPlatform) {
		return nil
	}
	return err
}

func clearTokenFromSecureStore() error {
	err := keyring.Delete(secureStoreService, secureStoreUser)
	if errors.Is(err, keyring.ErrNotFound) || errors.Is(err, keyring.ErrUnsupportedPlatform) {
		return nil
	}
	return err
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
	case ExportPDF:
		return "PDF"
	default:
		return string(f)
	}
}

// ScheduleInterval represents the interval for scheduled reports
type ScheduleInterval string

const (
	ScheduleDaily   ScheduleInterval = "daily"
	ScheduleWeekly  ScheduleInterval = "weekly"
	ScheduleMonthly ScheduleInterval = "monthly"
	ScheduleCustom  ScheduleInterval = "custom"
)

// CronExpression returns the cron expression for the given interval
func (si ScheduleInterval) CronExpression() string {
	switch si {
	case ScheduleDaily:
		return "0 9 * * *" // Daily at 9 AM
	case ScheduleWeekly:
		return "0 9 * * 1" // Weekly on Monday at 9 AM
	case ScheduleMonthly:
		return "0 9 1 * *" // Monthly on 1st at 9 AM
	default:
		return "0 9 * * *" // Default to daily
	}
}

// DisplayName returns a user-friendly name for the schedule interval
func (si ScheduleInterval) DisplayName() string {
	switch si {
	case ScheduleDaily:
		return "Daily"
	case ScheduleWeekly:
		return "Weekly"
	case ScheduleMonthly:
		return "Monthly"
	case ScheduleCustom:
		return "Custom"
	default:
		return string(si)
	}
}

// OutputDestination represents where scheduled reports are sent
type OutputDestination struct {
	Type       string `json:"type"`        // "local" or "webhook"
	LocalPath  string `json:"local_path"`  // Local directory path
	WebhookURL string `json:"webhook_url"` // Webhook URL for notifications
	Enabled    bool   `json:"enabled"`     // Whether this destination is enabled
}

// ScheduledJob represents a scheduled analysis report job
type ScheduledJob struct {
	ID             string            `json:"id"`              // Unique job identifier
	Owner          string            `json:"owner"`           // Repository owner
	Repo           string            `json:"repo"`            // Repository name
	Interval       ScheduleInterval  `json:"interval"`        // Schedule interval
	CronExpression string            `json:"cron_expression"` // Custom cron expression (if interval is custom)
	Format         ExportFormat      `json:"format"`          // Export format
	Destination    OutputDestination `json:"destination"`     // Output destination
	Enabled        bool              `json:"enabled"`         // Whether the job is enabled
	LastRun        time.Time         `json:"last_run"`        // Last execution time
	NextRun        time.Time         `json:"next_run"`        // Next scheduled execution
	CreatedAt      time.Time         `json:"created_at"`      // Job creation time
	AnalysisType   string            `json:"analysis_type"`   // Type of analysis to run
}

// GetRepoFullName returns the full repository name (owner/repo)
func (sj *ScheduledJob) GetRepoFullName() string {
	return sj.Owner + "/" + sj.Repo
}

// GetNextCronExpression returns the cron expression for the job
func (sj *ScheduledJob) GetCronExpression() string {
	if sj.Interval == ScheduleCustom && sj.CronExpression != "" {
		return sj.CronExpression
	}
	return sj.Interval.CronExpression()
}

// AddScheduledJob adds a new scheduled job and saves settings
func (s *AppSettings) AddScheduledJob(job ScheduledJob) error {
	// Generate ID if not provided
	if job.ID == "" {
		job.ID = generateJobID(job.Owner, job.Repo)
	}
	job.CreatedAt = time.Now()
	job.NextRun = calculateNextRun(job.GetCronExpression())

	s.ScheduledJobs = append(s.ScheduledJobs, job)
	return s.SaveSettings()
}

// RemoveScheduledJob removes a scheduled job by ID and saves settings
func (s *AppSettings) RemoveScheduledJob(jobID string) error {
	for i, job := range s.ScheduledJobs {
		if job.ID == jobID {
			s.ScheduledJobs = append(s.ScheduledJobs[:i], s.ScheduledJobs[i+1:]...)
			return s.SaveSettings()
		}
	}
	return nil
}

// GetScheduledJobs returns all scheduled jobs
func (s *AppSettings) GetScheduledJobs() []ScheduledJob {
	return s.ScheduledJobs
}

// GetScheduledJobByID returns a scheduled job by ID
func (s *AppSettings) GetScheduledJobByID(jobID string) *ScheduledJob {
	for _, job := range s.ScheduledJobs {
		if job.ID == jobID {
			return &job
		}
	}
	return nil
}

// UpdateScheduledJob updates an existing scheduled job and saves settings
func (s *AppSettings) UpdateScheduledJob(job ScheduledJob) error {
	for i, existingJob := range s.ScheduledJobs {
		if existingJob.ID == job.ID {
			s.ScheduledJobs[i] = job
			return s.SaveSettings()
		}
	}
	return nil
}

// EnableScheduledJob enables or disables a scheduled job
func (s *AppSettings) EnableScheduledJob(jobID string, enabled bool) error {
	for i, job := range s.ScheduledJobs {
		if job.ID == jobID {
			s.ScheduledJobs[i].Enabled = enabled
			return s.SaveSettings()
		}
	}
	return nil
}

// generateJobID generates a unique job ID
func generateJobID(owner, repo string) string {
	return owner + "-" + repo + "-" + time.Now().Format("20060102150405")
}

// calculateNextRun calculates the next run time based on cron expression
func calculateNextRun(cronExpr string) time.Time {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	schedule, err := parser.Parse(cronExpr)
	if err != nil {
		// Fallback if the expression is invalid
		return time.Now().Add(24 * time.Hour)
	}
	return schedule.Next(time.Now())
}
