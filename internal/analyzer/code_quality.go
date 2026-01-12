// Package analyzer provides functions for analyzing GitHub repository data.
// This file implements code quality metrics analysis.
package analyzer

import (
	"path/filepath"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// CodeQualityMetrics contains comprehensive code quality analysis
type CodeQualityMetrics struct {
	OverallScore      int                    `json:"overall_score"`       // 0-100
	Grade             string                 `json:"grade"`               // A, B, C, D, F
	DocumentationScore int                   `json:"documentation_score"` // 0-100
	TestingScore      int                    `json:"testing_score"`       // 0-100
	StructureScore    int                    `json:"structure_score"`     // 0-100
	MaintenanceScore  int                    `json:"maintenance_score"`   // 0-100
	HasReadme         bool                   `json:"has_readme"`
	HasContributing   bool                   `json:"has_contributing"`
	HasLicense        bool                   `json:"has_license"`
	HasChangelog      bool                   `json:"has_changelog"`
	HasCodeOfConduct  bool                   `json:"has_code_of_conduct"`
	HasTests          bool                   `json:"has_tests"`
	HasCI             bool                   `json:"has_ci"`
	HasDocker         bool                   `json:"has_docker"`
	HasEditorConfig   bool                   `json:"has_editorconfig"`
	HasGitignore      bool                   `json:"has_gitignore"`
	TestFrameworks    []string               `json:"test_frameworks"`
	CIProviders       []string               `json:"ci_providers"`
	FileStats         FileStatistics         `json:"file_stats"`
	CodeSmells        []CodeSmell            `json:"code_smells"`
	Recommendations   []string               `json:"recommendations"`
}

// FileStatistics contains file-related metrics
type FileStatistics struct {
	TotalFiles       int            `json:"total_files"`
	SourceFiles      int            `json:"source_files"`
	TestFiles        int            `json:"test_files"`
	DocFiles         int            `json:"doc_files"`
	ConfigFiles      int            `json:"config_files"`
	TestRatio        float64        `json:"test_ratio"`        // test files / source files
	AvgPathDepth     float64        `json:"avg_path_depth"`    // average directory depth
	FilesByExtension map[string]int `json:"files_by_extension"`
	LargestFiles     []string       `json:"largest_files"`     // files with deep paths (potential complexity)
}

// CodeSmell represents a potential code quality issue
type CodeSmell struct {
	Type        string `json:"type"`
	Severity    string `json:"severity"` // Low, Medium, High
	Description string `json:"description"`
	Location    string `json:"location,omitempty"`
}

// AnalyzeCodeQuality performs comprehensive code quality analysis
func AnalyzeCodeQuality(repo *github.Repo, fileTree []github.TreeEntry, languages map[string]int) *CodeQualityMetrics {
	metrics := &CodeQualityMetrics{
		FileStats: FileStatistics{
			FilesByExtension: make(map[string]int),
		},
		TestFrameworks: []string{},
		CIProviders:    []string{},
		CodeSmells:     []CodeSmell{},
	}

	if len(fileTree) == 0 {
		metrics.Grade = "N/A"
		metrics.Recommendations = []string{"No file tree data available"}
		return metrics
	}

	// Analyze file tree
	analyzeFileTree(metrics, fileTree)

	// Check for important files
	checkImportantFiles(metrics, fileTree)

	// Detect test frameworks
	detectTestFrameworks(metrics, fileTree, languages)

	// Detect CI providers
	detectCIProviders(metrics, fileTree)

	// Detect code smells
	detectCodeSmells(metrics, fileTree, repo)

	// Calculate scores
	calculateScores(metrics, repo)

	// Generate recommendations
	generateQualityRecommendations(metrics)

	return metrics
}

func analyzeFileTree(metrics *CodeQualityMetrics, fileTree []github.TreeEntry) {
	var totalDepth int
	
	for _, entry := range fileTree {
		if entry.Type != "blob" {
			continue
		}

		metrics.FileStats.TotalFiles++
		
		// Calculate path depth
		depth := strings.Count(entry.Path, "/")
		totalDepth += depth

		// Track deep paths (potential complexity)
		if depth > 5 {
			metrics.FileStats.LargestFiles = append(metrics.FileStats.LargestFiles, entry.Path)
		}

		// Get extension
		ext := strings.ToLower(filepath.Ext(entry.Path))
		if ext != "" {
			metrics.FileStats.FilesByExtension[ext]++
		}

		// Categorize files
		lowerPath := strings.ToLower(entry.Path)
		
		if isSourceFile(lowerPath) {
			metrics.FileStats.SourceFiles++
		}
		if isTestFile(lowerPath) {
			metrics.FileStats.TestFiles++
		}
		if isDocFile(lowerPath) {
			metrics.FileStats.DocFiles++
		}
		if isConfigFile(lowerPath) {
			metrics.FileStats.ConfigFiles++
		}
	}

	// Calculate averages
	if metrics.FileStats.TotalFiles > 0 {
		metrics.FileStats.AvgPathDepth = float64(totalDepth) / float64(metrics.FileStats.TotalFiles)
	}
	if metrics.FileStats.SourceFiles > 0 {
		metrics.FileStats.TestRatio = float64(metrics.FileStats.TestFiles) / float64(metrics.FileStats.SourceFiles)
	}

	// Limit largest files list
	if len(metrics.FileStats.LargestFiles) > 10 {
		metrics.FileStats.LargestFiles = metrics.FileStats.LargestFiles[:10]
	}
}

func checkImportantFiles(metrics *CodeQualityMetrics, fileTree []github.TreeEntry) {
	for _, entry := range fileTree {
		lowerPath := strings.ToLower(entry.Path)
		baseName := strings.ToLower(filepath.Base(entry.Path))

		// README
		if strings.HasPrefix(baseName, "readme") {
			metrics.HasReadme = true
		}

		// CONTRIBUTING
		if strings.HasPrefix(baseName, "contributing") {
			metrics.HasContributing = true
		}

		// LICENSE
		if strings.HasPrefix(baseName, "license") || baseName == "copying" {
			metrics.HasLicense = true
		}

		// CHANGELOG
		if strings.HasPrefix(baseName, "changelog") || strings.HasPrefix(baseName, "history") || baseName == "news" {
			metrics.HasChangelog = true
		}

		// CODE_OF_CONDUCT
		if strings.Contains(baseName, "code_of_conduct") || strings.Contains(baseName, "code-of-conduct") {
			metrics.HasCodeOfConduct = true
		}

		// .gitignore
		if baseName == ".gitignore" {
			metrics.HasGitignore = true
		}

		// .editorconfig
		if baseName == ".editorconfig" {
			metrics.HasEditorConfig = true
		}

		// Docker
		if strings.Contains(lowerPath, "dockerfile") || baseName == "docker-compose.yml" || baseName == "docker-compose.yaml" {
			metrics.HasDocker = true
		}
	}
}

func detectTestFrameworks(metrics *CodeQualityMetrics, fileTree []github.TreeEntry, languages map[string]int) {
	frameworks := make(map[string]bool)

	for _, entry := range fileTree {
		lowerPath := strings.ToLower(entry.Path)
		baseName := strings.ToLower(filepath.Base(entry.Path))

		// Check for test files
		if isTestFile(lowerPath) {
			metrics.HasTests = true
		}

		// JavaScript/TypeScript
		if baseName == "jest.config.js" || baseName == "jest.config.ts" || baseName == "jest.config.json" {
			frameworks["Jest"] = true
		}
		if baseName == "mocha.opts" || baseName == ".mocharc.js" || baseName == ".mocharc.json" {
			frameworks["Mocha"] = true
		}
		if baseName == "karma.conf.js" {
			frameworks["Karma"] = true
		}
		if baseName == "cypress.json" || baseName == "cypress.config.js" || baseName == "cypress.config.ts" {
			frameworks["Cypress"] = true
		}
		if baseName == "vitest.config.ts" || baseName == "vitest.config.js" {
			frameworks["Vitest"] = true
		}
		if baseName == "playwright.config.ts" || baseName == "playwright.config.js" {
			frameworks["Playwright"] = true
		}

		// Python
		if baseName == "pytest.ini" || baseName == "conftest.py" || baseName == "setup.cfg" {
			if strings.Contains(lowerPath, "pytest") || baseName == "conftest.py" {
				frameworks["pytest"] = true
			}
		}
		if baseName == "tox.ini" {
			frameworks["tox"] = true
		}

		// Go
		if strings.HasSuffix(lowerPath, "_test.go") {
			frameworks["Go testing"] = true
		}

		// Java
		if strings.Contains(lowerPath, "src/test/java") {
			frameworks["JUnit"] = true
		}

		// Ruby
		if baseName == ".rspec" || strings.Contains(lowerPath, "spec/") {
			frameworks["RSpec"] = true
		}

		// Rust
		if strings.Contains(lowerPath, "tests/") && strings.HasSuffix(lowerPath, ".rs") {
			frameworks["Rust tests"] = true
		}

		// PHP
		if baseName == "phpunit.xml" || baseName == "phpunit.xml.dist" {
			frameworks["PHPUnit"] = true
		}
	}

	// Also check by language
	if _, hasGo := languages["Go"]; hasGo {
		for _, entry := range fileTree {
			if strings.HasSuffix(entry.Path, "_test.go") {
				frameworks["Go testing"] = true
				break
			}
		}
	}

	for framework := range frameworks {
		metrics.TestFrameworks = append(metrics.TestFrameworks, framework)
	}
}

func detectCIProviders(metrics *CodeQualityMetrics, fileTree []github.TreeEntry) {
	providers := make(map[string]bool)

	for _, entry := range fileTree {
		lowerPath := strings.ToLower(entry.Path)
		baseName := strings.ToLower(filepath.Base(entry.Path))

		// GitHub Actions
		if strings.HasPrefix(lowerPath, ".github/workflows/") {
			providers["GitHub Actions"] = true
			metrics.HasCI = true
		}

		// Travis CI
		if baseName == ".travis.yml" {
			providers["Travis CI"] = true
			metrics.HasCI = true
		}

		// CircleCI
		if strings.HasPrefix(lowerPath, ".circleci/") {
			providers["CircleCI"] = true
			metrics.HasCI = true
		}

		// GitLab CI
		if baseName == ".gitlab-ci.yml" {
			providers["GitLab CI"] = true
			metrics.HasCI = true
		}

		// Jenkins
		if baseName == "jenkinsfile" {
			providers["Jenkins"] = true
			metrics.HasCI = true
		}

		// Azure Pipelines
		if baseName == "azure-pipelines.yml" {
			providers["Azure Pipelines"] = true
			metrics.HasCI = true
		}

		// Bitbucket Pipelines
		if baseName == "bitbucket-pipelines.yml" {
			providers["Bitbucket Pipelines"] = true
			metrics.HasCI = true
		}

		// AppVeyor
		if baseName == "appveyor.yml" || baseName == ".appveyor.yml" {
			providers["AppVeyor"] = true
			metrics.HasCI = true
		}

		// Drone
		if baseName == ".drone.yml" {
			providers["Drone"] = true
			metrics.HasCI = true
		}
	}

	for provider := range providers {
		metrics.CIProviders = append(metrics.CIProviders, provider)
	}
}

func detectCodeSmells(metrics *CodeQualityMetrics, fileTree []github.TreeEntry, repo *github.Repo) {
	// Check for deep nesting
	deepFiles := 0
	for _, entry := range fileTree {
		if strings.Count(entry.Path, "/") > 6 {
			deepFiles++
		}
	}
	if deepFiles > 10 {
		metrics.CodeSmells = append(metrics.CodeSmells, CodeSmell{
			Type:        "Deep Nesting",
			Severity:    "Medium",
			Description: "Many files with deep directory nesting (>6 levels)",
		})
	}

	// Check for missing documentation
	if !metrics.HasReadme {
		metrics.CodeSmells = append(metrics.CodeSmells, CodeSmell{
			Type:        "Missing Documentation",
			Severity:    "High",
			Description: "No README file found",
		})
	}

	// Check for missing license
	if !metrics.HasLicense {
		metrics.CodeSmells = append(metrics.CodeSmells, CodeSmell{
			Type:        "Missing License",
			Severity:    "High",
			Description: "No LICENSE file found - unclear usage rights",
		})
	}

	// Check for missing tests
	if !metrics.HasTests && metrics.FileStats.SourceFiles > 10 {
		metrics.CodeSmells = append(metrics.CodeSmells, CodeSmell{
			Type:        "No Tests",
			Severity:    "High",
			Description: "No test files detected in a project with many source files",
		})
	}

	// Check for low test ratio
	if metrics.FileStats.TestRatio < 0.1 && metrics.FileStats.SourceFiles > 20 {
		metrics.CodeSmells = append(metrics.CodeSmells, CodeSmell{
			Type:        "Low Test Coverage",
			Severity:    "Medium",
			Description: "Test to source file ratio is below 10%",
		})
	}

	// Check for missing CI
	if !metrics.HasCI && metrics.FileStats.SourceFiles > 5 {
		metrics.CodeSmells = append(metrics.CodeSmells, CodeSmell{
			Type:        "No CI/CD",
			Severity:    "Medium",
			Description: "No continuous integration configuration found",
		})
	}

	// Check for missing .gitignore
	if !metrics.HasGitignore {
		metrics.CodeSmells = append(metrics.CodeSmells, CodeSmell{
			Type:        "Missing .gitignore",
			Severity:    "Low",
			Description: "No .gitignore file - may commit unwanted files",
		})
	}

	// Check for too many file types (potential complexity)
	if len(metrics.FileStats.FilesByExtension) > 20 {
		metrics.CodeSmells = append(metrics.CodeSmells, CodeSmell{
			Type:        "High Complexity",
			Severity:    "Low",
			Description: "Many different file types - may indicate complex or unfocused project",
		})
	}

	// Check for stale repo
	if repo != nil && repo.OpenIssues > 100 {
		metrics.CodeSmells = append(metrics.CodeSmells, CodeSmell{
			Type:        "Issue Backlog",
			Severity:    "Medium",
			Description: "Large number of open issues may indicate maintenance challenges",
		})
	}
}

func calculateScores(metrics *CodeQualityMetrics, repo *github.Repo) {
	// Documentation Score (0-100)
	docScore := 0
	if metrics.HasReadme {
		docScore += 40
	}
	if metrics.HasContributing {
		docScore += 20
	}
	if metrics.HasChangelog {
		docScore += 15
	}
	if metrics.HasCodeOfConduct {
		docScore += 10
	}
	if metrics.FileStats.DocFiles > 0 {
		docScore += 15
	}
	metrics.DocumentationScore = min(docScore, 100)

	// Testing Score (0-100)
	testScore := 0
	if metrics.HasTests {
		testScore += 30
	}
	if len(metrics.TestFrameworks) > 0 {
		testScore += 20
	}
	if metrics.FileStats.TestRatio >= 0.5 {
		testScore += 30
	} else if metrics.FileStats.TestRatio >= 0.2 {
		testScore += 20
	} else if metrics.FileStats.TestRatio >= 0.1 {
		testScore += 10
	}
	if metrics.HasCI {
		testScore += 20
	}
	metrics.TestingScore = min(testScore, 100)

	// Structure Score (0-100)
	structScore := 50 // Base score
	if metrics.HasGitignore {
		structScore += 10
	}
	if metrics.HasEditorConfig {
		structScore += 10
	}
	if metrics.HasLicense {
		structScore += 15
	}
	if metrics.FileStats.AvgPathDepth < 4 {
		structScore += 15
	} else if metrics.FileStats.AvgPathDepth > 6 {
		structScore -= 10
	}
	metrics.StructureScore = max(0, min(structScore, 100))

	// Maintenance Score (0-100)
	maintScore := 50 // Base score
	if metrics.HasCI {
		maintScore += 20
	}
	if metrics.HasDocker {
		maintScore += 10
	}
	if len(metrics.CodeSmells) == 0 {
		maintScore += 20
	} else if len(metrics.CodeSmells) <= 2 {
		maintScore += 10
	} else if len(metrics.CodeSmells) > 5 {
		maintScore -= 20
	}
	if repo != nil && repo.OpenIssues < 20 {
		maintScore += 10
	}
	metrics.MaintenanceScore = max(0, min(maintScore, 100))

	// Overall Score (weighted average)
	metrics.OverallScore = (metrics.DocumentationScore*25 +
		metrics.TestingScore*30 +
		metrics.StructureScore*20 +
		metrics.MaintenanceScore*25) / 100

	// Grade
	switch {
	case metrics.OverallScore >= 90:
		metrics.Grade = "A"
	case metrics.OverallScore >= 80:
		metrics.Grade = "B"
	case metrics.OverallScore >= 70:
		metrics.Grade = "C"
	case metrics.OverallScore >= 60:
		metrics.Grade = "D"
	default:
		metrics.Grade = "F"
	}
}

func generateQualityRecommendations(metrics *CodeQualityMetrics) {
	// Documentation recommendations
	if !metrics.HasReadme {
		metrics.Recommendations = append(metrics.Recommendations, "📝 Add a README.md file to describe your project")
	}
	if !metrics.HasContributing {
		metrics.Recommendations = append(metrics.Recommendations, "🤝 Add CONTRIBUTING.md to guide contributors")
	}
	if !metrics.HasLicense {
		metrics.Recommendations = append(metrics.Recommendations, "⚖️ Add a LICENSE file to clarify usage rights")
	}
	if !metrics.HasChangelog {
		metrics.Recommendations = append(metrics.Recommendations, "📋 Add CHANGELOG.md to track version changes")
	}

	// Testing recommendations
	if !metrics.HasTests {
		metrics.Recommendations = append(metrics.Recommendations, "🧪 Add tests to improve code reliability")
	} else if metrics.FileStats.TestRatio < 0.2 {
		metrics.Recommendations = append(metrics.Recommendations, "📈 Increase test coverage (currently low)")
	}

	// CI/CD recommendations
	if !metrics.HasCI {
		metrics.Recommendations = append(metrics.Recommendations, "🔄 Set up CI/CD (GitHub Actions recommended)")
	}

	// Structure recommendations
	if !metrics.HasGitignore {
		metrics.Recommendations = append(metrics.Recommendations, "🚫 Add .gitignore to exclude build artifacts")
	}
	if !metrics.HasEditorConfig {
		metrics.Recommendations = append(metrics.Recommendations, "⚙️ Add .editorconfig for consistent formatting")
	}

	// Limit recommendations
	if len(metrics.Recommendations) > 5 {
		metrics.Recommendations = metrics.Recommendations[:5]
	}

	if len(metrics.Recommendations) == 0 {
		metrics.Recommendations = append(metrics.Recommendations, "✅ Great job! Your project follows best practices")
	}
}

// Helper functions
func isSourceFile(path string) bool {
	sourceExts := []string{".go", ".js", ".ts", ".jsx", ".tsx", ".py", ".java", ".rb", ".rs", ".c", ".cpp", ".h", ".cs", ".php", ".swift", ".kt", ".scala", ".ex", ".exs"}
	for _, ext := range sourceExts {
		if strings.HasSuffix(path, ext) && !isTestFile(path) {
			return true
		}
	}
	return false
}

func isTestFile(path string) bool {
	lowerPath := strings.ToLower(path)
	testPatterns := []string{"_test.", ".test.", ".spec.", "_spec.", "/test/", "/tests/", "/__tests__/", "/spec/", "test/", "tests/", "__tests__/"}
	for _, pattern := range testPatterns {
		if strings.Contains(lowerPath, pattern) {
			return true
		}
	}
	return false
}

func isDocFile(path string) bool {
	lowerPath := strings.ToLower(path)
	docPatterns := []string{".md", ".rst", ".txt", "/docs/", "/doc/", "readme", "changelog", "contributing", "license"}
	for _, pattern := range docPatterns {
		if strings.Contains(lowerPath, pattern) {
			return true
		}
	}
	return false
}

func isConfigFile(path string) bool {
	lowerPath := strings.ToLower(path)
	configPatterns := []string{".json", ".yaml", ".yml", ".toml", ".ini", ".cfg", ".conf", ".config", "dockerfile", "makefile", ".env"}
	for _, pattern := range configPatterns {
		if strings.Contains(lowerPath, pattern) {
			return true
		}
	}
	return false
}
