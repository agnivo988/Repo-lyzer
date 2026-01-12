// Package ui provides the terminal user interface for Repo-lyzer.
// This file implements accessibility features and keyboard navigation helpers.
package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// AccessibilityConfig holds accessibility settings
type AccessibilityConfig struct {
	HighContrast     bool   // Use high contrast colors
	LargeText        bool   // Use larger text where possible
	ReducedMotion    bool   // Reduce animations
	ScreenReader     bool   // Optimize for screen readers
	KeyRepeatDelay   int    // Milliseconds before key repeat
	FocusIndicator   string // Character(s) to indicate focus
	AnnounceChanges  bool   // Announce view changes
}

// DefaultAccessibilityConfig returns sensible defaults
func DefaultAccessibilityConfig() AccessibilityConfig {
	return AccessibilityConfig{
		HighContrast:    false,
		LargeText:       false,
		ReducedMotion:   false,
		ScreenReader:    false,
		KeyRepeatDelay:  250,
		FocusIndicator:  "▶",
		AnnounceChanges: true,
	}
}

// KeyBinding represents a keyboard shortcut with context
type KeyBinding struct {
	Key         string   // Primary key
	AltKeys     []string // Alternative keys (e.g., vim bindings)
	Action      string   // Action identifier
	Description string   // Human-readable description
	Context     string   // Where this binding is active
	Category    string   // Grouping category
}

// KeyBindings is the central registry of all keyboard shortcuts
var KeyBindings = []KeyBinding{
	// Navigation - Universal
	{Key: "↑", AltKeys: []string{"k"}, Action: "nav_up", Description: "Move up", Context: "global", Category: "Navigation"},
	{Key: "↓", AltKeys: []string{"j"}, Action: "nav_down", Description: "Move down", Context: "global", Category: "Navigation"},
	{Key: "←", AltKeys: []string{"h"}, Action: "nav_left", Description: "Move left / Previous", Context: "global", Category: "Navigation"},
	{Key: "→", AltKeys: []string{"l"}, Action: "nav_right", Description: "Move right / Next", Context: "global", Category: "Navigation"},
	{Key: "Enter", AltKeys: []string{"Space"}, Action: "select", Description: "Select / Confirm", Context: "global", Category: "Navigation"},
	{Key: "Esc", AltKeys: []string{"q"}, Action: "back", Description: "Go back / Cancel", Context: "global", Category: "Navigation"},
	{Key: "Tab", AltKeys: []string{}, Action: "next_focus", Description: "Next focusable element", Context: "global", Category: "Navigation"},
	{Key: "Shift+Tab", AltKeys: []string{}, Action: "prev_focus", Description: "Previous focusable element", Context: "global", Category: "Navigation"},

	// Quick Jump - Numbers
	{Key: "1", AltKeys: []string{}, Action: "jump_1", Description: "Jump to tab 1 (Overview)", Context: "dashboard", Category: "Quick Jump"},
	{Key: "2", AltKeys: []string{}, Action: "jump_2", Description: "Jump to tab 2 (Repo)", Context: "dashboard", Category: "Quick Jump"},
	{Key: "3", AltKeys: []string{}, Action: "jump_3", Description: "Jump to tab 3 (Languages)", Context: "dashboard", Category: "Quick Jump"},
	{Key: "4", AltKeys: []string{}, Action: "jump_4", Description: "Jump to tab 4 (Activity)", Context: "dashboard", Category: "Quick Jump"},
	{Key: "5", AltKeys: []string{}, Action: "jump_5", Description: "Jump to tab 5 (Contributors)", Context: "dashboard", Category: "Quick Jump"},
	{Key: "6", AltKeys: []string{}, Action: "jump_6", Description: "Jump to tab 6 (Insights)", Context: "dashboard", Category: "Quick Jump"},
	{Key: "7", AltKeys: []string{}, Action: "jump_7", Description: "Jump to tab 7 (Recruiter)", Context: "dashboard", Category: "Quick Jump"},
	{Key: "8", AltKeys: []string{}, Action: "jump_8", Description: "Jump to tab 8 (API Status)", Context: "dashboard", Category: "Quick Jump"},

	// Actions
	{Key: "e", AltKeys: []string{}, Action: "export", Description: "Export results", Context: "dashboard", Category: "Actions"},
	{Key: "f", AltKeys: []string{}, Action: "file_tree", Description: "Open file tree", Context: "dashboard", Category: "Actions"},
	{Key: "r", AltKeys: []string{"F5"}, Action: "refresh", Description: "Refresh data", Context: "dashboard", Category: "Actions"},
	{Key: "b", AltKeys: []string{}, Action: "bookmark", Description: "Toggle bookmark", Context: "dashboard", Category: "Actions"},
	{Key: "t", AltKeys: []string{}, Action: "theme", Description: "Cycle theme", Context: "global", Category: "Actions"},
	{Key: "/", AltKeys: []string{"Ctrl+F"}, Action: "search", Description: "Search", Context: "global", Category: "Actions"},

	// Help
	{Key: "?", AltKeys: []string{"F1"}, Action: "help", Description: "Show help", Context: "global", Category: "Help"},
	{Key: "Ctrl+?", AltKeys: []string{}, Action: "shortcuts", Description: "Show all shortcuts", Context: "global", Category: "Help"},

	// System
	{Key: "Ctrl+C", AltKeys: []string{}, Action: "quit", Description: "Quit application", Context: "global", Category: "System"},
	{Key: "Ctrl+L", AltKeys: []string{}, Action: "clear", Description: "Clear/redraw screen", Context: "global", Category: "System"},
	{Key: "Ctrl+R", AltKeys: []string{}, Action: "hard_refresh", Description: "Force refresh", Context: "global", Category: "System"},

	// Input editing
	{Key: "Ctrl+U", AltKeys: []string{}, Action: "clear_line", Description: "Clear input line", Context: "input", Category: "Editing"},
	{Key: "Ctrl+W", AltKeys: []string{}, Action: "delete_word", Description: "Delete word", Context: "input", Category: "Editing"},
	{Key: "Ctrl+A", AltKeys: []string{"Home"}, Action: "line_start", Description: "Go to line start", Context: "input", Category: "Editing"},
	{Key: "Ctrl+E", AltKeys: []string{"End"}, Action: "line_end", Description: "Go to line end", Context: "input", Category: "Editing"},
	{Key: "Backspace", AltKeys: []string{"Ctrl+H"}, Action: "delete_char", Description: "Delete character", Context: "input", Category: "Editing"},

	// History navigation
	{Key: "d", AltKeys: []string{"Delete"}, Action: "delete", Description: "Delete entry", Context: "history", Category: "History"},
	{Key: "c", AltKeys: []string{}, Action: "clear_all", Description: "Clear all history", Context: "history", Category: "History"},

	// File tree
	{Key: "Enter", AltKeys: []string{}, Action: "open_file", Description: "Open/view file", Context: "tree", Category: "File Tree"},
	{Key: "o", AltKeys: []string{}, Action: "expand", Description: "Expand folder", Context: "tree", Category: "File Tree"},
	{Key: "O", AltKeys: []string{}, Action: "expand_all", Description: "Expand all", Context: "tree", Category: "File Tree"},
	{Key: "c", AltKeys: []string{}, Action: "collapse", Description: "Collapse folder", Context: "tree", Category: "File Tree"},
	{Key: "C", AltKeys: []string{}, Action: "collapse_all", Description: "Collapse all", Context: "tree", Category: "File Tree"},
}

// GetBindingsForContext returns all key bindings for a specific context
func GetBindingsForContext(context string) []KeyBinding {
	var bindings []KeyBinding
	for _, b := range KeyBindings {
		if b.Context == context || b.Context == "global" {
			bindings = append(bindings, b)
		}
	}
	return bindings
}

// GetBindingsByCategory groups bindings by category
func GetBindingsByCategory(context string) map[string][]KeyBinding {
	categories := make(map[string][]KeyBinding)
	for _, b := range KeyBindings {
		if b.Context == context || b.Context == "global" {
			categories[b.Category] = append(categories[b.Category], b)
		}
	}
	return categories
}

// FormatKeyBindingHelp creates a formatted help string for key bindings
func FormatKeyBindingHelp(context string, width int) string {
	categories := GetBindingsByCategory(context)
	
	var sections []string
	
	// Define category order
	categoryOrder := []string{"Navigation", "Quick Jump", "Actions", "Editing", "Help", "System"}
	
	for _, cat := range categoryOrder {
		bindings, exists := categories[cat]
		if !exists || len(bindings) == 0 {
			continue
		}
		
		section := fmt.Sprintf("━━ %s ━━\n", cat)
		for _, b := range bindings {
			keys := b.Key
			if len(b.AltKeys) > 0 {
				keys += " / " + strings.Join(b.AltKeys, " / ")
			}
			section += fmt.Sprintf("  %-20s %s\n", keys, b.Description)
		}
		sections = append(sections, section)
	}
	
	return strings.Join(sections, "\n")
}

// FocusRing manages focus between UI elements
type FocusRing struct {
	elements     []string
	currentIndex int
}

// NewFocusRing creates a new focus ring with the given elements
func NewFocusRing(elements []string) *FocusRing {
	return &FocusRing{
		elements:     elements,
		currentIndex: 0,
	}
}

// Next moves focus to the next element
func (f *FocusRing) Next() string {
	if len(f.elements) == 0 {
		return ""
	}
	f.currentIndex = (f.currentIndex + 1) % len(f.elements)
	return f.elements[f.currentIndex]
}

// Previous moves focus to the previous element
func (f *FocusRing) Previous() string {
	if len(f.elements) == 0 {
		return ""
	}
	f.currentIndex--
	if f.currentIndex < 0 {
		f.currentIndex = len(f.elements) - 1
	}
	return f.elements[f.currentIndex]
}

// Current returns the currently focused element
func (f *FocusRing) Current() string {
	if len(f.elements) == 0 {
		return ""
	}
	return f.elements[f.currentIndex]
}

// SetFocus sets focus to a specific element
func (f *FocusRing) SetFocus(element string) bool {
	for i, e := range f.elements {
		if e == element {
			f.currentIndex = i
			return true
		}
	}
	return false
}

// SkipLink represents a skip navigation link for accessibility
type SkipLink struct {
	Label  string
	Target string
}

// GetSkipLinks returns skip links for the current view
func GetSkipLinks(view string) []SkipLink {
	switch view {
	case "dashboard":
		return []SkipLink{
			{Label: "Skip to content", Target: "content"},
			{Label: "Skip to navigation", Target: "tabs"},
			{Label: "Skip to actions", Target: "actions"},
		}
	case "menu":
		return []SkipLink{
			{Label: "Skip to menu", Target: "menu"},
		}
	default:
		return []SkipLink{
			{Label: "Skip to main content", Target: "main"},
		}
	}
}

// ARIARole represents ARIA roles for screen readers
type ARIARole string

const (
	RoleButton      ARIARole = "button"
	RoleTab         ARIARole = "tab"
	RoleTabPanel    ARIARole = "tabpanel"
	RoleMenu        ARIARole = "menu"
	RoleMenuItem    ARIARole = "menuitem"
	RoleNavigation  ARIARole = "navigation"
	RoleMain        ARIARole = "main"
	RoleStatus      ARIARole = "status"
	RoleAlert       ARIARole = "alert"
	RoleProgressbar ARIARole = "progressbar"
)

// ScreenReaderAnnouncement creates an announcement for screen readers
type ScreenReaderAnnouncement struct {
	Message  string
	Priority string // "polite" or "assertive"
	Role     ARIARole
}

// CreateAnnouncement creates a screen reader announcement
func CreateAnnouncement(message string, assertive bool) ScreenReaderAnnouncement {
	priority := "polite"
	if assertive {
		priority = "assertive"
	}
	return ScreenReaderAnnouncement{
		Message:  message,
		Priority: priority,
		Role:     RoleStatus,
	}
}

// RenderAccessibleLabel renders a label with accessibility info
func RenderAccessibleLabel(label string, shortcut string, focused bool) string {
	style := NormalStyle
	if focused {
		style = SelectedStyle
	}
	
	if shortcut != "" {
		return style.Render(fmt.Sprintf("%s [%s]", label, shortcut))
	}
	return style.Render(label)
}

// RenderFocusIndicator renders a focus indicator
func RenderFocusIndicator(focused bool, config AccessibilityConfig) string {
	if focused {
		return config.FocusIndicator + " "
	}
	return "  "
}

// HighContrastStyle returns high contrast version of a style
func HighContrastStyle(base lipgloss.Style) lipgloss.Style {
	return base.
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color("#000000")).
		Bold(true)
}

// GetContrastColors returns appropriate colors based on accessibility settings
func GetContrastColors(config AccessibilityConfig) (fg, bg lipgloss.Color) {
	if config.HighContrast {
		return lipgloss.Color("#FFFFFF"), lipgloss.Color("#000000")
	}
	return lipgloss.Color("#FAFAFA"), lipgloss.Color("#1a1a2e")
}

// KeyboardHelpOverlay renders a keyboard shortcuts overlay
func KeyboardHelpOverlay(context string, width, height int) string {
	title := TitleStyle.Render("⌨️ Keyboard Shortcuts")
	
	help := FormatKeyBindingHelp(context, width)
	
	footer := SubtleStyle.Render("\nPress ? or ESC to close")
	
	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		help,
		footer,
	)
	
	box := BoxStyle.
		Width(min(width-4, 60)).
		Render(content)
	
	return lipgloss.Place(
		width,
		height,
		lipgloss.Center,
		lipgloss.Center,
		box,
	)
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// QuickNavHint returns a hint string for quick navigation
func QuickNavHint(currentTab int, totalTabs int) string {
	hints := []string{}
	
	if currentTab > 1 {
		hints = append(hints, fmt.Sprintf("← %d", currentTab-1))
	}
	if currentTab < totalTabs {
		hints = append(hints, fmt.Sprintf("%d →", currentTab+1))
	}
	
	return strings.Join(hints, " • ")
}

// NavigationBreadcrumb creates a breadcrumb trail
func NavigationBreadcrumb(path []string) string {
	if len(path) == 0 {
		return ""
	}
	return SubtleStyle.Render(strings.Join(path, " > "))
}
