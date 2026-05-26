package ui

import (
	"fmt"

	"github.com/agnivo988/Repo-lyzer/internal/cache"
	"github.com/agnivo988/Repo-lyzer/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

type SettingsModel struct {
	cursor         int
	inTokenInput   bool
	tokenInput     string
	settingsOption string
	cache          *cache.Cache
	appConfig      *config.AppSettings
}

func NewSettingsModel() SettingsModel {
	return SettingsModel{
		cursor:         0,
		inTokenInput:   false,
		tokenInput:     "",
		settingsOption: "",
	}
}

func (m SettingsModel) Init() tea.Cmd {
	return nil
}

func (m SettingsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Handle token input mode separately
	if m.inTokenInput {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEnter:
				// Save the token
				if m.tokenInput != "" {
					// This would need access to appConfig, but for now just return
					m.inTokenInput = false
					m.tokenInput = ""
				}
			case tea.KeyEsc:
				m.inTokenInput = false
				m.tokenInput = ""
			case tea.KeyBackspace:
				if len(m.tokenInput) > 0 {
					m.tokenInput = m.tokenInput[:len(m.tokenInput)-1]
				}
			case tea.KeyRunes:
				m.tokenInput += string(msg.Runes)
			}
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, nil
		case "t":
			// Cycle through themes
			if m.settingsOption == "theme" || m.settingsOption == "" {
				CycleTheme()
			}
		case "1", "2", "3", "4", "5", "6", "7":
			// Select theme by number
			if m.settingsOption == "theme" || m.settingsOption == "" {
				idx := int(msg.String()[0] - '1')
				if idx >= 0 && idx < len(AvailableThemes) {
					SetThemeByIndex(idx)
				}
			}
		case "e":
			// Toggle cache enabled
			if m.settingsOption == "cache" {
				// Would need cache instance
			}
		case "a":
			// Toggle auto-cache
			if m.settingsOption == "cache" {
				// Would need cache instance
			}
		case "c":
			// Clear cache or clear token
			if m.settingsOption == "cache" {
				// Would need cache instance
			} else if m.settingsOption == "token" {
				// Would need appConfig
			}
		case "x":
			// Clean expired entries
			if m.settingsOption == "cache" {
				// Would need cache instance
			}
		case "f":
			// Cycle export format
			if m.settingsOption == "export" {
				// Would need appConfig
			}
		case "i":
			// Enter token input mode
			if m.settingsOption == "token" {
				m.inTokenInput = true
				m.tokenInput = ""
			}
		case "y":
			// Confirm reset
			if m.settingsOption == "reset" {
				// Would need to reset settings
			}
		}
	}
	return m, nil
}

func (m SettingsModel) View() string {
	var title string
	var content string

	switch m.settingsOption {
	case "theme":
		title = "🎨 Theme Settings"
		content = fmt.Sprintf(`
Current theme: %s

Available themes:
  [1] Default
  [2] Dark
  [3] Light
  [4] Monokai
  [5] Solarized
  [6] Dracula
  [7] Nord

Keybindings:
  • Press 1-7 to select a theme
  • Press 't' to cycle through themes
`, CurrentTheme.Name)
	case "cache":
		title = "💾 Cache Settings"
		status := "Disabled"
		autoCache := "Off"
		ttlStr := "24h"
		maxSize := "100 MB"
		totalRepos := 0
		validRepos := 0
		expiredRepos := 0
		cacheSize := 0.0

		if m.cache != nil {
			cfg := m.cache.GetConfig()
			if cfg.Enabled {
				status = "Enabled"
			}
			if cfg.AutoCache {
				autoCache = "On"
			}
			ttlStr = cache.FormatTTL(cfg.TTL)
			maxSize = fmt.Sprintf("%d MB", cfg.MaxSize)

			stats := m.cache.GetStats()
			totalRepos = stats.TotalRepos
			validRepos = stats.ValidRepos
			expiredRepos = stats.ExpiredRepos
			cacheSize = stats.TotalSizeMB
		}

		content = fmt.Sprintf(`
Status: %s
Auto-cache: %s
TTL: %s
Max Size: %s

Statistics:
  • Total repos cached: %d
  • Valid (not expired): %d
  • Expired: %d
  • Cache size: %.2f MB

Keybindings:
  • Press 'e' to toggle caching
  • Press 'a' to toggle auto-cache
  • Press 'c' to clear all cache
  • Press 'x' to clean expired entries
`, status, autoCache, ttlStr, maxSize, totalRepos, validRepos, expiredRepos, cacheSize)
	case "export":
		title = "📤 Export Options"
		formatStr := "JSON"
		exportDir := "~/Downloads/"
		if m.appConfig != nil {
			formatStr = m.appConfig.DefaultExportFormat.DisplayName()
			exportDir = m.appConfig.ExportDirectory
		}
		content = fmt.Sprintf(`
Current export format: %s
Export directory: %s

Available formats:
  [1] JSON
  [2] Markdown
  [3] CSV
  [4] HTML
  [5] PDF

Keybindings:
  • Press 'f' to cycle through formats
`, formatStr, exportDir)
	case "token":
		title = "🔑 GitHub Token"

		if m.inTokenInput {
			content = fmt.Sprintf(`
Enter GitHub Personal Access Token:

> %s█

Press Enter to save, ESC to cancel.
`, m.tokenInput)
		} else {
			status := "❌ Not configured"
			if m.appConfig != nil && m.appConfig.HasGitHubToken() {
				status = fmt.Sprintf("✅ Configured (%s)", m.appConfig.GetMaskedToken())
			}
			content = fmt.Sprintf(`
GitHub API Token Configuration:

Status: %s

Benefits of using a token:
  • Higher API rate limits (5000 vs 60 requests/hour)
  • Access to private repositories
  • More detailed analysis

Keybindings:
  • Press 'i' to input a new token
  • Press 'c' to clear saved token
`, status)
		}
	case "reset":
		title = "🔄 Reset to Defaults"
		content = `
Reset all settings to default values:

This will:
  • Clear all saved settings
  • Reset theme to default
  • Clear export preferences
  • Remove custom configurations

Warning: This action cannot be undone.

Press 'y' to confirm reset, or ESC to cancel.
`
	default:
		title = "⚙️ Settings"
		content = `
Select a settings option from the menu.
`
	}

	settingsContent := TitleStyle.Render(title) + "\n\n" + content + "\n\n" + SubtleStyle.Render("Press ESC or q to go back")

	box := BoxStyle.Render(settingsContent)

	return box
}
