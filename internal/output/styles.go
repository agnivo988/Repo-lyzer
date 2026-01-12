package output

import "github.com/charmbracelet/lipgloss"

var (
	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7CFF00"))

	SectionStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00E5FF"))

	SuccessStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00FF87"))

	WarningStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFB000"))

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF5F5F"))
)
