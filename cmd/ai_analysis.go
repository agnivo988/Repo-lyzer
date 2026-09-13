package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// AIAnalysisCmd represents the ai-analysis command
var AIAnalysisCmd = &cobra.Command{
	Use:   "ai-analyze [path]",
	Short: "Perform AI-driven code smell and anti-pattern detection",
	Long: `Integrates an LLM-based analysis layer to evaluate code blocks against known 
architectural anti-patterns (e.g. God Objects, duplicate logic) and provides contextual 
warnings and automated refactoring suggestions.`,
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		fmt.Printf("Starting AI-driven analysis on %s...\n", targetPath)
		fmt.Println("Connecting to LLM analysis engine (mocked for this iteration)...")
		
		findings, err := scanForAntiPatterns(targetPath)
		if err != nil {
			fmt.Printf("Error analyzing codebase: %v\n", err)
			return
		}

		if len(findings) == 0 {
			fmt.Println("No major architectural anti-patterns found.")
			return
		}

		fmt.Println("\n--- AI Analysis Report ---")
		for _, finding := range findings {
			fmt.Printf("\n[Warning] File: %s\n", finding.File)
			fmt.Printf("Anti-Pattern: %s\n", finding.PatternType)
			fmt.Printf("Suggestion: %s\n", finding.Suggestion)
		}
		fmt.Println("\nEnd of report.")
	},
}

func init() {
	rootCmd.AddCommand(AIAnalysisCmd)
}

type AntiPatternFinding struct {
	File        string
	PatternType string
	Suggestion  string
}

func scanForAntiPatterns(root string) ([]AntiPatternFinding, error) {
	findings := []AntiPatternFinding{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() && (info.Name() == "node_modules" || info.Name() == "vendor" || info.Name() == ".git") {
			return filepath.SkipDir
		}

		// Mock LLM analysis heuristics
		if !info.IsDir() && filepath.Ext(path) == ".go" && info.Size() > 10000 {
			findings = append(findings, AntiPatternFinding{
				File:        path,
				PatternType: "God Object / Bloated File",
				Suggestion:  "This file exceeds 10KB. Consider breaking down the structural logic into smaller, domain-specific packages.",
			})
		}
		
		if !info.IsDir() && info.Name() == "utils.go" {
			findings = append(findings, AntiPatternFinding{
				File:        path,
				PatternType: "Kitchen Sink Utility",
				Suggestion:  "Generic 'utils' files often become dumping grounds. Extract functions into specific, behavior-driven modules (e.g., 'stringutils', 'mathutils').",
			})
		}

		return nil
	})

	return findings, err
}
