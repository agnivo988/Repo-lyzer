package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HelpModel struct {
	helpContent string
}

func NewHelpModel() HelpModel {
	return HelpModel{}
}

func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			return m, func() tea.Msg { return BackToMenuMsg{} }
		}
	}
	return m, nil
}

func (m HelpModel) View(width, height int) string {
	var title string
	var content string

	switch m.helpContent {
	case "shortcuts":
		title = "❓ Keyboard Shortcuts"
		content = `
Global:
  Ctrl+H        Open analysis history
Main Menu:
  ↑↓/jk         Navigate menu
  Enter         Select option
  q             Quit application

Repository Input:
  Enter         Start analysis
  ESC           Back to menu
  Ctrl+U        Clear input
  Ctrl+W        Delete word
  Ctrl+A        Move to start
  Ctrl+E        Move to end

Dashboard Navigation:
  ←→/hl         Switch between views
  1-8           Jump to visible tabs
  e             Toggle export menu
  f             Open file tree
  z             Show repository trends
  r             Refresh data
  ?/h           Toggle help
  q/ESC         Go back
   .            Re-analyze current repository

File Tree:
  ↑↓/jk         Navigate files
  Enter         Open file details
  ESC           Back to dashboard

History:
  ↑↓/jk         Navigate entries
  Enter         Re-analyze repository
  d             Delete entry
  c             Clear all history
  q/ESC         Back to menu
`
	case "getting-started":
		title = "🚀 Getting Started"
		content = `
Welcome to Repo-lyzer!

1. Choose "Analyze Repository" from the main menu
2. Enter a repository in the format: owner/repo
   Example: microsoft/vscode
3. Select analysis type:
   - Quick: Fast overview
   - Detailed: Comprehensive analysis
   - Custom: Advanced options
4. Wait for analysis to complete
5. Navigate through the dashboard views
6. Export results if needed

For GitHub API access:
- Set GITHUB_TOKEN environment variable for higher rate limits
- Private repositories require authentication
`
	case "features":
		title = "✨ Features Guide"
		content = `
Repository Analysis:
  • Health Score: Overall repository health
  • Bus Factor: Risk of losing key contributors
  • Maturity Level: Project maturity assessment
  • Language Breakdown: Programming languages used
  • Commit Activity: Development activity over time
  • Top Contributors: Most active contributors
  • Recruiter Summary: Key insights for hiring

Export Options:
  • JSON: Structured data for further processing
  • Markdown: Human-readable reports

Additional Features:
  • Repository Comparison: Compare multiple repos
  • Analysis History: Re-analyze previous repos
  • File Tree: Explore repository structure
  • GitHub API Status: Monitor rate limit usage
`
	case "troubleshooting":
		title = "🔧 Troubleshooting"
		content = `
Common Issues:

Repository Not Found:
  • Check spelling: owner/repo format
  • Ensure repository is public or you have access
  • GitHub API might be rate limited

Analysis Fails:
  • Check internet connection
  • Verify GitHub API status
  • Try again later if rate limited

High Rate Limits:
  • Set GITHUB_TOKEN environment variable
  • Authenticated requests: 5000/hour
  • Unauthenticated: 60/hour

Private Repositories:
  • Require GITHUB_TOKEN with repo scope
  • Token must have access to the repository

Performance:
  • Detailed analysis takes longer
  • Large repositories may take several minutes
  • Use Quick analysis for fast results
`
	default:
		title = "❓ Help"
		content = `
Select a help topic from the menu above.
`
	}

	helpContent := TitleStyle.Render(title) + "\n\n" + content + "\n\n" + SubtleStyle.Render("Press ESC or q to go back")

	box := BoxStyle.Render(helpContent)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		box,
	)
}

func (m *HelpModel) SetHelpContent(content string) {
	m.helpContent = content
}
