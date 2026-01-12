package cache

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewCache(t *testing.T) {
	cache, err := NewCache()
	if err != nil {
		t.Fatalf("NewCache() error = %v", err)
	}
	if cache == nil {
		t.Fatal("NewCache() returned nil")
	}
	if cache.cacheDir == "" {
		t.Error("Cache directory not set")
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if !config.Enabled {
		t.Error("Default config should be enabled")
	}
	if config.TTL != 24*time.Hour {
		t.Errorf("Default TTL = %v, want 24h", config.TTL)
	}
	if config.MaxSize != 100 {
		t.Errorf("Default MaxSize = %d, want 100", config.MaxSize)
	}
	if !config.AutoCache {
		t.Error("Default AutoCache should be true")
	}
}

func TestCache_SetAndGet(t *testing.T) {
	cache, err := NewCache()
	if err != nil {
		t.Fatalf("NewCache() error = %v", err)
	}

	testRepo := "test/repo"
	testData := map[string]interface{}{
		"name":  "test",
		"score": 85,
	}

	// Set
	err = cache.Set(testRepo, testData)
	if err != nil {
		t.Fatalf("Set() error = %v", err)
	}

	// Get
	entry, found := cache.Get(testRepo)
	if !found {
		t.Fatal("Get() did not find cached entry")
	}
	if entry == nil {
		t.Fatal("Get() returned nil entry")
	}
	if entry.RepoName != testRepo {
		t.Errorf("Entry.RepoName = %s, want %s", entry.RepoName, testRepo)
	}

	// Cleanup
	cache.Delete(testRepo)
}

func TestCache_Delete(t *testing.T) {
	cache, err := NewCache()
	if err != nil {
		t.Fatalf("NewCache() error = %v", err)
	}

	testRepo := "test/delete-repo"
	cache.Set(testRepo, "test data")

	// Verify it exists
	_, found := cache.Get(testRepo)
	if !found {
		t.Fatal("Entry should exist before delete")
	}

	// Delete
	err = cache.Delete(testRepo)
	if err != nil {
		t.Fatalf("Delete() error = %v", err)
	}

	// Verify it's gone
	_, found = cache.Get(testRepo)
	if found {
		t.Error("Entry should not exist after delete")
	}
}

func TestCache_Clear(t *testing.T) {
	cache, err := NewCache()
	if err != nil {
		t.Fatalf("NewCache() error = %v", err)
	}

	// Add some entries
	cache.Set("test/repo1", "data1")
	cache.Set("test/repo2", "data2")

	// Clear
	err = cache.Clear()
	if err != nil {
		t.Fatalf("Clear() error = %v", err)
	}

	// Verify all gone
	stats := cache.GetStats()
	if stats.TotalRepos != 0 {
		t.Errorf("After Clear(), TotalRepos = %d, want 0", stats.TotalRepos)
	}
}

func TestCache_GetStats(t *testing.T) {
	cache, err := NewCache()
	if err != nil {
		t.Fatalf("NewCache() error = %v", err)
	}

	// Clear first
	cache.Clear()

	// Add entries
	cache.Set("test/stats1", "data1")
	cache.Set("test/stats2", "data2")

	stats := cache.GetStats()
	
	if stats.TotalRepos != 2 {
		t.Errorf("TotalRepos = %d, want 2", stats.TotalRepos)
	}
	if stats.CacheDir == "" {
		t.Error("CacheDir should not be empty")
	}

	// Cleanup
	cache.Clear()
}

func TestCache_HasCache(t *testing.T) {
	cache, err := NewCache()
	if err != nil {
		t.Fatalf("NewCache() error = %v", err)
	}

	testRepo := "test/has-cache"
	
	// Should not exist initially
	if cache.HasCache(testRepo) {
		t.Error("HasCache() should return false for non-existent repo")
	}

	// Add it
	cache.Set(testRepo, "data")

	// Should exist now
	if !cache.HasCache(testRepo) {
		t.Error("HasCache() should return true after Set()")
	}

	// Cleanup
	cache.Delete(testRepo)
}

func TestCache_SetEnabled(t *testing.T) {
	cache, err := NewCache()
	if err != nil {
		t.Fatalf("NewCache() error = %v", err)
	}

	// Disable
	cache.SetEnabled(false)
	config := cache.GetConfig()
	if config.Enabled {
		t.Error("Cache should be disabled")
	}

	// Re-enable
	cache.SetEnabled(true)
	config = cache.GetConfig()
	if !config.Enabled {
		t.Error("Cache should be enabled")
	}
}

func TestFormatTTL(t *testing.T) {
	tests := []struct {
		duration time.Duration
		want     string
	}{
		{24 * time.Hour, "1 day"},
		{48 * time.Hour, "2 days"},
		{time.Hour, "1 hour"},
		{2 * time.Hour, "2 hours"},
		{30 * time.Minute, "30 minutes"},
		{time.Minute, "1 minute"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			got := FormatTTL(tt.duration)
			if got != tt.want {
				t.Errorf("FormatTTL(%v) = %s, want %s", tt.duration, got, tt.want)
			}
		})
	}
}

func TestRepoToFilename(t *testing.T) {
	tests := []struct {
		repoName string
		want     string
	}{
		{"owner/repo", "owner_repo.json"},
		{"user/my-project", "user_my-project.json"},
		{"org/repo.name", "org_repo.name.json"},
	}

	for _, tt := range tests {
		t.Run(tt.repoName, func(t *testing.T) {
			got := repoToFilename(tt.repoName)
			if got != tt.want {
				t.Errorf("repoToFilename(%s) = %s, want %s", tt.repoName, got, tt.want)
			}
		})
	}
}

func TestGetCacheDir(t *testing.T) {
	dir, err := getCacheDir()
	if err != nil {
		t.Fatalf("getCacheDir() error = %v", err)
	}

	home, _ := os.UserHomeDir()
	expected := filepath.Join(home, ".repo-lyzer", "cache")
	
	if dir != expected {
		t.Errorf("getCacheDir() = %s, want %s", dir, expected)
	}
}
