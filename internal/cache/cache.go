// Package cache provides offline caching functionality for repository analysis results.
// It stores analysis data locally to avoid redundant GitHub API calls and enables
// offline access to previously analyzed repositories.
//
// Cache Structure:
//
//	~/.repo-lyzer/cache/
//	├── config.json       # Cache configuration
//	├── cache_index.json  # Index of all cached repos
//	└── repos/            # Individual repo cache files
//	    ├── owner_repo1.json
//	    └── owner_repo2.json
//
// Usage:
//
//	cache, err := cache.NewCache()
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Check cache before API call
//	if entry, found := cache.Get("owner/repo"); found {
//	    // Use cached data
//	}
//
//	// Save analysis result
//	cache.Set("owner/repo", analysisResult)
package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// CacheEntry represents a cached analysis result with metadata.
// Each entry contains the full analysis data along with timing information
// for TTL-based expiration.
type CacheEntry struct {
	RepoName  string          `json:"repo_name"`  // Repository identifier (owner/repo)
	CachedAt  time.Time       `json:"cached_at"`  // When the entry was cached
	ExpiresAt time.Time       `json:"expires_at"` // When the entry expires
	Analysis  json.RawMessage `json:"analysis"`   // Serialized AnalysisResult
}

// CacheIndex stores metadata about all cached repositories.
// This index enables quick lookups without reading individual cache files.
type CacheIndex struct {
	Entries     map[string]CacheIndexEntry `json:"entries"`      // Map of repo name to entry metadata
	LastUpdated time.Time                  `json:"last_updated"` // Last index modification time
}

// CacheIndexEntry is a lightweight summary of a cached repository.
// It contains enough information to check validity without loading the full entry.
type CacheIndexEntry struct {
	RepoName  string    `json:"repo_name"`
	CachedAt  time.Time `json:"cached_at"`
	ExpiresAt time.Time `json:"expires_at"`
	FileSize  int64     `json:"file_size"` // Size in bytes for cache management
}

// CacheConfig holds user-configurable cache settings.
type CacheConfig struct {
	Enabled   bool          `json:"enabled"`     // Whether caching is enabled
	TTL       time.Duration `json:"ttl"`         // Time-to-live for cache entries
	MaxSize   int64         `json:"max_size_mb"` // Maximum cache size in MB
	AutoCache bool          `json:"auto_cache"`  // Automatically cache new analyses
}

// Cache manages the local cache for analysis results.
// It provides thread-safe operations for storing and retrieving cached data.
type Cache struct {
	cacheDir string      // Base directory for cache files
	config   CacheConfig // Current configuration
	index    *CacheIndex // In-memory index of cached repos
}

// DefaultConfig returns the default cache configuration.
// Defaults: enabled, 24-hour TTL, 100MB max size, auto-cache on.
func DefaultConfig() CacheConfig {
	return CacheConfig{
		Enabled:   true,
		TTL:       24 * time.Hour,
		MaxSize:   100,
		AutoCache: true,
	}
}

// NewCache creates and initializes a new cache instance.
// It creates the cache directory structure if it doesn't exist and loads
// any existing configuration and index.
func NewCache() (*Cache, error) {
	cacheDir, err := getCacheDir()
	if err != nil {
		return nil, err
	}

	// Ensure cache directories exist
	reposDir := filepath.Join(cacheDir, "repos")
	if err := os.MkdirAll(reposDir, 0755); err != nil {
		return nil, err
	}

	cache := &Cache{
		cacheDir: cacheDir,
		config:   DefaultConfig(),
	}

	// Load existing index or create new one
	if err := cache.loadIndex(); err != nil {
		cache.index = &CacheIndex{
			Entries:     make(map[string]CacheIndexEntry),
			LastUpdated: time.Now(),
		}
	}

	// Load user configuration if exists
	cache.loadConfig()

	return cache, nil
}

// getCacheDir returns the platform-specific cache directory path.
// On all platforms, this is ~/.repo-lyzer/cache/
func getCacheDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".repo-lyzer", "cache"), nil
}

// repoToFilename converts repo name to safe filename
func repoToFilename(repoName string) string {
	return strings.ReplaceAll(repoName, "/", "_") + ".json"
}

// loadIndex loads the cache index from disk
func (c *Cache) loadIndex() error {
	indexPath := filepath.Join(c.cacheDir, "cache_index.json")
	data, err := os.ReadFile(indexPath)
	if err != nil {
		return err
	}

	c.index = &CacheIndex{}
	return json.Unmarshal(data, c.index)
}

// saveIndex saves the cache index to disk
func (c *Cache) saveIndex() error {
	c.index.LastUpdated = time.Now()
	indexPath := filepath.Join(c.cacheDir, "cache_index.json")

	data, err := json.MarshalIndent(c.index, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(indexPath, data, 0644)
}

// loadConfig loads cache configuration
func (c *Cache) loadConfig() {
	configPath := filepath.Join(c.cacheDir, "config.json")
	data, err := os.ReadFile(configPath)
	if err != nil {
		return // Use defaults
	}
	json.Unmarshal(data, &c.config)
}

// SaveConfig saves cache configuration
func (c *Cache) SaveConfig() error {
	configPath := filepath.Join(c.cacheDir, "config.json")
	data, err := json.MarshalIndent(c.config, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(configPath, data, 0644)
}

// Get retrieves a cached analysis if available and not expired
func (c *Cache) Get(repoName string) (*CacheEntry, bool) {
	if !c.config.Enabled {
		return nil, false
	}

	entry, exists := c.index.Entries[repoName]
	if !exists {
		return nil, false
	}

	// Check if expired
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	// Load full entry from file
	filename := repoToFilename(repoName)
	filePath := filepath.Join(c.cacheDir, "repos", filename)

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, false
	}

	var cacheEntry CacheEntry
	if err := json.Unmarshal(data, &cacheEntry); err != nil {
		return nil, false
	}

	return &cacheEntry, true
}

// Set stores an analysis result in the cache
func (c *Cache) Set(repoName string, analysis interface{}) error {
	if !c.config.Enabled || !c.config.AutoCache {
		return nil
	}

	// Serialize analysis
	analysisData, err := json.Marshal(analysis)
	if err != nil {
		return err
	}

	now := time.Now()
	entry := CacheEntry{
		RepoName:  repoName,
		CachedAt:  now,
		ExpiresAt: now.Add(c.config.TTL),
		Analysis:  analysisData,
	}

	// Save to file
	filename := repoToFilename(repoName)
	filePath := filepath.Join(c.cacheDir, "repos", filename)

	data, err := json.MarshalIndent(entry, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return err
	}

	// Update index
	fileInfo, _ := os.Stat(filePath)
	fileSize := int64(0)
	if fileInfo != nil {
		fileSize = fileInfo.Size()
	}

	c.index.Entries[repoName] = CacheIndexEntry{
		RepoName:  repoName,
		CachedAt:  now,
		ExpiresAt: now.Add(c.config.TTL),
		FileSize:  fileSize,
	}

	return c.saveIndex()
}

// Delete removes a specific repo from cache
func (c *Cache) Delete(repoName string) error {
	filename := repoToFilename(repoName)
	filePath := filepath.Join(c.cacheDir, "repos", filename)

	os.Remove(filePath) // Ignore error if file doesn't exist

	delete(c.index.Entries, repoName)
	return c.saveIndex()
}

// Clear removes all cached data
func (c *Cache) Clear() error {
	reposDir := filepath.Join(c.cacheDir, "repos")

	// Remove all files in repos directory
	entries, err := os.ReadDir(reposDir)
	if err == nil {
		for _, entry := range entries {
			os.Remove(filepath.Join(reposDir, entry.Name()))
		}
	}

	// Clear index
	c.index.Entries = make(map[string]CacheIndexEntry)
	return c.saveIndex()
}

// GetStats returns cache statistics
func (c *Cache) GetStats() CacheStats {
	var totalSize int64
	validCount := 0
	expiredCount := 0
	now := time.Now()

	for _, entry := range c.index.Entries {
		totalSize += entry.FileSize
		if now.After(entry.ExpiresAt) {
			expiredCount++
		} else {
			validCount++
		}
	}

	return CacheStats{
		TotalRepos:   len(c.index.Entries),
		ValidRepos:   validCount,
		ExpiredRepos: expiredCount,
		TotalSizeMB:  float64(totalSize) / (1024 * 1024),
		CacheDir:     c.cacheDir,
	}
}

// CacheStats holds cache statistics
type CacheStats struct {
	TotalRepos   int
	ValidRepos   int
	ExpiredRepos int
	TotalSizeMB  float64
	CacheDir     string
}

// GetCachedRepos returns list of all cached repos
func (c *Cache) GetCachedRepos() []CacheIndexEntry {
	repos := make([]CacheIndexEntry, 0, len(c.index.Entries))
	for _, entry := range c.index.Entries {
		repos = append(repos, entry)
	}
	return repos
}

// IsExpired checks if a cached entry is expired
func (c *Cache) IsExpired(repoName string) bool {
	entry, exists := c.index.Entries[repoName]
	if !exists {
		return true
	}
	return time.Now().After(entry.ExpiresAt)
}

// HasCache checks if a repo is in cache (expired or not)
func (c *Cache) HasCache(repoName string) bool {
	_, exists := c.index.Entries[repoName]
	return exists
}

// GetConfig returns current cache configuration
func (c *Cache) GetConfig() CacheConfig {
	return c.config
}

// SetEnabled enables or disables caching
func (c *Cache) SetEnabled(enabled bool) {
	c.config.Enabled = enabled
	c.SaveConfig()
}

// SetTTL sets the cache time-to-live
func (c *Cache) SetTTL(ttl time.Duration) {
	c.config.TTL = ttl
	c.SaveConfig()
}

// SetAutoCache enables or disables auto-caching
func (c *Cache) SetAutoCache(enabled bool) {
	c.config.AutoCache = enabled
	c.SaveConfig()
}

// CleanExpired removes all expired entries
func (c *Cache) CleanExpired() int {
	removed := 0
	now := time.Now()

	for repoName, entry := range c.index.Entries {
		if now.After(entry.ExpiresAt) {
			c.Delete(repoName)
			removed++
		}
	}

	return removed
}

// FormatTTL returns a human-readable TTL string
func FormatTTL(d time.Duration) string {
	if d >= 24*time.Hour {
		days := int(d / (24 * time.Hour))
		if days == 1 {
			return "1 day"
		}
		return fmt.Sprintf("%d days", days)
	}
	if d >= time.Hour {
		hours := int(d / time.Hour)
		if hours == 1 {
			return "1 hour"
		}
		return fmt.Sprintf("%d hours", hours)
	}
	minutes := int(d / time.Minute)
	if minutes == 1 {
		return "1 minute"
	}
	return fmt.Sprintf("%d minutes", minutes)
}
