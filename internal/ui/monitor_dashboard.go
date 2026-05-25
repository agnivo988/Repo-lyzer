package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/monitor"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// MonitorDashboardModel handles the monitoring dashboard view
type MonitorDashboardModel struct {
	monitor       *monitor.Monitor
	repoName      string
	owner         string
	repo          string
	interval      time.Duration
	notifications []monitor.Notification
	width         int
	height        int
	isMonitoring  bool
	err           error
	scrollOffset  int
	autoScroll    bool
	monitorRunID  int
}

// NewMonitorDashboardModel creates a new monitor dashboard model
func NewMonitorDashboardModel(owner, repo string, interval time.Duration) MonitorDashboardModel {
	return MonitorDashboardModel{
		owner:         owner,
		repo:          repo,
		repoName:      fmt.Sprintf("%s/%s", owner, repo),
		interval:      interval,
		notifications: []monitor.Notification{},
		autoScroll:    true,
		monitorRunID:  1,
	}
}

// MonitorUpdateMsg carries monitoring updates
type MonitorUpdateMsg struct {
	notification monitor.Notification
	runID        int
}

// MonitorErrorMsg carries monitoring errors
type MonitorErrorMsg struct {
	err error
}

type monitorStartedMsg struct {
	monitor *monitor.Monitor
	runID   int
}

type monitorStoppedMsg struct {
	runID int
}

func (m MonitorDashboardModel) Init() tea.Cmd {
	return m.startMonitoring(m.monitorRunID)
}

func (m MonitorDashboardModel) startMonitoring(runID int) tea.Cmd {
	return func() tea.Msg {
		mon, err := monitor.NewMonitor(m.owner, m.repo, m.interval)
		if err != nil {
			return MonitorErrorMsg{err: err}
		}

		return monitorStartedMsg{
			monitor: mon,
			runID:   runID,
		}
	}
}

func waitMonitorUpdate(mon *monitor.Monitor, runID int) tea.Cmd {
	return func() tea.Msg {
		notification, ok := <-mon.Notifications()
		if !ok {
			return monitorStoppedMsg{runID: runID}
		}
		return MonitorUpdateMsg{
			notification: notification,
			runID:        runID,
		}
	}
}

func (m *MonitorDashboardModel) StopMonitoring() {
	if m.monitor != nil {
		m.monitor.Stop()
		m.monitor = nil
	}
	m.isMonitoring = false
}

func (m MonitorDashboardModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case monitorStartedMsg:
		if msg.runID != m.monitorRunID {
			msg.monitor.Stop()
			return m, nil
		}
		if m.monitor != nil {
			m.monitor.Stop()
		}
		m.monitor = msg.monitor
		m.monitor.StartAsync()
		m.isMonitoring = true
		m.err = nil

		m.notifications = append([]monitor.Notification{{
			Type:      "info",
			Title:     "Monitoring Started",
			Message:   fmt.Sprintf("Monitoring %s every %v", m.repoName, m.interval),
			Timestamp: time.Now(),
			Severity:  "info",
		}}, m.notifications...)

		if len(m.notifications) > 50 {
			m.notifications = m.notifications[:50]
		}

		return m, waitMonitorUpdate(m.monitor, m.monitorRunID)

	case MonitorUpdateMsg:
		if msg.runID != m.monitorRunID {
			return m, nil
		}
		m.notifications = append([]monitor.Notification{msg.notification}, m.notifications...)

		// Keep only last 50 notifications
		if len(m.notifications) > 50 {
			m.notifications = m.notifications[:50]
		}

		// Add to global notifications
		AddMonitorNotification(m.repoName, msg.notification.Message)

		// Auto-scroll to top if enabled
		if m.autoScroll {
			m.scrollOffset = 0
		}

		return m, waitMonitorUpdate(m.monitor, m.monitorRunID)

	case monitorStoppedMsg:
		if msg.runID != m.monitorRunID {
			return m, nil
		}
		m.isMonitoring = false

	case MonitorErrorMsg:
		m.err = msg.err
		m.isMonitoring = false

	case tea.KeyMsg:
		maxNotifications := len(m.notifications)
		visibleLines := m.height - 12
		if visibleLines < 1 {
			visibleLines = 1
		}

		switch msg.String() {
		case "up", "k":
			if m.scrollOffset > 0 {
				m.scrollOffset--
			}
		case "down", "j":
			if m.scrollOffset < maxNotifications-visibleLines {
				m.scrollOffset++
			}
		case "home", "g":
			m.scrollOffset = 0
		case "end", "G":
			m.scrollOffset = maxInt(0, maxNotifications-visibleLines)
		case "a":
			// Toggle auto-scroll
			m.autoScroll = !m.autoScroll
		case "r":
			// Refresh/restart monitoring
			m.StopMonitoring()
			m.monitorRunID++
			return m, m.startMonitoring(m.monitorRunID)
		case "c":
			// Clear notifications
			m.notifications = []monitor.Notification{}
			m.scrollOffset = 0
		}

	}

	return m, nil
}

func (m MonitorDashboardModel) View() string {
	return m.ViewWithSize(m.width, m.height)
}

func (m MonitorDashboardModel) ViewWithSize(width, height int) string {
	m.width = width
	m.height = height

	header := TitleStyle.Render(fmt.Sprintf("👀 MONITORING: %s", m.repoName))

	// Status bar
	statusParts := []string{
		fmt.Sprintf("Interval: %v", m.interval),
		fmt.Sprintf("Notifications: %d", len(m.notifications)),
	}

	if m.autoScroll {
		statusParts = append(statusParts, "Auto-scroll: ON")
	} else {
		statusParts = append(statusParts, "Auto-scroll: OFF")
	}

	statusBar := SubtleStyle.Render(strings.Join(statusParts, " • "))

	// Error display
	if m.err != nil {
		errorBox := ErrorStyle.Render(fmt.Sprintf("Error: %v", m.err))
		footer := SubtleStyle.Render("\nr: retry • q/ESC: back to menu")

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			"\n",
			statusBar,
			"\n",
			BoxStyle.Render(errorBox),
			footer,
		)

		return lipgloss.Place(
			width, height,
			lipgloss.Center, lipgloss.Center,
			content,
		)
	}

	// No notifications yet
	if len(m.notifications) == 0 {
		waitingMsg := SubtleStyle.Render("Waiting for updates...\n\nMonitoring for:\n• New commits\n• Issues and pull requests\n• Contributor changes\n• Repository health metrics")
		footer := SubtleStyle.Render("\na: toggle auto-scroll • c: clear • r: refresh • q/ESC: back")

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			"\n",
			statusBar,
			"\n",
			BoxStyle.Render(waitingMsg),
			footer,
		)

		return lipgloss.Place(
			width, height,
			lipgloss.Center, lipgloss.Center,
			content,
		)
	}

	// Build notification list
	var items []string
	visibleLines := height - 12
	if visibleLines < 1 {
		visibleLines = 1
	}
	start := m.scrollOffset
	end := minInt(start+visibleLines, len(m.notifications))

	for i := start; i < end; i++ {
		notif := m.notifications[i]

		// Icon based on severity
		var icon string
		var style lipgloss.Style
		switch notif.Severity {
		case "error":
			icon = "❌"
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
		case "warning":
			icon = "⚠️"
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("11"))
		default:
			icon = "ℹ️"
			style = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
		}

		timestamp := notif.Timestamp.Format("15:04:05")
		line := fmt.Sprintf("[%s] %s %s: %s",
			timestamp,
			icon,
			style.Render(notif.Title),
			notif.Message,
		)

		items = append(items, line)
	}

	listContent := strings.Join(items, "\n")

	// Scroll indicator
	scrollInfo := ""
	if len(m.notifications) > visibleLines {
		scrollInfo = SubtleStyle.Render(fmt.Sprintf("\n[%d-%d of %d]", start+1, end, len(m.notifications)))
	}

	footer := SubtleStyle.Render("↑↓/jk: scroll • a: toggle auto-scroll • c: clear • r: refresh • q/ESC: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		"\n",
		statusBar,
		"\n",
		BoxStyle.Render(listContent),
		scrollInfo,
		"\n"+footer,
	)

	return lipgloss.Place(
		width, height,
		lipgloss.Center, lipgloss.Center,
		content,
	)
}
