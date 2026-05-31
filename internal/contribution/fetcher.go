package contribution

import (
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/github"
)

// CheckContributingFile checks if a contributing guide exists in the file tree.
func CheckContributingFile(fileTree []github.TreeEntry) bool {
	contributingNames := []string{
		"CONTRIBUTING",
		"CONTRIBUTING.MD",
		"CONTRIBUTING.TXT",
		"CONTRIBUTING.RST",
	}

	for _, entry := range fileTree {
		if entry.Type != "blob" {
			continue
		}

		// Get filename from path and uppercase it
		parts := strings.Split(entry.Path, "/")
		filename := strings.ToUpper(parts[len(parts)-1])

		for _, name := range contributingNames {
			if filename == name {
				return true
			}
		}
	}
	return false
}

// FindReadmePath searches for a readme file in the file tree.
func FindReadmePath(fileTree []github.TreeEntry) string {
	readmeNames := []string{
		"README",
		"README.MD",
		"README.TXT",
		"README.RST",
	}

	for _, entry := range fileTree {
		if entry.Type != "blob" {
			continue
		}

		// Get filename from path and uppercase it
		parts := strings.Split(entry.Path, "/")
		filename := strings.ToUpper(parts[len(parts)-1])

		for _, name := range readmeNames {
			if filename == name {
				return entry.Path
			}
		}
	}
	return ""
}
