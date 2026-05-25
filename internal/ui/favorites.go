package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

type Favorite struct {
	RepoName string    `json:"repo_name"`
	UseCount int       `json:"use_count"`
	LastUsed time.Time `json:"last_used"`
	AddedAt  time.Time `json:"added_at"`
	Notes    string    `json:"notes"`
}

type FavoriteItem = Favorite

type Favorites struct {
	Items []Favorite `json:"items"`
}

type FavoritesModel struct {
	favorites *Favorites
	cursor    int
	width     int
	height    int
}

func NewFavoritesModel() *FavoritesModel {
	return &FavoritesModel{
		favorites: &Favorites{Items: []Favorite{}},
		cursor:    0,
	}
}

func (f *Favorites) Add(repoName string) {
	for i, item := range f.Items {
		if item.RepoName == repoName {
			f.Items[i].UseCount++
			f.Items[i].LastUsed = time.Now()
			return
		}
	}
	// Add new favorite
	f.Items = append(f.Items, Favorite{
		RepoName: repoName,
		UseCount: 1,
		LastUsed: time.Now(),
		AddedAt:  time.Now(),
	})
}

func (f *Favorites) Remove(repoName string) {
	for i, item := range f.Items {
		if item.RepoName == repoName {
			f.Items = append(f.Items[:i], f.Items[i+1:]...)
			return
		}
	}
}

func (f *Favorites) UpdateUsage(repoName string) {
	for i, item := range f.Items {
		if item.RepoName == repoName {
			f.Items[i].UseCount++
			f.Items[i].LastUsed = time.Now()
			return
		}
	}
}

func (f *FavoritesModel) Remove(repoName string) {
	f.favorites.Remove(repoName)
}

func (f *FavoritesModel) UpdateUsage(repoName string) {
	f.favorites.UpdateUsage(repoName)
}

func (f *FavoritesModel) IsFavorite(repoName string) bool {
	for _, item := range f.favorites.Items {
		if item.RepoName == repoName {
			return true
		}
	}
	return false
}

func (f *FavoritesModel) GetTopFavorites(limit int) []FavoriteItem {
	if limit <= 0 {
		return []FavoriteItem{}
	}
	if limit >= len(f.favorites.Items) {
		return f.favorites.Items
	}
	return f.favorites.Items[:limit]
}

func (f *FavoritesModel) Clear() {
	f.favorites.Items = []FavoriteItem{}
}

func (f *FavoritesModel) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	favoritesPath := filepath.Join(home, ".repo-lyzer", "favorites.json")
	data, err := json.MarshalIndent(f.favorites, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(favoritesPath, data, 0644)
}

func LoadFavorites() (*FavoritesModel, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return NewFavoritesModel(), err
	}

	filePath := filepath.Join(home, ".repo-lyzer", "favorites.json")
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return NewFavoritesModel(), nil
		}
		return NewFavoritesModel(), err
	}

	var favorites FavoritesModel
	if err := json.Unmarshal(data, &favorites); err != nil {
		return NewFavoritesModel(), err
	}

	// Sort by last used (most recent first)
	sort.Slice(favorites.favorites.Items, func(i, j int) bool {
		return favorites.favorites.Items[i].LastUsed.After(favorites.favorites.Items[j].LastUsed)
	})

	return &favorites, nil
}

func (m FavoritesModel) View(width, height int) string {
	header := TitleStyle.Render("⭐ Favorite Repositories")

	if m.favorites == nil || len(m.favorites.Items) == 0 {
		content := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			BoxStyle.Render("No favorites yet!\n\nAnalyze a repository and press 'b' to bookmark it."),
			SubtleStyle.Render("a: add new • q/ESC: back to menu"),
		)

		return lipgloss.Place(
			width, height,
			lipgloss.Center, lipgloss.Center,
			content,
		)
	}

	// Build favorites list
	var lines []string
	lines = append(lines, fmt.Sprintf("%-35s │ %-10s │ %s", "Repository", "Uses", "Last Used"))
	lines = append(lines, lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Render(strings.Repeat("─", 65)))

	for i, fav := range m.favorites.Items {
		prefix := "  "
		if i == m.cursor {
			prefix = "▶ "
		}
		line := fmt.Sprintf("%s%-33s │ %-10d │ %s",
			prefix,
			fav.RepoName,
			fav.UseCount,
			fav.LastUsed.Format("2006-01-02"),
		)
		if i == 0 {
			lines = append(lines, SelectedStyle.Render(line))
		} else {
			lines = append(lines, line)
		}
	}

	tableBox := BoxStyle.Render(strings.Join(lines, "\n"))
	footer := SubtleStyle.Render("↑↓: navigate • Enter: analyze • d: remove • a: add new • q/ESC: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		tableBox,
		footer,
	)

	if m.width == 0 {
		return content
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		content,
	)
}
