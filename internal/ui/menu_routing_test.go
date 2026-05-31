package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func keyRune(r rune) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{r}}
}

func TestMenuModelCapturesSubmenuSelection(t *testing.T) {
	menu := NewMenuModel()
	menu.cursor = 7
	menu.enterSubmenu()

	model, _ := menu.Update(keyRune('4'))
	updated := model.(MenuModel)

	if !updated.Done {
		t.Fatal("expected submenu selection to complete menu choice")
	}
	if updated.SelectedOption != 7 {
		t.Fatalf("SelectedOption = %d, want 7", updated.SelectedOption)
	}
	if updated.SelectedSubmenuType != "settings" {
		t.Fatalf("SelectedSubmenuType = %q, want settings", updated.SelectedSubmenuType)
	}
	if updated.SelectedSubmenuOption != 3 {
		t.Fatalf("SelectedSubmenuOption = %d, want 3", updated.SelectedSubmenuOption)
	}
}

func TestMainModelRoutesSettingsSubmenu(t *testing.T) {
	model := NewMainModel(nil, nil)

	next, _ := model.Update(keyRune('8'))
	model = next.(MainModel)
	next, _ = model.Update(keyRune('4'))
	model = next.(MainModel)

	if model.state != stateSettings {
		t.Fatalf("state = %v, want stateSettings", model.state)
	}
	if model.settingsOption != "token" {
		t.Fatalf("settingsOption = %q, want token", model.settingsOption)
	}
	if model.settings.settingsOption != "token" {
		t.Fatalf("settings model option = %q, want token", model.settings.settingsOption)
	}
	if !strings.Contains(model.View(), "GitHub API Token Configuration") {
		t.Fatalf("settings view did not render token page:\n%s", model.View())
	}

	next, _ = model.Update(keyRune('i'))
	model = next.(MainModel)

	if !model.settings.inTokenInput {
		t.Fatal("expected token input state to sync into settings view model")
	}
	if !strings.Contains(model.View(), "Enter GitHub Personal Access Token") {
		t.Fatalf("settings view did not render token input mode:\n%s", model.View())
	}
}

func TestMainModelRoutesHelpSubmenu(t *testing.T) {
	model := NewMainModel(nil, nil)

	next, _ := model.Update(keyRune('9'))
	model = next.(MainModel)
	next, _ = model.Update(keyRune('2'))
	model = next.(MainModel)

	if model.state != stateHelp {
		t.Fatalf("state = %v, want stateHelp", model.state)
	}
	if model.helpContent != "getting-started" {
		t.Fatalf("helpContent = %q, want getting-started", model.helpContent)
	}
	if !strings.Contains(model.View(), "Getting Started") {
		t.Fatalf("help view did not render getting started page:\n%s", model.View())
	}
}

func TestMainModelPreservesAnalyzeSubmenuChoice(t *testing.T) {
	model := NewMainModel(nil, nil)

	next, _ := model.Update(keyRune('1'))
	model = next.(MainModel)
	next, _ = model.Update(keyRune('2'))
	model = next.(MainModel)

	if model.state != stateInput {
		t.Fatalf("state = %v, want stateInput", model.state)
	}
	if model.analysisType != "detailed" {
		t.Fatalf("analysisType = %q, want detailed", model.analysisType)
	}
	if model.loading.analysisType != "detailed" {
		t.Fatalf("loading.analysisType = %q, want detailed", model.loading.analysisType)
	}
}
