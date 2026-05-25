package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type LoadingModel struct {
	spinner      spinner.Model
	animTick     int
	progress     *ProgressTracker
	repoName     string
	analysisType string
}

func NewLoadingModel() LoadingModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return LoadingModel{
		spinner: s,
	}
}

func (m LoadingModel) Update(msg tea.Msg) (LoadingModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case struct{}:
		m.animTick++
		return m, TickProgressCmd()
	case AnalysisResult:
		return m, nil // Handled by parent
	case CachedAnalysisResult:
		return m, nil // Handled by parent
	case error:
		return m, nil // Handled by parent
	case tea.KeyMsg:
		if msg.String() == "esc" || msg.String() == "q" {
			return m, func() tea.Msg { return BackToMenuMsg{} }
		}
	}

	m.spinner, cmd = m.spinner.Update(msg)
	return m, cmd
}

func (m LoadingModel) View(width, height int) string {
	loadMsg := fmt.Sprintf("📊 Analyzing %s", m.repoName)
	if m.analysisType != "" {
		loadMsg += fmt.Sprintf(" (%s mode)", m.analysisType)
	}

	statusView := fmt.Sprintf("%s %s...", m.spinner.View(), loadMsg)

	if len(SatelliteFrames) > 0 {
		frame := SatelliteFrames[m.animTick%len(SatelliteFrames)]
		statusView += "\n" + lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF")).Render(frame)
	}

	// Show progress stages if available
	if m.progress != nil {
		stages := m.progress.GetAllStages()
		statusView += "\n\n"
		for _, stage := range stages {
			prefix := "⏳ "
			if stage.IsComplete {
				prefix = "✅ "
			} else if stage.IsActive {
				prefix = "⚙️  "
			}
			statusView += prefix + stage.Name + "\n"
		}

		// Add elapsed time
		elapsed := m.progress.GetElapsedTime()
		statusView += fmt.Sprintf("\n⏱️  %ds elapsed", int(elapsed.Seconds()))
	}

	statusView += "\n\n" + SubtleStyle.Render("Press ESC to cancel")

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		statusView,
	)
}

func (m *LoadingModel) SetRepoName(name string) {
	m.repoName = name
}

func (m *LoadingModel) SetAnalysisType(analysisType string) {
	m.analysisType = analysisType
}

func (m *LoadingModel) SetProgress(progress *ProgressTracker) {
	m.progress = progress
}

func (m *LoadingModel) GetProgress() *ProgressTracker {
	return m.progress
}
