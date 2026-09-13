package pathutil

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// CleanAndResolvePath cleans a path, expands tilde prefix (~), and resolves relative paths to absolute paths.
// It also corrects absolute Unix paths missing a leading slash if that makes the path valid.
func CleanAndResolvePath(path string) (string, error) {
	path = strings.TrimSpace(path)
	if path == "" {
		return "", fmt.Errorf("path cannot be empty")
	}

	if strings.HasPrefix(path, "~") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get user home directory: %w", err)
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~"))
	}

	// Fix common user typos where they omit the leading slash on absolute Unix paths
	if !strings.HasPrefix(path, "/") && !strings.HasPrefix(path, ".") && !strings.HasPrefix(path, "~") {
		slashed := "/" + path
		if stat, err := os.Stat(slashed); err == nil && stat.IsDir() {
			path = slashed
		}
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("failed to resolve absolute path: %w", err)
	}

	return filepath.Clean(absPath), nil
}
