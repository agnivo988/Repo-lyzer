package ui

import (
	"fmt"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type ExportMenuOption int

const (
	ExportJSON ExportMenuOption = iota
	ExportMarkdown
	ExportCSV
	ExportHTML
	ExportCancel
)

type DashboardModel struct {
	data              AnalysisResult
	bridge            *AnalyzerDataBridge
	BackToMenu        bool
	width             int
	height            int
	showExport        bool
	statusMsg         string
	exportCursor      int
	exportOptions     []string
	exportMenuVisible bool
	selectedMetric    string
	metricsScroll     int
}

func NewDashboardModel() DashboardModel {
	return DashboardModel{
		exportOptions: []string{
			"💾 Export as JSON",
			"📄 Export as Markdown",
			"📊 Export as CSV",
			"🌐 Export as HTML",
			"❌ Cancel",
		},
	}
}

func (m *DashboardModel) SetData(data AnalysisResult) {
	m.data = data
	m.bridge = NewAnalyzerDataBridge(data)
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
		// Clear status after 3 seconds
		return m, tea.Tick(3*time.Second, func(t time.Time) tea.Msg { return "clear_status" })

	case string:
		if msg == "clear_status" {
			m.statusMsg = ""
		}

	case tea.KeyMsg:
		if m.exportMenuVisible {
			switch msg.String() {
			case "up", "k":
				if m.exportCursor > 0 {
					m.exportCursor--
				}
			case "down", "j":
				if m.exportCursor < len(m.exportOptions)-1 {
					m.exportCursor++
				}
			case "enter":
				return m.handleExportSelection()
			case "esc":
				m.exportMenuVisible = false
			}
		} else {
			switch msg.String() {
			case "esc", "q":
				m.BackToMenu = true
			case "e":
				m.exportMenuVisible = !m.exportMenuVisible
				m.exportCursor = 0
			case "f":
				// Switch to tree view - handled in app.go
				return m, func() tea.Msg { return "switch_to_tree" }()
			}
		}
	}
	return m, nil
}

func (m *DashboardModel) handleExportSelection() (tea.Model, tea.Cmd) {
	switch m.exportCursor {
	case 0: // JSON
		return m, func() tea.Msg {
			err := ExportJSON(m.data, "analysis.json")
			return exportMsg{err, "✅ Exported to analysis.json"}
		}
	case 1: // Markdown
		return m, func() tea.Msg {
			err := ExportMarkdown(m.data, "analysis.md")
			return exportMsg{err, "✅ Exported to analysis.md"}
		}
	case 2: // CSV
		return m, func() tea.Msg {
			err := ExportCSV(m.data, "analysis.csv")
			return exportMsg{err, "✅ Exported to analysis.csv"}
		}
	case 3: // HTML
		return m, func() tea.Msg {
			err := ExportHTML(m.data, "analysis.html")
			return exportMsg{err, "✅ Exported to analysis.html"}
		}
	case 4: // Cancel
		m.exportMenuVisible = false
	}
	return m, nil
}

func (m DashboardModel) View() string {
	if m.data.Repo == nil {
		return "No data"
	}

	// Header
	header := TitleStyle.Render(fmt.Sprintf("📊 Analysis for %s", m.data.Repo.FullName))

	// Repository Info
	repoInfo := fmt.Sprintf(
		"Stars: %d  •  Forks: %d  •  Contributors: %d  •  Commits: %d",
		m.data.Repo.StargazersCount,
		m.data.Repo.ForksCount,
		len(m.data.Contributors),
		len(m.data.Commits),
	)

	// Metrics Column
	metrics := fmt.Sprintf(
		"🏥 Health Score: %d/100\n🚌 Bus Factor: %d (%s)\n📈 Maturity: %s (Score: %d)",
		m.data.HealthScore,
		m.data.BusFactor, m.data.BusRisk,
		m.data.MaturityLevel, m.data.MaturityScore,
	)
	metricsBox := BoxStyle.Render(metrics)

	// Charts
	activityData := analyzer.CommitsPerDay(m.data.Commits)
	chart := RenderCommitActivity(activityData, 10) // Show last 10 days
	chartBox := BoxStyle.Render(chart)

	// Layout
	content := lipgloss.JoinHorizontal(lipgloss.Top, metricsBox, chartBox)
	content = lipgloss.JoinVertical(lipgloss.Left, header, SubtleStyle.Render(repoInfo), content)

	if m.exportMenuVisible {
		exportMenu := "📥 EXPORT OPTIONS:\n\n"
		for i, opt := range m.exportOptions {
			cursor := "  "
			style := NormalStyle
			if m.exportCursor == i {
				cursor = "▶ "
				style = SelectedStyle
			}
			exportMenu += fmt.Sprintf("%s%s\n", cursor, style.Render(opt))
		}
		exportBox := BoxStyle.Render(exportMenu)
		content = lipgloss.JoinVertical(lipgloss.Left, content, exportBox)
	}

	if m.statusMsg != "" {
		content = lipgloss.JoinVertical(lipgloss.Left, content, SuccessStyle.Render(m.statusMsg))
	}

	footer := SubtleStyle.Render("e: export • f: file tree")
	if !m.exportMenuVisible {
		footer += SubtleStyle.Render(" • q: back")
	} else {
		footer += SubtleStyle.Render(" • ↑ ↓ select • Enter confirm • ESC close")
	}
	content += "\n" + footer

	if m.width == 0 {
		return BoxStyle.Render(content)
	}

	return lipgloss.Place(
		m.width, m.height,
		lipgloss.Center, lipgloss.Center,
		BoxStyle.Render(content),
	)
}
		content,
	)
}
