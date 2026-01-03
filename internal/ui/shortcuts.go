package ui

// KeyboardShortcut represents a single keyboard shortcut
type KeyboardShortcut struct {
	Key         string
	Description string
}

// GetMainMenuShortcuts returns shortcuts for the main menu screen
func GetMainMenuShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "↑/↓ or j/k", Description: "Navigate menu"},
		{Key: "Enter", Description: "Select option"},
		{Key: "h", Description: "Show help"},
		{Key: "s", Description: "Settings"},
		{Key: "ESC", Description: "Go back"},
		{Key: "q or Ctrl+C", Description: "Quit application"},
	}
}

// GetInputShortcuts returns shortcuts for input screen
func GetInputShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "Enter", Description: "Analyze repository"},
		{Key: "Backspace", Description: "Delete character"},
		{Key: "Ctrl+U", Description: "Clear input"},
		{Key: "Ctrl+A", Description: "Go to start"},
		{Key: "Ctrl+E", Description: "Go to end"},
		{Key: "ESC", Description: "Cancel"},
	}
}

// GetDashboardShortcuts returns shortcuts for dashboard screen
func GetDashboardShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "e", Description: "Export results"},
		{Key: "↑/↓ or j/k", Description: "Navigate export menu"},
		{Key: "Enter", Description: "Select export format"},
		{Key: "t", Description: "Toggle theme"},
		{Key: "f", Description: "Show file tree"},
		{Key: "q or ESC", Description: "Back to menu"},
	}
}

// GetSettingsShortcuts returns shortcuts for settings screen
func GetSettingsShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "↑/↓ or j/k", Description: "Navigate settings"},
		{Key: "Enter/Space", Description: "Toggle option"},
		{Key: "→/← or h/l", Description: "Change value"},
		{Key: "r", Description: "Reset to defaults"},
		{Key: "s", Description: "Save settings"},
		{Key: "ESC", Description: "Go back"},
	}
}

// GetHistoryShortcuts returns shortcuts for history screen
func GetHistoryShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "↑/↓ or j/k", Description: "Navigate history"},
		{Key: "Enter", Description: "Re-analyze repository"},
		{Key: "d", Description: "Delete entry"},
		{Key: "c", Description: "Clear history"},
		{Key: "ESC", Description: "Go back"},
	}
}

// GetHelpShortcuts returns shortcuts for help screen
func GetHelpShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "↑/↓ or j/k", Description: "Navigate topics"},
		{Key: "→/← or h/l", Description: "Next/Previous topic"},
		{Key: "Enter", Description: "Select topic"},
		{Key: "/", Description: "Search help"},
		{Key: "ESC", Description: "Go back"},
	}
}

// GetFileTreeShortcuts returns shortcuts for file tree viewer
func GetFileTreeShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "↑/↓ or j/k", Description: "Navigate files"},
		{Key: "→ or l", Description: "Expand folder"},
		{Key: "← or h", Description: "Collapse folder"},
		{Key: "Enter", Description: "View file details"},
		{Key: "Ctrl+S", Description: "Search files"},
		{Key: "ESC", Description: "Go back"},
	}
}

// GetUniversalShortcuts returns shortcuts that work everywhere
func GetUniversalShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{Key: "?", Description: "Show this help"},
		{Key: "Ctrl+C", Description: "Quit application"},
		{Key: "Ctrl+L", Description: "Clear screen"},
		{Key: "Ctrl+R", Description: "Refresh"},
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

	// Always append universal shortcuts
	shortcuts = append(shortcuts, GetUniversalShortcuts()...)
	return shortcuts
}

// FormatShortcutsForDisplay returns formatted shortcuts as a string
func FormatShortcutsForDisplay(shortcuts []KeyboardShortcut, maxWidth int) string {
	if len(shortcuts) == 0 {
		return ""
	}

	result := ""
	for i, sc := range shortcuts {
		if i > 0 && i%3 == 0 { // 3 shortcuts per line
			result += "\n"
		} else if i > 0 {
			result += " • "
		}
		result += sc.Key + ": " + sc.Description
	}

	return result
}
