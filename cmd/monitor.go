package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/cache"
	"github.com/agnivo988/Repo-lyzer/internal/config"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/agnivo988/Repo-lyzer/internal/monitor"
	"github.com/agnivo988/Repo-lyzer/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
)

var monitorCmd = &cobra.Command{
	Use:   "monitor owner/repo",
	Short: "Monitor a GitHub repository in real-time",
	Long: `Monitor a GitHub repository for real-time updates including:
  • New commits
  • Issues and pull requests
  • Contributor changes
  • Repository health metrics

The monitoring runs continuously with configurable intervals and provides
notifications within the interactive TUI.

Examples:
  # Monitor with default 5-minute interval
  repo-lyzer monitor kubernetes/kubernetes

  # Monitor with custom 10-minute interval
  repo-lyzer monitor golang/go --interval 10m

  # Monitor critical production repo every minute
  repo-lyzer monitor company/production-api --interval 1m
  
  # Monitor with TUI dashboard
  repo-lyzer monitor facebook/react --dashboard`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// Validate the repository URL format
		owner, repo, err := validateRepoURL(args[0])
		if err != nil {
			return fmt.Errorf("invalid repository URL: %w", err)
		}

		// Get monitoring configuration
		interval, _ := cmd.Flags().GetDuration("interval")
		if interval == 0 {
			interval = 5 * time.Minute // Default 5 minutes
		}

		// Check if dashboard mode is requested
		useDashboard, _ := cmd.Flags().GetBool("dashboard")

		if useDashboard {
			// Use TUI dashboard
			cache, err := cache.NewCache()
			if err != nil {
				return fmt.Errorf("failed to initialize cache: %w", err)
			}

			config, err := config.LoadSettings()
			if err != nil {
				return fmt.Errorf("failed to load settings: %w", err)
			}

			model := ui.NewMainModel(cache, config)
			model.SetStateMonitorDashboard(owner, repo, interval)

			p := tea.NewProgram(model, tea.WithAltScreen())
			_, err = p.Run()
			return err
		}

		// Create monitor instance (CLI mode)
		// Check repository accessibility and prompt for token if private
		client := github.NewClient()
		_, err = client.GetRepo(owner, repo)
		if err != nil {
			// Check if it's a private repo error and no token is set
			if strings.Contains(err.Error(), "repository not found") && !client.HasToken() {
				fmt.Print("This appears to be a private repository. Please enter your GitHub access token: ")
				scanner := bufio.NewScanner(os.Stdin)
				if scanner.Scan() {
					token := strings.TrimSpace(scanner.Text())
					if token != "" {
						client.SetToken(token)
						// Retry fetching the repo with the token
						_, err = client.GetRepo(owner, repo)
						if err != nil {
							return fmt.Errorf("failed to access repository even with token: %w", err)
						}
						// If successful, set env variable for the monitor to inherit
						os.Setenv("GITHUB_TOKEN", token)
					} else {
						return fmt.Errorf("no token provided, cannot access private repository")
					}
				} else {
					return fmt.Errorf("failed to read token input")
				}
			} else {
				return fmt.Errorf("failed to access repository: %w", err)
			}
		}

		mon, err := monitor.NewMonitor(owner, repo, interval)
		if err != nil {
			return fmt.Errorf("failed to create monitor: %w", err)
		}

		// Start monitoring
		fmt.Printf("🔍 Starting real-time monitoring for %s/%s\n", owner, repo)
		fmt.Printf("📊 Check interval: %v\n", interval)
		fmt.Println("Press Ctrl+C to stop monitoring")

		return mon.Start()
	},
}

func init() {
	rootCmd.AddCommand(monitorCmd)
	monitorCmd.Flags().Duration("interval", 5*time.Minute, "Monitoring check interval (e.g., 5m, 10m, 1h)")
	monitorCmd.Flags().Bool("dashboard", false, "Use interactive TUI dashboard for monitoring")
}
