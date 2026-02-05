package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "repo-lyzer",
	Short: "Analyze GitHub repositories from the terminal",
	Long: `Repo-lyzer - Professional GitHub Repository Analysis Tool

A modern, terminal-based CLI tool written in Go that analyzes GitHub repositories 
and presents comprehensive insights in a beautifully formatted, interactive dashboard.

Perfect for:
  • Developers evaluating open-source projects
  • Recruiters assessing repository health and activity
  • Contributors exploring project structure and engagement
  • Teams monitoring code quality and security

Features:
  ✓ Repository Health Score & Maturity Analysis
  ✓ Commit Activity & Contributor Insights
  ✓ Code Quality Dashboard with Problem Hotspots
  ✓ Security Vulnerability Scanning
  ✓ Bus Factor & Risk Assessment
  ✓ Language Breakdown & File Tree Viewer
  ✓ Enhanced PDF Reports with Charts & Branding
  ✓ Export to JSON, Markdown, CSV, HTML, PDF
  ✓ Side-by-side Repository Comparison
  ✓ Offline Caching & Real-time Monitoring
  ✓ Interactive TUI with Keyboard Navigation

Examples:
  # Interactive menu mode
  repo-lyzer

  # Quick analysis
  repo-lyzer analyze golang/go

  # With recruiter summary
  repo-lyzer summary charmbracelet/bubbletea

  # Compare repositories
  repo-lyzer compare facebook/react vuejs/vue

  # Monitor repository changes
  repo-lyzer monitor kubernetes/kubernetes

  # View cached analyses
  repo-lyzer cache list

Documentation: https://github.com/agnivo988/Repo-lyzer`,
	Version: "v1.0.0",
}

// Execute is used for cobra commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
