package ui

import (
	"strings"
	"testing"
)

func TestDefaultAccessibilityConfig(t *testing.T) {
	config := DefaultAccessibilityConfig()

	if config.HighContrast {
		t.Error("HighContrast should be false by default")
	}
	if config.LargeText {
		t.Error("LargeText should be false by default")
	}
	if config.FocusIndicator != "▶" {
		t.Errorf("FocusIndicator = %s, want ▶", config.FocusIndicator)
	}
	if config.KeyRepeatDelay != 250 {
		t.Errorf("KeyRepeatDelay = %d, want 250", config.KeyRepeatDelay)
	}
}

func TestGetBindingsForContext(t *testing.T) {
	// Dashboard context should include dashboard and global bindings
	bindings := GetBindingsForContext("dashboard")

	if len(bindings) == 0 {
		t.Error("Expected bindings for dashboard context")
	}

	// Check that global bindings are included
	hasQuit := false
	for _, b := range bindings {
		if b.Action == "quit" {
			hasQuit = true
			break
		}
	}
	if !hasQuit {
		t.Error("Global 'quit' binding should be included in dashboard context")
	}

	// Check dashboard-specific bindings
	hasExport := false
	for _, b := range bindings {
		if b.Action == "export" {
			hasExport = true
			break
		}
	}
	if !hasExport {
		t.Error("Dashboard 'export' binding should be present")
	}
}

func TestGetBindingsByCategory(t *testing.T) {
	categories := GetBindingsByCategory("dashboard")

	if len(categories) == 0 {
		t.Error("Expected categories for dashboard context")
	}

	// Check that Navigation category exists
	if _, exists := categories["Navigation"]; !exists {
		t.Error("Navigation category should exist")
	}

	// Check that Quick Jump category exists for dashboard
	if _, exists := categories["Quick Jump"]; !exists {
		t.Error("Quick Jump category should exist for dashboard")
	}
}

func TestFormatKeyBindingHelp(t *testing.T) {
	help := FormatKeyBindingHelp("dashboard", 80)

	if help == "" {
		t.Error("Expected non-empty help text")
	}

	// Should contain category headers
	if !strings.Contains(help, "Navigation") {
		t.Error("Help should contain Navigation category")
	}

	// Should contain key descriptions
	if !strings.Contains(help, "Move up") || !strings.Contains(help, "Move down") {
		t.Error("Help should contain navigation descriptions")
	}
}

func TestFocusRing(t *testing.T) {
	elements := []string{"tab1", "tab2", "tab3", "tab4"}
	ring := NewFocusRing(elements)

	// Initial focus should be first element
	if ring.Current() != "tab1" {
		t.Errorf("Initial focus = %s, want tab1", ring.Current())
	}

	// Next should move to second element
	next := ring.Next()
	if next != "tab2" {
		t.Errorf("Next() = %s, want tab2", next)
	}

	// Continue to end
	ring.Next() // tab3
	ring.Next() // tab4

	// Next should wrap to first
	wrapped := ring.Next()
	if wrapped != "tab1" {
		t.Errorf("Wrapped Next() = %s, want tab1", wrapped)
	}

	// Previous should wrap to last
	ring.Previous() // tab4
	if ring.Current() != "tab4" {
		t.Errorf("Previous() from tab1 = %s, want tab4", ring.Current())
	}
}

func TestFocusRing_SetFocus(t *testing.T) {
	elements := []string{"a", "b", "c"}
	ring := NewFocusRing(elements)

	// Set focus to existing element
	found := ring.SetFocus("b")
	if !found {
		t.Error("SetFocus should return true for existing element")
	}
	if ring.Current() != "b" {
		t.Errorf("Current() = %s, want b", ring.Current())
	}

	// Set focus to non-existing element
	found = ring.SetFocus("z")
	if found {
		t.Error("SetFocus should return false for non-existing element")
	}
}

func TestFocusRing_Empty(t *testing.T) {
	ring := NewFocusRing([]string{})

	if ring.Current() != "" {
		t.Error("Empty ring should return empty string")
	}
	if ring.Next() != "" {
		t.Error("Empty ring Next() should return empty string")
	}
	if ring.Previous() != "" {
		t.Error("Empty ring Previous() should return empty string")
	}
}

func TestGetSkipLinks(t *testing.T) {
	// Dashboard should have skip links
	links := GetSkipLinks("dashboard")
	if len(links) == 0 {
		t.Error("Dashboard should have skip links")
	}

	// Check for content skip link
	hasContent := false
	for _, link := range links {
		if link.Target == "content" {
			hasContent = true
			break
		}
	}
	if !hasContent {
		t.Error("Dashboard should have 'Skip to content' link")
	}
}

func TestCreateAnnouncement(t *testing.T) {
	// Polite announcement
	polite := CreateAnnouncement("Tab changed", false)
	if polite.Priority != "polite" {
		t.Errorf("Priority = %s, want polite", polite.Priority)
	}
	if polite.Message != "Tab changed" {
		t.Errorf("Message = %s, want 'Tab changed'", polite.Message)
	}

	// Assertive announcement
	assertive := CreateAnnouncement("Error occurred", true)
	if assertive.Priority != "assertive" {
		t.Errorf("Priority = %s, want assertive", assertive.Priority)
	}
}

func TestRenderFocusIndicator(t *testing.T) {
	config := DefaultAccessibilityConfig()

	// Focused
	focused := RenderFocusIndicator(true, config)
	if !strings.Contains(focused, config.FocusIndicator) {
		t.Error("Focused indicator should contain focus character")
	}

	// Not focused
	notFocused := RenderFocusIndicator(false, config)
	if strings.Contains(notFocused, config.FocusIndicator) {
		t.Error("Not focused indicator should not contain focus character")
	}
}

func TestQuickNavHint(t *testing.T) {
	// Middle tab
	hint := QuickNavHint(3, 5)
	if !strings.Contains(hint, "← 2") || !strings.Contains(hint, "4 →") {
		t.Errorf("Middle tab hint = %s, should show both directions", hint)
	}

	// First tab
	hint = QuickNavHint(1, 5)
	if strings.Contains(hint, "←") {
		t.Error("First tab should not show left arrow")
	}
	if !strings.Contains(hint, "2 →") {
		t.Error("First tab should show right arrow")
	}

	// Last tab
	hint = QuickNavHint(5, 5)
	if strings.Contains(hint, "→") {
		t.Error("Last tab should not show right arrow")
	}
	if !strings.Contains(hint, "← 4") {
		t.Error("Last tab should show left arrow")
	}
}

func TestNavigationBreadcrumb(t *testing.T) {
	// Empty path
	crumb := NavigationBreadcrumb([]string{})
	if crumb != "" {
		t.Error("Empty path should return empty string")
	}

	// Single item
	crumb = NavigationBreadcrumb([]string{"Home"})
	if !strings.Contains(crumb, "Home") {
		t.Error("Single item breadcrumb should contain item")
	}

	// Multiple items
	crumb = NavigationBreadcrumb([]string{"Home", "Dashboard", "Overview"})
	if !strings.Contains(crumb, ">") {
		t.Error("Multi-item breadcrumb should contain separator")
	}
}

func TestKeyBindings_Coverage(t *testing.T) {
	// Ensure all contexts have bindings
	contexts := []string{"global", "dashboard", "input", "history", "tree"}

	for _, ctx := range contexts {
		bindings := GetBindingsForContext(ctx)
		if len(bindings) == 0 {
			t.Errorf("Context %s should have bindings", ctx)
		}
	}
}
