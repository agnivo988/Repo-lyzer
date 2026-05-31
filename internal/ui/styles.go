package ui

import "github.com/charmbracelet/lipgloss"

var (
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 2).
			MarginBottom(1)

	LogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00E5FF")).
			MarginBottom(1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 3).
			Margin(0, 0)

	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			Padding(1, 2).
			Margin(0, 1)

	SelectedStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 2).
			MarginLeft(2)

	NormalStyle = lipgloss.NewStyle().
			Padding(0, 2).
			MarginLeft(2)

	InputStyle = lipgloss.NewStyle().
			Bold(true)

	PlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))

	SubtleStyle = lipgloss.NewStyle().
			Italic(true)

	SuccessStyle = lipgloss.NewStyle().
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Bold(true)

	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 2).
			MarginRight(1).
			Border(lipgloss.RoundedBorder(), true, true, false, true)

	InactiveTabStyle = lipgloss.NewStyle().
				Padding(0, 2).
				MarginRight(1).
				Border(lipgloss.RoundedBorder(), true, true, false, true)
)
