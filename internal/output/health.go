package output

import (
	"fmt"
	"os"

	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/charmbracelet/lipgloss"
)

func PrintHealth(score int) {
	color := "#FF5F5F"
	label := "🔴 Poor"

	if score >= 80 {
		color = "#00FF87"
		label = "🟢 Excellent"
	} else if score >= 60 {
		color = "#FFB000"
		label = "🟡 Good"
	}

	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(color))

	fmt.Println(style.Render(
		fmt.Sprintf("\n🏆 Repo Health Score : %d/100 (%s)\n", score, label),
	))
}
func PrintGitHubAPIStatus(client *github.Client) {
	rateLimit, err := client.GetRateLimit()
	if err != nil {
		fmt.Println("⚠️ Unable to fetch GitHub API status")
		return
	}

	mode := "Unauthenticated"
	if os.Getenv("GITHUB_TOKEN") != "" {
		mode = "Authenticated"
	}

	style := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#7AE7C7"))

	fmt.Println(style.Render("🔐 GitHub API Status"))
	fmt.Printf("Mode        : %s\n", mode)
	fmt.Printf(
		"Requests    : %d / %d\n",
		rateLimit.Resources.Core.Remaining,
		rateLimit.Resources.Core.Limit,
	)
	fmt.Printf(
		"Resets At   : %s\n\n",
		rateLimit.ResetTime().Format("15:04"),
	)
}
