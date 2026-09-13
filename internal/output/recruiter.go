package output

import (
	"fmt"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/charmbracelet/lipgloss"
)

func PrintRecruiterSummary(s analyzer.RecruiterSummary) {
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#00E5FF"))

	fmt.Println(title.Render("\n👔 Recruiter Summary"))
	fmt.Println("Repository:", s.RepoName)
	fmt.Println("⭐ Stars:", s.Stars)
	fmt.Println("🍴 Forks:", s.Forks)
	fmt.Println("📦 Commits (1y):", s.CommitsLastYear)
	fmt.Println("👥 Contributors:", s.Contributors)
	fmt.Println("🏗️ Maturity:", s.MaturityLevel, "(", s.MaturityScore, ")")
	fmt.Println("⚠️ Bus Factor:", s.BusFactor, "-", s.BusRisk)
	fmt.Println("🐛 Issue Health:", s.IssueHealth)
	fmt.Println("📥 PR Health:", s.PRHealth)
	fmt.Println("🔥 Activity:", s.ActivityLevel)
}
