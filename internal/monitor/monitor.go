package monitor

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/cache"
	"github.com/agnivo988/Repo-lyzer/internal/github"
)

const monitorStateTTL = 100 * 365 * 24 * time.Hour

// MonitorState represents the current state of a monitored repository
type MonitorState struct {
	Owner                string    `json:"owner"`
	Repo                 string    `json:"repo"`
	LastCommitSHA        string    `json:"last_commit_sha"`
	LastIssueCount       int       `json:"last_issue_id"`
	LastPRCount          int       `json:"last_pr_id"`
	LastContributorCount int       `json:"last_contributor_count"`
	LastUpdated          time.Time `json:"last_updated"`
}

// Monitor manages real-time monitoring of a GitHub repository
type Monitor struct {
	client        *github.Client
	cache         *cache.Cache
	owner         string
	repo          string
	interval      time.Duration
	state         *MonitorState
	stateMutex    sync.RWMutex
	ctx           context.Context
	cancel        context.CancelFunc
	wg            sync.WaitGroup
	notifications chan Notification
	startOnce     sync.Once
}

// Notification represents a monitoring notification
type Notification struct {
	Type      string    `json:"type"` // "commit", "issue", "pr", "contributor"
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Severity  string    `json:"severity"` // "info", "warning", "error"
}

// NewMonitor creates a new repository monitor
func NewMonitor(owner, repo string, interval time.Duration) (*Monitor, error) {
	client := github.NewClient()
	cache, err := cache.NewCache()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize cache: %w", err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &Monitor{
		client:        client,
		cache:         cache,
		owner:         owner,
		repo:          repo,
		interval:      interval,
		state:         &MonitorState{Owner: owner, Repo: repo},
		ctx:           ctx,
		cancel:        cancel,
		notifications: make(chan Notification, 100),
	}, nil
}

// Start begins the monitoring process
func (m *Monitor) Start() error {
	// Start monitoring loop (non-blocking setup)
	m.StartAsync()

	// Start CLI notification handler
	go m.handleNotifications()

	// Wait for interrupt signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-sigChan:
		fmt.Println("\n🛑 Stopping monitoring...")
		m.Stop()
	case <-m.ctx.Done():
		// Context cancelled
	}

	m.wg.Wait()
	return nil
}

// StartAsync begins monitoring without waiting for OS signals.
func (m *Monitor) StartAsync() {
	m.startOnce.Do(func() {
		m.loadState()
		m.wg.Add(1)
		go m.monitorLoop()
	})
}

// Notifications exposes monitoring notifications as a read-only channel.
func (m *Monitor) Notifications() <-chan Notification {
	return m.notifications
}

// Stop stops the monitoring process
func (m *Monitor) Stop() {
	m.cancel()
}

// monitorLoop runs the main monitoring loop
func (m *Monitor) monitorLoop() {
	defer m.wg.Done()
	defer close(m.notifications)

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	// Initial check
	m.checkForUpdates()

	for {
		select {
		case <-m.ctx.Done():
			return
		case <-ticker.C:
			m.checkForUpdates()
		}
	}
}

// checkForUpdates performs the actual monitoring checks
func (m *Monitor) checkForUpdates() {
	m.stateMutex.Lock()
	defer m.stateMutex.Unlock()

	// Check for new commits
	m.checkCommits()

	// Check for new issues
	m.checkIssues()

	// Check for new pull requests
	m.checkPullRequests()

	// Check for contributor changes
	m.checkContributors()

	// Save state
	m.saveState()
}

// checkCommits monitors for new commits
func (m *Monitor) checkCommits() {
	commits, err := m.client.GetCommits(m.owner, m.repo, 1) // Get latest commit
	if err != nil {
		log.Printf("Failed to get commits: %v", err)
		return
	}

	if len(commits) > 0 {
		latestCommit := commits[0]
		if latestCommit.SHA != m.state.LastCommitSHA {
			// New commit detected
			notification := Notification{
				Type:      "commit",
				Title:     "New Commit",
				Message:   fmt.Sprintf("New commit: %s", latestCommit.SHA[:8]),
				Timestamp: time.Now(),
				Severity:  "info",
			}
			m.notifications <- notification

			m.state.LastCommitSHA = latestCommit.SHA
			m.state.LastUpdated = time.Now()
		}
	}
}

// countIssues returns the number of true issues (PullRequest == nil) in the list.
func countIssues(items []github.Issue) int {
	count := 0
	for _, item := range items {
		if item.PullRequest == nil {
			count++
		}
	}
	return count
}

// countPRs returns the number of pull requests (PullRequest != nil) in the list.
func countPRs(items []github.Issue) int {
	count := 0
	for _, item := range items {
		if item.PullRequest != nil {
			count++
		}
	}
	return count
}

// checkIssues monitors for new issues (excludes pull requests)
func (m *Monitor) checkIssues() {
	issues, err := m.client.GetIssues(m.owner, m.repo, "open")
	if err != nil {
		log.Printf("Failed to get issues: %v", err)
		return
	}

	actualIssuesCount := countIssues(issues)

	// Only notify when the count actually changes
	if actualIssuesCount != m.state.LastIssueCount {
		notification := Notification{
			Type:      "issue",
			Title:     "Issues Update",
			Message:   fmt.Sprintf("Repository has %d open issues", actualIssuesCount),
			Timestamp: time.Now(),
			Severity:  "info",
		}
		m.notifications <- notification
		m.state.LastIssueCount = actualIssuesCount
	}
}

// checkPullRequests monitors for new pull requests (filtered and state-tracked)
func (m *Monitor) checkPullRequests() {
	issues, err := m.client.GetIssues(m.owner, m.repo, "open")
	if err != nil {
		log.Printf("Failed to get pull requests: %v", err)
		return
	}

	actualPRCount := countPRs(issues)

	// Only notify when the PR count actually changes
	if actualPRCount != m.state.LastPRCount {
		notification := Notification{
			Type:      "pr",
			Title:     "Pull Requests Update",
			Message:   fmt.Sprintf("Repository has %d open pull requests", actualPRCount),
			Timestamp: time.Now(),
			Severity:  "info",
		}
		m.notifications <- notification
		m.state.LastPRCount = actualPRCount
	}
}

// checkContributors monitors for contributor changes
func (m *Monitor) checkContributors() {
	contributors, err := m.client.GetContributors(m.owner, m.repo)
	if err != nil {
		log.Printf("Failed to get contributors: %v", err)
		return
	}

	// Check if the contributor count has changed
	if len(contributors) != m.state.LastContributorCount {
		notification := Notification{
			Type:      "contributor",
			Title:     "Contributor Update",
			Message:   fmt.Sprintf("Repository has %d contributors", len(contributors)),
			Timestamp: time.Now(),
			Severity:  "info",
		}
		m.notifications <- notification
		m.state.LastContributorCount = len(contributors)
	}
}

// handleNotifications processes incoming notifications
func (m *Monitor) handleNotifications() {
	for notification := range m.notifications {
		m.displayNotification(notification)
	}
}

// displayNotification displays a notification to the user
func (m *Monitor) displayNotification(n Notification) {
	timestamp := n.Timestamp.Format("15:04:05")
	var icon string

	switch n.Severity {
	case "error":
		icon = "❌"
	case "warning":
		icon = "⚠️"
	default:
		icon = "ℹ️"
	}

	fmt.Printf("[%s] %s %s: %s\n", timestamp, icon, n.Title, n.Message)
}

// loadState loads the monitoring state from cache
func (m *Monitor) loadState() {
	m.stateMutex.Lock()
	defer m.stateMutex.Unlock()

	key := fmt.Sprintf("%s/%s", m.owner, m.repo)
	entry, found := m.cache.Get(key)
	if !found || entry == nil || len(entry.Analysis) == 0 {
		return
	}

	var cachedState MonitorState
	if err := json.Unmarshal(entry.Analysis, &cachedState); err != nil {
		log.Printf("Failed to restore monitoring state for %s: %v", key, err)
		return
	}

	// Ensure monitor identity remains consistent even with older cache payloads.
	if cachedState.Owner == "" {
		cachedState.Owner = m.owner
	}
	if cachedState.Repo == "" {
		cachedState.Repo = m.repo
	}
	if cachedState.Owner != m.owner || cachedState.Repo != m.repo {
		log.Printf("Ignoring monitoring state for mismatched repository: %s/%s", cachedState.Owner, cachedState.Repo)
		return
	}

	m.state = &cachedState
}

// saveState saves the monitoring state to cache
func (m *Monitor) saveState() {
	key := fmt.Sprintf("%s/%s", m.owner, m.repo)
	if m.state == nil {
		return
	}

	stateCopy := *m.state
	if stateCopy.Owner == "" {
		stateCopy.Owner = m.owner
	}
	if stateCopy.Repo == "" {
		stateCopy.Repo = m.repo
	}

	if err := m.cache.SetWithTTL(key, stateCopy, monitorStateTTL); err != nil {
		log.Printf("Failed to persist monitoring state for %s: %v", key, err)
	}
}
