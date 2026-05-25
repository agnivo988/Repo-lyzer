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
			"⭐ Favorites",
			"🔄 Compare Repositories",
			"📜 View History",
			"📥 Clone Repository",
			"🔔 Notifications",
			"👀 Monitoring",
			"⚙️ Settings",
			"❓ Help",
			"🚪 Exit",
		},
		inSubmenu: false,
	}
}

// Helper functions for cursor navigation
func (m *MenuModel) moveCursorUp(choices []string, cursor *int) {
	if *cursor > 0 {
		*cursor--
	} else {
		*cursor = len(choices) - 1
	}
}

func (m *MenuModel) moveCursorDown(choices []string, cursor *int) {
	if *cursor < len(choices)-1 {
		*cursor++
	} else {
		*cursor = 0
	}
}

func (m *MenuModel) moveCursorHome(cursor *int) {
	*cursor = 0
}

func (m *MenuModel) moveCursorEnd(choices []string, cursor *int) {
	*cursor = len(choices) - 1
}

func (m MenuModel) Init() tea.Cmd { return nil }

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k", "w", "W":
			if m.inSubmenu {
				m.moveCursorUp(m.submenuChoices, &m.submenuCursor)
			} else {
				m.moveCursorUp(m.choices, &m.cursor)
			}
		case "down", "j", "S":
			if m.inSubmenu {
				m.moveCursorDown(m.submenuChoices, &m.submenuCursor)
			} else {
				m.moveCursorDown(m.choices, &m.cursor)
			}
		case "home", "g":
			if m.inSubmenu {
				m.submenuCursor = 0
			} else {
				m.cursor = 0
			}
		case "end", "G":
			if m.inSubmenu {
				m.submenuCursor = len(m.submenuChoices) - 1
			} else {
				m.cursor = len(m.choices) - 1
			}
		case "1", "2", "3", "4", "5", "6", "7", "8", "9":
			idx := int(msg.String()[0] - '1')
			if !m.inSubmenu && idx < len(m.choices) {
				m.cursor = idx
				m.enterSubmenu()
			} else if m.inSubmenu && idx < len(m.submenuChoices) {
				m.submenuCursor = idx
				m.SelectedOption = m.cursor
				m.Done = true
				m.inSubmenu = false
			}
		case "enter", " ":
			if m.inSubmenu {
				m.SelectedOption = m.cursor
				m.Done = true
				m.inSubmenu = false
			} else {
				m.enterSubmenu()
			}
		case "esc", "backspace":
			if m.inSubmenu {
				m.inSubmenu = false
				m.submenuCursor = 0
				m.submenuChoices = nil
				m.submenuType = ""
			}
		case "q":
			if !m.inSubmenu {
				m.SelectedOption = 9
				m.Done = true
			} else {
				m.inSubmenu = false
				m.submenuCursor = 0
				m.submenuChoices = nil
				m.submenuType = ""
			}
		case "?":
			if !m.inSubmenu {
				m.cursor = 8
				m.enterSubmenu()
			}
		case "a":
			if !m.inSubmenu {
				m.cursor = 0
				m.enterSubmenu()
			}
		case "f":
			if !m.inSubmenu {
				m.cursor = 1 // Favorites
				m.enterSubmenu()
			}
		case "c":
			if !m.inSubmenu {
				m.cursor = 2
				m.enterSubmenu()
			}
		case "h", "H":
			// Quick access: History
			if !m.inSubmenu {
				m.cursor = 3
				m.enterSubmenu()
			}
		case "d":
			if !m.inSubmenu {
				m.cursor = 4 // Clone
				m.enterSubmenu()
			}
		case "n":
			if !m.inSubmenu {
				m.cursor = 5 // Notification
				m.enterSubmenu()
			}
		case "m":
			if !m.inSubmenu {
				m.cursor = 6 // Monitor
				m.enterSubmenu()
			}
		case "s":
			if !m.inSubmenu {
				m.cursor = 7 //Setting
				m.enterSubmenu()
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
			"Quick Analysis",
			"Detailed Analysis",
			"Custom Analysis",
		}
		m.inSubmenu = true
		m.submenuCursor = 0
	case 1: // Favorites
		m.SelectedOption = 1
		m.Done = true
	case 2: // Compare Repositories
		m.SelectedOption = 2
		m.Done = true
	case 3: // View History
		m.SelectedOption = 3
		m.Done = true
	case 4: // Clone Repository
		m.SelectedOption = 4
		m.Done = true
	case 5: // Notification
		m.SelectedOption = 5
		m.Done = true
	case 6: // Monitoring
		m.SelectedOption = 6
		m.Done = true
	case 7: // Settings
		m.submenuType = "settings"
		m.submenuChoices = []string{
			"Theme Settings",
			"Cache Settings",
			"Export Options",
			"GitHub Token",
			"Reset to Defaults",
		}
		m.inSubmenu = true
		m.submenuCursor = 0
	case 8: // Help
		m.submenuType = "help"
		m.submenuChoices = []string{
			"Keyboard Shortcuts",
			"Getting Started",
			"Features Guide",
			"Troubleshooting",
		}
		m.inSubmenu = true
		m.submenuCursor = 0
	case 9: // Exit
		m.SelectedOption = 9
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
	logoView := LogoStyle.Render(logo)

	if m.inSubmenu {
		return m.submenuView(logoView)
	}
	// Menu items with keyboard shortcuts
	shortcuts := []string{"a", "f", "c", "h", "d", "n", "m", "s", "?", "q"}
	var menuItems []string

	for i, choice := range m.choices {
		shortcut := ""
		if i < len(shortcuts) {
			shortcut = fmt.Sprintf("[%s] ", shortcuts[i])
		}

		if m.cursor == i {
			item := fmt.Sprintf("%s%s", shortcut, choice)
			menuItems = append(menuItems, SelectedStyle.Render(item))
		} else {
			item := fmt.Sprintf("%s%s", shortcut, choice)
			menuItems = append(menuItems, NormalStyle.Render(item))
		}
	}

	menuContent := lipgloss.JoinVertical(lipgloss.Left, menuItems...)

	footer := SubtleStyle.Render("\n↑↓: navigate • Enter: select")

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logoView,
		"\n",
		BoxStyle.Render(menuContent),
		footer,
	)

	if m.width == 0 {
		return content
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

func (m MenuModel) submenuView(logoView string) string {
	var title string
	var hint string

	switch m.submenuType {
	case "analyze":
		title = "📊 ANALYSIS TYPE"
		hint = "↑↓/jk/ws: navigate • 1-3: jump • Enter/Space: select • Esc/q: back"
	case "settings":
		title = "⚙️ SETTINGS"
		hint = "↑↓/jk/ws: navigate • 1-5: jump • Enter/Space: select • Esc/q: back"
	case "help":
		title = "❓ HELP MENU"
		hint = "↑↓/jk/ws: navigate • 1-4: jump • Enter/Space: select • Esc/q: back"
	default:
		title = "SUBMENU"
		hint = "↑↓/jk/ws: navigate • Enter/Space: select • Esc/q: back"
	}

	header := TitleStyle.Render(title)

	var menuItems []string

	for i, choice := range m.submenuChoices {
		shortcut := fmt.Sprintf("[%d] ", i+1)

		if m.submenuCursor == i {
			item := fmt.Sprintf("%s%s", shortcut, choice)
			menuItems = append(menuItems, SelectedStyle.Render(item))
		} else {
			item := fmt.Sprintf("%s%s", shortcut, choice)
			menuItems = append(menuItems, NormalStyle.Render(item))
		}
	}

	menuContent := lipgloss.JoinVertical(lipgloss.Left, menuItems...)
	footer := SubtleStyle.Render("\n" + hint)

	content := lipgloss.JoinVertical(
		lipgloss.Center,
		logoView,
		"\n",
		header,
		BoxStyle.Render(menuContent),
		footer,
	)

	if m.width == 0 {
		return content
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}
