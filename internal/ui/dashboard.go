
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
	viewDependencies
	viewSecurity
	viewRecruiter
	viewAPIStatus
)

type DashboardModel struct {
 HEAD
	data        AnalysisResult
	BackToMenu  bool
	width       int
	height      int
	showExport  bool
	statusMsg   string
	currentView dashboardView
	showHelp    bool
feat/code-search-filter-by-filetype

	data       AnalysisResult
	err        error // explicit error state
	BackToMenu bool
	width      int
	height     int
	showExport bool
	statusMsg  string
552a131 (fix: remove duplicate tree definitions and unused types (#58))

	cacheStatus string // "fresh", "cached", or ""

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

func (m *DashboardModel) SetCacheStatus(status string) {
	m.cacheStatus = status
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
			m.err = msg.err
			m.statusMsg = msg.err.Error()
		} else {
			m.statusMsg = msg.msg
		}
HEAD
 feat/empty-state-error-handling-58
		return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg { return "clear_status" })

		return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg {
			return "clear_status"
		})
 552a131 (fix: remove duplicate tree definitions and unused types (#58))

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

 HEAD
		case "?", "h":
			m.showHelp = !m.showHelp


 552a131 (fix: remove duplicate tree definitions and unused types (#58))
		case "e":
			m.showExport = !m.showExport

		case "j":
			if m.showExport {
				return m, func() tea.Msg {
 feat/code-search-filter-by-filetype
 HEAD
					_,err := ExportJSON(m.data, "analysis.json")

					_, err := ExportJSON(m.data, "analysis.json")

					if err != nil {
						return exportMsg{err, ""}
					}
					return exportMsg{nil, "✓ Exported to analysis.json"}

					err := ExportJSON(m.data, "analysis.json")
					return exportMsg{err: err, msg: "Exported to analysis.json"}
 552a131 (fix: remove duplicate tree definitions and unused types (#58))
				}
			}

		case "m":
			if m.showExport {
				return m, func() tea.Msg {
 feat/code-search-filter-by-filetype
HEAD
					_,err := ExportMarkdown(m.data, "analysis.md")

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
			m.currentView = viewDependencies
			m.showHelp = false
			m.showExport = false
		case "7":
			m.currentView = viewSecurity
			m.showHelp = false
			m.showExport = false
		case "8":
			m.currentView = viewRecruiter
			m.showHelp = false
			m.showExport = false
		case "9":
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

					err := ExportMarkdown(m.data, "analysis.md")
					return exportMsg{err: err, msg: "Exported to analysis.md"}
 552a131 (fix: remove duplicate tree definitions and unused types (#58))
				}
			}

		case "t":
			// Toggle theme
			theme := CycleTheme()
			m.statusMsg = fmt.Sprintf("Theme: %s", theme.Name)
			return m, tea.Tick(2*time.Second, func(time.Time) tea.Msg {
				return "clear_status"
			})
		}
	}

	return m, nil
}

func (m DashboardModel) View() string {
 feat/empty-state-error-handling-58

	// ❌ Error state (explicit)
	if m.err != nil {
		return errorStateView(m.err.Error())
	}

	// 📭 Empty state (single source of truth)
	if m.data.IsEmpty() {
		return emptyStateView()

	if m.data.Repo == nil {
		return "No data loaded"
	}

 HEAD
	// Show help overlay
	if m.showHelp {
		return m.helpView()
	}

	// Header
	header := TitleStyle.Render(
		fmt.Sprintf("Analysis for %s", m.data.Repo.FullName),
 552a131 (fix: remove duplicate tree definitions and unused types (#58))

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
	views := []string{"Overview", "Repo", "Langs", "Activity", "Contribs", "Deps", "Security", "Recruiter", "API"}
	var tabs []string

	for i, name := range views {
		tab := fmt.Sprintf(" %d:%s ", i+1, name)
		if dashboardView(i) == m.currentView {
			tabs = append(tabs, SelectedStyle.Render(tab))
		} else {
			tabs = append(tabs, SubtleStyle.Render(tab))
		}

	}

	return BoxStyle.Render(strings.Join(tabs, "│"))
}

func (m DashboardModel) overviewView() string {
	// Cache status indicator
	cacheIndicator := ""
	switch m.cacheStatus {
	case "fresh":
		cacheIndicator = " 🟢 Fresh"
	case "cached":
		cacheIndicator = " 🟡 Cached"
	case "expired":
		cacheIndicator = " 🔴 Expired"
	}

	header := TitleStyle.Render(
		fmt.Sprintf("📊 Analysis for %s%s", m.data.Repo.FullName, cacheIndicator),
	)

feat/empty-state-error-handling-58
	// Metrics

	metrics := fmt.Sprintf(
		"Health Score: %d\nBus Factor: %d (%s)\nMaturity: %s (%d)",
		m.data.HealthScore,
		m.data.BusFactor,
		m.data.BusRisk,
		m.data.MaturityLevel,
		m.data.MaturityScore,
	)
	metricsBox := BoxStyle.Render(metrics)
 HEAD
 feat/empty-state-error-handling-58
	// Charts

	// Commit activity chart
 552a131 (fix: remove duplicate tree definitions and unused types (#58))
	activityData := analyzer.CommitsPerDay(m.data.Commits)
	chart := RenderCommitActivity(activityData, 10)
	chartBox := BoxStyle.Render(chart)

	// File tree (safe)
	treeContent := "📂 Files (Top 10):\n"
	limit := 10
	if len(m.data.FileTree) < limit {
		limit = len(m.data.FileTree)

	activity := analyzer.CommitsPerDay(m.data.Commits)
	chart := RenderCommitActivity(activity, 10)
	chartBox := BoxStyle.Render(chart)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		lipgloss.JoinHorizontal(lipgloss.Top, metricsBox, chartBox),
	)
}

func (m DashboardModel) repoView() string {
	header := TitleStyle.Render("📦 Repository Details")

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
		m.data.Repo.Description,
		m.data.Repo.Stars,
		m.data.Repo.Forks,
		m.data.Repo.OpenIssues,
		m.data.Repo.CreatedAt.Format("2006-01-02"),
		m.data.Repo.PushedAt.Format("2006-01-02"),
		m.data.Repo.DefaultBranch,
		m.data.Repo.HTMLURL,
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render(info))
}

func (m DashboardModel) languagesView() string {
	header := TitleStyle.Render("💻 Languages")

	if len(m.data.Languages) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render("No language data available"))

	}
 HEAD
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

	for i := 0; i < limit; i++ {
		icon := "📄"
		if m.data.FileTree[i].Type == "dir" {
			icon = "📁"
		}
		treeContent += fmt.Sprintf(
			"%s %s\n",
			icon,
			m.data.FileTree[i].Path,
		)
	}

	if len(m.data.FileTree) > limit {
		treeContent += fmt.Sprintf(
			"... and %d more",
			len(m.data.FileTree)-limit,
		)
	}

	treeBox := BoxStyle.Render(treeContent)

	// Layout
	row := lipgloss.JoinHorizontal(
		lipgloss.Top,
		metricsBox,
		chartBox,
	)

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		row,
		treeBox,
	)

	if m.showExport {
		exportMenu := BoxStyle.Render(
			"Export Options:\n[J] JSON\n[M] Markdown",
		)
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			exportMenu,
		)
 552a131 (fix: remove duplicate tree definitions and unused types (#58))
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

 feat/empty-state-error-handling-58
	if m.statusMsg != "" {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("205")).
				Render(m.statusMsg),
		)

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

func (m DashboardModel) dependenciesView() string {
	header := TitleStyle.Render("📦 Dependencies")

	if m.data.Dependencies == nil || len(m.data.Dependencies.Files) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render("No dependency files found (package.json, go.mod, requirements.txt, etc.)"))
	}

	deps := m.data.Dependencies
	var sections []string

	// Summary
	summary := fmt.Sprintf(
		"Total Dependencies: %d\n"+
			"Package Managers: %s\n"+
			"Lock File: %s",
		deps.TotalDeps,
		strings.Join(deps.Languages, ", "),
		boolToYesNo(deps.HasLockFile),
	)
	sections = append(sections, BoxStyle.Render("📊 Summary\n"+summary))

	// Show dependencies by file
	for _, file := range deps.Files {
		var lines []string
		lines = append(lines, fmt.Sprintf("📄 %s (%s) - %d deps", file.Filename, file.FileType, file.TotalCount))
		lines = append(lines, strings.Repeat("─", 50))

		// Group by type
		prodDeps := []string{}
		devDeps := []string{}
		otherDeps := []string{}

		for _, dep := range file.Dependencies {
			depLine := fmt.Sprintf("  %-30s %s", dep.Name, dep.Version)
			switch dep.Type {
			case "production":
				prodDeps = append(prodDeps, depLine)
			case "dev":
				devDeps = append(devDeps, depLine)
			default:
				otherDeps = append(otherDeps, depLine)
			}
		}

		if len(prodDeps) > 0 {
			lines = append(lines, "\n🔹 Production:")
			maxShow := 10
			if len(prodDeps) < maxShow {
				maxShow = len(prodDeps)
			}
			lines = append(lines, prodDeps[:maxShow]...)
			if len(prodDeps) > maxShow {
				lines = append(lines, fmt.Sprintf("  ... and %d more", len(prodDeps)-maxShow))
			}
		}

		if len(devDeps) > 0 {
			lines = append(lines, "\n🔸 Dev:")
			maxShow := 5
			if len(devDeps) < maxShow {
				maxShow = len(devDeps)
			}
			lines = append(lines, devDeps[:maxShow]...)
			if len(devDeps) > maxShow {
				lines = append(lines, fmt.Sprintf("  ... and %d more", len(devDeps)-maxShow))
			}
		}

		if len(otherDeps) > 0 {
			lines = append(lines, "\n🔻 Other:")
			lines = append(lines, otherDeps...)
		}

		sections = append(sections, BoxStyle.Render(strings.Join(lines, "\n")))
	}

	content := lipgloss.JoinVertical(lipgloss.Left, sections...)
	return lipgloss.JoinVertical(lipgloss.Left, header, content)
}

func boolToYesNo(b bool) string {
	if b {
		return "✓ Yes"
	}
	return "✗ No"
}

func (m DashboardModel) securityView() string {
	header := TitleStyle.Render("�  Security Scan")

	if m.data.Security == nil {
		return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render("No security scan data available"))
	}

	sec := m.data.Security

	// Security score and grade
	grade := analyzer.GetSecurityGrade(sec.SecurityScore)
	scoreColor := "green"
	if sec.SecurityScore < 70 {
		scoreColor = "red"
	} else if sec.SecurityScore < 90 {
		scoreColor = "yellow"
	}

	summary := fmt.Sprintf(
		"Security Score: %d/100 (Grade: %s)\n"+
			"Packages Scanned: %d\n"+
			"Vulnerabilities Found: %d\n\n"+
			"  🔴 Critical: %d\n"+
			"  🟠 High: %d\n"+
			"  🟡 Medium: %d\n"+
			"  🟢 Low: %d",
		sec.SecurityScore, grade,
		sec.ScannedPackages,
		sec.TotalCount,
		sec.CriticalCount,
		sec.HighCount,
		sec.MediumCount,
		sec.LowCount,
	)
	_ = scoreColor // Used for styling if needed

	summaryBox := BoxStyle.Render("📊 Summary\n" + summary)

	// List vulnerabilities
	var vulnLines []string
	if len(sec.Vulnerabilities) == 0 {
		vulnLines = append(vulnLines, "✅ No known vulnerabilities found!")
	} else {
		vulnLines = append(vulnLines, "⚠️ Vulnerabilities Detected:\n")
		maxShow := 10
		if len(sec.Vulnerabilities) < maxShow {
			maxShow = len(sec.Vulnerabilities)
		}

		for i := 0; i < maxShow; i++ {
			v := sec.Vulnerabilities[i]
			emoji := analyzer.GetSeverityEmoji(v.Severity)
			line := fmt.Sprintf("%s %s - %s@%s", emoji, v.ID, v.Package, v.Version)
			vulnLines = append(vulnLines, line)

			if v.Summary != "" {
				// Truncate summary if too long
				summ := v.Summary
				if len(summ) > 60 {
					summ = summ[:57] + "..."
				}
				vulnLines = append(vulnLines, "   "+summ)
			}

			if v.FixedIn != "" {
				vulnLines = append(vulnLines, fmt.Sprintf("   Fix: upgrade to %s", v.FixedIn))
			}
			vulnLines = append(vulnLines, "")
		}

		if len(sec.Vulnerabilities) > maxShow {
			vulnLines = append(vulnLines, fmt.Sprintf("... and %d more vulnerabilities", len(sec.Vulnerabilities)-maxShow))
		}
	}

	vulnBox := BoxStyle.Render(strings.Join(vulnLines, "\n"))

	return lipgloss.JoinVertical(lipgloss.Left, header, summaryBox, vulnBox)
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
			"🔥 Activity: %s\n"+
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
  
Views:
  1  Overview     - Health, Bus Factor, Maturity
  2  Repo         - Repository details
  3  Languages    - Language breakdown
  4  Activity     - Commit activity chart
  5  Contributors - Top contributors
  6  Recruiter    - Summary for recruiters
  7  API Status   - GitHub API rate limits

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
 HEAD
		lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render(help)),

		content,
 552a131 (fix: remove duplicate tree definitions and unused types (#58))
	)
}

 feat/empty-state-error-handling-58
/* ---------- Empty & Error Views ---------- */

func emptyStateView() string {
	return lipgloss.Place(
		60,
		10,
		lipgloss.Center,
		lipgloss.Center,
		BoxStyle.Render(
			"📭 No analysis data available\n\n"+
				"This repository does not contain enough data to analyze.\n"+
				"Try another repository.",
		),
	)
}

func errorStateView(msg string) string {
	return lipgloss.Place(
		60,
		10,
		lipgloss.Center,
		lipgloss.Center,
		BoxStyle.Render(
			"❌ Analysis failed\n\n"+msg+"\n\nPress q to return.",
		),
	)

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
