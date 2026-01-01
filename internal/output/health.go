package output

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

func PrintHealth(score int) {
    color := "#FF5F5F"
	label:= "🔴 Poor"

	if score >= 80 {
		color = "#00FF87"
		label = "🟢 Excellent"
	} else if score >= 60 {
		color = "#FFB000"
		label = "🟡 Good"
	 }

	 style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(color))

     fmt.Println(style.Render(
		fmt.Sprintf("\n🏆 Repo Health Score : %d/100 (%s)\n",score,label),
	 ))
}
