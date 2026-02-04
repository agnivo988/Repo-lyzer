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

// MonitorState represents the current state of a monitored repository
type MonitorState struct {
	Owner               string    `json:"owner"`
	Repo                string    `json:"repo"`
	LastCommitSHA       string    `json:"last_commit_sha"`
	LastIssueID         int       `json:"last_issue_id"`
	LastIssueCount      int       `json:"last_issue_count"`
	LastPRID            int       `json:"last_pr_id"`
	LastPRCount         int       `json:"last_pr_count"`
	LastContributorCount int      `json:"last_contributor_count"`
	LastUpdated         time.Time `json:"last_updated"`
}

// Monitor manages real-time monitoring of a GitHub repository
type Monitor struct {
	client                    *github.Client
	cache                     *cache.Cache
	owner                     string
	repo                      string
	interval                  time.Duration
	state                     *MonitorState
	stateMutex                sync.RWMutex
	ctx                       context.Context
	cancel                    context.CancelFunc
	wg                        sync.WaitGroup
	notifications             chan Notification
	notificationsCloseOnce    sync.Once
	prevContribCount          int
}

// Notification represents a monitoring notification
type Notification struct {
	Type      string    `json:"type"`      // "commit", "issue", "pr", "contributor"
	Title     string    `json:"title"`
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
	Severity  string    `json:"severity"` // "info", "warning", "error"
}

// NewMonitor creates a new repository monitor
func NewMonitor(owner, repo string, interval time.Duration) (*Monitor, error) {
	// Validate interval parameter
	if interval <= 0 {
		return nil, fmt.Errorf("invalid interval: must be > 0")
	}

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

// Start begins the monitoring process and spawns goroutines
// Use Wait() to block until the monitor is stopped
func (m *Monitor) Start() error {
	// Load previous state
	m.loadState()

	// Start notification handler
	m.wg.Add(1)
	go m.handleNotifications()

	// Start monitoring loop
	m.wg.Add(1)
	go m.monitorLoop()

	// Start signal handler in a separate goroutine
	m.wg.Add(1)
	go m.handleSignals()

	return nil
}

// Wait blocks until the monitoring process is stopped via signal or context cancellation
func (m *Monitor) Wait() {
	m.wg.Wait()
}

// handleSignals waits for interrupt signals and gracefully stops monitoring
func (m *Monitor) handleSignals() {
	defer m.wg.Done()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer signal.Stop(sigChan)

	<-sigChan
	fmt.Println("\n🛑 Stopping monitoring...")
	// Call cancel directly instead of Stop to avoid deadlock (Stop calls wg.Wait())
	m.cancel()
}

// Stop stops the monitoring process
func (m *Monitor) Stop() {
	m.cancel()
	// Use sync.Once to safely close the notifications channel
	m.notificationsCloseOnce.Do(func() {
		close(m.notifications)
	})
	m.wg.Wait()
}

// monitorLoop runs the main monitoring loop
func (m *Monitor) monitorLoop() {
	defer m.wg.Done()

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
	select {
	case <-m.ctx.Done():
		return
	default:
	}

	// Perform all checks without holding the mutex (to avoid deadlock with network calls)
	var notifications []Notification

	// Check for new commits
	commitNotifs := m.checkCommits()
	notifications = append(notifications, commitNotifs...)

	// Check for new issues
	issueNotifs := m.checkIssues()
	notifications = append(notifications, issueNotifs...)

	// Check for new pull requests
	prNotifs := m.checkPullRequests()
	notifications = append(notifications, prNotifs...)

	// Check for contributor changes
	contribNotifs := m.checkContributors()
	notifications = append(notifications, contribNotifs...)

	// Now acquire mutex to save state and send notifications
	m.stateMutex.Lock()
	m.saveState()
	m.stateMutex.Unlock()

	// Send notifications after releasing mutex
	for _, notif := range notifications {
		select {
		case <-m.ctx.Done():
			return
		case m.notifications <- notif:
		}
	}
}

// checkCommits monitors for new commits
func (m *Monitor) checkCommits() []Notification {
	var notifs []Notification

	commits, err := m.client.GetCommits(m.owner, m.repo, 1) // Get latest commit
	if err != nil {
		log.Printf("Failed to get commits: %v", err)
		return notifs
	}

	if len(commits) > 0 {
		latestCommit := commits[0]
		m.stateMutex.RLock()
		lastSHA := m.state.LastCommitSHA
		m.stateMutex.RUnlock()

		if latestCommit.SHA != lastSHA {
			// New commit detected
			notification := Notification{
				Type:      "commit",
				Title:     "New Commit",
				Message:   fmt.Sprintf("New commit: %s", truncateSHA(latestCommit.SHA, 8)),
				Timestamp: time.Now(),
				Severity:  "info",
			}
			notifs = append(notifs, notification)

			m.stateMutex.Lock()
			m.state.LastCommitSHA = latestCommit.SHA
			m.state.LastUpdated = time.Now()
			m.stateMutex.Unlock()
		}
	}

	return notifs
}

// checkIssues monitors for new issues
func (m *Monitor) checkIssues() []Notification {
	var notifs []Notification

	issues, err := m.client.GetIssues(m.owner, m.repo, "open")
	if err != nil {
		log.Printf("Failed to get issues: %v", err)
		return notifs
	}

	// Get current count
	currentIssueCount := len(issues)

	// Check if the number of issues has changed
	m.stateMutex.RLock()
	lastIssueCount := m.state.LastIssueCount
	m.stateMutex.RUnlock()

	if currentIssueCount != lastIssueCount {
		notification := Notification{
			Type:      "issue",
			Title:     "Issues Update",
			Message:   fmt.Sprintf("Repository has %d open issues", currentIssueCount),
			Timestamp: time.Now(),
			Severity:  "info",
		}
		notifs = append(notifs, notification)

		m.stateMutex.Lock()
		m.state.LastIssueCount = currentIssueCount
		m.stateMutex.Unlock()
	}

	return notifs
}

// checkPullRequests monitors for new pull requests
func (m *Monitor) checkPullRequests() []Notification {
	var notifs []Notification

	// Use dedicated GetPullRequests API method to avoid duplicate API calls
	prs, err := m.client.GetPullRequests(m.owner, m.repo, "open")
	if err != nil {
		log.Printf("Failed to get pull requests: %v", err)
		return notifs
	}

	// Check if PR count changed
	currentPRCount := len(prs)

	m.stateMutex.RLock()
	lastPRCount := m.state.LastPRCount
	m.stateMutex.RUnlock()

	if currentPRCount != lastPRCount {
		notification := Notification{
			Type:      "pr",
			Title:     "Pull Requests Update",
			Message:   fmt.Sprintf("Repository has %d open pull requests", currentPRCount),
			Timestamp: time.Now(),
			Severity:  "info",
		}
		notifs = append(notifs, notification)

		m.stateMutex.Lock()
		m.state.LastPRCount = currentPRCount
		m.stateMutex.Unlock()
	}

	return notifs
}

// checkContributors monitors for contributor changes
func (m *Monitor) checkContributors() []Notification {
	var notifs []Notification

	contributors, err := m.client.GetContributors(m.owner, m.repo)
	if err != nil {
		log.Printf("Failed to get contributors: %v", err)
		return notifs
	}

	// Get current count
	currentCount := len(contributors)

	// Check if the contributor count has changed
	m.stateMutex.RLock()
	lastContributorCount := m.state.LastContributorCount
	m.stateMutex.RUnlock()

	if currentCount != lastContributorCount {
		notification := Notification{
			Type:      "contributor",
			Title:     "Contributor Update",
			Message:   fmt.Sprintf("Repository has %d contributors", currentCount),
			Timestamp: time.Now(),
			Severity:  "info",
		}
		notifs = append(notifs, notification)

		m.stateMutex.Lock()
		m.state.LastContributorCount = currentCount
		m.stateMutex.Unlock()
	}

	return notifs
}

// handleNotifications processes incoming notifications
func (m *Monitor) handleNotifications() {
	defer m.wg.Done()
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
	cachedEntry, found := m.cache.Get(key)
	if found {
		// Deserialize the cached state
		var cachedState MonitorState
		err := json.Unmarshal(cachedEntry.Analysis, &cachedState)
		if err != nil {
			log.Printf("Failed to unmarshal cached state: %v, using fresh state", err)
			m.state.LastUpdated = time.Now()
			return
		}
		// Restore the cached state
		m.state = &cachedState
	} else {
		// Cache miss, initialize with current time
		m.state.LastUpdated = time.Now()
	}
}

// saveState saves the monitoring state to cache
func (m *Monitor) saveState() {
	key := fmt.Sprintf("%s/%s", m.owner, m.repo)
	jsonData, err := json.Marshal(m.state)
	if err != nil {
		log.Printf("Failed to marshal state: %v", err)
		return
	}
	m.cache.Set(key, jsonData)
}

// truncateSHA truncates a SHA to a maximum length
func truncateSHA(sha string, maxLen int) string {
	if len(sha) <= maxLen {
		return sha
	}
	return sha[:maxLen]
}
