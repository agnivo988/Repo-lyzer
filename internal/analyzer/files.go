package analyzer
 feat/code-search-filter-by-filetype

import (
	"path/filepath"
	"strings"
)

// FilterFilesByExtension filters file paths by given extensions.
// If no extensions are provided, it returns the original file list.
func FilterFilesByExtension(files []string, extensions []string) []string {
	// Default behavior: no filtering
	if len(extensions) == 0 {
		return files
	}

	// Normalize extensions (lowercase, ensure leading dot)
	extSet := make(map[string]struct{})
	for _, ext := range extensions {
		ext = strings.ToLower(strings.TrimSpace(ext))
		if ext == "" {
			continue
		}
		if !strings.HasPrefix(ext, ".") {
			ext = "." + ext
		}
		extSet[ext] = struct{}{}
	}

	var filtered []string
	for _, file := range files {
		fileExt := strings.ToLower(filepath.Ext(file))
		if _, ok := extSet[fileExt]; ok {
			filtered = append(filtered, file)
		}
	}

	return filtered
}

