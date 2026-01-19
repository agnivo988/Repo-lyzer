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
		case "up", "k", "w", "W":
			if m.inSubmenu {
				if m.submenuCursor > 0 {
					m.submenuCursor--
				} else {
					m.submenuCursor = len(m.submenuChoices) - 1
				}
			} else {
				if m.cursor > 0 {
					m.cursor--
				} else {
					m.cursor = len(m.choices) - 1
				}
			}
		case "down", "j", "s", "S":
			if m.inSubmenu {
				if m.submenuCursor < len(m.submenuChoices)-1 {
					m.submenuCursor++
				} else {
					m.submenuCursor = 0
				}
			} else {
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				} else {
					m.cursor = 0
				}
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
		case "1", "2", "3", "4", "5", "6", "7":
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
				m.SelectedOption = 7 // Exit
				m.Done = true
			} else {
				m.inSubmenu = false
				m.submenuCursor = 0
				m.submenuChoices = nil
				m.submenuType = ""
			}
		case "?":
			if !m.inSubmenu {
				m.cursor = 5 // Help
				m.enterSubmenu()
			}
		case "a":
			if !m.inSubmenu {
				m.cursor = 0
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
			// Quick access: Clone
			if !m.inSubmenu {
				m.cursor = 4
				m.enterSubmenu()
			}
		case "s":
			if !m.inSubmenu {
				m.cursor = 5
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
	case 5: // Settings
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
	case 6: // Help
		m.submenuType = "help"
		m.submenuChoices = []string{
			"Keyboard Shortcuts",
			"Getting Started",
			"Features Guide",
			"Troubleshooting",
		}
		m.inSubmenu = true
		m.submenuCursor = 0
	case 7: // Exit
		m.SelectedOption = 7
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
	shortcuts := []string{"a", "c", "h", "d", "s", "?", "q"}
	
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
	footer := SubtleStyle.Render("\n↑↓: navigate • Enter: select • Esc: back")

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