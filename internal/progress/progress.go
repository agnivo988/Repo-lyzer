// Package progress provides CLI progress indicators and spinners for async user feedback during long-running operations.
package progress

import (
	"fmt"
	"sync"
	"time"
)

// Spinner represents a CLI spinner for displaying loading states
type Spinner struct {
	frames    []string
	delay     time.Duration
	message   string
	active    bool
	mu        sync.Mutex
	stopChan  chan bool
	completed bool
}

// NewSpinner creates a new spinner with default settings
func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{
			"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏",
		},
		delay:     100 * time.Millisecond,
		message:   "",
		stopChan:  make(chan bool),
		completed: false,
	}
}

// Start begins the spinner animation with a custom message
func (s *Spinner) Start(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.active {
		return
	}

	s.message = message
	s.active = true
	s.completed = false

	go s.animate()
}

// Stop stops the spinner animation
func (s *Spinner) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active {
		return
	}

	s.active = false
	s.stopChan <- true

	// Clear the line
	fmt.Print("\r" + clearLine())
}

// StopWithMessage stops the spinner and displays a completion message
func (s *Spinner) StopWithMessage(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.active {
		return
	}

	s.active = false
	s.stopChan <- true

	// Display the completion message
	fmt.Printf("\r✓ %s\033[K\n", message)
	s.completed = true
}

// Update changes the spinner message without restarting it
func (s *Spinner) Update(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.message = message
}

// animate runs the spinner animation loop
func (s *Spinner) animate() {
	frameIndex := 0

	for {
		select {
		case <-s.stopChan:
			return
		default:
			s.mu.Lock()
			if s.active {
				frame := s.frames[frameIndex%len(s.frames)]
				fmt.Printf("\r%s %s\033[K", frame, s.message)
			}
			s.mu.Unlock()

			frameIndex++
			time.Sleep(s.delay)
		}
	}
}

// ProgressBar represents a simple progress bar for tracking item completion
type ProgressBar struct {
	current   int
	total     int
	message   string
	active    bool
	mu        sync.Mutex
	startTime time.Time
}

// NewProgressBar creates a new progress bar
func NewProgressBar(total int, message string) *ProgressBar {
	return &ProgressBar{
		current:   0,
		total:     total,
		message:   message,
		active:    false,
		startTime: time.Now(),
	}
}

// Start initializes the progress bar display
func (p *ProgressBar) Start() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.active = true
	p.startTime = time.Now()
	p.display()
}

// Increment increases the progress by 1
func (p *ProgressBar) Increment() {
	p.mu.Lock()
	defer p.mu.Unlock()

	if p.current < p.total {
		p.current++
		p.display()
	}
}

// IncrementBy increases the progress by a specific amount
func (p *ProgressBar) IncrementBy(amount int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current = p.current + amount
	if p.current > p.total {
		p.current = p.total
	}
	p.display()
}

// SetCurrent sets the progress to a specific value
func (p *ProgressBar) SetCurrent(current int) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current = current
	if p.current > p.total {
		p.current = p.total
	}
	p.display()
}

// Stop completes the progress bar
func (p *ProgressBar) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.current = p.total
	p.display()
	fmt.Println()
	p.active = false
}

// display renders the progress bar
func (p *ProgressBar) display() {
	if !p.active {
		return
	}

	percentage := 0
	if p.total > 0 {
		percentage = int((p.current * 100) / p.total)
	}

	barWidth := 30
	filledWidth := (p.current * barWidth) / p.total
	if p.current > 0 && filledWidth == 0 {
		filledWidth = 1
	}

	bar := "["
	for i := 0; i < barWidth; i++ {
		if i < filledWidth {
			bar += "="
		} else if i == filledWidth && p.current > 0 && p.current < p.total {
			bar += ">"
		} else {
			bar += " "
		}
	}
	bar += "]"

	elapsed := time.Since(p.startTime).Seconds()
	eta := 0.0
	if p.current > 0 && percentage < 100 {
		eta = (elapsed / float64(p.current)) * float64(p.total-p.current)
	}

	fmt.Printf("\r%s %s %d%% (%.0fs/~%.0fs)", bar, p.message, percentage, elapsed, eta)
}

// Status represents a simple status message with emoji indicators
type Status struct {
	mu sync.Mutex
}

// NewStatus creates a new status display
func NewStatus() *Status {
	return &Status{}
}

// Info displays an informational message
func (s *Status) Info(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Printf("ℹ️  %s\n", message)
}

// Success displays a success message
func (s *Status) Success(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Printf("✓ %s\n", message)
}

// Error displays an error message
func (s *Status) Error(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Printf("✗ %s\n", message)
}

// Warning displays a warning message
func (s *Status) Warning(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Printf("⚠️  %s\n", message)
}

// Message displays a generic message
func (s *Status) Message(message string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	fmt.Printf("→ %s\n", message)
}

// clearLine returns ANSI code to clear current line
func clearLine() string {
	return "\033[K"
}

// OverallProgress tracks the overall progress of a multi-step analysis operation
// It displays both the current step spinner and an overall percentage progress
type OverallProgress struct {
	stepDone       chan struct{}
	totalSteps     int
	completedSteps int
	currentStep    string
	mu             sync.Mutex
	spinner        *Spinner
	startTime      time.Time
	stepStartTime  time.Time
}

// NewOverallProgress creates a new overall progress tracker
// Parameters:
//   - totalSteps: Total number of steps in the analysis
func NewOverallProgress(totalSteps int) *OverallProgress {
	return &OverallProgress{
		totalSteps:     totalSteps,
		completedSteps: 0,
		currentStep:    "",
		spinner:        NewSpinner(),
		startTime:      time.Now(),
		stepStartTime:  time.Now(),
	}
}

// StartStep begins a new analysis step with a message
// This updates the current step display and starts the spinner
func (o *OverallProgress) StartStep(stepMessage string) {
	o.mu.Lock()
	if o.stepDone != nil {
		close(o.stepDone)
	}

	o.stepDone = make(chan struct{})

	o.currentStep = stepMessage
	o.stepStartTime = time.Now()

	percentage := o.calculatePercentage()

	message := fmt.Sprintf(
		"%s [%d/%d - %d%%] (0s elapsed)",
		stepMessage,
		o.completedSteps+1,
		o.totalSteps,
		percentage,
	)

	o.spinner.Start(message)
	o.mu.Unlock()

	// Continuously update elapsed time while spinner is active
	go func(step string, done <-chan struct{}) {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return

			case <-ticker.C:
				o.mu.Lock()

				elapsed := int(time.Since(o.stepStartTime).Seconds())

				updatedMessage := fmt.Sprintf(
					"%s [%d/%d - %d%%] (%ds elapsed)",
					step,
					o.completedSteps+1,
					o.totalSteps,
					o.calculatePercentage(),
					elapsed,
				)

				o.spinner.Update(updatedMessage)
				o.mu.Unlock()
			}
		}
	}(stepMessage, o.stepDone)
}

// CompleteStep marks the current step as complete and moves to the next
func (o *OverallProgress) CompleteStep(stepMessage string) {
	o.mu.Lock()
	if o.stepDone != nil {
		close(o.stepDone)
		o.stepDone = nil
	}
	defer o.mu.Unlock()

	o.completedSteps++
	percentage := o.calculatePercentage()

	o.spinner.StopWithMessage(
		fmt.Sprintf(
			"%s [%d/%d - %d%%]",
			stepMessage,
			o.completedSteps,
			o.totalSteps,
			percentage,
		),
	)
}

// UpdateStep updates the current step message without stopping the spinner
func (o *OverallProgress) UpdateStep(stepMessage string) {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.currentStep = stepMessage
	percentage := o.calculatePercentage()
	elapsed := int(time.Since(o.stepStartTime).Seconds())

	message := fmt.Sprintf(
		"%s [%d/%d - %d%%] (%ds elapsed)",
		stepMessage,
		o.completedSteps+1,
		o.totalSteps,
		percentage,
		elapsed,
	)

	o.spinner.Update(message)
}

// Finish completes the analysis operation and displays final summary
func (o *OverallProgress) Finish() {
	o.mu.Lock()
	defer o.mu.Unlock()

	o.spinner.Stop()
	elapsed := time.Since(o.startTime).Seconds()
	fmt.Printf("✨ Analysis completed in %.2f seconds\n", elapsed)
}

// calculatePercentage returns the current completion percentage
func (o *OverallProgress) calculatePercentage() int {
	if o.totalSteps == 0 {
		return 0
	}
	return (o.completedSteps * 100) / o.totalSteps
}

// GetProgress returns current progress information
// Returns completedSteps, totalSteps, and percentage
func (o *OverallProgress) GetProgress() (int, int, int) {
	o.mu.Lock()
	defer o.mu.Unlock()

	return o.completedSteps, o.totalSteps, o.calculatePercentage()
}
