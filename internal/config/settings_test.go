package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func withTempConfigPath(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	configPath := filepath.Join(dir, "settings.json")
	t.Setenv("REPO_LYZER_CONFIG_PATH", configPath)
	t.Setenv("REPO_LYZER_GITHUB_TOKEN", "test-token")

	return configPath
}

func TestDefaultSettingsIncludesCacheDefaults(t *testing.T) {
	settings := DefaultSettings()

	if settings.CacheTTL != 24*time.Hour {
		t.Fatalf("CacheTTL = %v, want %v", settings.CacheTTL, 24*time.Hour)
	}
	if !settings.CacheAutoRefresh {
		t.Fatal("CacheAutoRefresh = false, want true")
	}
}

func TestSettingsSaveAndLoadPreservesCacheFields(t *testing.T) {
	configPath := withTempConfigPath(t)

	settings := DefaultSettings()
	settings.ThemeName = "Nord"
	settings.CacheTTL = 6 * time.Hour
	settings.CacheAutoRefresh = false

	if err := settings.SaveSettings(); err != nil {
		t.Fatalf("SaveSettings() error = %v", err)
	}

	loaded, err := LoadSettings()
	if err != nil {
		t.Fatalf("LoadSettings() error = %v", err)
	}

	if loaded.ThemeName != "Nord" {
		t.Fatalf("ThemeName = %q, want %q", loaded.ThemeName, "Nord")
	}
	if loaded.CacheTTL != 6*time.Hour {
		t.Fatalf("CacheTTL = %v, want %v", loaded.CacheTTL, 6*time.Hour)
	}
	if loaded.CacheAutoRefresh {
		t.Fatal("CacheAutoRefresh = true, want false")
	}

	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("expected settings file to exist at %s: %v", configPath, err)
	}
}

func TestLoadSettingsPreservesDefaultsForLegacyFiles(t *testing.T) {
	configPath := withTempConfigPath(t)

	legacyJSON := []byte(`{
		"theme_name": "Solarized",
		"default_export_format": "json",
		"export_directory": "/tmp/exports",
		"default_analysis_type": "quick",
		"log_level": "debug",
		"monitoring_enabled": true,
		"default_monitor_interval": 300000000000,
		"notification_enabled": true,
		"scheduled_jobs": []
	}`)

	if err := os.WriteFile(configPath, legacyJSON, 0o600); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}

	loaded, err := LoadSettings()
	if err != nil {
		t.Fatalf("LoadSettings() error = %v", err)
	}

	if loaded.ThemeName != "Solarized" {
		t.Fatalf("ThemeName = %q, want %q", loaded.ThemeName, "Solarized")
	}
	if loaded.CacheTTL != 24*time.Hour {
		t.Fatalf("CacheTTL = %v, want default %v for legacy files", loaded.CacheTTL, 24*time.Hour)
	}
	if !loaded.CacheAutoRefresh {
		t.Fatal("CacheAutoRefresh = false, want default true for legacy files")
	}
}
