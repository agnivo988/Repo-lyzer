package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// CompactConfig holds the measurements used to build the compact JSON export.
type CompactConfig struct {
	Repo            *github.Repo
	HealthScore     int
	BusFactor       int
	BusRisk         string
	MaturityScore   int
	MaturityLevel   string
	CommitsLastYear int
	Contributors    int
	Duration        time.Duration
	Languages       map[string]int
}

// PrintCompactJSON writes the compact analysis summary to stdout.
func PrintCompactJSON(cfg CompactConfig) error {
	if cfg.Repo == nil {
		cfg.Repo = &github.Repo{}
	}

	topLangs := buildTopLanguages(cfg.Languages, 3)
	primaryLanguage := cfg.Repo.Language
	if primaryLanguage == "" && len(topLangs) > 0 {
		primaryLanguage = topLangs[0].Name
	}

	summary := compactAnalysis{
		Repository: compactRepository{
			FullName:        cfg.Repo.FullName,
			Description:     cfg.Repo.Description,
			URL:             cfg.Repo.HTMLURL,
			PrimaryLanguage: primaryLanguage,
			Stars:           cfg.Repo.Stars,
			Forks:           cfg.Repo.Forks,
			OpenIssues:      cfg.Repo.OpenIssues,
		},
		Metrics: compactMetrics{
			HealthScore:     cfg.HealthScore,
			BusFactor:       cfg.BusFactor,
			BusRisk:         cfg.BusRisk,
			MaturityScore:   cfg.MaturityScore,
			MaturityLevel:   cfg.MaturityLevel,
			CommitsLastYear: cfg.CommitsLastYear,
			Contributors:    cfg.Contributors,
		},
		Metadata: compactMetadata{
			DurationSeconds: cfg.Duration.Seconds(),
			TopLanguages:    topLangs,
		},
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	return encoder.Encode(summary)
}

type compactAnalysis struct {
	Repository compactRepository `json:"repository"`
	Metrics    compactMetrics    `json:"metrics"`
	Metadata   compactMetadata   `json:"metadata"`
}

type compactRepository struct {
	FullName        string `json:"full_name"`
	Description     string `json:"description,omitempty"`
	URL             string `json:"url"`
	PrimaryLanguage string `json:"primary_language,omitempty"`
	Stars           int    `json:"stars"`
	Forks           int    `json:"forks"`
	OpenIssues      int    `json:"open_issues"`
}

type compactMetrics struct {
	HealthScore     int    `json:"health_score"`
	BusFactor       int    `json:"bus_factor"`
	BusRisk         string `json:"bus_risk"`
	MaturityScore   int    `json:"maturity_score"`
	MaturityLevel   string `json:"maturity_level"`
	CommitsLastYear int    `json:"commit_count_1y"`
	Contributors    int    `json:"contributors"`
}

type compactMetadata struct {
	DurationSeconds float64           `json:"analysis_duration_seconds"`
	TopLanguages    []compactLanguage `json:"top_languages,omitempty"`
}

type compactLanguage struct {
	Name       string  `json:"name"`
	Percentage float64 `json:"percentage"`
}

func buildTopLanguages(langs map[string]int, limit int) []compactLanguage {
	if limit <= 0 || len(langs) == 0 {
		return nil
	}

	type langEntry struct {
		name string
		size int
	}

	entries := make([]langEntry, 0, len(langs))
	total := 0
	for name, size := range langs {
		if size <= 0 {
			continue
		}
		total += size
		entries = append(entries, langEntry{name: name, size: size})
	}

	if total == 0 || len(entries) == 0 {
		return nil
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].size > entries[j].size
	})

	if limit > len(entries) {
		limit = len(entries)
	}

	top := make([]compactLanguage, 0, limit)
	for i := 0; i < limit; i++ {
		entry := entries[i]
		percentage := float64(entry.size) / float64(total) * 100
		top = append(top, compactLanguage{Name: entry.name, Percentage: percentage})
	}

	return top
}

// SaveCompactJSON writes the compact analysis summary to a file.
func SaveCompactJSON(filePath string, cfg CompactConfig) error {
	if cfg.Repo == nil {
		cfg.Repo = &github.Repo{}
	}

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if dir != "." && dir != "" {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}

	topLangs := buildTopLanguages(cfg.Languages, 3)
	primaryLanguage := cfg.Repo.Language
	if primaryLanguage == "" && len(topLangs) > 0 {
		primaryLanguage = topLangs[0].Name
	}

	summary := compactAnalysis{
		Repository: compactRepository{
			FullName:        cfg.Repo.FullName,
			Description:     cfg.Repo.Description,
			URL:             cfg.Repo.HTMLURL,
			PrimaryLanguage: primaryLanguage,
			Stars:           cfg.Repo.Stars,
			Forks:           cfg.Repo.Forks,
			OpenIssues:      cfg.Repo.OpenIssues,
		},
		Metrics: compactMetrics{
			HealthScore:     cfg.HealthScore,
			BusFactor:       cfg.BusFactor,
			BusRisk:         cfg.BusRisk,
			MaturityScore:   cfg.MaturityScore,
			MaturityLevel:   cfg.MaturityLevel,
			CommitsLastYear: cfg.CommitsLastYear,
			Contributors:    cfg.Contributors,
		},
		Metadata: compactMetadata{
			DurationSeconds: cfg.Duration.Seconds(),
			TopLanguages:    topLangs,
		},
	}

	// Create the file
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Write JSON to file
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(summary); err != nil {
		return fmt.Errorf("failed to encode JSON: %w", err)
	}

	return nil
}
