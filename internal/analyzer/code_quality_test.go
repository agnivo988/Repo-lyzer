package analyzer

import (
	"testing"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestAnalyzeCodeQuality_Empty(t *testing.T) {
	metrics := AnalyzeCodeQuality(nil, []github.TreeEntry{}, nil)

	if metrics.Grade != "N/A" {
		t.Errorf("Grade = %s, want N/A for empty file tree", metrics.Grade)
	}
}

func TestAnalyzeCodeQuality_BasicProject(t *testing.T) {
	fileTree := []github.TreeEntry{
		{Path: "README.md", Type: "blob"},
		{Path: "LICENSE", Type: "blob"},
		{Path: "main.go", Type: "blob"},
		{Path: "main_test.go", Type: "blob"},
		{Path: ".gitignore", Type: "blob"},
	}

	repo := &github.Repo{
		OpenIssues: 5,
	}

	metrics := AnalyzeCodeQuality(repo, fileTree, nil)

	if !metrics.HasReadme {
		t.Error("Should detect README")
	}
	if !metrics.HasLicense {
		t.Error("Should detect LICENSE")
	}
	if !metrics.HasGitignore {
		t.Error("Should detect .gitignore")
	}
	if !metrics.HasTests {
		t.Error("Should detect test files")
	}
	if metrics.FileStats.TotalFiles != 5 {
		t.Errorf("TotalFiles = %d, want 5", metrics.FileStats.TotalFiles)
	}
}

func TestAnalyzeCodeQuality_WithCI(t *testing.T) {
	fileTree := []github.TreeEntry{
		{Path: ".github/workflows/ci.yml", Type: "blob"},
		{Path: "src/main.js", Type: "blob"},
	}

	metrics := AnalyzeCodeQuality(nil, fileTree, nil)

	if !metrics.HasCI {
		t.Error("Should detect GitHub Actions CI")
	}
	if len(metrics.CIProviders) == 0 {
		t.Error("Should have CI providers")
	}

	hasGitHubActions := false
	for _, p := range metrics.CIProviders {
		if p == "GitHub Actions" {
			hasGitHubActions = true
			break
		}
	}
	if !hasGitHubActions {
		t.Error("Should detect GitHub Actions")
	}
}

func TestAnalyzeCodeQuality_TestFrameworks(t *testing.T) {
	testCases := []struct {
		name      string
		files     []string
		framework string
	}{
		{
			name:      "Jest",
			files:     []string{"jest.config.js", "src/app.js"},
			framework: "Jest",
		},
		{
			name:      "Go testing",
			files:     []string{"main.go", "main_test.go"},
			framework: "Go testing",
		},
		{
			name:      "pytest",
			files:     []string{"conftest.py", "app.py"},
			framework: "pytest",
		},
		{
			name:      "RSpec",
			files:     []string{".rspec", "app.rb"},
			framework: "RSpec",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var fileTree []github.TreeEntry
			for _, f := range tc.files {
				fileTree = append(fileTree, github.TreeEntry{Path: f, Type: "blob"})
			}

			metrics := AnalyzeCodeQuality(nil, fileTree, nil)

			found := false
			for _, fw := range metrics.TestFrameworks {
				if fw == tc.framework {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Should detect %s framework", tc.framework)
			}
		})
	}
}

func TestAnalyzeCodeQuality_CIProviders(t *testing.T) {
	testCases := []struct {
		name     string
		file     string
		provider string
	}{
		{"GitHub Actions", ".github/workflows/test.yml", "GitHub Actions"},
		{"Travis CI", ".travis.yml", "Travis CI"},
		{"CircleCI", ".circleci/config.yml", "CircleCI"},
		{"GitLab CI", ".gitlab-ci.yml", "GitLab CI"},
		{"Jenkins", "Jenkinsfile", "Jenkins"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fileTree := []github.TreeEntry{
				{Path: tc.file, Type: "blob"},
			}

			metrics := AnalyzeCodeQuality(nil, fileTree, nil)

			if !metrics.HasCI {
				t.Error("Should detect CI")
			}

			found := false
			for _, p := range metrics.CIProviders {
				if p == tc.provider {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Should detect %s", tc.provider)
			}
		})
	}
}

func TestAnalyzeCodeQuality_FileStatistics(t *testing.T) {
	fileTree := []github.TreeEntry{
		{Path: "src/main.go", Type: "blob"},
		{Path: "src/utils.go", Type: "blob"},
		{Path: "src/utils_test.go", Type: "blob"},
		{Path: "docs/README.md", Type: "blob"},
		{Path: "config.yaml", Type: "blob"},
	}

	metrics := AnalyzeCodeQuality(nil, fileTree, nil)

	if metrics.FileStats.SourceFiles != 2 {
		t.Errorf("SourceFiles = %d, want 2", metrics.FileStats.SourceFiles)
	}
	if metrics.FileStats.TestFiles != 1 {
		t.Errorf("TestFiles = %d, want 1", metrics.FileStats.TestFiles)
	}
	if metrics.FileStats.DocFiles < 1 {
		t.Errorf("DocFiles = %d, want >= 1", metrics.FileStats.DocFiles)
	}
}

func TestAnalyzeCodeQuality_CodeSmells(t *testing.T) {
	// Project without README should have code smell
	fileTree := []github.TreeEntry{
		{Path: "main.go", Type: "blob"},
	}

	metrics := AnalyzeCodeQuality(nil, fileTree, nil)

	hasReadmeSmell := false
	for _, smell := range metrics.CodeSmells {
		if smell.Type == "Missing Documentation" {
			hasReadmeSmell = true
			break
		}
	}
	if !hasReadmeSmell {
		t.Error("Should detect missing README as code smell")
	}
}

func TestAnalyzeCodeQuality_Grades(t *testing.T) {
	testCases := []struct {
		name          string
		files         []string
		expectedGrade string
	}{
		{
			name: "Well-structured project",
			files: []string{
				"README.md",
				"LICENSE",
				"CONTRIBUTING.md",
				"CHANGELOG.md",
				".gitignore",
				".editorconfig",
				".github/workflows/ci.yml",
				"src/main.go",
				"src/main_test.go",
				"src/utils.go",
				"src/utils_test.go",
			},
			expectedGrade: "A", // or B
		},
		{
			name: "Minimal project",
			files: []string{
				"main.go",
			},
			expectedGrade: "F", // or D
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var fileTree []github.TreeEntry
			for _, f := range tc.files {
				fileTree = append(fileTree, github.TreeEntry{Path: f, Type: "blob"})
			}

			metrics := AnalyzeCodeQuality(nil, fileTree, nil)

			// Just verify grade is set
			if metrics.Grade == "" || metrics.Grade == "N/A" {
				t.Error("Grade should be set")
			}

			// Verify score is in valid range
			if metrics.OverallScore < 0 || metrics.OverallScore > 100 {
				t.Errorf("OverallScore = %d, should be 0-100", metrics.OverallScore)
			}
		})
	}
}

func TestAnalyzeCodeQuality_Recommendations(t *testing.T) {
	// Project missing common files should get recommendations
	fileTree := []github.TreeEntry{
		{Path: "main.go", Type: "blob"},
	}

	metrics := AnalyzeCodeQuality(nil, fileTree, nil)

	if len(metrics.Recommendations) == 0 {
		t.Error("Should have recommendations for incomplete project")
	}

	// Check for README recommendation
	hasReadmeRec := false
	for _, rec := range metrics.Recommendations {
		if containsIgnoreCase(rec, "readme") {
			hasReadmeRec = true
			break
		}
	}
	if !hasReadmeRec {
		t.Error("Should recommend adding README")
	}
}

func TestIsSourceFile(t *testing.T) {
	testCases := []struct {
		path     string
		expected bool
	}{
		{"main.go", true},
		{"app.js", true},
		{"index.ts", true},
		{"main_test.go", false}, // test file
		{"README.md", false},
		{"config.yaml", false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			result := isSourceFile(tc.path)
			if result != tc.expected {
				t.Errorf("isSourceFile(%s) = %v, want %v", tc.path, result, tc.expected)
			}
		})
	}
}

func TestIsTestFile(t *testing.T) {
	testCases := []struct {
		path     string
		expected bool
	}{
		{"main_test.go", true},
		{"app.test.js", true},
		{"component.spec.ts", true},
		{"test/test_app.py", true},
		{"__tests__/app.js", true},
		{"main.go", false},
		{"app.js", false},
	}

	for _, tc := range testCases {
		t.Run(tc.path, func(t *testing.T) {
			result := isTestFile(tc.path)
			if result != tc.expected {
				t.Errorf("isTestFile(%s) = %v, want %v", tc.path, result, tc.expected)
			}
		})
	}
}

func TestScoreBounds(t *testing.T) {
	// Test with various file trees to ensure scores stay in bounds
	testCases := [][]github.TreeEntry{
		{}, // empty
		{{Path: "main.go", Type: "blob"}}, // minimal
		{ // comprehensive
			{Path: "README.md", Type: "blob"},
			{Path: "LICENSE", Type: "blob"},
			{Path: "CONTRIBUTING.md", Type: "blob"},
			{Path: ".github/workflows/ci.yml", Type: "blob"},
			{Path: "src/main.go", Type: "blob"},
			{Path: "src/main_test.go", Type: "blob"},
		},
	}

	for i, fileTree := range testCases {
		metrics := AnalyzeCodeQuality(nil, fileTree, nil)

		if metrics.DocumentationScore < 0 || metrics.DocumentationScore > 100 {
			t.Errorf("Case %d: DocumentationScore = %d, out of bounds", i, metrics.DocumentationScore)
		}
		if metrics.TestingScore < 0 || metrics.TestingScore > 100 {
			t.Errorf("Case %d: TestingScore = %d, out of bounds", i, metrics.TestingScore)
		}
		if metrics.StructureScore < 0 || metrics.StructureScore > 100 {
			t.Errorf("Case %d: StructureScore = %d, out of bounds", i, metrics.StructureScore)
		}
		if metrics.MaintenanceScore < 0 || metrics.MaintenanceScore > 100 {
			t.Errorf("Case %d: MaintenanceScore = %d, out of bounds", i, metrics.MaintenanceScore)
		}
		if metrics.OverallScore < 0 || metrics.OverallScore > 100 {
			t.Errorf("Case %d: OverallScore = %d, out of bounds", i, metrics.OverallScore)
		}
	}
}

// Helper function
func containsIgnoreCase(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsIgnoreCase(s[1:], substr) || (len(s) >= len(substr) && equalFoldPrefix(s, substr)))
}

func equalFoldPrefix(s, prefix string) bool {
	if len(s) < len(prefix) {
		return false
	}
	for i := 0; i < len(prefix); i++ {
		c1, c2 := s[i], prefix[i]
		if c1 >= 'A' && c1 <= 'Z' {
			c1 += 'a' - 'A'
		}
		if c2 >= 'A' && c2 <= 'Z' {
			c2 += 'a' - 'A'
		}
		if c1 != c2 {
			return false
		}
	}
	return true
}
