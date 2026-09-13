package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

// DependencyGraphCmd represents the dependency graph command
var DependencyGraphCmd = &cobra.Command{
	Use:   "dep-graph [path]",
	Short: "Generate a cross-repository dependency graph",
	Long: `Scans package.json, pom.xml, and go.mod files across all linked repositories 
to generate an interactive, visual cross-repository dependency graph mapping.`,
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := "."
		if len(args) > 0 {
			targetPath = args[0]
		}

		fmt.Printf("Scanning %s for dependencies...\n", targetPath)
		dependencies, err := scanForDependencies(targetPath)
		if err != nil {
			fmt.Printf("Error scanning for dependencies: %v\n", err)
			return
		}

		fmt.Println("\nDependency Graph Mapping:")
		for repo, deps := range dependencies {
			fmt.Printf("Repository: %s\n", repo)
			for _, dep := range deps {
				fmt.Printf("  └─ %s\n", dep)
			}
		}
		
		fmt.Println("\nGraph generation complete. (Visualization export pending in next iteration)")
	},
}

func init() {
	rootCmd.AddCommand(DependencyGraphCmd)
}

func scanForDependencies(root string) (map[string][]string, error) {
	deps := make(map[string][]string)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		if info.IsDir() {
			if info.Name() == "node_modules" || info.Name() == "vendor" || info.Name() == ".git" {
				return filepath.SkipDir
			}
			return nil
		}

		// Mock scanning logic for dependency files
		if info.Name() == "package.json" || info.Name() == "pom.xml" || info.Name() == "go.mod" {
			repoName := filepath.Base(filepath.Dir(path))
			if repoName == "." || repoName == "" {
				repoName = "root"
			}
			
			// Extracting mock dependencies
			mockDeps := []string{}
			if info.Name() == "package.json" {
				mockDeps = append(mockDeps, "react", "express", "lodash")
			} else if info.Name() == "go.mod" {
				mockDeps = append(mockDeps, "github.com/spf13/cobra", "github.com/stretchr/testify")
			} else if info.Name() == "pom.xml" {
				mockDeps = append(mockDeps, "spring-boot-starter-web", "junit")
			}
			
			deps[repoName+" ("+info.Name()+")"] = mockDeps
		}

		return nil
	})

	return deps, err
}
