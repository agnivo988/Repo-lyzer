package cmd

import (
	"github.com/agnivo988/Repo-lyzer/internal/cache"
	"github.com/agnivo988/Repo-lyzer/internal/config"
	"github.com/agnivo988/Repo-lyzer/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var notificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "View analysis and export notifications",
	Long: `Display a list of all notifications from repository analyses, exports, and monitoring.

Notifications include:
  • Repository analysis completions
  • Export operations (JSON, PDF, Markdown, etc.)
  • Monitoring alerts and updates
  • Error notifications

The notifications view provides an interactive interface to:
  • Browse all notifications chronologically
  • Delete individual notifications
  • Clear all notifications at once
  • View detailed information about each event

Examples:
  # View all notifications
  repo-lyzer notifications

  # View notifications (alias)
  repo-lyzer notif`,
	Aliases: []string{"notif", "notify"},
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize cache and config
		cache, err := cache.NewCache()
		if err != nil {
			return err
		}

		appConfig, err := config.LoadSettings()
		if err != nil {
			// If settings file doesn't exist, use defaults
			appConfig = config.DefaultSettings()
		}

		// Create main model and set to notifications state
		model := ui.NewMainModel(cache, appConfig)
		model.SetStateNotifications()

		// Run the TUI
		p := tea.NewProgram(model, tea.WithAltScreen())
		_, err = p.Run()
		return err
	},
}

func init() {
	rootCmd.AddCommand(notificationsCmd)
}
