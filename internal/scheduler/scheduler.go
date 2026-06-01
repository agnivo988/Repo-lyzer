// Package scheduler provides automated report scheduling functionality for Repo-lyzer.
// It enables periodic analysis reports that run automatically and export results
// to various formats (JSON, PDF, Markdown) and destinations (local path, webhook).
package scheduler

import (
	"bytes"
	"context"

	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/agnivo988/Repo-lyzer/internal/analyzer"
	"github.com/agnivo988/Repo-lyzer/internal/config"
	"github.com/agnivo988/Repo-lyzer/internal/github"
	"github.com/agnivo988/Repo-lyzer/internal/monitor"
	"github.com/agnivo988/Repo-lyzer/internal/output"
	"github.com/robfig/cron/v3"
)

// Scheduler manages scheduled analysis report jobs
type Scheduler struct {
	cron           *cron.Cron
	settings       *config.AppSettings
	jobEntries     map[string]cron.EntryID
	reportExporter *ReportExporter

	// New fields for stabilization
	taskTimeout    time.Duration             // default 2m
	quotaTicker    *time.Ticker              // 1s ticker for adaptive quota
	stopChan       chan struct{}             // for graceful shutdown of ticker
	maxWorkers     int32                     // maximum concurrent workers
	activeWorkers  int32                     // current active workers (atomic)
	quotaTokens    int64                     // adaptive quota tokens
	maxQuota       int64                     // max quota tokens
	retryAttempts  map[string]int            // track retry count per job ID
	metrics        *monitor.SchedulerMetrics // telemetry counters
	wg             sync.WaitGroup            // wait for active jobs
	failureStreak  int                       // consecutive failures
	cooldownActive bool                      // indicates cooldown state
	cooldownUntil  time.Time                 // cooldown expiry
	mutex          sync.Mutex
}

// ReportExporter handles exporting analysis reports to various formats
type ReportExporter struct{}

// NewScheduler creates a new scheduler instance
func NewScheduler() (*Scheduler, error) {
	settings, err := config.LoadSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to load settings: %w", err)
	}

	// Default values
	defaultTimeout := 2 * time.Minute
	maxWorkers := int32(10) // sensible default, can be overridden via config later

	maxQuota := int64(maxWorkers)

	return &Scheduler{
		cron:           cron.New(),
		settings:       settings,
		jobEntries:     make(map[string]cron.EntryID),
		reportExporter: &ReportExporter{},
		taskTimeout:    defaultTimeout,
		quotaTicker:    nil,
		stopChan:       make(chan struct{}),
		maxWorkers:     maxWorkers,
		activeWorkers:  0,
		quotaTokens:    maxQuota,
		maxQuota:       maxQuota,
		retryAttempts:  make(map[string]int),
		metrics:        monitor.NewSchedulerMetrics(),
	}, nil
}

// Start initializes and starts the scheduler with all registered jobs
func (s *Scheduler) Start() {
	// Start quota recovery
	s.quotaTicker = time.NewTicker(1 * time.Second)
	go s.runQuotaRecovery()

	// Schedule existing jobs
	for _, job := range s.settings.GetScheduledJobs() {
		if job.Enabled {
			_ = s.scheduleJob(job)
		}
	}
	s.cron.Start()
	log.Println("Scheduler started")
}

// Stop stops the scheduler gracefully
func (s *Scheduler) Stop() {
	log.Println("Stopping scheduler...")
	s.cron.Stop()
	// Stop quota ticker gracefully
	if s.quotaTicker != nil {
		s.quotaTicker.Stop()
		close(s.stopChan)
	}
	// Wait for all running jobs to finish
	done := make(chan struct{})
	go func() { s.wg.Wait(); close(done) }()
	select {
	case <-done:
		log.Println("All jobs completed during shutdown")
	case <-time.After(10 * time.Second):
		log.Println("Timeout waiting for jobs, proceeding with shutdown")
	}
	log.Println("Scheduler stopped")
}

// scheduleJob adds a job to the cron scheduler
func (s *Scheduler) scheduleJob(job config.ScheduledJob) error {
	spec := job.GetCronExpression()

	jobFunc := func() {
		// Check cooldown before execution
		if s.isInCooldown() {
			log.Printf("Scheduler in cooldown, skipping job %s", job.ID)
			return
		}
		log.Printf("Executing scheduled job: %s for %s/%s", job.ID, job.Owner, job.Repo)

		s.wg.Add(1)
		defer s.wg.Done()

		if err := s.executeJobWithRetry(job); err != nil {
			log.Printf("Job execution failed after retries: %v", err)
		}
	}

	entryID, err := s.cron.AddFunc(spec, jobFunc)
	if err != nil {
		return fmt.Errorf("failed to add cron job: %w", err)
	}

	s.jobEntries[job.ID] = entryID
	log.Printf("Job %s scheduled with spec: %s", job.ID, spec)
	return nil
}

// executeJob runs the analysis and exports the report
func (s *Scheduler) executeJob(job config.ScheduledJob) error {
	startTime := time.Now()

	// Context with timeout for the whole job execution
	ctx, cancel := context.WithTimeout(context.Background(), s.taskTimeout)
	defer cancel()

	// Increment active workers counter
	atomic.AddInt32(&s.activeWorkers, 1)
	defer func() {
		atomic.AddInt32(&s.activeWorkers, -1)
	}()

	// Increment telemetry metric for active workers
	if s.metrics != nil {
		s.metrics.IncActiveWorkers()
	}

	// Initialize GitHub client
	client := github.NewClientWithContext(ctx)

	// Fetch repository information
	repoInfo, err := client.GetRepo(job.Owner, job.Repo)
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	// Fetch languages
	langs, err := client.GetLanguages(job.Owner, job.Repo)
	if err != nil {
		return fmt.Errorf("failed to get languages: %w", err)
	}

	// Fetch commits
	commits, err := client.GetCommits(job.Owner, job.Repo, 365)
	if err != nil {
		return fmt.Errorf("failed to get commits: %w", err)
	}

	// Fetch contributors
	contributors, err := client.GetContributors(job.Owner, job.Repo)
	if err != nil {
		return fmt.Errorf("failed to get contributors: %w", err)
	}

	// Calculate metrics
	healthScore := analyzer.CalculateHealth(repoInfo, commits)
	busFactor, busRisk := analyzer.BusFactor(contributors)
	maturityScore, maturityLevel := analyzer.RepoMaturityScore(repoInfo, len(commits), len(contributors), false)

	// Build compact config for export
	compactCfg := output.CompactConfig{
		Repo:            repoInfo,
		HealthScore:     healthScore,
		BusFactor:       busFactor,
		BusRisk:         busRisk,
		MaturityScore:   maturityScore,
		MaturityLevel:   maturityLevel,
		CommitsLastYear: len(commits),
		Contributors:    len(contributors),
		Duration:        time.Since(startTime),
		Languages:       langs,
	}

	// Export based on format
	var reportData []byte
	var filename string

	switch job.Format {
	case config.ExportJSON:
		reportData, err = s.reportExporter.exportJSON(compactCfg)
		filename = fmt.Sprintf("%s_%s_report.json", job.Owner+"-"+job.Repo, time.Now().Format("20060102"))
	case config.ExportMarkdown:
		reportData, err = s.reportExporter.exportMarkdown(compactCfg, repoInfo)
		filename = fmt.Sprintf("%s_%s_report.md", job.Owner+"-"+job.Repo, time.Now().Format("20060102"))
	case config.ExportPDF:
		// For PDF, we'll save as JSON and note that PDF requires additional implementation
		reportData, err = s.reportExporter.exportJSON(compactCfg)
		filename = fmt.Sprintf("%s_%s_report.json", job.Owner+"-"+job.Repo, time.Now().Format("20060102"))
		log.Println("PDF export not fully implemented, exporting as JSON instead")
	default:
		reportData, err = s.reportExporter.exportJSON(compactCfg)
		filename = fmt.Sprintf("%s_%s_report.json", job.Owner+"-"+job.Repo, time.Now().Format("20060102"))
	}

	if err != nil {
		return fmt.Errorf("failed to export report: %w", err)
	}

	// Save to destination
	if err := s.saveReport(job.Destination, filename, reportData); err != nil {
		return fmt.Errorf("failed to save report: %w", err)
	}

	// Update last run time
	job.LastRun = time.Now()
	job.NextRun = s.calculateNextRunTime(job.GetCronExpression())
	s.settings.UpdateScheduledJob(job)

	// Telemetry: record latency
	if s.metrics != nil {
		s.metrics.RecordLatency(time.Since(startTime))
	}

	log.Printf("Job %s completed successfully", job.ID)
	// Increment successful job metric
	if s.metrics != nil {
		s.metrics.IncSuccess()
	}
	return nil
}

// exportJSON exports the report as JSON
func (s *ReportExporter) exportJSON(cfg output.CompactConfig) ([]byte, error) {
	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(cfg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// exportMarkdown exports the report as Markdown
func (s *ReportExporter) exportMarkdown(cfg output.CompactConfig, repoInfo *github.Repo) ([]byte, error) {
	var buf bytes.Buffer

	buf.WriteString("# Repository Analysis Report\n\n")
	buf.WriteString(fmt.Sprintf("**Repository:** %s\n\n", cfg.Repo.FullName))
	buf.WriteString(fmt.Sprintf("**Generated:** %s\n\n", time.Now().Format("2006-01-02 15:04:05")))

	buf.WriteString("## Metrics Summary\n\n")
	buf.WriteString(fmt.Sprintf("- **Health Score:** %d/100\n", cfg.HealthScore))
	buf.WriteString(fmt.Sprintf("- **Bus Factor:** %d\n", cfg.BusFactor))
	buf.WriteString(fmt.Sprintf("- **Bus Risk:** %s\n", cfg.BusRisk))
	buf.WriteString(fmt.Sprintf("- **Maturity Score:** %d\n", cfg.MaturityScore))
	buf.WriteString(fmt.Sprintf("- **Maturity Level:** %s\n", cfg.MaturityLevel))
	buf.WriteString(fmt.Sprintf("- **Commits (1 year):** %d\n", cfg.CommitsLastYear))
	buf.WriteString(fmt.Sprintf("- **Contributors:** %d\n\n", cfg.Contributors))

	buf.WriteString("## Repository Info\n\n")
	buf.WriteString(fmt.Sprintf("- **Stars:** %d\n", cfg.Repo.Stars))
	buf.WriteString(fmt.Sprintf("- **Forks:** %d\n", cfg.Repo.Forks))
	buf.WriteString(fmt.Sprintf("- **Open Issues:** %d\n", cfg.Repo.OpenIssues))
	if cfg.Repo.Language != "" {
		buf.WriteString(fmt.Sprintf("- **Primary Language:** %s\n", cfg.Repo.Language))
	}

	if cfg.Repo.Description != "" {
		buf.WriteString(fmt.Sprintf("\n**Description:** %s\n", cfg.Repo.Description))
	}

	return buf.Bytes(), nil
}

// saveReport saves the report to the specified destination
func (s *Scheduler) saveReport(dest config.OutputDestination, filename string, data []byte) error {
	if !dest.Enabled {
		return fmt.Errorf("destination is not enabled")
	}

	switch dest.Type {
	case "local":
		return s.saveToLocalPath(dest.LocalPath, filename, data)
	case "webhook":
		return s.sendToWebhook(dest.WebhookURL, filename, data)
	default:
		// Default to local if not specified
		return s.saveToLocalPath(dest.LocalPath, filename, data)
	}
}

// saveToLocalPath saves the report to a local directory
func (s *Scheduler) saveToLocalPath(localPath, filename string, data []byte) error {
	if localPath == "" {
		// Use default export directory
		localPath = s.settings.ExportDirectory
	}

	// Ensure directory exists
	if err := os.MkdirAll(localPath, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	filePath := filepath.Join(localPath, filename)
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	log.Printf("Report saved to: %s", filePath)
	return nil
}

// sendToWebhook sends the report to a webhook URL
func (s *Scheduler) sendToWebhook(webhookURL, filename string, data []byte) error {
	if webhookURL == "" {
		return fmt.Errorf("webhook URL is not configured")
	}

	payload := map[string]interface{}{
		"filename":  filename,
		"content":   string(data),
		"timestamp": time.Now().Format(time.RFC3339),
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		return fmt.Errorf("webhook returned status: %d", resp.StatusCode)
	}

	log.Printf("Report sent to webhook: %s", webhookURL)
	return nil
}

// AddJob adds a new scheduled job
func (s *Scheduler) AddJob(job config.ScheduledJob) error {
	if err := s.settings.AddScheduledJob(job); err != nil {
		return fmt.Errorf("failed to add job: %w", err)
	}

	if job.Enabled {
		return s.scheduleJob(job)
	}

	return nil
}

// RemoveJob removes a scheduled job
func (s *Scheduler) RemoveJob(jobID string) error {
	// Remove from cron if scheduled
	if entryID, ok := s.jobEntries[jobID]; ok {
		s.cron.Remove(entryID)
		delete(s.jobEntries, jobID)
	}

	return s.settings.RemoveScheduledJob(jobID)
}

// ListJobs returns all scheduled jobs
func (s *Scheduler) ListJobs() []config.ScheduledJob {
	return s.settings.GetScheduledJobs()
}

// GetJob returns a specific job by ID
func (s *Scheduler) GetJob(jobID string) *config.ScheduledJob {
	return s.settings.GetScheduledJobByID(jobID)
}

// EnableJob enables or disables a job
func (s *Scheduler) EnableJob(jobID string, enabled bool) error {
	job := s.settings.GetScheduledJobByID(jobID)
	if job == nil {
		return fmt.Errorf("job not found: %s", jobID)
	}

	job.Enabled = enabled

	if enabled {
		if err := s.scheduleJob(*job); err != nil {
			return fmt.Errorf("failed to schedule job: %w", err)
		}
	} else {
		if entryID, ok := s.jobEntries[jobID]; ok {
			s.cron.Remove(entryID)
			delete(s.jobEntries, jobID)
		}
	}

	return s.settings.EnableScheduledJob(jobID, enabled)
}

// calculateNextRunTime calculates the next run time based on cron expression
func (s *Scheduler) calculateNextRunTime(cronExpr string) time.Time {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	sched, err := parser.Parse(cronExpr)
	if err != nil {
		return time.Now().Add(24 * time.Hour) // fallback
	}
	return sched.Next(time.Now())
}

// runQuotaRecovery handles adaptive quota recovery logic.
func (s *Scheduler) runQuotaRecovery() {
	for {
		select {
		case <-s.quotaTicker.C:
			// Adaptive refill scaling based on current load
			load := float64(atomic.LoadInt32(&s.activeWorkers)) / float64(s.maxWorkers)
			// Refill rate decreases as load approaches capacity
			refill := int64(math.Max(1, float64(s.maxQuota)*(1.0-load)))
			// Add tokens up to maxQuota
			newTokens := atomic.AddInt64(&s.quotaTokens, refill)
			if newTokens > s.maxQuota {
				atomic.StoreInt64(&s.quotaTokens, s.maxQuota)
			}
			if s.metrics != nil {
				s.metrics.IncQuotaRecovery()
			}
		case <-s.stopChan:
			return
		}
	}
}

// ValidateCronExpression validates a cron expression
func ValidateCronExpression(expr string) error {
	parser := cron.NewParser(cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow | cron.Descriptor)
	_, err := parser.Parse(expr)
	return err
}

// isInCooldown checks if the scheduler is currently in a cooldown period.
func (s *Scheduler) isInCooldown() bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.cooldownActive && time.Now().Before(s.cooldownUntil) {
		return true
	}
	// Reset cooldown if expired
	if s.cooldownActive && time.Now().After(s.cooldownUntil) {
		s.cooldownActive = false
		s.failureStreak = 0
	}
	return false
}

// recordFailure increments failure streak and triggers cooldown if threshold exceeded.
func (s *Scheduler) recordFailure() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.failureStreak++
	if s.failureStreak >= 3 { // threshold
		if s.metrics != nil {
			s.metrics.IncCooldownTrigger()
		}
		s.cooldownActive = true
		s.cooldownUntil = time.Now().Add(10 * time.Second) // cooldown duration
	}
}

// resetFailureStreak resets failure counters after a successful job.
func (s *Scheduler) resetFailureStreak() {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.failureStreak = 0
	s.cooldownActive = false
}

// executeJobWithRetry wraps executeJob with bounded retries and backoff.
func (s *Scheduler) executeJobWithRetry(job config.ScheduledJob) error {
	const maxRetries = 5
	baseBackoff := 200 * time.Millisecond
	maxBackoff := 5 * time.Second

	var attempt int
	backoff := baseBackoff

	for {
		err := s.executeJob(job)
		if err == nil {
			// success
			s.resetFailureStreak()
			return nil
		}

		// Check for timeout
		if errors.Is(err, context.DeadlineExceeded) {
			if s.metrics != nil {
				s.metrics.IncTimeoutCount()
			}
			s.recordFailure()
			return err
		}

		attempt++
		if attempt > maxRetries {
			s.recordFailure()
			return err
		}

		// Record retry metric
		if s.metrics != nil {
			s.metrics.IncRetryCount()
		}
		// Backoff with jitter
		jitter := time.Duration(rand.Int63n(int64(backoff)))
		time.Sleep(backoff + jitter)
		backoff *= 2
		if backoff > maxBackoff {
			backoff = maxBackoff
		}
	}
}

// FormatScheduleInterval returns available schedule intervals
func FormatScheduleInterval() string {
	var intervals []string
	intervals = append(intervals, string(config.ScheduleDaily))
	intervals = append(intervals, string(config.ScheduleWeekly))
	intervals = append(intervals, string(config.ScheduleMonthly))
	intervals = append(intervals, string(config.ScheduleCustom))
	return strings.Join(intervals, ", ")
}

// GetCronExpressionForInterval returns the cron expression for a given interval
func GetCronExpressionForInterval(interval config.ScheduleInterval) string {
	return interval.CronExpression()
}
