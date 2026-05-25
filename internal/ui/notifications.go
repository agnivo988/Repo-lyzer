package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// NotificationEntry represents a single notification
type NotificationEntry struct {
	Type      string    `json:"type"` // "analysis", "export", "monitor"
	RepoName  string    `json:"repo_name"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"` // "success", "error", "info"
	Details   string    `json:"details"`
}

// NotificationStore manages notifications
type NotificationStore struct {
	Entries []NotificationEntry `json:"entries"`
}

// NotificationsModel handles the notifications view
type NotificationsModel struct {
	notifications *NotificationStore
	cursor        int
	width         int
	height        int
	err           error
	scrollOffset  int
}

// NewNotificationsModel creates a new notifications model
func NewNotificationsModel() NotificationsModel {
	notifications, _ := LoadNotifications()
	return NotificationsModel{
		notifications: notifications,
		cursor:        0,
		scrollOffset:  0,
	}
}

func (m NotificationsModel) Init() tea.Cmd {
	return nil
}

func (m NotificationsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

	case tea.KeyMsg:
		maxEntries := len(m.notifications.Entries)
		visibleLines := m.height - 10 // Account for header and footer
		if visibleLines < 1 {
			visibleLines = 1
		}

		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
				if m.cursor < m.scrollOffset {
					m.scrollOffset = m.cursor
				}
			}
		case "down", "j":
			if m.cursor < maxEntries-1 {
				m.cursor++
				if m.cursor >= m.scrollOffset+visibleLines {
					m.scrollOffset = m.cursor - visibleLines + 1
				}
			}
		case "home", "g":
			m.cursor = 0
			m.scrollOffset = 0
		case "end", "G":
			if maxEntries > 0 {
				m.cursor = maxEntries - 1
				m.scrollOffset = maxInt(0, maxEntries-visibleLines)
			}
		case "d":
			// Delete selected notification
			if maxEntries > 0 {
				m.notifications.Delete(m.cursor)
				m.notifications.Save()
				if m.cursor >= len(m.notifications.Entries) && m.cursor > 0 {
					m.cursor--
				}
			}
		case "c":
			// Clear all notifications
			m.notifications.Clear()
			m.notifications.Save()
			m.cursor = 0
			m.scrollOffset = 0
		}
	}

	return m, nil
}

func (m NotificationsModel) View() string {
	return m.ViewWithSize(m.width, m.height)
}

func (m NotificationsModel) ViewWithSize(width, height int) string {
	m.width = width
	m.height = height

	header := TitleStyle.Render("🔔 NOTIFICATIONS")

	if len(m.notifications.Entries) == 0 {
		emptyMsg := SubtleStyle.Render("No notifications yet.\n\nNotifications will appear here when you:\n• Analyze repositories\n• Export reports\n• Monitor repository changes")
		footer := SubtleStyle.Render("\nq/ESC: back to menu")

		content := lipgloss.JoinVertical(
			lipgloss.Left,
			header,
			"\n",
			BoxStyle.Render(emptyMsg),
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
	visibleLines := height - 10
	start := m.scrollOffset
	end := minInt(start+visibleLines, len(m.notifications.Entries))

	for i := start; i < end; i++ {
		entry := m.notifications.Entries[i]

		// Icon based on type and status
		var icon string
		switch entry.Type {
		case "analysis":
			icon = "📊"
		case "export":
			icon = "📤"
		case "monitor":
			icon = "👀"
		default:
			icon = "ℹ️"
		}

		// Status indicator
		var statusStyle lipgloss.Style
		switch entry.Status {
		case "success":
			statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
		case "error":
			statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))
		default:
			statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
		}

		timestamp := entry.Timestamp.Format("Jan 02, 15:04")
		line := fmt.Sprintf("%s %s %s - %s",
			icon,
			statusStyle.Render(entry.Status),
			entry.RepoName,
			timestamp,
		)

		if entry.Message != "" {
			line += fmt.Sprintf("\n   %s", SubtleStyle.Render(entry.Message))
		}

		if i == m.cursor {
			items = append(items, SelectedStyle.Render(line))
		} else {
			items = append(items, NormalStyle.Render(line))
		}
		items = append(items, "") // Add spacing
	}

	listContent := strings.Join(items, "\n")

	// Scroll indicator
	scrollInfo := ""
	if len(m.notifications.Entries) > visibleLines {
		scrollInfo = SubtleStyle.Render(fmt.Sprintf("\n[%d-%d of %d]", start+1, end, len(m.notifications.Entries)))
	}

	footer := SubtleStyle.Render("↑↓/jk: navigate • d: delete • c: clear all • q/ESC: back")

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
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

// LoadNotifications loads notifications from disk
func LoadNotifications() (*NotificationStore, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return &NotificationStore{Entries: []NotificationEntry{}}, err
	}

	notifPath := filepath.Join(home, ".repo-lyzer", "notifications.json")
	data, err := os.ReadFile(notifPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &NotificationStore{Entries: []NotificationEntry{}}, nil
		}
		return &NotificationStore{Entries: []NotificationEntry{}}, err
	}

	var store NotificationStore
	if err := json.Unmarshal(data, &store); err != nil {
		return &NotificationStore{Entries: []NotificationEntry{}}, err
	}

	// Sort by timestamp (newest first)
	sort.Slice(store.Entries, func(i, j int) bool {
		return store.Entries[i].Timestamp.After(store.Entries[j].Timestamp)
	})

	return &store, nil
}

// Save saves notifications to disk
func (ns *NotificationStore) Save() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	repoDir := filepath.Join(home, ".repo-lyzer")
	if err := os.MkdirAll(repoDir, 0755); err != nil {
		return err
	}

	notifPath := filepath.Join(repoDir, "notifications.json")
	data, err := json.MarshalIndent(ns, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(notifPath, data, 0644)
}

// AddNotification adds a new notification
func (ns *NotificationStore) AddNotification(notif NotificationEntry) {
	notif.Timestamp = time.Now()
	ns.Entries = append([]NotificationEntry{notif}, ns.Entries...)

	// Keep only last 100 notifications
	if len(ns.Entries) > 100 {
		ns.Entries = ns.Entries[:100]
	}
}

// Delete removes a notification at the given index
func (ns *NotificationStore) Delete(index int) {
	if index >= 0 && index < len(ns.Entries) {
		ns.Entries = append(ns.Entries[:index], ns.Entries[index+1:]...)
	}
}

// Clear removes all notifications
func (ns *NotificationStore) Clear() {
	ns.Entries = []NotificationEntry{}
}

// Helper functions
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// AddAnalysisNotification is a helper to add analysis notifications
func AddAnalysisNotification(repoName string, success bool) {
	store, _ := LoadNotifications()

	status := "success"
	message := "Repository analyzed successfully"
	if !success {
		status = "error"
		message = "Analysis failed"
	}

	store.AddNotification(NotificationEntry{
		Type:     "analysis",
		RepoName: repoName,
		Message:  message,
		Status:   status,
	})
	store.Save()
}

// AddExportNotification is a helper to add export notifications
func AddExportNotification(repoName, format string, success bool) {
	store, _ := LoadNotifications()

	status := "success"
	message := fmt.Sprintf("Exported to %s", format)
	if !success {
		status = "error"
		message = fmt.Sprintf("Export to %s failed", format)
	}

	store.AddNotification(NotificationEntry{
		Type:     "export",
		RepoName: repoName,
		Message:  message,
		Status:   status,
	})
	store.Save()
}

// AddMonitorNotification is a helper to add monitoring notifications
func AddMonitorNotification(repoName, message string) {
	store, _ := LoadNotifications()

	store.AddNotification(NotificationEntry{
		Type:     "monitor",
		RepoName: repoName,
		Message:  message,
		Status:   "info",
	})
	store.Save()
}
