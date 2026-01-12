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

// NewDashboardModel creates a new dashboard model for displaying analysis results.
// It initializes an empty dashboard that can be populated with analysis data.
// Returns the initialized DashboardModel.
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
			// Refresh - re-analyze current repo
			if m.data.Repo != nil {
				return m, func() tea.Msg { return "refresh_data" }
			}

		case "b":
			// Add to favorites
			if m.data.Repo != nil {
				return m, func() tea.Msg { return "add_to_favorites" }
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
			BoxStyle.Render("📥 Export:\n[J] JSON  [M] Markdown  [P] PDF"),
		)
	}

	if m.statusMsg != "" {
		content += "\n" + SubtleStyle.Render(m.statusMsg)
	}

	// Navigation tabs
	tabs := m.renderTabs()
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
	views := []string{"Overview", "Repo", "Langs", "Activity", "Contribs", "Insights", "Deps", "Security", "Recruiter", "API"}
	var tabs []string

	for i, name := range views {
		if dashboardView(i) == m.currentView {
			tabs = append(tabs, SelectedStyle.Render(" "+name+" "))
		} else {
			tabs = append(tabs, SubtleStyle.Render(" "+name+" "))
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

	metrics := fmt.Sprintf(
		"Health Score: %d\nBus Factor: %d (%s)\nMaturity: %s (%d)",
		m.data.HealthScore,
		m.data.BusFactor,
		m.data.BusRisk,
		m.data.MaturityLevel,
		m.data.MaturityScore,
	)
	metricsBox := BoxStyle.Render(metrics)

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

func boolToYesNo(b bool) string {
	if b {
		return "✓ Yes"
	}
	return "✗ No"
}

func (m DashboardModel) dependenciesView() string {
	header := TitleStyle.Render("📦 Dependencies")

	if m.data.Dependencies == nil || len(m.data.Dependencies.Files) == 0 {
		return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render("No dependency files found"))
	}

	deps := m.data.Dependencies
	summary := fmt.Sprintf(
		"Total Dependencies: %d\nPackage Managers: %s\nLock File: %s",
		deps.TotalDeps,
		strings.Join(deps.Languages, ", "),
		boolToYesNo(deps.HasLockFile),
	)

	var depLines []string
	for _, file := range deps.Files {
		depLines = append(depLines, fmt.Sprintf("\n📄 %s (%d deps)", file.Filename, file.TotalCount))
		maxShow := 5
		if len(file.Dependencies) < maxShow {
			maxShow = len(file.Dependencies)
		}
		for i := 0; i < maxShow; i++ {
			d := file.Dependencies[i]
			depLines = append(depLines, fmt.Sprintf("  • %s %s", d.Name, d.Version))
		}
		if len(file.Dependencies) > maxShow {
			depLines = append(depLines, fmt.Sprintf("  ... and %d more", len(file.Dependencies)-maxShow))
		}
	}

	content := BoxStyle.Render(summary) + "\n" + BoxStyle.Render(strings.Join(depLines, "\n"))
	return lipgloss.JoinVertical(lipgloss.Left, header, content)
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

func (m DashboardModel) securityView() string {
	header := TitleStyle.Render("🔒 Security Scan")

	quality := m.data.CodeQuality
	if quality == nil {
		// Generate on the fly if not pre-computed
		quality = analyzer.AnalyzeCodeQuality(m.data.Repo, m.data.FileTree, m.data.Languages)
	}

	if m.data.Security == nil {
		return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render("No security scan data"))
	}

	sec := m.data.Security
	grade := analyzer.GetSecurityGrade(sec.SecurityScore)

	summary := fmt.Sprintf(
		"Security Score: %d/100 (Grade: %s)\nPackages Scanned: %d\nVulnerabilities: %d\n\n🔴 Critical: %d  🟠 High: %d  🟡 Medium: %d  🟢 Low: %d",
		sec.SecurityScore, grade, sec.ScannedPackages, sec.TotalCount,
		sec.CriticalCount, sec.HighCount, sec.MediumCount, sec.LowCount,
	)

	var vulnLines []string
	if len(sec.Vulnerabilities) == 0 {
		vulnLines = append(vulnLines, "✅ No known vulnerabilities!")
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
			vulnLines = append(vulnLines, fmt.Sprintf("... and %d more", len(sec.Vulnerabilities)-maxShow))
		}
	}

	content := BoxStyle.Render(summary) + "\n" + BoxStyle.Render(strings.Join(vulnLines, "\n"))
	return lipgloss.JoinVertical(lipgloss.Left, header, content)
}

func (m DashboardModel) licenseView() string {
	header := TitleStyle.Render("📜 License")

	if m.data.License == nil {
		return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render("No license data"))
	}

	lic := m.data.License
	grade := analyzer.GetLicenseGrade(lic.LicenseScore)

	var mainLic string
	if lic.MainLicense != nil {
		emoji := analyzer.GetLicenseEmoji(lic.MainLicense.Category)
		mainLic = fmt.Sprintf(
			"License: %s %s\nSPDX: %s\nCategory: %s\n\nPermissions:\n  Commercial: %s\n  Modify: %s\n  Distribute: %s\n  Patent Grant: %s",
			emoji, lic.MainLicense.Name, lic.MainLicense.SPDX, lic.MainLicense.Category,
			boolToYesNo(lic.MainLicense.Commercial), boolToYesNo(lic.MainLicense.Modify),
			boolToYesNo(lic.MainLicense.Distribute), boolToYesNo(lic.MainLicense.Patent),
		)
	} else {
		mainLic = "⚠️ No license detected"
	}

	score := fmt.Sprintf("\nLicense Score: %d/100 (Grade: %s)\nCompatibility: %s", lic.LicenseScore, grade, lic.Compatibility)

	var warnings string
	if len(lic.Warnings) > 0 {
		warnings = "\n\n⚠️ Warnings:\n• " + strings.Join(lic.Warnings, "\n• ")
	}

	content := BoxStyle.Render(mainLic + score + warnings)
	return lipgloss.JoinVertical(lipgloss.Left, header, content)
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
			" Activity: %s\n"+
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
  1-0           Jump to specific view
  
Views:
  1  Overview     - Health, Bus Factor, Maturity
  2  Repo         - Repository details
  3  Languages    - Language breakdown
  4  Activity     - Commit activity chart
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
  p             Export to PDF (when export menu open)
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

	// Get rate limit info
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
			"Rate Limit Status: %s\n"+
				"Requests: %d / %d (%.1f%% used)\n"+
				"Resets in: %s",
			status,
			rateLimit.Resources.Core.Limit-rateLimit.Resources.Core.Remaining,
			rateLimit.Resources.Core.Limit,
			usage,
			resetTime,
		)
	}

	// Check authentication mode
	mode := "🔴 Unauthenticated (60 req/hour)"
	if client.HasToken() {
		mode = "🟢 Authenticated (5000 req/hour)"
	}

	info := fmt.Sprintf(
		"Authentication: %s\n\n"+
			"%s\n\n"+
			"Data Fetched:\n"+
			"  • Repository info: ✓\n"+
			"  • Commits (1 year): %d\n"+
			"  • Contributors: %d\n"+
			"  • Languages: %d\n"+
			"  • File tree: %d entries\n\n"+
			"💡 Tip: Set GITHUB_TOKEN env variable\n"+
			"   for higher rate limits (5000/hour)",
		mode,
		rateLimitInfo,
		len(m.data.Commits),
		len(m.data.Contributors),
		len(m.data.Languages),
		len(m.data.FileTree),
	)

	return lipgloss.JoinVertical(lipgloss.Left, header, BoxStyle.Render(info))
}