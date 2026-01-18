// Package main is the entry point for the Repo-lyzer application.
// Repo-lyzer is a tool for analyzing GitHub repositories, providing insights
// into code health, contributor activity, and project maturity.
package main

import (
	"os"

	"github.com/agnivo988/Repo-lyzer/cmd"
)

// main initializes and runs the Repo-lyzer application.
// It starts the interactive menu interface for repository analysis or executes CLI commands.
func main() {
	if len(os.Args) > 1 {
		cmd.Execute()
	} else {
		cmd.RunMenu()
	}
}
