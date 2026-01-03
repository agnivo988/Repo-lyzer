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
}

func NewMenuModel() MenuModel {
	return MenuModel{
		choices: []string{
			"Analyze a repository",
			"Compare repositories",
			"History",
			"Exit",
		},
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
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.SelectedOption = m.cursor
			m.Done = true
		}
	}
	return m, nil
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
