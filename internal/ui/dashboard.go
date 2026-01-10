package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Dashboard view modes
type dashboardView int

const (
	viewOverview dashboardView = iota
	viewRepo
	viewLanguages
	viewActivity
	viewContributors
	viewContributorInsights
	viewDependencies
	viewSecurity
	viewRecruiter
	viewAPIStatus
)

type DashboardModel struct {
	data        AnalysisResult
	BackToMenu  bool
	width       int
	height      int
	showExport  bool
	statusMsg   string
	currentView dashboardView
	showHelp    bool
}

func NewDashboardModel() DashboardModel {
	return DashboardModel{
		currentView: viewOverview,
	}
}

func (m DashboardModel) Init() tea.Cmd { return nil }

func (m *DashboardModel) SetData(data AnalysisResult) {
	m.data = data
}

type exportMsg struct {
	err error
	msg string
}

func (m DashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case exportMsg:
		if msg.err != nil {
			m.statusMsg = fmt.Sprintf("Export failed: %v", msg.err)
		} else {
			m.statusMsg = msg.msg
		}
		return m, tea.Tick(3*time.Second, func(time.Time) tea.Msg {
			return "clear_status"
		})

	case string:
		if msg == "clear_status" {
			m.statusMsg = ""
		}

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc":
			if m.showHelp {
				m.showHelp = false
			} else if m.showExport {
				m.showExport = false
			} else if m.currentView != viewOverview {
				m.currentView = viewOverview
			} else {
				m.BackToMenu = true
			}

		case "?", "h":
			m.showHelp = !m.showHelp

		case "e":
			m.showExport = !m.showExport

		case "j":
			if m.showExport {
				return m, func() tea.Msg {
					_, err := ExportJSON(m.data, "analysis.json")
					if err != nil {
						return exportMsg{err, ""}
					}
					return exportMsg{nil, "✓ Exported to analysis.json"}
				}
			}

		case "m":
			if m.showExport {
				return m, func() tea.Msg {
					_, err := ExportMarkdown(m.data, "analysis.md")
					if err != nil {
						return exportMsg{err, ""}
					}
					return exportMsg{nil, "✓ Exported to analysis.md"}
				}
			}

		case "f":
			return m, func() tea.Msg { return "switch_to_tree" }

		case "r":
			// Refresh - re-analyze current repo
			if m.data.Repo != nil {
				return m, func() tea.Msg { return "refresh_data" }
			}

		// View switching keybindings
		case "1":
			m.currentView = viewOverview
			m.showHelp = false
			m.showExport = false
		case "2":
			m.currentView = viewRepo
			m.showHelp = false
			m.showExport = false
		case "3":
			m.currentView = viewLanguages
			m.showHelp = false
			m.showExport = false
		case "4":
			m.currentView = viewActivity
			m.showHelp = false
			m.showExport = false
		case "5":
			m.currentView = viewContributors
			m.showHelp = false
			m.showExport = false
		case "6":
			m.currentView = viewRecruiter
			m.showHelp = false
			m.showExport = false
		case "7":
			m.currentView = viewContributorInsights
			m.showHelp = false
			m.showExport = false
		case "7":
			m.currentView = viewDependencies
			m.showHelp = false
			m.showExport = false
		case "8":
			m.currentView = viewSecurity
			m.showHelp = false
			m.showExport = false
		case "9":
			m.currentView = viewRecruiter
			m.showHelp = false
			m.showExport = false
		case "0":
			m.currentView = viewAPIStatus
			m.showHelp = false
			m.showExport = false

		// Arrow key navigation between views
		case "right", "l":
			if !m.showHelp && !m.showExport {
				if m.currentView < viewAPIStatus {
					m.currentView++
				}
			}
		case "left":
			if !m.showHelp && !m.showExport {
				if m.currentView > viewOverview {
					m.currentView--
				}
			}
		}
	}

	return m, nil
}

func (m DashboardModel) View() string {
	if m.data.Repo == nil {
		return "No data loaded"
	}

	// Show help overlay
	if m.showHelp {
		return m.helpView()
	}

	var content string

	switch m.currentView {
	case viewOverview:
		content = m.overviewView()
	case viewRepo:
		content = m.repoView()
	case viewLanguages:
		content = m.languagesView()
	case viewActivity:
		content = m.activityView()
	case viewContributors:
		content = m.contributorsView()
	case viewContributorInsights:
		content = m.contributorInsightsView()
	case viewDependencies:
		content = m.dependenciesView()
	case viewSecurity:
		content = m.securityView()
	case viewRecruiter:
		content = m.recruiterView()
	case viewAPIStatus:
		content = m.apiStatusView()
	}

	// Add export panel if shown
	if m.showExport {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			BoxStyle.Render("📥 Export:\n[J] JSON  [M] Markdown"),
		)
	}

	if m.statusMsg != "" {
		content += "\n" + SubtleStyle.Render(m.statusMsg)
	}

	// Navigation tabs
	tabs := m.renderTabs()
	footer := SubtleStyle.Render("←→/hl: switch view • 1-6: jump to view • e: export • f: file tree • ?: help • q: back")
	footer := SubtleStyle.Render("←→/hl: switch view • 1-0: jump to view • e: export • f: file tree • ?: help • q: back")

	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		tabs,
		content,
		footer,
	)

	if m.width == 0 {
		return fullContent
	}

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		fullContent,
	)
}

func (m DashboardModel) renderTabs() string {
	views := []string{"Overview", "Repo", "Languages", "Activity", "Contributors", "Recruiter", "API"}
	views := []string{"Overview", "Repo", "Langs", "Activity", "Contribs", "Insights", "Deps", "Security", "Recruiter", "API"}
	var tabs []string

	for i, name := range views {
		tab := fmt.Sprintf(" %d:%s ", (i+1)%10, name)
		if dashboardView(i) == m.currentView {
			tabs = append(tabs, SelectedStyle.Render(tab))
		} else {
			tabs = append(tabs, SubtleStyle.Render(tab))
		}
	}

	return BoxStyle.Render(strings.Join(tabs, "│"))
}

// Helper to determine rank based on score
func getGamifiedRank(score int) (string, lipgloss.Style) {
	switch {
	case score >= 90:
		return "🏆 S-TIER (Legendary)", RankGoldStyle
	case score >= 75:
		return "⚔️ A-RANK (Master)", SelectedStyle
	case score >= 50:
		return "🛡️ B-RANK (Warrior)", NormalStyle
	default:
		return "💀 DEADLY (Needs Healing)", ErrorStyle
	}
}

func (m DashboardModel) overviewView() string {
	// 1. Get the Gamified Rank
	rankTitle, rankStyle := getGamifiedRank(m.data.HealthScore)

	// 2. Create the Header
	header := TitleStyle.Render(fmt.Sprintf("📊 ANALYSIS TARGET: %s", m.data.Repo.FullName))

	// 3. Create "Stat Cards"
	healthCard := CardStyle.Render(fmt.Sprintf(
		"HEALTH POINTS\n\n%s",
		rankStyle.Render(fmt.Sprintf("%d/100", m.data.HealthScore)),
	))

	rankCard := CardStyle.Render(fmt.Sprintf(
		"CURRENT RANK\n\n%s",
		rankStyle.Render(rankTitle),
	))

	busIcon := "🚌"
	if m.data.BusFactor > 5 {
		busIcon = "🏰"
	}
	busCard := CardStyle.Render(fmt.Sprintf(
		"DEFENSE (BUS)\n\n%s %d (%s)",
		busIcon, m.data.BusFactor, m.data.BusRisk,
	))

	maturityCard := CardStyle.Render(fmt.Sprintf(
		"MATURITY LVL\n\n%s (%d)",
		m.data.MaturityLevel, m.data.MaturityScore,
	))

	// 4. FIXED LAYOUT (Removed Responsive Logic)
	// Always use the horizontal layout
	statsLayout := lipgloss.JoinHorizontal(
		lipgloss.Top,
		healthCard,
		rankCard,
		busCard,
		maturityCard,
	)

	// 5. Add the chart
	activity := analyzer.CommitsPerDay(m.data.Commits)
	// Standard chart width
	chartWidth := 20
	chart := RenderCommitActivity(activity, chartWidth)
	chartBox := BoxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, "⚔️  COMBAT LOG (Activity)", chart))

	// 6. Join everything
	return lipgloss.JoinVertical(
		lipgloss.Center,
		header,
		"\n",
		statsLayout,
		"\n",
		chartBox,
	)
}

func (m DashboardModel) repoView() string {
	header := TitleStyle.Render("📦 Repository Details")

	// 1. Calculate safe width (Terminal width - padding)
	safeWidth := m.width - 10
	if safeWidth < 40 {
		safeWidth = 40 // Minimum width
	}

	// 2. Wrap the description text specifically
	desc := m.data.Repo.Description
	// Simple manual wrap or use lipgloss to wrap
	wrappedDesc := lipgloss.NewStyle().Width(safeWidth).Render(desc)

	info := fmt.Sprintf(
		"Name: %s\n"+
			"Description: %s\n"+
			"⭐ Stars: %d\n"+
			"🍴 Forks: %d\n"+
			"🐛 Open Issues: %d\n"+
			"📅 Created: %s\n"+
			"🔄 Last Push: %s\n"+
			"🌿 Default Branch: %s\n"+
			"🔗 URL: %s",
		m.data.Repo.FullName,
		wrappedDesc, // Use the wrapped description
		m.data.Repo.Stars,
		m.data.Repo.Forks,
		m.data.Repo.OpenIssues,
		m.data.Repo.CreatedAt.Format("2006-01-02"),
		m.data.Repo.PushedAt.Format("2006-01-02"),
		m.data.Repo.DefaultBranch,
		m.data.Repo.HTMLURL,
	)

	// 3. Render the box with a max width limit
	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		BoxStyle.Width(safeWidth).Render(info), // Force the box to fit
	)
}

func (m DashboardModel) languagesView() string {
	header := TitleStyle.Render("💻 Languages")

	if len(m.data.Languages) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render("No language data available"))
	}

	// Calculate total bytes
	total := 0
	for _, bytes := range m.data.Languages {
		total += bytes
	}

	// Sort languages by bytes
	type langStat struct {
		name  string
		bytes int
	}
	var langs []langStat
	for name, bytes := range m.data.Languages {
		langs = append(langs, langStat{name, bytes})
	}
	sort.Slice(langs, func(i, j int) bool {
		return langs[i].bytes > langs[j].bytes
	})

	var lines []string
	for _, lang := range langs {
		pct := float64(lang.bytes) / float64(total) * 100
		barLen := int(pct / 5) // 20 chars max
		if barLen < 1 && lang.bytes > 0 {
			barLen = 1
		}
		bar := strings.Repeat("█", barLen)
		lines = append(lines, fmt.Sprintf("%-15s %s %.1f%%", lang.name, bar, pct))
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render(strings.Join(lines, "\n")))
}

func (m DashboardModel) activityView() string {
	header := TitleStyle.Render("📈 Commit Activity (Last 30 Days)")

	activity := analyzer.CommitsPerDay(m.data.Commits)
	chart := RenderCommitActivity(activity, 30)

	totalCommits := len(m.data.Commits)
	stats := fmt.Sprintf("\nTotal Commits (1 year): %d", totalCommits)

	return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render(chart+stats))
}

func (m DashboardModel) contributorsView() string {
	header := TitleStyle.Render("👥 Top Contributors")

	if len(m.data.Contributors) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render("No contributor data available"))
	}

	var lines []string
	maxShow := 15
	if len(m.data.Contributors) < maxShow {
		maxShow = len(m.data.Contributors)
	}

	// Find max contributions for bar scaling
	maxContribs := m.data.Contributors[0].Commits

	for i := 0; i < maxShow; i++ {
		c := m.data.Contributors[i]
		barLen := int(float64(c.Commits) / float64(maxContribs) * 20)
		if barLen < 1 {
			barLen = 1
		}
		bar := strings.Repeat("█", barLen)
		lines = append(lines, fmt.Sprintf("%2d. %-20s %s %d", i+1, c.Login, bar, c.Commits))
	}

	summary := fmt.Sprintf("\nTotal Contributors: %d", len(m.data.Contributors))
	lines = append(lines, summary)

	return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render(strings.Join(lines, "\n")))
}

func (m DashboardModel) contributorInsightsView() string {
	header := TitleStyle.Render("🔍 Contributor Insights")

	insights := m.data.ContributorInsights
	if insights == nil {
		// Generate insights on the fly if not pre-computed
		insights = analyzer.AnalyzeContributors(m.data.Contributors)
	}

	if insights.TotalContributors == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render("No contributor data available"))
	}

	// Overview section
	overview := fmt.Sprintf(
		"📊 Overview\n"+
			"   Total Contributors: %d\n"+
			"   Active Contributors: %d (>1%% commits)\n"+
			"   Team Size: %s\n"+
			"   Diversity Score: %.1f/100\n"+
			"   Concentration Risk: %s",
		insights.TotalContributors,
		insights.ActiveContributors,
		insights.TeamSize,
		insights.DiversityScore,
		insights.ConcentrationRisk,
	)

	// Top contributor
	topContrib := ""
	if insights.TopContributor != nil {
		topContrib = fmt.Sprintf(
			"\n\n👑 Top Contributor\n"+
				"   %s: %d commits (%.1f%%)\n"+
				"   Type: %s",
			insights.TopContributor.Login,
			insights.TopContributor.Commits,
			insights.TopContributor.Percentage,
			insights.TopContributor.ContributorType,
		)
	}

	// Distribution stats
	dist := insights.CommitDistribution
	distribution := fmt.Sprintf(
		"\n\n📈 Commit Distribution\n"+
			"   Top 1%%:  %.1f%% of commits\n"+
			"   Top 10%%: %.1f%% of commits\n"+
			"   Top 50%%: %.1f%% of commits\n"+
			"   Gini Index: %.2f (0=equal, 1=unequal)",
		dist.Top1Percent,
		dist.Top10Percent,
		dist.Top50Percent,
		dist.GiniCoefficient,
	)

	// Contributor breakdown
	breakdown := fmt.Sprintf(
		"\n\n👥 Contributor Breakdown\n"+
			"   Veterans (>100 commits): %d\n"+
			"   New (<10 commits): %d",
		insights.VeteranContributors,
		insights.NewContributors,
	)

	// Recommendations
	recs := "\n\n💡 Recommendations\n"
	for _, rec := range insights.Recommendations {
		recs += fmt.Sprintf("   %s\n", rec)
	}

	content := overview + topContrib + distribution + breakdown + recs

	return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render(content))
}

func (m DashboardModel) dependenciesView() string {
	header := TitleStyle.Render("📦 Dependencies")

	insights := m.data.ContributorInsights
	if insights == nil {
		// Generate insights on the fly if not pre-computed
		insights = analyzer.AnalyzeContributors(m.data.Contributors)
	}

	if insights.TotalContributors == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render("No contributor data available"))
	}

	// Overview section
	overview := fmt.Sprintf(
		"📊 Overview\n"+
			"   Total Contributors: %d\n"+
			"   Active Contributors: %d (>1%% commits)\n"+
			"   Team Size: %s\n"+
			"   Diversity Score: %.1f/100\n"+
			"   Concentration Risk: %s",
		insights.TotalContributors,
		insights.ActiveContributors,
		insights.TeamSize,
		insights.DiversityScore,
		insights.ConcentrationRisk,
	)

	// Top contributor
	topContrib := ""
	if insights.TopContributor != nil {
		topContrib = fmt.Sprintf(
			"\n\n👑 Top Contributor\n"+
				"   %s: %d commits (%.1f%%)\n"+
				"   Type: %s",
			insights.TopContributor.Login,
			insights.TopContributor.Commits,
			insights.TopContributor.Percentage,
			insights.TopContributor.ContributorType,
		)
	}

	// Distribution stats
	dist := insights.CommitDistribution
	distribution := fmt.Sprintf(
		"\n\n📈 Commit Distribution\n"+
			"   Top 1%%:  %.1f%% of commits\n"+
			"   Top 10%%: %.1f%% of commits\n"+
			"   Top 50%%: %.1f%% of commits\n"+
			"   Gini Index: %.2f (0=equal, 1=unequal)",
		dist.Top1Percent,
		dist.Top10Percent,
		dist.Top50Percent,
		dist.GiniCoefficient,
	)

	// Contributor breakdown
	breakdown := fmt.Sprintf(
		"\n\n👥 Contributor Breakdown\n"+
			"   Veterans (>100 commits): %d\n"+
			"   New (<10 commits): %d",
		insights.VeteranContributors,
		insights.NewContributors,
	)

	// Recommendations
	recs := "\n\n💡 Recommendations\n"
	for _, rec := range insights.Recommendations {
		recs += fmt.Sprintf("   %s\n", rec)
	}

	content := overview + topContrib + distribution + breakdown + recs

	return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render(content))
}

func (m DashboardModel) securityView() string {
	header := TitleStyle.Render("🔒 Security Scan")

	quality := m.data.CodeQuality
	if quality == nil {
		// Generate on the fly if not pre-computed
		quality = analyzer.AnalyzeCodeQuality(m.data.Repo, m.data.FileTree, m.data.Languages)
	}

	if quality.Grade == "N/A" {
		return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render("No file tree data available for analysis"))
	}

	// Grade and overall score
	gradeColor := getGradeColor(quality.Grade)
	overview := fmt.Sprintf(
		"🎯 Overall Score: %d/100  Grade: %s\n\n"+
			"📊 Score Breakdown:\n"+
			"   Documentation: %s %d/100\n"+
			"   Testing:       %s %d/100\n"+
			"   Structure:     %s %d/100\n"+
			"   Maintenance:   %s %d/100",
		quality.OverallScore,
		gradeColor(quality.Grade),
		getScoreBar(quality.DocumentationScore), quality.DocumentationScore,
		getScoreBar(quality.TestingScore), quality.TestingScore,
		getScoreBar(quality.StructureScore), quality.StructureScore,
		getScoreBar(quality.MaintenanceScore), quality.MaintenanceScore,
	)

	// Project files checklist
	checklist := "\n\n📁 Project Files:\n"
	checklist += fmt.Sprintf("   %s README\n", checkMark(quality.HasReadme))
	checklist += fmt.Sprintf("   %s LICENSE\n", checkMark(quality.HasLicense))
	checklist += fmt.Sprintf("   %s CONTRIBUTING\n", checkMark(quality.HasContributing))
	checklist += fmt.Sprintf("   %s CHANGELOG\n", checkMark(quality.HasChangelog))
	checklist += fmt.Sprintf("   %s .gitignore\n", checkMark(quality.HasGitignore))
	checklist += fmt.Sprintf("   %s CI/CD\n", checkMark(quality.HasCI))
	checklist += fmt.Sprintf("   %s Tests\n", checkMark(quality.HasTests))

	// File statistics
	stats := fmt.Sprintf(
		"\n\n📈 File Statistics:\n"+
			"   Total Files: %d\n"+
			"   Source Files: %d\n"+
			"   Test Files: %d\n"+
			"   Test Ratio: %.1f%%",
		quality.FileStats.TotalFiles,
		quality.FileStats.SourceFiles,
		quality.FileStats.TestFiles,
		quality.FileStats.TestRatio*100,
	)

	// CI/Test frameworks
	frameworks := ""
	if len(quality.CIProviders) > 0 {
		frameworks += fmt.Sprintf("\n\n🔄 CI: %s", strings.Join(quality.CIProviders, ", "))
	}
	if len(quality.TestFrameworks) > 0 {
		frameworks += fmt.Sprintf("\n🧪 Tests: %s", strings.Join(quality.TestFrameworks, ", "))
	}

	// Code smells
	smells := ""
	if len(quality.CodeSmells) > 0 {
		smells = "\n\n⚠️ Issues Found:\n"
		for _, smell := range quality.CodeSmells {
			icon := "⚪"
			if smell.Severity == "High" {
				icon = "🔴"
			} else if smell.Severity == "Medium" {
				icon = "🟡"
			}
			smells += fmt.Sprintf("   %s %s\n", icon, smell.Description)
		}
	}

	// Recommendations
	recs := "\n\n💡 Recommendations:\n"
	for _, rec := range quality.Recommendations {
		recs += fmt.Sprintf("   %s\n", rec)
	}

	content := overview + checklist + stats + frameworks + smells + recs

	return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render(content))
}

// Helper functions for code quality view
func getGradeColor(grade string) func(string) string {
	return func(g string) string {
		switch grade {
		case "A":
			return "🟢 " + g
		case "B":
			return "🟢 " + g
		case "C":
			return "🟡 " + g
		case "D":
			return "🟠 " + g
		default:
			return "🔴 " + g
		}
	}
}

func getScoreBar(score int) string {
	filled := score / 10
	empty := 10 - filled
	return "[" + strings.Repeat("█", filled) + strings.Repeat("░", empty) + "]"
}

func checkMark(has bool) string {
	if has {
		return "✅"
	}
	return "❌"
}

func (m DashboardModel) recruiterView() string {
	header := TitleStyle.Render("👔 Recruiter Summary")

	// Determine activity level
	activityLevel := "Low"
	if len(m.data.Commits) > 500 {
		activityLevel = "Very High"
	} else if len(m.data.Commits) > 200 {
		activityLevel = "High"
	} else if len(m.data.Commits) > 50 {
		activityLevel = "Medium"
	}

	summary := fmt.Sprintf(
		"Repository: %s\n"+
			"⭐ Stars: %d\n"+
			"🍴 Forks: %d\n"+
			"📦 Commits (1y): %d\n"+
			"👥 Contributors: %d\n"+
			"🏗️ Maturity: %s (%d)\n"+
			"⚠️ Bus Factor: %d - %s\n"+
			"� Activity: %s\n"+
			"💚 Health Score: %d/100",
		m.data.Repo.FullName,
		m.data.Repo.Stars,
		m.data.Repo.Forks,
		len(m.data.Commits),
		len(m.data.Contributors),
		m.data.MaturityLevel, m.data.MaturityScore,
		m.data.BusFactor, m.data.BusRisk,
		activityLevel,
		m.data.HealthScore,
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render(summary))
}

func (m DashboardModel) helpView() string {
	header := TitleStyle.Render("❓ Keyboard Shortcuts")

	help := `
Dashboard Navigation:
  ←/→ or h/l    Switch between views
  1-7           Jump to specific view
  1-0           Jump to specific view
  
Views:
  1  Overview     - Health, Bus Factor, Maturity
  2  Repo         - Repository details
  3  Languages    - Language breakdown
  4  Activity     - Commit activity chart
  5  Contributors - Top contributors
  6  Recruiter    - Summary for recruiters
  7  API Status   - GitHub API rate limits
  5  Contributors - Top contributors list
  6  Insights     - Detailed contributor insights
  7  Dependencies - Project dependencies
  8  Security     - Security vulnerability scan
  9  Recruiter    - Summary for recruiters
  0  API Status   - GitHub API rate limits

Actions:
  e             Toggle export menu
  j             Export to JSON (when export menu open)
  m             Export to Markdown (when export menu open)
  f             Open file tree
  r             Refresh data
  ?/h           Toggle this help
  q/ESC         Go back / Close overlay
  Ctrl+C        Quit application
`

	return lipgloss.Place(
		m.width,
		m.height,
		lipgloss.Center,
		lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render(help)),
	)
}

func (m DashboardModel) apiStatusView() string {
	header := TitleStyle.Render("🔐 GitHub API Status")

	// Check if authenticated
	mode := "Unauthenticated (60 req/hour)"
	if m.data.Repo != nil && m.data.Repo.Private {
		mode = "Authenticated (5000 req/hour)"
	} else {
		// Simple check - if we got detailed data, likely authenticated
		if len(m.data.Contributors) > 30 {
			mode = "Authenticated (5000 req/hour)"
		}
	}

	info := fmt.Sprintf(
		"Mode: %s\n\n"+
			"Data Fetched:\n"+
			"  • Repository info: ✓\n"+
			"  • Commits (1 year): %d\n"+
			"  • Contributors: %d\n"+
			"  • Languages: %d\n"+
			"  • File tree: %d entries\n\n"+
			"Tip: Set GITHUB_TOKEN env variable\n"+
			"for higher rate limits (5000/hour)",
		mode,
		len(m.data.Commits),
		len(m.data.Contributors),
		len(m.data.Languages),
		len(m.data.FileTree),
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render(info))
}