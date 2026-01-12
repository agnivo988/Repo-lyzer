package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const historyFile = "exports/history.json"
const maxHistoryItems = 50

// HistoryEntry represents a saved analysis
type HistoryEntry struct {
	RepoName      string    `json:"repo_name"`
	AnalyzedAt    time.Time `json:"analyzed_at"`
	HealthScore   int       `json:"health_score"`
	Stars         int       `json:"stars"`
	Forks         int       `json:"forks"`
	MaturityLevel string    `json:"maturity_level"`
}

// History holds all history entries
type History struct {
	Entries []HistoryEntry `json:"entries"`
}

// LoadHistory loads history from file
func LoadHistory() (*History, error) {
	history := &History{Entries: []HistoryEntry{}}

	data, err := os.ReadFile(historyFile)
	if err != nil {
		if os.IsNotExist(err) {
			return history, nil
		}
		return nil, err
	}

	if err := json.Unmarshal(data, history); err != nil {
		return history, nil // Return empty on parse error
	}

	return history, nil
}

// SaveHistory saves history to file
func (h *History) Save() error {
	if err := os.MkdirAll(filepath.Dir(historyFile), 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(h, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(historyFile, data, 0644)
}

// AddEntry adds a new entry to history
func (h *History) AddEntry(data AnalysisResult) {
	entry := HistoryEntry{
		RepoName:      data.Repo.FullName,
		AnalyzedAt:    time.Now(),
		HealthScore:   data.HealthScore,
		Stars:         data.Repo.Stars,
		Forks:         data.Repo.Forks,
		MaturityLevel: data.MaturityLevel,
	}

	// Remove duplicate if exists
	for i, e := range h.Entries {
		if e.RepoName == entry.RepoName {
			h.Entries = append(h.Entries[:i], h.Entries[i+1:]...)
			break
		}
	}

	// Add to front
	h.Entries = append([]HistoryEntry{entry}, h.Entries...)

	// Trim to max size
	if len(h.Entries) > maxHistoryItems {
		h.Entries = h.Entries[:maxHistoryItems]
	}
}

// GetRecent returns the most recent entries
func (h *History) GetRecent(count int) []HistoryEntry {
	if count > len(h.Entries) {
		count = len(h.Entries)
	}
	return h.Entries[:count]
}

// Clear removes all history
func (h *History) Clear() {
	h.Entries = []HistoryEntry{}
}

// Delete removes a specific entry
func (h *History) Delete(index int) {
	if index >= 0 && index < len(h.Entries) {
		h.Entries = append(h.Entries[:index], h.Entries[index+1:]...)
	}
}

// SortByDate sorts entries by date (newest first)
func (h *History) SortByDate() {
	sort.Slice(h.Entries, func(i, j int) bool {
		return h.Entries[i].AnalyzedAt.After(h.Entries[j].AnalyzedAt)
	})
}

// FormatEntry formats a history entry for display
func (e HistoryEntry) Format() string {
	return fmt.Sprintf("%-30s │ ⭐%-6d │ 💚%-3d │ %s │ %s",
		e.RepoName,
		e.Stars,
		e.HealthScore,
		e.MaturityLevel,
		e.AnalyzedAt.Format("2006-01-02 15:04"),
	)
}
