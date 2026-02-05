package output

import (
	"fmt"
	"strings"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/charmbracelet/lipgloss"
)

// PrintHotspots prints the identify hotspots in a table format using lipgloss
func PrintHotspots(hotspots []analyzer.Hotspot) {
	if len(hotspots) == 0 {
		return
	}

	fmt.Println("\n🔥 Hotspot files (Complex & Frequently Changed)")

	// Define styles
	borderColor := lipgloss.Color("240")
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("205")).
		Padding(0, 1)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	borderStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor)

	// Calculate column widths
	headers := []string{"File", "Score", "Churn", "Size", "Complexity", "Issues"}
	widths := make([]int, len(headers))

	// Set minimum widths for each column to prevent wrapping
	minWidths := []int{35, 5, 5, 5, 10, 5} // Minimum widths for each column
	for i, h := range headers {
		widths[i] = len(h)
		if widths[i] < minWidths[i] {
			widths[i] = minWidths[i]
		}
	}

	// Calculate max widths from data
	maxFileWidth := 50 // Cap file path width to prevent extremely wide tables
	for _, h := range hotspots {
		// For file paths, truncate if too long
		filePath := h.FilePath
		if len(filePath) > maxFileWidth {
			filePath = "..." + filePath[len(filePath)-maxFileWidth+3:]
		}
		updateWidth(widths, 0, filePath)
		updateWidth(widths, 1, fmt.Sprintf("%d", h.Score))
		updateWidth(widths, 2, fmt.Sprintf("%d", h.ChurnScore))
		updateWidth(widths, 3, fmt.Sprintf("%d", h.SizeScore))
		updateWidth(widths, 4, fmt.Sprintf("%d", h.Complexity))
		updateWidth(widths, 5, h.Reason)
	}

	// Cap file width
	if widths[0] > maxFileWidth {
		widths[0] = maxFileWidth
	}

	// Render Header - use simple string formatting for consistent alignment
	var headerParts []string
	for i, h := range headers {
		// Pad header to column width
		padded := h + strings.Repeat(" ", widths[i]-len(h))
		headerParts = append(headerParts, headerStyle.Render(padded))
	}
	headerStr := strings.Join(headerParts, lipgloss.NewStyle().Foreground(borderColor).Render(" │ "))

	// Render Rows
	var rows []string
	for _, h := range hotspots {
		// Truncate file path if too long
		filePath := h.FilePath
		if len(filePath) > maxFileWidth {
			filePath = "..." + filePath[len(filePath)-maxFileWidth+3:]
		}

		data := []string{
			filePath,
			fmt.Sprintf("%d", h.Score),
			fmt.Sprintf("%d", h.ChurnScore),
			fmt.Sprintf("%d", h.SizeScore),
			fmt.Sprintf("%d", h.Complexity),
			h.Reason,
		}

		// Build row with consistent padding
		var rowParts []string
		for i, d := range data {
			// Truncate if data is longer than column width
			if len(d) > widths[i] {
				d = d[:widths[i]-3] + "..."
			}
			// Pad data to column width
			padded := d + strings.Repeat(" ", widths[i]-len(d))
			rowParts = append(rowParts, cellStyle.Render(padded))
		}
		rows = append(rows, strings.Join(rowParts, lipgloss.NewStyle().Foreground(borderColor).Render(" │ ")))
	}

	// Combine
	tableContent := headerStr + "\n" +
		lipgloss.NewStyle().Foreground(borderColor).Render(strings.Repeat("─", lipgloss.Width(headerStr))) + "\n" +
		strings.Join(rows, "\n")

	// Print with border
	fmt.Println(borderStyle.Render(tableContent))
}

func updateWidth(widths []int, idx int, content string) {
	w := len(content)
	if w > widths[idx] {
		widths[idx] = w
	}
}
