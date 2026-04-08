package output

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

func TestSaveCompactJSON(t *testing.T) {
	// Create a temporary directory for test files
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "test-analysis.json")

	// Create test data
	testRepo := &github.Repo{
		FullName:    "test/repo",
		Description: "Test repository",
		HTMLURL:     "https://github.com/test/repo",
		Language:    "Go",
		Stars:       100,
		Forks:       50,
		OpenIssues:  5,
	}

	cfg := CompactConfig{
		Repo:            testRepo,
		HealthScore:     85,
		BusFactor:       3,
		BusRisk:         "Medium",
		MaturityScore:   75,
		MaturityLevel:   "Mature",
		CommitsLastYear: 250,
		Contributors:    10,
		Duration:        2 * time.Second,
		Languages: map[string]int{
			"Go":         80000,
			"JavaScript": 15000,
			"Python":     5000,
		},
	}

	// Test saving JSON
	err := SaveCompactJSON(testFile, cfg)
	if err != nil {
		t.Fatalf("SaveCompactJSON failed: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(testFile); os.IsNotExist(err) {
		t.Fatal("File was not created")
	}

	// Read and verify the JSON content
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]interface{}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}

	// Verify key fields exist
	if _, ok := result["repository"]; !ok {
		t.Error("Missing 'repository' field")
	}
	if _, ok := result["metrics"]; !ok {
		t.Error("Missing 'metrics' field")
	}
	if _, ok := result["metadata"]; !ok {
		t.Error("Missing 'metadata' field")
	}

	t.Logf("Successfully saved and verified JSON to: %s", testFile)
}

func TestSaveCompactJSON_CreateDirectory(t *testing.T) {
	// Create a temporary directory
	tempDir := t.TempDir()
	
	// Try to save to a nested path that doesn't exist
	nestedPath := filepath.Join(tempDir, "subdir1", "subdir2", "analysis.json")

	cfg := CompactConfig{
		Repo: &github.Repo{
			FullName: "test/repo",
		},
		HealthScore: 80,
		Duration:    1 * time.Second,
	}

	// Test saving to nested directory
	err := SaveCompactJSON(nestedPath, cfg)
	if err != nil {
		t.Fatalf("SaveCompactJSON failed to create directories: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(nestedPath); os.IsNotExist(err) {
		t.Fatal("File was not created in nested directory")
	}

	t.Logf("Successfully created nested directories and saved file: %s", nestedPath)
}

func TestSaveCompactJSON_Overwrite(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "overwrite-test.json")

	cfg1 := CompactConfig{
		Repo:        &github.Repo{FullName: "test/repo1"},
		HealthScore: 70,
		Duration:    1 * time.Second,
	}

	cfg2 := CompactConfig{
		Repo:        &github.Repo{FullName: "test/repo2"},
		HealthScore: 90,
		Duration:    1 * time.Second,
	}

	// Save first time
	if err := SaveCompactJSON(testFile, cfg1); err != nil {
		t.Fatalf("First save failed: %v", err)
	}

	// Save second time (overwrite)
	if err := SaveCompactJSON(testFile, cfg2); err != nil {
		t.Fatalf("Overwrite failed: %v", err)
	}

	// Read and verify the content is from cfg2
	data, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	var result compactAnalysis
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Invalid JSON: %v", err)
	}

	if result.Repository.FullName != "test/repo2" {
		t.Errorf("Expected repo name 'test/repo2', got '%s'", result.Repository.FullName)
	}

	if result.Metrics.HealthScore != 90 {
		t.Errorf("Expected health score 90, got %d", result.Metrics.HealthScore)
	}

	t.Log("Successfully verified file overwrite")
}
