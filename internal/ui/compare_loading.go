package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type CompareLoadingModel struct {
	spinner  spinner.Model
	animTick int
	repo1    string
	repo2    string
}

func NewCompareLoadingModel() CompareLoadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return CompareLoadingModel{
		spinner: s,
	}
}

func (m CompareLoadingModel) Update(msg tea.Msg) (CompareLoadingModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case struct{}:
		m.animTick++
		return m, TickProgressCmd()
	case CompareResult:
		return m, nil // Handled by parent
	case error:
		return m, nil // Handled by parent
	case tea.KeyMsg:
		if msg.String() == "esc" {
			return m, func() tea.Msg { return BackToMenuMsg{} }
		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m CompareLoadingModel) View(width, height int) string {
	loadMsg := fmt.Sprintf("📊 Comparing %s vs %s", m.repo1, m.repo2)
	statusView := fmt.Sprintf("%s %s...", m.spinner.View(), loadMsg)

	if len(SatelliteFrames) > 0 {
		frame := SatelliteFrames[m.animTick%len(SatelliteFrames)]
		statusView += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF")).Render(frame)
	}

	statusView += "\n\n" + SubtleStyle.Render("Press ESC to cancel")

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		statusView,
	)
}

func (m *CompareLoadingModel) SetRepos(repo1, repo2 string) {
	m.repo1 = repo1
	m.repo2 = repo2
}
