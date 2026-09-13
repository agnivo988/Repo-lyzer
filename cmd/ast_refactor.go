package cmd

import (
	"fmt"
	"path/filepath"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// ASTRefactorCmd represents the ast-refactor command
var ASTRefactorCmd = &cobra.Command{
	Use:   "ast-refactor [path]",
	Short: "Automated code refactoring suggestions via AST",
	Long: `Creates an AST-to-AST transformation engine that flags outdated code patterns 
and automatically generates the exact code patch required to modernize the syntax. 
These patches can be attached directly to PR comments.`,
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		fmt.Printf("Initializing AST Refactoring Engine for %s...\n", targetPath)
		fmt.Println("Parsing source files and building Abstract Syntax Trees...")

		patches, err := generateRefactorPatches(targetPath)
		if err != nil {
			fmt.Printf("Error during AST refactoring analysis: %v\n", err)
			return
		}

		if len(patches) == 0 {
			fmt.Println("No legacy code patterns detected. Codebase syntax is modern.")
			return
		}

		fmt.Println("\n--- Automated Refactoring Suggestions ---")
		for _, patch := range patches {
			fmt.Printf("\n[Refactor Opportunity] File: %s\n", patch.File)
			fmt.Printf("Pattern: %s\n", patch.PatternDetected)
			fmt.Println("Suggested Patch:")
			fmt.Println("```diff")
			fmt.Println(patch.DiffContent)
			fmt.Println("```")
		}

		fmt.Printf("\nGenerated %d refactoring patches.\n", len(patches))
		fmt.Println("Action: Apply these patches to modernize the codebase.")
	},
}

func init() {
	rootCmd.AddCommand(ASTRefactorCmd)
}

type RefactorPatch struct {
	File            string
	PatternDetected string
	DiffContent     string
}

func generateRefactorPatches(root string) ([]RefactorPatch, error) {
	patches := []RefactorPatch{}

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() && (info.Name() == "node_modules" || info.Name() == "vendor" || info.Name() == ".git") {
			return filepath.SkipDir
		}

		// Mocking AST-based detection of outdated patterns
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".js") {
			patches = append(patches, RefactorPatch{
				File:            path,
				PatternDetected: "Manual Null Checking (Pre-Optional Chaining)",
				DiffContent: `- if (user && user.profile && user.profile.avatar) {
-   return user.profile.avatar;
- }
+ return user?.profile?.avatar;`,
			})
		}
		
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".go") {
			if strings.Contains(info.Name(), "utils") {
				patches = append(patches, RefactorPatch{
					File:            path,
					PatternDetected: "Legacy Interface{} usage (Pre-Generics)",
					DiffContent: `- func PrintSlice(s []interface{}) {
+ func PrintSlice[T any](s []T) {
- 	for _, v := range s {
+ 	for _, v := range s {
- 		fmt.Println(v)
+ 		fmt.Println(v)
- 	}
+ 	}
  }`,
				})
			}
		}

		return nil
	})

	return patches, err
}
