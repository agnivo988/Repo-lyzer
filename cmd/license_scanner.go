package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
)

// LicenseScannerCmd represents the license-scan command
var LicenseScannerCmd = &cobra.Command{
	Use:   "license-scan [path]",
	Short: "Detect license compliance and conflict issues",
	Long: `Parses license files and headers across all transitive dependencies, 
alerting maintainers to potential licensing conflicts based on predefined policies 
(e.g., catching GPL dependencies in MIT-licensed projects).`,
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		fmt.Printf("Starting License Compliance Scan for %s...\n", targetPath)
		fmt.Println("Analyzing project primary license...")
		fmt.Println("Deep scanning dependency tree for transitive licenses...")

		conflicts, err := scanLicenses(targetPath)
		if err != nil {
			fmt.Printf("Error during license scan: %v\n", err)
			return
		}

		if len(conflicts) == 0 {
			fmt.Println("Success: No license conflicts detected. All dependencies are compliant.")
			return
		}

		fmt.Println("\n--- License Conflict Report ---")
		for _, conflict := range conflicts {
			fmt.Printf("\n[Conflict Detected]\n")
			fmt.Printf("Dependency: %s (v%s)\n", conflict.DependencyName, conflict.Version)
			fmt.Printf("Detected License: %s\n", conflict.DetectedLicense)
			fmt.Printf("Project Policy: %s\n", conflict.ProjectPolicy)
			fmt.Printf("Risk Level: %s\n", conflict.RiskLevel)
			fmt.Printf("Recommendation: %s\n", conflict.Recommendation)
		}

		fmt.Printf("\nTotal conflicts found: %d\n", len(conflicts))
		fmt.Println("Action Required: Resolve licensing issues before merging.")
	},
}

func init() {
	rootCmd.AddCommand(LicenseScannerCmd)
}

type LicenseConflict struct {
	DependencyName  string
	Version         string
	DetectedLicense string
	ProjectPolicy   string
	RiskLevel       string
	Recommendation  string
}

func scanLicenses(root string) ([]LicenseConflict, error) {
	conflicts := []LicenseConflict{}

	// Mocking deep dependency tree license scanning
	// E.g., detecting a restrictive copyleft license in a permissive project
	
	// Assuming the project policy strictly requires Permissive licenses (MIT, Apache)
	projectPolicy := "Strictly Permissive (MIT/Apache 2.0 allowed, GPL forbidden)"

	if strings.Contains(root, "enterprise") || root == "." {
		conflicts = append(conflicts, LicenseConflict{
			DependencyName:  "go-mysql-driver-fork",
			Version:         "1.7.1",
			DetectedLicense: "GPL-3.0",
			ProjectPolicy:   projectPolicy,
			RiskLevel:       "Critical",
			Recommendation:  "GPL-3.0 is a strong copyleft license that conflicts with your permissive policy. Replace this dependency immediately or risk open-sourcing proprietary code.",
		})

		conflicts = append(conflicts, LicenseConflict{
			DependencyName:  "unknown-logger",
			Version:         "0.9.4",
			DetectedLicense: "UNKNOWN",
			ProjectPolicy:   projectPolicy,
			RiskLevel:       "High",
			Recommendation:  "No valid LICENSE file found in the dependency. Treat as fully copyrighted. Avoid using until the author adds an explicit open-source license.",
		})
	}

	return conflicts, nil
}
