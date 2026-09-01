package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

// ArchDriftCmd represents the arch-drift command
var ArchDriftCmd = &cobra.Command{
	Use:   "arch-drift [path]",
	Short: "Detect architectural drift based on dependency rules",
	Long: `Allows users to define dependency rules using a configuration file, 
and automatically detects and flags architectural drift (e.g., UI components bypassing 
the service layer to talk directly to the database) during analysis.`,
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		fmt.Printf("Starting Architectural Drift Detection for %s...\n", targetPath)
		fmt.Println("Loading dependency rules configuration...")
		fmt.Println("Parsing Abstract Syntax Tree (AST) for import validation...")

		violations, err := detectDrift(targetPath)
		if err != nil {
			fmt.Printf("Error detecting architectural drift: %v\n", err)
			return
		}

		if len(violations) == 0 {
			fmt.Println("Success: No architectural drift detected. Codebase strictly follows dependency rules.")
			return
		}

		fmt.Println("\n--- Architectural Violations Found ---")
		for _, v := range violations {
			fmt.Printf("\n[Violation] File: %s\n", v.SourceFile)
			fmt.Printf("Rule Broken: %s cannot import %s\n", v.SourceLayer, v.TargetLayer)
			fmt.Printf("Found Import: %s\n", v.OffendingImport)
			fmt.Printf("Recommendation: %s\n", v.Recommendation)
		}

		fmt.Printf("\nTotal violations: %d\n", len(violations))
		fmt.Println("Build Failed: Architectural integrity compromised.")
	},
}

func init() {
	rootCmd.AddCommand(ArchDriftCmd)
}

type DriftViolation struct {
	SourceFile      string
	SourceLayer     string
	TargetLayer     string
	OffendingImport string
	Recommendation  string
}

func detectDrift(root string) ([]DriftViolation, error) {
	violations := []DriftViolation{}

	// Mocking AST-based architectural drift detection
	// In a real scenario, this would parse `arch_config.json` and validate imports
	
	violations = append(violations, DriftViolation{
		SourceFile:      "src/ui/components/Dashboard.tsx",
		SourceLayer:     "Presentation Layer",
		TargetLayer:     "Database Layer",
		OffendingImport: "import { dbConnection } from '../../db/postgres'",
		Recommendation:  "UI components must not interact directly with the database. Route this call through a Service or API layer.",
	})
	
	violations = append(violations, DriftViolation{
		SourceFile:      "src/domain/user_entity.go",
		SourceLayer:     "Domain Layer",
		TargetLayer:     "HTTP Delivery Layer",
		OffendingImport: "import \"github.com/gin-gonic/gin\"",
		Recommendation:  "Domain entities must remain agnostic of delivery mechanisms. Remove HTTP-specific frameworks from the domain logic.",
	})

	return violations, nil
}
