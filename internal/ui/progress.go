package ui

import (
	"sync"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// ProgressStage represents a step in the analysis process
type ProgressStage struct {
	Name       string
	IsComplete bool
	IsActive   bool
}

// ProgressTracker manages multi-step analysis progress
type ProgressTracker struct {
	mu        sync.Mutex
	stages    []ProgressStage
	current   int
	startTime time.Time
}

var SatelliteFrames = []string{
	`
        .
       / \
      | . |
       \ /
        '
    Scanning...
	`,
	`
        o
       / \
      | o |
       \ /
        o
    Scanning...
	`,
	`
        O
       / \
      | O |
       \ /
        O
    Scanning...
	`,
	`
        @
       / \
      | @ |
       \ /
        @
    Scanning...
	`,
	`
        O
       / \
      | O |
       \ /
        O
    Scanning...
	`,
	`
        o
       / \
      | o |
       \ /
        o
    Scanning...
	`,
}

// NewProgressTracker creates a tracker with default analysis stages
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		stages: []ProgressStage{
			{Name: "🔗 Fetching repository metadata", IsComplete: false, IsActive: true},
			{Name: "📝 Fetching commit history", IsComplete: false, IsActive: false},
			{Name: "👥 Fetching contributors & activity", IsComplete: false, IsActive: false},
			{Name: "📂 Parsing directory structure & languages", IsComplete: false, IsActive: false},
			{Name: "📊 Calculating repository health metrics", IsComplete: false, IsActive: false},
			{Name: "🔍 Scanning dependencies & security", IsComplete: false, IsActive: false},
			{Name: "📈 Identifying hotspots & anomalies", IsComplete: false, IsActive: false},
			{Name: "📋 Gathering issues & pull requests", IsComplete: false, IsActive: false},
			{Name: "✅ Analysis complete", IsComplete: false, IsActive: false},
		},
		current:   0,
		startTime: time.Now(),
	}
}

// NextStage moves to the next analysis stage
func (pt *ProgressTracker) NextStage() {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	if pt.current < len(pt.stages) {
		pt.stages[pt.current].IsComplete = true
		pt.stages[pt.current].IsActive = false
		pt.current++
		if pt.current < len(pt.stages) {
			pt.stages[pt.current].IsActive = true
		}
	}
}

// GetCurrentStage returns the current stage information
func (pt *ProgressTracker) GetCurrentStage() ProgressStage {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	if pt.current < len(pt.stages) {
		return pt.stages[pt.current]
	}
	return ProgressStage{Name: "Complete", IsComplete: true, IsActive: false}
}

// GetAllStages returns all stages with their status
func (pt *ProgressTracker) GetAllStages() []ProgressStage {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	stagesCopy := make([]ProgressStage, len(pt.stages))
	copy(stagesCopy, pt.stages)
	return stagesCopy
}

// GetProgress returns completed stages / total stages
func (pt *ProgressTracker) GetProgress() (completed int, total int) {
	pt.mu.Lock()
	defer pt.mu.Unlock()
	total = len(pt.stages)
	for _, stage := range pt.stages {
		if stage.IsComplete {
			completed++
		}
	}
	return
}

// GetProgressBar returns a visual progress bar with a shimmering skeleton effect
func (pt *ProgressTracker) GetProgressBar(width int) string {
	completed, total := pt.GetProgress()
	if width < 10 {
		width = 10
	}

	fillWidth := (completed * width) / total
	emptyWidth := width - fillWidth

	fill := ""
	for i := 0; i < fillWidth; i++ {
		fill += "█"
	}

	// SKELETON EFFECT: Create a shimmering effect in the empty area
	elapsedMs := time.Since(pt.startTime).Milliseconds()
	shimmerPos := int((elapsedMs / 150) % int64(emptyWidth+1))

	empty := ""
	for i := 0; i < emptyWidth; i++ {
		if i == shimmerPos || i == shimmerPos-1 {
			empty += "▒" // The "shimmer" highlight
		} else {
			empty += "░" // The standard background
		}
	}

	return "[" + fill + empty + "] "
}

// GetElapsedTime returns how long the analysis has been running
func (pt *ProgressTracker) GetElapsedTime() time.Duration {
	return time.Since(pt.startTime)
}

// TickProgressCmd returns a command that ticks every 150ms for smoother skeleton animation
func TickProgressCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*150, func(t time.Time) tea.Msg {
		return struct{}{} // Progress tick message
	})
}
