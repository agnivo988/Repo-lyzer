// Package analyzer provides analysis functions for GitHub repositories.
// This file contains dependency analysis functionality that parses various
// package manager files to extract dependency information.
package analyzer

import (
	"encoding/base64"
	"encoding/json"
	"regexp"
	"sort"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// Dependency represents a single project dependency with its metadata.
// It captures the essential information about a package dependency
// regardless of the package manager used.
type Dependency struct {
	Name    string `json:"name"`    // Package name (e.g., "react", "github.com/gin-gonic/gin")
	Version string `json:"version"` // Version constraint (e.g., "^1.0.0", "v1.9.1")
	Type    string `json:"type"`    // Dependency type: "production", "dev", "peer", "indirect"
}

// DependencyFile represents all dependencies extracted from a single file.
// A repository may contain multiple dependency files (e.g., monorepo with
// multiple package.json files).
type DependencyFile struct {
	Filename     string       `json:"filename"`  // Full path to the file (e.g., "packages/web/package.json")
	FileType     string       `json:"file_type"` // Package manager type: "npm", "go", "python", "rust", "ruby"
	Dependencies []Dependency `json:"dependencies"`
	TotalCount   int          `json:"total_count"` // Total number of dependencies in this file
}

// DependencyAnalysis holds the complete dependency analysis for a repository.
// It aggregates information from all dependency files found in the repo.
type DependencyAnalysis struct {
	Files       []DependencyFile `json:"files"`         // All parsed dependency files
	TotalDeps   int              `json:"total_deps"`    // Total dependencies across all files
	Languages   []string         `json:"languages"`     // Detected package managers/languages
	HasLockFile bool             `json:"has_lock_file"` // Whether a lock file exists
}

// AnalyzeDependencies fetches and parses dependency files from a repository.
// It supports multiple package managers and handles monorepos with multiple
// dependency files.
//
// Supported package managers:
//   - npm (package.json)
//   - Go (go.mod)
//   - Python (requirements.txt, Pipfile, pyproject.toml)
//   - Rust (Cargo.toml)
//   - Ruby (Gemfile)
//
// Parameters:
//   - client: GitHub API client for fetching file contents
//   - owner: Repository owner (e.g., "facebook")
//   - repo: Repository name (e.g., "react")
//   - branch: Branch name to analyze (e.g., "main")
//   - fileTree: Pre-fetched file tree from the repository
//
// Returns:
//   - *DependencyAnalysis: Aggregated dependency information
//   - error: Any error encountered during analysis
func AnalyzeDependencies(client *github.Client, owner, repo, branch string, fileTree []github.TreeEntry) (*DependencyAnalysis, error) {
	analysis := &DependencyAnalysis{
		Files:     []DependencyFile{},
		Languages: []string{},
	}

	// Find all dependency files in the repository tree
	depFiles := findDependencyFiles(fileTree)

	for _, df := range depFiles {
		// Fetch file content from GitHub API
		content, err := client.GetFileContent(owner, repo, df.path)
		if err != nil {
			continue // Skip files we can't read (permissions, size limits, etc.)
		}

		// GitHub API returns base64 encoded content
		decoded, err := base64.StdEncoding.DecodeString(content)
		if err != nil {
			continue
		}

		var deps []Dependency
		var fileType string

		// Parse based on file type
		switch df.fileType {
		case "npm":
			deps, fileType = parsePackageJSON(decoded)
		case "go":
			deps, fileType = parseGoMod(decoded)
		case "python":
			deps, fileType = parseRequirementsTxt(decoded)
		case "rust":
			deps, fileType = parseCargoToml(decoded)
		case "ruby":
			deps, fileType = parseGemfile(decoded)
		}

		if len(deps) > 0 {
			analysis.Files = append(analysis.Files, DependencyFile{
				Filename:     df.path,
				FileType:     fileType,
				Dependencies: deps,
				TotalCount:   len(deps),
			})
			analysis.TotalDeps += len(deps)

			// Track unique languages/package managers
			if !contains(analysis.Languages, fileType) {
				analysis.Languages = append(analysis.Languages, fileType)
			}
		}
	}

	// Check for lock files (indicates reproducible builds)
	analysis.HasLockFile = hasLockFile(fileTree)

	return analysis, nil
}

// depFileInfo holds metadata about a dependency file found in the repo.
type depFileInfo struct {
	path     string // Full path to the file
	fileType string // Package manager type
}

// findDependencyFiles scans the file tree for known dependency file patterns.
// It returns a list of files that should be parsed for dependencies.
func findDependencyFiles(tree []github.TreeEntry) []depFileInfo {
	var files []depFileInfo

	// Map of filename to package manager type
	depFilePatterns := map[string]string{
		"package.json":     "npm",
		"go.mod":           "go",
		"requirements.txt": "python",
		"Pipfile":          "python",
		"pyproject.toml":   "python",
		"Cargo.toml":       "rust",
		"Gemfile":          "ruby",
	}

	for _, entry := range tree {
		// Only process files (blobs), not directories (trees)
		if entry.Type != "blob" {
			continue
		}

		// Extract filename from full path
		parts := strings.Split(entry.Path, "/")
		filename := parts[len(parts)-1]

		if fileType, ok := depFilePatterns[filename]; ok {
			files = append(files, depFileInfo{
				path:     entry.Path,
				fileType: fileType,
			})
		}
	}

	return files
}

// hasLockFile checks if the repository contains any lock files.
// Lock files indicate that the project uses reproducible dependency resolution.
func hasLockFile(tree []github.TreeEntry) bool {
	lockFiles := []string{
		"package-lock.json", // npm
		"yarn.lock",         // Yarn
		"pnpm-lock.yaml",    // pnpm
		"go.sum",            // Go modules
		"Pipfile.lock",      // Pipenv
		"poetry.lock",       // Poetry
		"Cargo.lock",        // Cargo (Rust)
		"Gemfile.lock",      // Bundler (Ruby)
	}

	for _, entry := range tree {
		parts := strings.Split(entry.Path, "/")
		filename := parts[len(parts)-1]

		for _, lockFile := range lockFiles {
			if filename == lockFile {
				return true
			}
		}
	}
	return false
}

// parsePackageJSON parses an npm package.json file and extracts dependencies.
// It handles dependencies, devDependencies, and peerDependencies sections.
//
// Example package.json structure:
//
//	{
//	  "dependencies": { "react": "^18.0.0" },
//	  "devDependencies": { "jest": "^29.0.0" },
//	  "peerDependencies": { "react-dom": "^18.0.0" }
//	}
func parsePackageJSON(content []byte) ([]Dependency, string) {
	var pkg struct {
		Dependencies    map[string]string `json:"dependencies"`
		DevDependencies map[string]string `json:"devDependencies"`
		PeerDeps        map[string]string `json:"peerDependencies"`
	}

	if err := json.Unmarshal(content, &pkg); err != nil {
		return nil, "npm"
	}

	var deps []Dependency

	// Production dependencies
	for name, version := range pkg.Dependencies {
		deps = append(deps, Dependency{
			Name:    name,
			Version: cleanVersion(version),
			Type:    "production",
		})
	}

	// Development dependencies
	for name, version := range pkg.DevDependencies {
		deps = append(deps, Dependency{
			Name:    name,
			Version: cleanVersion(version),
			Type:    "dev",
		})
	}

	// Peer dependencies
	for name, version := range pkg.PeerDeps {
		deps = append(deps, Dependency{
			Name:    name,
			Version: cleanVersion(version),
			Type:    "peer",
		})
	}

	// Sort alphabetically for consistent output
	sort.Slice(deps, func(i, j int) bool {
		return deps[i].Name < deps[j].Name
	})

	return deps, "npm"
}

// parseGoMod parses a Go go.mod file and extracts module dependencies.
// It handles both single-line requires and require blocks.
//
// Example go.mod structure:
//
//	require (
//	    github.com/gin-gonic/gin v1.9.1
//	    golang.org/x/text v0.14.0 // indirect
//	)
func parseGoMod(content []byte) ([]Dependency, string) {
	var deps []Dependency
	lines := strings.Split(string(content), "\n")

	inRequire := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Start of require block
		if strings.HasPrefix(line, "require (") {
			inRequire = true
			continue
		}
		// End of require block
		if line == ")" {
			inRequire = false
			continue
		}

		// Single line require (e.g., "require github.com/pkg/errors v0.9.1")
		if strings.HasPrefix(line, "require ") && !strings.Contains(line, "(") {
			parts := strings.Fields(line)
			if len(parts) >= 3 {
				deps = append(deps, Dependency{
					Name:    parts[1],
					Version: parts[2],
					Type:    "production",
				})
			}
			continue
		}

		// Inside require block
		if inRequire && line != "" && !strings.HasPrefix(line, "//") {
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				depType := "production"
				// Indirect dependencies are transitive deps
				if strings.Contains(line, "// indirect") {
					depType = "indirect"
				}
				deps = append(deps, Dependency{
					Name:    parts[0],
					Version: parts[1],
					Type:    depType,
				})
			}
		}
	}

	return deps, "go"
}

// parseRequirementsTxt parses a Python requirements.txt file.
// It handles various version specifier formats.
//
// Example requirements.txt:
//
//	flask>=2.0.0
//	requests==2.28.0
//	numpy
func parseRequirementsTxt(content []byte) ([]Dependency, string) {
	var deps []Dependency
	lines := strings.Split(string(content), "\n")

	// Pattern for package with version specifier (e.g., "flask>=2.0.0")
	versionPattern := regexp.MustCompile(`^([a-zA-Z0-9_-]+)\s*([=<>!~]+.*)$`)
	// Pattern for package name only (e.g., "flask")
	simplePattern := regexp.MustCompile(`^([a-zA-Z0-9_-]+)\s*$`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip comments, empty lines, and flags (like -r, -e)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "-") {
			continue
		}

		// Try versioned pattern first
		if matches := versionPattern.FindStringSubmatch(line); len(matches) >= 3 {
			deps = append(deps, Dependency{
				Name:    matches[1],
				Version: strings.TrimSpace(matches[2]),
				Type:    "production",
			})
			continue
		}

		// Try simple pattern (just package name)
		if matches := simplePattern.FindStringSubmatch(line); len(matches) >= 2 {
			deps = append(deps, Dependency{
				Name:    matches[1],
				Version: "*",
				Type:    "production",
			})
		}
	}

	return deps, "python"
}

// parseCargoToml parses a Rust Cargo.toml file.
// It handles [dependencies] and [dev-dependencies] sections.
//
// Example Cargo.toml:
//
//	[dependencies]
//	serde = "1.0"
//
//	[dev-dependencies]
//	tokio-test = "0.4"
func parseCargoToml(content []byte) ([]Dependency, string) {
	var deps []Dependency
	lines := strings.Split(string(content), "\n")

	inDeps := false
	inDevDeps := false

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Section headers
		if line == "[dependencies]" {
			inDeps = true
			inDevDeps = false
			continue
		}
		if line == "[dev-dependencies]" {
			inDeps = false
			inDevDeps = true
			continue
		}
		// Any other section ends dependency parsing
		if strings.HasPrefix(line, "[") {
			inDeps = false
			inDevDeps = false
			continue
		}

		// Parse dependency line (e.g., 'serde = "1.0"')
		if (inDeps || inDevDeps) && strings.Contains(line, "=") {
			parts := strings.SplitN(line, "=", 2)
			if len(parts) == 2 {
				name := strings.TrimSpace(parts[0])
				version := strings.Trim(strings.TrimSpace(parts[1]), "\"'")

				depType := "production"
				if inDevDeps {
					depType = "dev"
				}

				deps = append(deps, Dependency{
					Name:    name,
					Version: version,
					Type:    depType,
				})
			}
		}
	}

	return deps, "rust"
}

// parseGemfile parses a Ruby Gemfile.
// It extracts gem declarations with optional version constraints.
//
// Example Gemfile:
//
//	gem 'rails', '~> 7.0'
//	gem 'puma'
func parseGemfile(content []byte) ([]Dependency, string) {
	var deps []Dependency
	lines := strings.Split(string(content), "\n")

	// Pattern matches: gem 'name' or gem 'name', 'version'
	gemPattern := regexp.MustCompile(`gem\s+['"]([^'"]+)['"](?:\s*,\s*['"]([^'"]+)['"])?`)

	for _, line := range lines {
		line = strings.TrimSpace(line)

		// Skip comments and empty lines
		if strings.HasPrefix(line, "#") || line == "" {
			continue
		}

		if matches := gemPattern.FindStringSubmatch(line); len(matches) >= 2 {
			version := "*"
			if len(matches) >= 3 && matches[2] != "" {
				version = matches[2]
			}

			deps = append(deps, Dependency{
				Name:    matches[1],
				Version: version,
				Type:    "production",
			})
		}
	}

	return deps, "ruby"
}

// cleanVersion normalizes version strings for display.
// It preserves the original format including prefixes like ^, ~, >=.
func cleanVersion(v string) string {
	return strings.TrimSpace(v)
}

// contains checks if a string slice contains a specific item.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
