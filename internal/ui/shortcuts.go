package ui

import (
	"fmt"
	"strings"
)

// KeyboardShortcut represents a single keyboard shortcut
type KeyboardShortcut struct {
	Key         string
	AltKey      string // Alternative key (vim-style)
	Description string
	Category    string
}

// GetMainMenuShortcuts returns shortcuts for the main menu screen
func GetMainMenuShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "↑/↓", AltKey: "j/k", Description: "Navigate menu", Category: "Navigation"},
		{Key: "Home", AltKey: "g", Description: "Jump to first", Category: "Navigation"},
		{Key: "End", AltKey: "G", Description: "Jump to last", Category: "Navigation"},
		{Key: "1-7", AltKey: "", Description: "Quick jump to item", Category: "Navigation"},
		{Key: "Enter", AltKey: "Space", Description: "Select option", Category: "Actions"},
		{Key: "a", AltKey: "", Description: "Quick: Analyze", Category: "Quick Access"},
		{Key: "c", AltKey: "", Description: "Quick: Compare", Category: "Quick Access"},
		{Key: "h", AltKey: "", Description: "Quick: History", Category: "Quick Access"},
		{Key: "s", AltKey: "", Description: "Quick: Settings", Category: "Quick Access"},
		{Key: "?", AltKey: "", Description: "Show help", Category: "Help"},
		{Key: "ESC", AltKey: "q", Description: "Go back / Quit", Category: "System"},
		{Key: "Ctrl+C", AltKey: "", Description: "Force quit", Category: "System"},
	}
}

// GetInputShortcuts returns shortcuts for input screen
func GetInputShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "Enter", AltKey: "", Description: "Submit input", Category: "Actions"},
		{Key: "Backspace", AltKey: "Ctrl+H", Description: "Delete character", Category: "Editing"},
		{Key: "Ctrl+U", AltKey: "", Description: "Clear entire line", Category: "Editing"},
		{Key: "Ctrl+W", AltKey: "", Description: "Delete word", Category: "Editing"},
		{Key: "Ctrl+A", AltKey: "Home", Description: "Go to start", Category: "Editing"},
		{Key: "Ctrl+E", AltKey: "End", Description: "Go to end", Category: "Editing"},
		{Key: "ESC", AltKey: "", Description: "Cancel", Category: "System"},
	}
}

// GetDashboardShortcuts returns shortcuts for dashboard screen
func GetDashboardShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "←/→", AltKey: "h/l", Description: "Switch tabs", Category: "Navigation"},
		{Key: "1-8", AltKey: "", Description: "Jump to tab", Category: "Navigation"},
		{Key: "e", AltKey: "", Description: "Export menu", Category: "Actions"},
		{Key: "j", AltKey: "", Description: "Export JSON", Category: "Actions"},
		{Key: "m", AltKey: "", Description: "Export Markdown", Category: "Actions"},
		{Key: "f", AltKey: "", Description: "File tree", Category: "Actions"},
		{Key: "r", AltKey: "F5", Description: "Refresh data", Category: "Actions"},
		{Key: "b", AltKey: "", Description: "Toggle bookmark", Category: "Actions"},
		{Key: "t", AltKey: "", Description: "Cycle theme", Category: "Display"},
		{Key: "?", AltKey: "", Description: "Show help", Category: "Help"},
		{Key: "ESC", AltKey: "q", Description: "Back to menu", Category: "System"},
	}
}

// GetSettingsShortcuts returns shortcuts for settings screen
func GetSettingsShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "↑/↓", AltKey: "j/k", Description: "Navigate settings", Category: "Navigation"},
		{Key: "Enter", AltKey: "Space", Description: "Toggle option", Category: "Actions"},
		{Key: "1-7", AltKey: "", Description: "Select theme", Category: "Theme"},
		{Key: "t", AltKey: "", Description: "Cycle theme", Category: "Theme"},
		{Key: "e", AltKey: "", Description: "Toggle cache", Category: "Cache"},
		{Key: "a", AltKey: "", Description: "Toggle auto-cache", Category: "Cache"},
		{Key: "c", AltKey: "", Description: "Clear cache", Category: "Cache"},
		{Key: "x", AltKey: "", Description: "Clean expired", Category: "Cache"},
		{Key: "ESC", AltKey: "q", Description: "Go back", Category: "System"},
	}
}

// GetHistoryShortcuts returns shortcuts for history screen
func GetHistoryShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "↑/↓", AltKey: "j/k", Description: "Navigate history", Category: "Navigation"},
		{Key: "Enter", AltKey: "", Description: "Re-analyze repo", Category: "Actions"},
		{Key: "d", AltKey: "Delete", Description: "Delete entry", Category: "Actions"},
		{Key: "c", AltKey: "", Description: "Clear all history", Category: "Actions"},
		{Key: "ESC", AltKey: "q", Description: "Go back", Category: "System"},
	}
}

// GetHelpShortcuts returns shortcuts for help screen
func GetHelpShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "↑/↓", AltKey: "j/k", Description: "Navigate topics", Category: "Navigation"},
		{Key: "←/→", AltKey: "h/l", Description: "Previous/Next topic", Category: "Navigation"},
		{Key: "Enter", AltKey: "", Description: "Select topic", Category: "Actions"},
		{Key: "/", AltKey: "Ctrl+F", Description: "Search help", Category: "Actions"},
		{Key: "ESC", AltKey: "q", Description: "Go back", Category: "System"},
	}
}

// GetFileTreeShortcuts returns shortcuts for file tree viewer
func GetFileTreeShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "↑/↓", AltKey: "j/k", Description: "Navigate files", Category: "Navigation"},
		{Key: "→", AltKey: "l/o", Description: "Expand folder", Category: "Navigation"},
		{Key: "←", AltKey: "h/c", Description: "Collapse folder", Category: "Navigation"},
		{Key: "O", AltKey: "", Description: "Expand all", Category: "Navigation"},
		{Key: "C", AltKey: "", Description: "Collapse all", Category: "Navigation"},
		{Key: "Enter", AltKey: "", Description: "View file", Category: "Actions"},
		{Key: "/", AltKey: "Ctrl+F", Description: "Search files", Category: "Actions"},
		{Key: "ESC", AltKey: "q", Description: "Go back", Category: "System"},
	}
}

// GetUniversalShortcuts returns shortcuts that work everywhere
func GetUniversalShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "?", AltKey: "F1", Description: "Show help", Category: "Help"},
		{Key: "Ctrl+C", AltKey: "", Description: "Quit application", Category: "System"},
		{Key: "Ctrl+L", AltKey: "", Description: "Redraw screen", Category: "System"},
	}
}

// GetShortcutsForScreen returns appropriate shortcuts for a screen
func GetShortcutsForScreen(screenName string) []KeyboardShortcut {
	var shortcuts []KeyboardShortcut

	switch screenName {
	case "menu":
		shortcuts = GetMainMenuShortcuts()
	case "input":
		shortcuts = GetInputShortcuts()
	case "dashboard":
		shortcuts = GetDashboardShortcuts()
	case "settings":
		shortcuts = GetSettingsShortcuts()
	case "history":
		shortcuts = GetHistoryShortcuts()
	case "help":
		shortcuts = GetHelpShortcuts()
	case "tree":
		shortcuts = GetFileTreeShortcuts()
	default:
		shortcuts = GetMainMenuShortcuts()
	}

	return shortcuts
}

// FormatShortcutsForDisplay returns formatted shortcuts as a string
func FormatShortcutsForDisplay(shortcuts []KeyboardShortcut, maxWidth int) string {
	if len(shortcuts) == 0 {
		return ""
	}

	// Group by category
	categories := make(map[string][]KeyboardShortcut)
	categoryOrder := []string{}
	
	for _, sc := range shortcuts {
		if _, exists := categories[sc.Category]; !exists {
			categoryOrder = append(categoryOrder, sc.Category)
		}
		categories[sc.Category] = append(categories[sc.Category], sc)
	}

	var sections []string
	for _, cat := range categoryOrder {
		scs := categories[cat]
		section := fmt.Sprintf("━━ %s ━━\n", cat)
		for _, sc := range scs {
			keys := sc.Key
			if sc.AltKey != "" {
				keys += " / " + sc.AltKey
			}
			section += fmt.Sprintf("  %-18s %s\n", keys, sc.Description)
		}
		sections = append(sections, section)
	}

	return strings.Join(sections, "\n")
}

// FormatShortcutsCompact returns a compact one-line hint
func FormatShortcutsCompact(shortcuts []KeyboardShortcut) string {
	var hints []string
	for _, sc := range shortcuts {
		if sc.Category == "Navigation" || sc.Category == "Actions" {
			hint := sc.Key
			if sc.AltKey != "" {
				hint += "/" + sc.AltKey
			}
			hints = append(hints, hint+": "+sc.Description)
		}
	}
	
	// Limit to first 5 hints
	if len(hints) > 5 {
		hints = hints[:5]
		hints = append(hints, "?: more")
	}
	
	return strings.Join(hints, " • ")
}
