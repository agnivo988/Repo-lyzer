package analyzer

import (
	"github.com/agnivo988/Repo-lyzer/internal/cache"
	gh "github.com/agnivo988/Repo-lyzer/internal/github"
)

type IncrementalStats struct {
	CacheHits    int
	CacheMisses  int
	FilesSkipped int
	FilesScanned int
}

func DetectChangedFiles(
	cached map[string]cache.FileMetadata,
	current []gh.TreeEntry,
) ([]gh.TreeEntry, IncrementalStats) {

	changed := []gh.TreeEntry{}
	stats := IncrementalStats{}

	for _, file := range current {

		// Skip directories
		if file.Type != "blob" {
			continue
		}

		stats.FilesScanned++

		old, exists := cached[file.Path]

		if exists && old.SHA == file.Sha {
			stats.CacheHits++
			stats.FilesSkipped++
			continue
		}

		stats.CacheMisses++
		changed = append(changed, file)
	}

	return changed, stats
}
