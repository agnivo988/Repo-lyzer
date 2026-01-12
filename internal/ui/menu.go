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

// NewMenuModel creates a new menu model with default options for repository analysis.
// It initializes the menu with choices for analyzing a repository, comparing repositories, and exiting.
// Returns the initialized MenuModel with cursor at the first option.
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
		case "up", "k":
			if m.inSubmenu {
				if m.submenuCursor > 0 {
					m.submenuCursor--
				} else {
					// Wrap to bottom
					m.submenuCursor = len(m.submenuChoices) - 1
				}
			} else {
				if m.cursor > 0 {
					m.cursor--
				} else {
					// Wrap to bottom
					m.cursor = len(m.choices) - 1
				}
			}
		case "down", "j":
			if m.inSubmenu {
				if m.submenuCursor < len(m.submenuChoices)-1 {
					m.submenuCursor++
				} else {
					// Wrap to top
					m.submenuCursor = 0
				}
			} else {
				if m.cursor < len(m.choices)-1 {
					m.cursor++
				} else {
					// Wrap to top
					m.cursor = 0
				}
			}
		case "home", "g":
			// Jump to first item
			if m.inSubmenu {
				m.submenuCursor = 0
			} else {
				m.cursor = 0
			}
		case "end", "G":
			// Jump to last item
			if m.inSubmenu {
				m.submenuCursor = len(m.submenuChoices) - 1
			} else {
				m.cursor = len(m.choices) - 1
			}
		case "1", "2", "3", "4", "5", "6", "7":
			// Quick jump to menu item by number
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
				m.SelectedOption = 6 // Exit
				m.Done = true
			} else {
				m.inSubmenu = false
				m.submenuCursor = 0
				m.submenuChoices = nil
				m.submenuType = ""
			}
		case "?":
			// Show help - jump to help menu
			if !m.inSubmenu {
				m.cursor = 5 // Help
				m.enterSubmenu()
			}
		case "a":
			// Quick access: Analyze
			if !m.inSubmenu {
				m.cursor = 0
				m.enterSubmenu()
			}
		case "c":
			// Quick access: Compare
			if !m.inSubmenu {
				m.cursor = 1
				m.enterSubmenu()
			}
		case "h":
			// Quick access: History
			if !m.inSubmenu {
				m.cursor = 2
				m.enterSubmenu()
			}
		case "s":
			// Quick access: Settings
			if !m.inSubmenu {
				m.cursor = 4
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
			"Quick Analysis (⚡ fast)",
			"Detailed Analysis (🔍 comprehensive)",
			"Custom Analysis (⚙️ custom)",
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
	content := TitleStyle.Render(logo) + "\n\n"

	if m.inSubmenu {
		return m.submenuView()
	}

	// Menu items with keyboard shortcuts
	shortcuts := []string{"a", "c", "h", "d", "s", "?", "q"}
	
	for i, choice := range m.choices {
		cursor := "  "
		style := NormalStyle
		shortcut := ""
		
		if i < len(shortcuts) {
			shortcut = fmt.Sprintf("[%s] ", shortcuts[i])
		}

		if m.cursor == i {
			cursor = "▶ "
			style = SelectedStyle
		}

		content += fmt.Sprintf("%s%s%s\n", cursor, SubtleStyle.Render(shortcut), style.Render(choice))
	}

	content += "\n" + SubtleStyle.Render("↑↓/jk: navigate • 1-7: jump • Enter/Space: select • ?: help • q: quit")

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
		hint = "↑↓/jk: navigate • 1-3: jump • Enter/Space: select • Esc/q: back"
	case "settings":
		title = "⚙️ SETTINGS"
		hint = "↑↓/jk: navigate • 1-5: jump • Enter/Space: select • Esc/q: back"
	case "help":
		title = "❓ HELP MENU"
		hint = "↑↓/jk: navigate • 1-4: jump • Enter/Space: select • Esc/q: back"
	default:
		title = "SUBMENU"
		hint = "↑↓/jk: navigate • Enter/Space: select • Esc/q: back"
	}

	content := TitleStyle.Render(title) + "\n\n"

	for i, choice := range m.submenuChoices {
		cursor := "  "
		style := NormalStyle
		shortcut := fmt.Sprintf("[%d] ", i+1)

		if m.submenuCursor == i {
			cursor = "▶ "
			style = SelectedStyle
		}

		content += fmt.Sprintf("%s%s%s\n", cursor, SubtleStyle.Render(shortcut), style.Render(choice))
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
