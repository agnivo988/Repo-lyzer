package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

// ProgressStage represents a step in the analysis process
type ProgressStage struct {
	Name       string
	IsComplete bool
	IsActive   bool
}

// ProgressTracker manages multi-step analysis progress
type ProgressTracker struct {
	stages   []ProgressStage
	current  int
	startTime time.Time
}

// ProgressUpdateMsg is sent to update progress
type ProgressUpdateMsg struct {
	StageIndex int
	IsComplete bool
}

// NewProgressTracker creates a tracker with default analysis stages
func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		stages: []ProgressStage{
			{Name: "🔗 Fetching repository data", IsComplete: false, IsActive: true},
			{Name: "📝 Analyzing commits", IsComplete: false, IsActive: false},
			{Name: "👥 Analyzing contributors", IsComplete: false, IsActive: false},
			{Name: "🗣️  Analyzing languages", IsComplete: false, IsActive: false},
			{Name: "📊 Computing metrics", IsComplete: false, IsActive: false},
			{Name: "✅ Analysis complete", IsComplete: false, IsActive: false},
		},
		current:   0,
		startTime: time.Now(),
	}
}

// NextStage moves to the next analysis stage
func (pt *ProgressTracker) NextStage() {
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
	if pt.current < len(pt.stages) {
		return pt.stages[pt.current]
	}
	return ProgressStage{Name: "Complete", IsComplete: true, IsActive: false}
}

// GetAllStages returns all stages with their status
func (pt *ProgressTracker) GetAllStages() []ProgressStage {
	return pt.stages
}

// GetProgress returns completed stages / total stages
func (pt *ProgressTracker) GetProgress() (completed int, total int) {
	total = len(pt.stages)
	for _, stage := range pt.stages {
		if stage.IsComplete {
			completed++
		}
	}
	return
}

// GetProgressBar returns a visual progress bar string
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

	empty := ""
	for i := 0; i < emptyWidth; i++ {
		empty += "░"
	}

	percentage := (completed * 100) / total

	return "[" + fill + empty + "] " + string(rune(percentage)) + "%"
}

// GetElapsedTime returns how long the analysis has been running
func (pt *ProgressTracker) GetElapsedTime() time.Duration {
	return time.Since(pt.startTime)
}

// TickProgressCmd returns a command that ticks every 300ms to update progress display
func TickProgressCmd() tea.Cmd {
	return tea.Tick(time.Millisecond*300, func(t time.Time) tea.Msg {
		return struct{}{} // Progress tick message
	})
}
