package ui

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/github"
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
					_,err := ExportJSON(m.data, "analysis.json")
					if err != nil {
						return exportMsg{err, ""}
					}
					return exportMsg{nil, "✓ Exported to analysis.json"}
				}
			}

		case "m":
			if m.showExport {
				return m, func() tea.Msg {
					_,err := ExportMarkdown(m.data, "analysis.md")
					if err != nil {
						return exportMsg{err, ""}
					}
					return exportMsg{nil, "✓ Exported to analysis.md"}
				}
			}
			
		case "p":
			if m.showExport {
				return m, func() tea.Msg {
					_,err := ExportPDF(m.data, "analysis.pdf")
					if err != nil {
						return exportMsg{err, ""}
					}
					return exportMsg{nil, "✓ Exported to analysis.pdf"}
				}
			}

		case "f":
			return m, func() tea.Msg { return "switch_to_tree" }

		case "r":
			if m.data.Repo != nil {
				return m, func() tea.Msg { return "refresh_data" }
			}

		case "b":
			if m.data.Repo != nil {
				return m, func() tea.Msg { return "add_to_favorites" }
			}

		case "1":
			m.currentView = viewOverview
		case "2":
			m.currentView = viewRepo
		case "3":
			m.currentView = viewLanguages
		case "4":
			m.currentView = viewActivity
		case "5":
			m.currentView = viewContributors
		case "6":
			m.currentView = viewContributorInsights
		case "7":
			m.currentView = viewDependencies
		case "8":
			m.currentView = viewSecurity
		case "9":
			m.currentView = viewRecruiter
		case "0":
			m.currentView = viewAPIStatus

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

		case "t":
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
	if m.data.Repo == nil {
		return "No data loaded"
	}

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

	if m.showExport {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			CardStyle.Render("📥 Export Options:\n[J] JSON  [M] Markdown  [P] PDF"),
		)
	}

	if m.statusMsg != "" {
		content += "\n" + SubtleStyle.Render(m.statusMsg)
	}

	tabs := m.renderTabs()
	
	footer := SubtleStyle.Render("←→: switch view • f: files • e: export • ?: help • q: back")

	fullContent := lipgloss.JoinVertical(
		lipgloss.Left,
		tabs,
		"\n",
		content,
		"\n",
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
	views := []string{"Overview", "Repo", "Langs", "Activity", "Contribs", "Insights", "Deps", "Security", "Recruiter", "API"}
	var renderedTabs []string

	for i, name := range views {
		if dashboardView(i) == m.currentView {
			renderedTabs = append(renderedTabs, ActiveTabStyle.Render(name))
		} else {
			renderedTabs = append(renderedTabs, InactiveTabStyle.Render(name))
		}
	}

	return lipgloss.JoinHorizontal(lipgloss.Top, renderedTabs...)
}

func (m DashboardModel) overviewView() string {
	cacheIndicator := ""
	switch m.cacheStatus {
	case "fresh":
		cacheIndicator = "🟢 Fresh"
	case "cached":
		cacheIndicator = "🟡 Cached"
	case "expired":
		cacheIndicator = "🔴 Expired"
	}

	header := TitleStyle.Render(fmt.Sprintf(" %s ", m.data.Repo.FullName))
	subHeader := SubtleStyle.Render(fmt.Sprintf(" %s", cacheIndicator))

	metrics := fmt.Sprintf(
		"💚 Health:   %d/100\n"+
		"🚌 Bus Risk: %d (%s)\n"+
		"🏗️ Maturity: %s",
		m.data.HealthScore,
		m.data.BusFactor,
		m.data.BusRisk,
		m.data.MaturityLevel,
	)
	
	metricsBox := CardStyle.Render(lipgloss.JoinVertical(lipgloss.Left, 
		lipgloss.NewStyle().Bold(true).Render("Key Metrics"), 
		"\n"+metrics,
	))

	activity := analyzer.CommitsPerDay(m.data.Commits)
	chart := RenderCommitActivity(activity, 15)
	chartBox := CardStyle.Render(lipgloss.JoinVertical(lipgloss.Left,
		lipgloss.NewStyle().Bold(true).Render("Activity Trend"),
		"\n"+chart,
	))

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center, header, subHeader),
		"\n",
		lipgloss.JoinHorizontal(lipgloss.Top, metricsBox, chartBox),
	)
}

func (m DashboardModel) repoView() string {
	header := TitleStyle.Render(" Repository Details ")

	info := fmt.Sprintf(
		"Name:           %s\n"+
		"Description:    %s\n\n"+
		"⭐ Stars:        %d\n"+
		"🍴 Forks:        %d\n"+
		"🐛 Open Issues:  %d\n\n"+
		"📅 Created:      %s\n"+
		"🔄 Last Push:    %s\n"+
		"🌿 Branch:       %s\n"+
		"🔗 URL:          %s",
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

	return lipgloss.JoinVertical(lipgloss.Left, header, CardStyle.Render(info))
}

func (m DashboardModel) languagesView() string {
	header := TitleStyle.Render(" Languages ")

	if len(m.data.Languages) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, header, CardStyle.Render("No language data available"))
	}

	total := 0
	for _, bytes := range m.data.Languages {
		total += bytes
	}

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
		barLen := int(pct / 2)
		if barLen < 1 && lang.bytes > 0 {
			barLen = 1
		}
		bar := strings.Repeat("█", barLen)
		lines = append(lines, fmt.Sprintf("%-15s %s %.1f%%", lang.name, bar, pct))
	}

	return lipgloss.JoinVertical(lipgloss.Left, header, CardStyle.Render(strings.Join(lines, "\n")))
}

func (m DashboardModel) activityView() string {
	header := TitleStyle.Render(" Commit Activity ")
	activity := analyzer.CommitsPerDay(m.data.Commits)
	chart := RenderCommitActivity(activity, 40) // Wider chart
	totalCommits := len(m.data.Commits)
	stats := fmt.Sprintf("\nTotal Commits (Last Year): %d", totalCommits)

	return lipgloss.JoinVertical(lipgloss.Left, header, CardStyle.Render(chart+stats))
}

func (m DashboardModel) contributorsView() string {
	header := TitleStyle.Render(" Top Contributors ")

	if len(m.data.Contributors) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, header, CardStyle.Render("No contributor data available"))
	}

	var lines []string
	maxShow := 15
	if len(m.data.Contributors) < maxShow {
		maxShow = len(m.data.Contributors)
	}

	maxContribs := m.data.Contributors[0].Commits

	for i := 0; i < maxShow; i++ {
		c := m.data.Contributors[i]
		barLen := int(float64(c.Commits) / float64(maxContribs) * 30)
		if barLen < 1 {
			barLen = 1
		}
		bar := strings.Repeat("█", barLen)
		lines = append(lines, fmt.Sprintf("%2d. %-20s %s %d", i+1, c.Login, bar, c.Commits))
	}

	summary := fmt.Sprintf("\nTotal Contributors: %d", len(m.data.Contributors))
	lines = append(lines, summary)

	return lipgloss.JoinVertical(lipgloss.Left, header, CardStyle.Render(strings.Join(lines, "\n")))
}

func boolToYesNo(b bool) string {
	if b {
		return "✓"
	}
	return "✗"
}

func (m DashboardModel) dependenciesView() string {
	header := TitleStyle.Render(" Dependencies ")

	if m.data.Dependencies == nil || len(m.data.Dependencies.Files) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, header, CardStyle.Render("No dependency files found"))
	}

	deps := m.data.Dependencies
	summary := fmt.Sprintf(
		"Total Deps:       %d\nPackage Managers: %s\nLock File:        %s",
		deps.TotalDeps,
		strings.Join(deps.Languages, ", "),
		boolToYesNo(deps.HasLockFile),
	)

	var depLines []string
	for _, file := range deps.Files {
		depLines = append(depLines, fmt.Sprintf("\n📄 %s (%d)", file.Filename, file.TotalCount))
		maxShow := 5
		if len(file.Dependencies) < maxShow {
			maxShow = len(file.Dependencies)
		}
		for i := 0; i < maxShow; i++ {
			d := file.Dependencies[i]
			depLines = append(depLines, fmt.Sprintf("  • %s %s", d.Name, d.Version))
		}
		if len(file.Dependencies) > maxShow {
			depLines = append(depLines, fmt.Sprintf("  ... %d more", len(file.Dependencies)-maxShow))
		}
	}

	content := CardStyle.Render(summary) + "\n" + CardStyle.Render(strings.Join(depLines, "\n"))
	return lipgloss.JoinVertical(lipgloss.Left, header, content)
}

func (m DashboardModel) contributorInsightsView() string {
	header := TitleStyle.Render(" Insights ")

	insights := m.data.ContributorInsights
	if insights == nil {
		insights = analyzer.AnalyzeContributors(m.data.Contributors)
	}

	if insights.TotalContributors == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, header, CardStyle.Render("No contributor data available"))
	}

	col1 := fmt.Sprintf(
		"📊 OVERVIEW\n"+
			"Total:        %d\n"+
			"Active:       %d\n"+
			"Team Size:    %s\n"+
			"Diversity:    %.1f\n"+
			"Risk:         %s",
		insights.TotalContributors,
		insights.ActiveContributors,
		insights.TeamSize,
		insights.DiversityScore,
		insights.ConcentrationRisk,
	)

	col2 := ""
	if insights.TopContributor != nil {
		col2 = fmt.Sprintf(
			"👑 TOP CONTRIBUTOR\n"+
				"%s\n"+
				"%d commits (%.1f%%)\n"+
				"Type: %s",
			insights.TopContributor.Login,
			insights.TopContributor.Commits,
			insights.TopContributor.Percentage,
			insights.TopContributor.ContributorType,
		)
	}

	// Recommendations
	recs := "\n💡 RECOMMENDATIONS\n"
	for _, rec := range insights.Recommendations {
		recs += fmt.Sprintf("• %s\n", rec)
	}

	content := lipgloss.JoinVertical(lipgloss.Left, 
		lipgloss.JoinHorizontal(lipgloss.Top, CardStyle.Render(col1), CardStyle.Render(col2)),
		CardStyle.Render(recs),
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, content)
}

func (m DashboardModel) securityView() string {
	header := TitleStyle.Render(" Security ")

	if m.data.Security == nil {
		return lipgloss.JoinVertical(lipgloss.Left, header, CardStyle.Render("No security scan data"))
	}

	sec := m.data.Security
	grade := analyzer.GetSecurityGrade(sec.SecurityScore)

	summary := fmt.Sprintf(
		"Score: %d/100 (Grade: %s)\nScanned: %d packages\nTotal Vulns: %d\n\n🔴 %d  🟠 %d  🟡 %d  🟢 %d",
		sec.SecurityScore, grade, sec.ScannedPackages, sec.TotalCount,
		sec.CriticalCount, sec.HighCount, sec.MediumCount, sec.LowCount,
	)

	var vulnLines []string
	if len(sec.Vulnerabilities) == 0 {
		vulnLines = append(vulnLines, "✅ No known vulnerabilities found")
	} else {
		maxShow := 5
		if len(sec.Vulnerabilities) < maxShow {
			maxShow = len(sec.Vulnerabilities)
		}
		for i := 0; i < maxShow; i++ {
			v := sec.Vulnerabilities[i]
			vulnLines = append(vulnLines, fmt.Sprintf("%s %s - %s", analyzer.GetSeverityEmoji(v.Severity), v.ID, v.Package))
		}
		if len(sec.Vulnerabilities) > maxShow {
			vulnLines = append(vulnLines, fmt.Sprintf("... %d more", len(sec.Vulnerabilities)-maxShow))
		}
	}

	content := CardStyle.Render(summary) + "\n" + CardStyle.Render(strings.Join(vulnLines, "\n"))
	return lipgloss.JoinVertical(lipgloss.Left, header, content)
}

func (m DashboardModel) recruiterView() string {
	header := TitleStyle.Render(" Recruiter Summary ")

	activityLevel := "Low"
	if len(m.data.Commits) > 500 {
		activityLevel = "Very High"
	} else if len(m.data.Commits) > 200 {
		activityLevel = "High"
	} else if len(m.data.Commits) > 50 {
		activityLevel = "Medium"
	}

	summary := fmt.Sprintf(
		"REPO:     %s\n"+
		"STARS:    %d\n"+
		"COMMITS:  %d (Last Year)\n"+
		"CONTRIBS: %d\n"+
		"ACTIVITY: %s\n"+
		"MATURITY: %s (%d)\n"+
		"HEALTH:   %d/100\n",
		m.data.Repo.FullName,
		m.data.Repo.Stars,
		len(m.data.Commits),
		len(m.data.Contributors),
		activityLevel,
		m.data.MaturityLevel, m.data.MaturityScore,
		m.data.HealthScore,
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, CardStyle.Render(summary))
}

func (m DashboardModel) apiStatusView() string {
	header := TitleStyle.Render(" API Status ")

	client := github.NewClient()
	rateLimit, err := client.GetRateLimit()

	var rateLimitInfo string
	if err != nil {
		rateLimitInfo = "⚠️ Could not fetch rate limit info"
	} else {
		status := rateLimit.GetRateLimitStatus()
		resetTime := rateLimit.FormatResetTime()
		usage := rateLimit.UsagePercent()

		rateLimitInfo = fmt.Sprintf(
			"Status:    %s\n"+
			"Remaining: %d / %d\n"+
			"Used:      %.1f%%\n"+
			"Reset:     %s",
			status,
			rateLimit.Resources.Core.Limit-rateLimit.Resources.Core.Remaining,
			rateLimit.Resources.Core.Limit,
			usage,
			resetTime,
		)
	}

	mode := "🔴 Unauthenticated"
	if client.HasToken() {
		mode = "🟢 Authenticated"
	}

	info := fmt.Sprintf(
		"Mode: %s\n\n%s",
		mode,
		rateLimitInfo,
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, CardStyle.Render(info))
}

func (m DashboardModel) helpView() string {
	header := TitleStyle.Render(" Keyboard Shortcuts ")

	help := `
NAVIGATION
  ←/→       Switch view
  1-0       Jump to view
  
ACTIONS
  e         Export menu
  f         File tree
  r         Refresh
  ?         Toggle help
  q         Go back
`
	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.JoinVertical(lipgloss.Left, header, CardStyle.Render(help)),
	)
}