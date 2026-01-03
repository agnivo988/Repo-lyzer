package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ResponsiveLayout handles responsive design for different terminal sizes
type ResponsiveLayout struct {
	Width  int
	Height int
}

// NewResponsiveLayout creates a layout for given dimensions
func NewResponsiveLayout(width, height int) *ResponsiveLayout {
	return &ResponsiveLayout{Width: width, Height: height}
}

// IsSmallTerminal returns true if terminal is very small
func (rl *ResponsiveLayout) IsSmallTerminal() bool {
	return rl.Width < 80 || rl.Height < 24
}

// IsMobileTerminal returns true if terminal is mobile-sized
func (rl *ResponsiveLayout) IsMobileTerminal() bool {
	return rl.Width < 60
}

// GetMaxContentWidth returns safe content width considering padding
func (rl *ResponsiveLayout) GetMaxContentWidth() int {
	padding := 6 // 3px on each side
	width := rl.Width - padding
	if width < 40 {
		width = 40
	}
	return width
}

// GetMaxContentHeight returns safe content height considering padding
func (rl *ResponsiveLayout) GetMaxContentHeight() int {
	padding := 4 // 2 rows padding
	height := rl.Height - padding
	if height < 10 {
		height = 10
	}
	return height
}

// CenterText centers text horizontally and vertically
func (rl *ResponsiveLayout) CenterText(text string) string {
	if rl.Width == 0 || rl.Height == 0 {
		return text
	}

	return lipgloss.Place(
		rl.Width, rl.Height,
		lipgloss.Center, lipgloss.Center,
		text,
	)
}

// CenterContent centers content with proper margin
func (rl *ResponsiveLayout) CenterContent(content string) string {
	if rl.Width == 0 {
		return content
	}

	contentWidth := GetStringWidth(content)
	if contentWidth > rl.GetMaxContentWidth() {
		content = TruncateString(content, rl.GetMaxContentWidth())
	}

	return lipgloss.Place(
		rl.Width, rl.Height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}

// WrapText wraps text to fit terminal width
func (rl *ResponsiveLayout) WrapText(text string, padding int) string {
	maxWidth := rl.Width - padding
	if maxWidth < 20 {
		maxWidth = 20
	}

	words := strings.Fields(text)
	var lines []string
	var currentLine strings.Builder

	for _, word := range words {
		if currentLine.Len()+len(word)+1 > maxWidth {
			if currentLine.Len() > 0 {
				lines = append(lines, currentLine.String())
				currentLine.Reset()
			}
		}
		if currentLine.Len() > 0 {
			currentLine.WriteString(" ")
		}
		currentLine.WriteString(word)
	}

	if currentLine.Len() > 0 {
		lines = append(lines, currentLine.String())
	}

	return strings.Join(lines, "\n")
}

// GetStringWidth returns display width of string (accounting for ANSI codes)
func GetStringWidth(s string) int {
	// Remove ANSI escape sequences
	cleaned := strings.Builder{}
	inEscape := false
	for _, r := range s {
		if r == '\x1b' {
			inEscape = true
		} else if inEscape && r == 'm' {
			inEscape = false
		} else if !inEscape {
			cleaned.WriteRune(r)
		}
	}
	return len(cleaned.String())
}

// TruncateString truncates string to specified width with ellipsis
func TruncateString(s string, maxWidth int) string {
	if GetStringWidth(s) <= maxWidth {
		return s
	}

	if maxWidth < 3 {
		maxWidth = 3
	}

	// Simple truncation (doesn't handle ANSI codes perfectly)
	if len(s) > maxWidth {
		s = s[:maxWidth-3]
	}

	return s + "..."
}

// FormatMenuForDisplay formats menu items to fit terminal
func (rl *ResponsiveLayout) FormatMenuForDisplay(items []string) []string {
	maxWidth := rl.GetMaxContentWidth()
	formatted := make([]string, len(items))

	for i, item := range items {
		if GetStringWidth(item) > maxWidth {
			formatted[i] = TruncateString(item, maxWidth)
		} else {
			formatted[i] = item
		}
	}

	return formatted
}

// GetMinimumWarning returns warning message if terminal is too small
func (rl *ResponsiveLayout) GetMinimumWarning() string {
	if rl.Width >= 80 && rl.Height >= 24 {
		return ""
	}

	warning := "⚠️  Recommended minimum size: 80x24\n"
	warning += fmt.Sprintf("Current size: %dx%d\n", rl.Width, rl.Height)
	warning += "Some features may not display correctly"

	return warning
}

// PadContent adds padding to content
func (rl *ResponsiveLayout) PadContent(content string, horizontal, vertical int) string {
	padH := strings.Repeat(" ", horizontal)
	padV := strings.Repeat("\n", vertical)

	lines := strings.Split(content, "\n")
	for i, line := range lines {
		lines[i] = padH + line + padH
	}

	return padV + strings.Join(lines, "\n") + padV
}

// RenderResponsiveBox renders a box that adapts to terminal size
func (rl *ResponsiveLayout) RenderResponsiveBox(title, content string) string {
	maxWidth := rl.GetMaxContentWidth()

	// Truncate content if too wide
	if GetStringWidth(content) > maxWidth-4 {
		content = rl.WrapText(content, 4)
	}

	// For mobile terminals, use minimal styling
	if rl.IsMobileTerminal() {
		return fmt.Sprintf("=== %s ===\n%s", title, content)
	}

	return BoxStyle.Render(fmt.Sprintf("%s\n%s", TitleStyle.Render(title), content))
}

// ShouldShowSidebar returns whether sidebar should be displayed
func (rl *ResponsiveLayout) ShouldShowSidebar() bool {
	return rl.Width > 120
}

// ShouldShowPreview returns whether preview pane should be shown
func (rl *ResponsiveLayout) ShouldShowPreview() bool {
	return rl.Width > 100
}

// GetLayoutMode returns the current layout mode
func (rl *ResponsiveLayout) GetLayoutMode() string {
	if rl.IsMobileTerminal() {
		return "mobile"
	} else if rl.IsSmallTerminal() {
		return "compact"
	} else if rl.Width > 120 && rl.Height > 30 {
		return "wide"
	}
	return "default"
}

// AdjustSpacing returns appropriate spacing for current layout
func (rl *ResponsiveLayout) AdjustSpacing() (vertical, horizontal int) {
	switch rl.GetLayoutMode() {
	case "mobile":
		return 0, 0
	case "compact":
		return 0, 1
	case "wide":
		return 1, 2
	default:
		return 1, 1
	}
}
