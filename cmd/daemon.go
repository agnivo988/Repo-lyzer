// Package cmd provides command-line interface commands for the Repo-lyzer application.
// It includes the daemon command for running the scheduler continuously.
package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/scheduler"
	"github.com/spf13/cobra"
)

// daemonCmd defines the "daemon" command for running the scheduler continuously
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run the scheduler daemon continuously",
	Long: `Run the Repo-lyzer scheduler as a daemon process that continuously
monitors and executes scheduled analysis reports.

The daemon will:
  • Load all scheduled jobs from configuration
  • Execute jobs according to their schedule (daily, weekly, monthly)
  • Export reports to configured destinations (local path or webhook)
  • Handle graceful shutdown on SIGINT/SIGTERM

Examples:
  # Start the scheduler daemon
  repo-lyzer daemon

  # Start with verbose output
  repo-lyzer daemon --verbose`,
	RunE: func(cmd *cobra.Command, args []string) error {
		verbose, _ := cmd.Flags().GetBool("verbose")

		if verbose {
			fmt.Println("🔄 Starting Repo-lyzer scheduler daemon...")
		}

		// Create new scheduler
		sched, err := scheduler.NewScheduler()
		if err != nil {
			return fmt.Errorf("failed to create scheduler: %w", err)
		}

		// Start the scheduler
		sched.Start()

		// Get job count
		jobs := sched.ListJobs()
		enabledJobs := 0
		for _, job := range jobs {
			if job.Enabled {
				enabledJobs++
			}
		}

		fmt.Println("📅 Scheduler daemon started")
		fmt.Printf("   Total jobs: %d\n", len(jobs))
		fmt.Printf("   Enabled jobs: %d\n", enabledJobs)
		fmt.Println("\nPress Ctrl+C to stop the daemon")

		// Set up signal handling for graceful shutdown
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Wait for shutdown signal
		sig := <-sigChan
		fmt.Printf("\n🛑 Received signal: %v\n", sig)
		fmt.Println("   Shutting down scheduler...")

		// Stop the scheduler gracefully
		sched.Stop()

		fmt.Println("✅ Scheduler daemon stopped gracefully")
		return nil
	},
}

// daemonStatusCmd defines the "daemon status" command
var daemonStatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check daemon status",
	RunE: func(cmd *cobra.Command, args []string) error {
		sched, err := scheduler.NewScheduler()
		if err != nil {
			return fmt.Errorf("failed to create scheduler: %w", err)
		}

		jobs := sched.ListJobs()

		fmt.Println("📊 Scheduler Status")
		fmt.Println(strings.Repeat("─", 31))

		enabledCount := 0
		disabledCount := 0

		for _, job := range jobs {
			if job.Enabled {
				enabledCount++
			} else {
				disabledCount++
			}
		}

		fmt.Printf("Total Jobs: %d\n", len(jobs))
		fmt.Printf("Enabled: %d\n", enabledCount)
		fmt.Printf("Disabled: %d\n", disabledCount)

		if len(jobs) > 0 {
			fmt.Println("\n📋 Job Summary:")
			for _, job := range jobs {
				status := "✅"
				if !job.Enabled {
					status = "❌"
				}
				fmt.Printf("   %s %s - %s (%s)\n", status, job.ID, job.GetRepoFullName(), job.Interval.DisplayName())
			}
		}

		fmt.Println()
		return nil
	},
}

// daemonValidateCmd defines the "daemon validate" command to validate cron expressions
var daemonValidateCmd = &cobra.Command{
	Use:   "validate <cron-expression>",
	Short: "Validate a cron expression",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cronExpr := args[0]

		err := scheduler.ValidateCronExpression(cronExpr)
		if err != nil {
			return fmt.Errorf("invalid cron expression: %w", err)
		}

		fmt.Printf("✅ Valid cron expression: %s\n", cronExpr)
		return nil
	},
}

// daemonNextRunCmd defines the "daemon next-run" command to show next run times
var daemonNextRunCmd = &cobra.Command{
	Use:   "next-run",
	Short: "Show next scheduled run times",
	RunE: func(cmd *cobra.Command, args []string) error {
		sched, err := scheduler.NewScheduler()
		if err != nil {
			return fmt.Errorf("failed to create scheduler: %w", err)
		}

		jobs := sched.ListJobs()

		if len(jobs) == 0 {
			fmt.Println("No scheduled jobs found.")
			return nil
		}

		fmt.Println("⏰ Next Scheduled Runs")
		fmt.Println(strings.Repeat("─", 51))

		for _, job := range jobs {
			if job.Enabled {
				nextRun := job.NextRun
				if nextRun.IsZero() {
					nextRun = time.Now().Add(24 * time.Hour) // Default
				}

				duration := nextRun.Sub(time.Now())
				if duration < 0 {
					duration = 0
				}

				fmt.Printf("\n%s\n", job.GetRepoFullName())
				fmt.Printf("   Next Run: %s\n", nextRun.Format("2006-01-02 15:04:05"))
				fmt.Printf("   In: %s\n", duration.Round(time.Minute))
				fmt.Printf("   Interval: %s\n", job.Interval.DisplayName())
			}
		}

		fmt.Println()
		return nil
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)

	// Add subcommands
	daemonCmd.AddCommand(daemonStatusCmd)
	daemonCmd.AddCommand(daemonValidateCmd)
	daemonCmd.AddCommand(daemonNextRunCmd)

	// Add flags
	daemonCmd.Flags().BoolP("verbose", "v", false, "Enable verbose output")
}

// RunDaemon starts the scheduler daemon programmatically
func RunDaemon() error {
	sched, err := scheduler.NewScheduler()
	if err != nil {
		return fmt.Errorf("failed to create scheduler: %w", err)
	}

	sched.Start()

	jobs := sched.ListJobs()
	enabledJobs := 0
	for _, job := range jobs {
		if job.Enabled {
			enabledJobs++
		}
	}

	fmt.Println("📅 Scheduler daemon started")
	fmt.Printf("   Total jobs: %d\n", len(jobs))
	fmt.Printf("   Enabled jobs: %d\n", enabledJobs)
	fmt.Println("\nPress Ctrl+C to stop the daemon")

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for shutdown signal
	sig := <-sigChan
	fmt.Printf("\n🛑 Received signal: %v\n", sig)
	fmt.Println("   Shutting down scheduler...")

	// Stop the scheduler gracefully
	sched.Stop()

	fmt.Println("✅ Scheduler daemon stopped gracefully")
	return nil
}
