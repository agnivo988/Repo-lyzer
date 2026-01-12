package output

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func PrintLanguages(langs map[string]int) {

	fmt.Println(SectionStyle.Render("\n⛳ Language Breakdown"))

	total := 0
	for _, v := range langs {
		total += v
	}

	for lang, size := range langs {
		percent := float64(size) / float64(total) * 100
		bar := lipgloss.NewStyle().Foreground(lipgloss.Color("#7CFF00")).Render(strings.Repeat("🟩", int(percent/5)))

		fmt.Printf("%-10s %s %.1f%%\n", lang, bar, percent)
	}
}
