package ui

import (
	"fmt"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	highThreshold = 0.67
	midThreshold  = 0.34
	maxBarWidth   = 20
)

var (
	dateStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#00E5FF"))
	countStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#FFB000"))
)

func barColor(count, max int) lipgloss.Style {
	if max == 0 {
		return lipgloss.NewStyle()
	}

	ratio := float64(count) / float64(max)

	switch {
	case ratio >= highThreshold:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#FF6E00"))
	case ratio >= midThreshold:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#E89149"))
	default:
		return lipgloss.NewStyle().Foreground(lipgloss.Color("#292C7B"))
	}
}

func RenderCommitActivity(data map[string]int, maxDays int) string {
	var sb strings.Builder

	sb.WriteString(TitleStyle.Render("📈 Commit Activity") + "\n")

	if len(data) == 0 {
		sb.WriteString("No commit activity found.\n")
		return sb.String()
	}

	dates := make([]string, 0, len(data))
	for d := range data {
		dates = append(dates, d)
	}

	sort.Strings(dates)

	if len(dates) > maxDays {
		dates = dates[len(dates)-maxDays:]
	}

	max := 0
	for _, d := range dates {
		count := data[d]
		if count > max {
			max = count
		}
	}

	for _, d := range dates {
		count := data[d]

		barLen := 0
		if max > 0 {
			barLen = int(float64(count) / float64(max) * maxBarWidth)

			// Ensure small non-zero values are still visible
			if count > 0 && barLen == 0 {
				barLen = 1
			}
		}

		bar := "░"
		if barLen > 0 {
			bar = barColor(count, max).Render(strings.Repeat("█", barLen))
		}

		sb.WriteString(fmt.Sprintf(
			"%s | %-20s %s\n",
			dateStyle.Render(d),
			bar,
			countStyle.Render(fmt.Sprintf("%d", count)),
		))
	}

	return sb.String()
}