package ui

import "github.com/charmbracelet/lipgloss"

var (
	AppStyle = lipgloss.NewStyle().
			Padding(1, 2)

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Background(lipgloss.Color("62")).
			Foreground(lipgloss.Color("230")).
			MarginBottom(1)

	LogoStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#00E5FF")).
			MarginBottom(1)

	BoxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1, 2).
			Margin(0, 0)

	CardStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder(), true).
			BorderForeground(lipgloss.Color("#bd93f9")).
			Padding(1, 2).
			Margin(0, 1)

	
	SelectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282a36")).
			Background(lipgloss.Color("#bd93f9")).
			Bold(true).
			Padding(0, 2).
			MarginLeft(2)

	NormalStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f8f8f2")).
			Padding(0, 2).
			MarginLeft(2)

	InputStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFD700")).
			Bold(true)

	PlaceholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("240"))

	SubtleStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272a4")).
			Italic(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50fa7b")).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff5555")).
			Bold(true)

	ActiveTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#282a36")).
			Background(lipgloss.Color("#bd93f9")).
			Bold(true).
			Padding(0, 1).
			MarginRight(1)

	InactiveTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f8f8f2")).
			Background(lipgloss.Color("#44475a")).
			Padding(0, 1).
			MarginRight(1)
)