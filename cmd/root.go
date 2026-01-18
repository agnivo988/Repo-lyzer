package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "Repo-lyzer",
	Short:   "Analyze GitHub repositories from the terminal",
	Long:    "Repo-lyzer is a fast CLI tool written in Go to analyze GitHub repositories.",
	Version: "v1.0.0",
}

// Execute is used for cobra commands
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
