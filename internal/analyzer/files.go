package analyzer

import (
	"path/filepath"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// ExtensionFilter restricts file analysis to a set of allowed extensions.
type ExtensionFilter struct {
	allowed map[string]struct{}
}

// ParseExtensionFilter converts a comma-separated list of extensions into a filter.
// Both "go" and ".go" formats are supported, and matching is case-insensitive.
func ParseExtensionFilter(value string) ExtensionFilter {
	filter := ExtensionFilter{allowed: make(map[string]struct{})}
	for _, fragment := range strings.Split(value, ",") {
		ext := strings.ToLower(strings.TrimSpace(fragment))
		if ext == "" {
			continue
		}
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		filter.allowed[ext] = struct{}{}
	}
	return filter
}

// IsActive reports whether the filter has any extensions defined.
func (f ExtensionFilter) IsActive() bool {
	return len(f.allowed) > 0
}

// Matches returns true when the file path ends in one of the allowed extensions.
func (f ExtensionFilter) Matches(path string) bool {
	if !f.IsActive() {
		return true
	}
	ext := strings.ToLower(filepath.Ext(path))
	_, ok := f.allowed[ext]
	return ok
}

// FilterTreeEntriesByExtensions keeps only blob entries that match the filter.
func FilterTreeEntriesByExtensions(entries []github.TreeEntry, filter ExtensionFilter) []github.TreeEntry {
	if !filter.IsActive() {
		return entries
	}

	filtered := make([]github.TreeEntry, 0, len(entries))
	for _, entry := range entries {
		if entry.Type != "blob" {
			continue
		}
		if filter.Matches(entry.Path) {
			filtered = append(filtered, entry)
		}
	}

	return filtered
}

// LanguageSizesFromTree estimates language sizes by summing blob sizes per extension.
// Files without an extension are counted under "unknown".
func LanguageSizesFromTree(entries []github.TreeEntry) map[string]int {
	sizes := make(map[string]int)
	for _, entry := range entries {
		if entry.Type != "blob" {
			continue
		}
		ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(entry.Path), "."))
		if ext == "" {
			ext = "unknown"
		}
		sizes[ext] += entry.Size
	}
	return sizes
}
