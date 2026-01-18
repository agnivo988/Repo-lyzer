// Package cmd provides command-line interface commands for the Repo-lyzer application.
// It includes commands for analyzing repositories, comparing repositories, cache management, and running the interactive menu.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/cache"
	"github.com/spf13/cobra"
)

// cacheCmd defines the "cache" command group for cache management operations.
// It provides subcommands for clearing cached data and managing cache settings.
// Usage example:
//   repo-lyzer cache clear
var cacheCmd = &cobra.Command{
	Use:   "cache",
	Short: "Manage cached analysis data",
	Long:  "Manage cached repository analysis data, including clearing stale or corrupted cache entries.",
}

// clearCmd defines the "cache clear" subcommand.
// It safely removes all cached analysis data with user confirmation.
// This command provides a way to resolve cache corruption issues and free up disk space.
var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear all cached analysis data",
	Long: `Clear all cached repository analysis data from the local cache directory.

This command will permanently remove all cached analysis results. This can be useful for:
- Resolving cache corruption issues
- Freeing up disk space
- Ensuring fresh data for subsequent analyses

The command will prompt for confirmation before proceeding.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Initialize cache
		c, err := cache.NewCache()
		if err != nil {
			return fmt.Errorf("failed to initialize cache: %w", err)
		}

		// Get cache statistics
		stats := c.GetStats()

		if stats.TotalRepos == 0 {
			fmt.Println("ℹ️  Cache is already empty - no data to clear.")
			return nil
		}

		// Display current cache status
		fmt.Printf("📊 Current cache status:\n")
		fmt.Printf("   • Total cached repositories: %d\n", stats.TotalRepos)
		fmt.Printf("   • Valid entries: %d\n", stats.ValidRepos)
		fmt.Printf("   • Expired entries: %d\n", stats.ExpiredRepos)
		fmt.Printf("   • Total size: %.2f MB\n", stats.TotalSizeMB)
		fmt.Printf("   • Cache directory: %s\n\n", stats.CacheDir)

		// Prompt for confirmation
		fmt.Print("⚠️  This will permanently delete all cached data. Continue? (y/N): ")
		reader := bufio.NewReader(os.Stdin)
		response, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("failed to read user input: %w", err)
		}

		response = strings.TrimSpace(strings.ToLower(response))
		if response != "y" && response != "yes" {
			fmt.Println("❌ Operation cancelled.")
			return nil
		}

		// Clear the cache
		fmt.Println("🧹 Clearing cache...")
		if err := c.Clear(); err != nil {
			return fmt.Errorf("failed to clear cache: %w", err)
		}

		fmt.Println("✅ Cache cleared successfully!")
		fmt.Printf("   • Removed %d cached repositories\n", stats.TotalRepos)
		fmt.Printf("   • Freed approximately %.2f MB of disk space\n", stats.TotalSizeMB)

		return nil
	},
}

func init() {
	rootCmd.AddCommand(cacheCmd)
	cacheCmd.AddCommand(clearCmd)
}