package ui

import "github.com/charmbracelet/lipgloss"

// Theme represents a color theme for the UI
type Theme struct {
	Name       string
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Background lipgloss.Color
	Text       lipgloss.Color
	TextMuted  lipgloss.Color
	Success    lipgloss.Color
	Error      lipgloss.Color
	Warning    lipgloss.Color
	Border     lipgloss.Color
}

// Available themes
var (
	// Catppuccin Mocha - warm, cozy dark theme
	CatppuccinMocha = Theme{
		Name:       "Catppuccin Mocha",
		Primary:    lipgloss.Color("#cba6f7"), // Mauve
		Secondary:  lipgloss.Color("#89b4fa"), // Blue
		Accent:     lipgloss.Color("#f5c2e7"), // Pink
		Background: lipgloss.Color("#1e1e2e"), // Base
		Text:       lipgloss.Color("#cdd6f4"), // Text
		TextMuted:  lipgloss.Color("#6c7086"), // Overlay0
		Success:    lipgloss.Color("#a6e3a1"), // Green
		Error:      lipgloss.Color("#f38ba8"), // Red
		Warning:    lipgloss.Color("#f9e2af"), // Yellow
		Border:     lipgloss.Color("#89b4fa"), // Blue
	}

	// Catppuccin Latte - light theme
	CatppuccinLatte = Theme{
		Name:       "Catppuccin Latte",
		Primary:    lipgloss.Color("#8839ef"), // Mauve
		Secondary:  lipgloss.Color("#1e66f5"), // Blue
		Accent:     lipgloss.Color("#ea76cb"), // Pink
		Background: lipgloss.Color("#eff1f5"), // Base
		Text:       lipgloss.Color("#4c4f69"), // Text
		TextMuted:  lipgloss.Color("#9ca0b0"), // Overlay0
		Success:    lipgloss.Color("#40a02b"), // Green
		Error:      lipgloss.Color("#d20f39"), // Red
		Warning:    lipgloss.Color("#df8e1d"), // Yellow
		Border:     lipgloss.Color("#1e66f5"), // Blue
	}

	// Dracula - popular dark theme
	Dracula = Theme{
		Name:       "Dracula",
		Primary:    lipgloss.Color("#bd93f9"), // Purple
		Secondary:  lipgloss.Color("#8be9fd"), // Cyan
		Accent:     lipgloss.Color("#ff79c6"), // Pink
		Background: lipgloss.Color("#282a36"), // Background
		Text:       lipgloss.Color("#f8f8f2"), // Foreground
		TextMuted:  lipgloss.Color("#6272a4"), // Comment
		Success:    lipgloss.Color("#50fa7b"), // Green
		Error:      lipgloss.Color("#ff5555"), // Red
		Warning:    lipgloss.Color("#f1fa8c"), // Yellow
		Border:     lipgloss.Color("#bd93f9"), // Purple
	}

	// Nord - cool, arctic theme
	Nord = Theme{
		Name:       "Nord",
		Primary:    lipgloss.Color("#88c0d0"), // Nord8
		Secondary:  lipgloss.Color("#81a1c1"), // Nord9
		Accent:     lipgloss.Color("#b48ead"), // Nord15
		Background: lipgloss.Color("#2e3440"), // Nord0
		Text:       lipgloss.Color("#eceff4"), // Nord6
		TextMuted:  lipgloss.Color("#4c566a"), // Nord3
		Success:    lipgloss.Color("#a3be8c"), // Nord14
		Error:      lipgloss.Color("#bf616a"), // Nord11
		Warning:    lipgloss.Color("#ebcb8b"), // Nord13
		Border:     lipgloss.Color("#5e81ac"), // Nord10
	}

	// Tokyo Night - modern dark theme
	TokyoNight = Theme{
		Name:       "Tokyo Night",
		Primary:    lipgloss.Color("#7aa2f7"), // Blue
		Secondary:  lipgloss.Color("#bb9af7"), // Purple
		Accent:     lipgloss.Color("#7dcfff"), // Cyan
		Background: lipgloss.Color("#1a1b26"), // Background
		Text:       lipgloss.Color("#c0caf5"), // Foreground
		TextMuted:  lipgloss.Color("#565f89"), // Comment
		Success:    lipgloss.Color("#9ece6a"), // Green
		Error:      lipgloss.Color("#f7768e"), // Red
		Warning:    lipgloss.Color("#e0af68"), // Yellow
		Border:     lipgloss.Color("#7aa2f7"), // Blue
	}

	// Gruvbox Dark - retro groove theme
	GruvboxDark = Theme{
		Name:       "Gruvbox Dark",
		Primary:    lipgloss.Color("#fe8019"), // Orange
		Secondary:  lipgloss.Color("#83a598"), // Aqua
		Accent:     lipgloss.Color("#d3869b"), // Purple
		Background: lipgloss.Color("#282828"), // Background
		Text:       lipgloss.Color("#ebdbb2"), // Foreground
		TextMuted:  lipgloss.Color("#928374"), // Gray
		Success:    lipgloss.Color("#b8bb26"), // Green
		Error:      lipgloss.Color("#fb4934"), // Red
		Warning:    lipgloss.Color("#fabd2f"), // Yellow
		Border:     lipgloss.Color("#d65d0e"), // Orange
	}

	// One Dark - Atom-inspired theme
	OneDark = Theme{
		Name:       "One Dark",
		Primary:    lipgloss.Color("#61afef"), // Blue
		Secondary:  lipgloss.Color("#c678dd"), // Purple
		Accent:     lipgloss.Color("#56b6c2"), // Cyan
		Background: lipgloss.Color("#282c34"), // Background
		Text:       lipgloss.Color("#abb2bf"), // Foreground
		TextMuted:  lipgloss.Color("#5c6370"), // Comment
		Success:    lipgloss.Color("#98c379"), // Green
		Error:      lipgloss.Color("#e06c75"), // Red
		Warning:    lipgloss.Color("#e5c07b"), // Yellow
		Border:     lipgloss.Color("#61afef"), // Blue
	}

	// All available themes
	AvailableThemes = []Theme{
		CatppuccinMocha,
		CatppuccinLatte,
		Dracula,
		Nord,
		TokyoNight,
		GruvboxDark,
		OneDark,
	}

	// Current active theme
	CurrentTheme      = CatppuccinMocha
	CurrentThemeIndex = 0
)

// ApplyTheme updates all style variables with the given theme
func ApplyTheme(theme Theme) {
	CurrentTheme = theme

	TitleStyle = lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Primary)

	BoxStyle = lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(theme.Border).
		Padding(1, 4)

	SelectedStyle = lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	NormalStyle = lipgloss.NewStyle().
		Foreground(theme.Text)

	InputStyle = lipgloss.NewStyle().
		Foreground(theme.Warning).
		Bold(true)

	SubtleStyle = lipgloss.NewStyle().
		Foreground(theme.TextMuted)

	SuccessStyle = lipgloss.NewStyle().
		Foreground(theme.Success).
		Bold(true)

	ErrorStyle = lipgloss.NewStyle().
		Foreground(theme.Error).
		Bold(true)
}

// CycleTheme switches to the next theme
func CycleTheme() Theme {
	CurrentThemeIndex = (CurrentThemeIndex + 1) % len(AvailableThemes)
	theme := AvailableThemes[CurrentThemeIndex]
	ApplyTheme(theme)
	return theme
}

// SetThemeByIndex sets theme by index
func SetThemeByIndex(index int) Theme {
	if index >= 0 && index < len(AvailableThemes) {
		CurrentThemeIndex = index
		theme := AvailableThemes[index]
		ApplyTheme(theme)
		return theme
	}
	return CurrentTheme
}

// GetThemeNames returns list of theme names
func GetThemeNames() []string {
	names := make([]string, len(AvailableThemes))
	for i, t := range AvailableThemes {
		names[i] = t.Name
	}
	return names
}

func init() {
	// Apply default theme on startup
	ApplyTheme(CurrentTheme)
}

// SetThemeByName sets the theme by its name
func SetThemeByName(name string) bool {
	for i, theme := range AvailableThemes {
		if theme.Name == name {
			CurrentThemeIndex = i
			ApplyTheme(theme)
			return true
		}
	}
	return false
}

// GetCurrentThemeName returns the name of the current theme
func GetCurrentThemeName() string {
	return CurrentTheme.Name
}
