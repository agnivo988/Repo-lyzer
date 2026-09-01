package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PushingModel struct {
	spinner  spinner.Model
	repoName string
}

func NewPushingModel() PushingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return PushingModel{
		spinner: s,
	}
}

func (m PushingModel) Update(msg tea.Msg) (PushingModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg.(type) {
	case pushResult:
		return m, nil // Handled by parent
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m PushingModel) View(width, height int) string {
	header := TitleStyle.Render("📤 PUSHING REPOSITORY")

	content := fmt.Sprintf(
		"%s Pushing code to GitHub repo %s...\n\n"+
			"Please wait while the repository is being prepared and pushed.",
		m.spinner.View(),
		m.repoName,
	)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			BoxStyle.Render(content),
		),
	)
}

func (m *PushingModel) SetRepoName(repoName string) {
	m.repoName = repoName
}
