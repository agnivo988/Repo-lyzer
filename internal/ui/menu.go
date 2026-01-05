package ui

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type MenuModel struct {
	cursor         int
	choices        []string
	SelectedOption int
	Done           bool
	width          int
	height         int
	inSubmenu      bool
	submenuType    string
	submenuCursor  int
	submenuChoices []string
	parentCursor   int
}

type SubmenuOption struct {
	Label  string
	Action string
}

func NewMenuModel() MenuModel {
	return MenuModel{
		choices: []string{
			"📊 Analyze Repository",
			"🔄 Compare Repositories",
			"📜 View History",
			"⚙️ Settings",
			"❓ Help",
			"🚪 Exit",
		},
		inSubmenu: false,
	}
}

func (m MenuModel) Init() tea.Cmd { return nil }

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.inSubmenu {
				if m.submenuCursor > 0 {
					m.submenuCursor--
				}
			} else {
				if m.cursor > 0 {
					m.cursor--
				}
			}
		case "down", "j":
			if m.inSubmenu {
				if m.submenuCursor < len(m.submenuChoices)-1 {
					m.submenuCursor++
				}
			} else {
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				}
			}
		case "enter":
			if m.inSubmenu {
				m.SelectedOption = m.cursor
				m.Done = true
				m.inSubmenu = false
			} else {
				m.enterSubmenu()
			}
		case "esc":
			if m.inSubmenu {
				m.inSubmenu = false
				m.submenuCursor = 0
				m.submenuChoices = nil
				m.submenuType = ""
			}
		case "q":
			if !m.inSubmenu {
				m.SelectedOption = 5 // Exit
				m.Done = true
			}
		}
	}
	return m, nil
}

func (m *MenuModel) enterSubmenu() {
	switch m.cursor {
	case 0: // Analyze Repository
		m.submenuType = "analyze"
		m.submenuChoices = []string{
			"Quick Analysis (⚡ fast)",
			"Detailed Analysis (🔍 comprehensive)",
			"Custom Analysis (⚙️ custom)",
		}
		m.inSubmenu = true
		m.submenuCursor = 0
	case 1: // Compare Repositories
		m.SelectedOption = 1
		m.Done = true
	case 2: // View History
		m.SelectedOption = 2
		m.Done = true
	case 3: // Settings
		m.submenuType = "settings"
		m.submenuChoices = []string{
			"Theme Settings",
			"Export Options",
			"GitHub Token",
			"Reset to Defaults",
		}
		m.inSubmenu = true
		m.submenuCursor = 0
	case 4: // Help
		m.submenuType = "help"
		m.submenuChoices = []string{
			"Keyboard Shortcuts",
			"Getting Started",
			"Features Guide",
			"Troubleshooting",
		}
		m.inSubmenu = true
		m.submenuCursor = 0
	case 5: // Exit
		m.SelectedOption = 5
		m.Done = true
	}
}

func (m MenuModel) View() string {
	logo := `
 ██████╗ ███████╗██████╗  ██████╗      ██╗     ██╗   ██╗ ███████╗  ███████╗██████╗ 
 ██╔══██╗██╔════╝██╔══██╗██╔═══██╗     ██║     ╚██╗ ██╔╝ ╚════██║  ██╔════╝██╔══██╗
 ██████╔╝█████╗  ██████╔╝██║   ██║█████╗██║      ╚████╔╝     ██╔╝  █████╗  ██████╔╝
 ██╔══██╗██╔══╝  ██╔═══╝ ██║   ██║╚════╝██║       ╚██╔╝     ██╔╝   ██╔══╝  ██╔══██╗
 ██║  ██║███████╗██║     ╚██████╔╝     ███████╗   ██║      ██╔╝     ███████╗██║  ██║   
 ╚═╝  ╚═╝╚══════╝╚═╝      ╚═════╝      ╚══════╝   ╚═╝     ███████╗ ╚══════╝╚═╝  ╚═╝     
`
	content := TitleStyle.Render(logo) + "\n\n"

	if m.inSubmenu {
		return m.submenuView()
	}

	for i, choice := range m.choices {
		cursor := "  "
		style := NormalStyle

		if m.cursor == i {
			cursor = "▶ "
			style = SelectedStyle
		}

		content += fmt.Sprintf("%s%s\n", cursor, style.Render(choice))
	}

	content += "\n" + SubtleStyle.Render("↑ ↓ navigate • Enter select • q quit")

	box := BoxStyle.Render(content)

	if m.width == 0 {
		return box
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}

func (m MenuModel) submenuView() string {
	var title string
	var hint string

	switch m.submenuType {
	case "analyze":
		title = "📊 ANALYSIS TYPE"
		hint = "↑ ↓ navigate • Enter select • ESC back"
	case "settings":
		title = "⚙️ SETTINGS"
		hint = "↑ ↓ navigate • Enter select • ESC back"
	case "help":
		title = "❓ HELP MENU"
		hint = "↑ ↓ navigate • Enter select • ESC back"
	default:
		title = "SUBMENU"
		hint = "↑ ↓ navigate • Enter select • ESC back"
	}

	content := TitleStyle.Render(title) + "\n\n"

	for i, choice := range m.submenuChoices {
		cursor := "  "
		style := NormalStyle

		if m.submenuCursor == i {
			cursor = "▶ "
			style = SelectedStyle
		}

		content += fmt.Sprintf("%s%s\n", cursor, style.Render(choice))
	}

	content += "\n" + SubtleStyle.Render(hint)

	box := BoxStyle.Render(content)

	if m.width == 0 {
		return box
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}
