// Package ui provides favorites/bookmarks functionality for repositories.
package ui

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// Favorite represents a bookmarked repository
type Favorite struct {
	RepoName  string    `json:"repo_name"`
	AddedAt   time.Time `json:"added_at"`
	LastUsed  time.Time `json:"last_used"`
	UseCount  int       `json:"use_count"`
	Notes     string    `json:"notes,omitempty"`
}

// Favorites manages the list of favorite repositories
type Favorites struct {
	Items []Favorite `json:"items"`
}

// LoadFavorites loads favorites from disk
func LoadFavorites() (*Favorites, error) {
	path, err := getFavoritesPath()
	if err != nil {
		return &Favorites{Items: []Favorite{}}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return &Favorites{Items: []Favorite{}}, nil // Return empty if file doesn't exist
	}

	var favs Favorites
	if err := json.Unmarshal(data, &favs); err != nil {
		return &Favorites{Items: []Favorite{}}, err
	}

	return &favs, nil
}

// Save saves favorites to disk
func (f *Favorites) Save() error {
	path, err := getFavoritesPath()
	if err != nil {
		return err
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// Add adds a repository to favorites
func (f *Favorites) Add(repoName string) {
	// Check if already exists
	for i, fav := range f.Items {
		if fav.RepoName == repoName {
			// Update existing
			f.Items[i].LastUsed = time.Now()
			f.Items[i].UseCount++
			return
		}
	}

	// Add new
	f.Items = append([]Favorite{{
		RepoName: repoName,
		AddedAt:  time.Now(),
		LastUsed: time.Now(),
		UseCount: 1,
	}}, f.Items...)
}

// Remove removes a repository from favorites
func (f *Favorites) Remove(repoName string) {
	for i, fav := range f.Items {
		if fav.RepoName == repoName {
			f.Items = append(f.Items[:i], f.Items[i+1:]...)
			return
		}
	}
}

// IsFavorite checks if a repository is in favorites
func (f *Favorites) IsFavorite(repoName string) bool {
	for _, fav := range f.Items {
		if fav.RepoName == repoName {
			return true
		}
	}
	return false
}

// UpdateUsage updates the last used time and count for a favorite
func (f *Favorites) UpdateUsage(repoName string) {
	for i, fav := range f.Items {
		if fav.RepoName == repoName {
			f.Items[i].LastUsed = time.Now()
			f.Items[i].UseCount++
			return
		}
	}
}

// GetTopFavorites returns the most used favorites
func (f *Favorites) GetTopFavorites(limit int) []Favorite {
	if len(f.Items) <= limit {
		return f.Items
	}
	return f.Items[:limit]
}

// Clear removes all favorites
func (f *Favorites) Clear() {
	f.Items = []Favorite{}
}

func getFavoritesPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".repo-lyzer", "favorites.json"), nil
}
