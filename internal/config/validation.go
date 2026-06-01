// Package config provides application settings and configuration management.
package config

import (
	"errors"
	"fmt"
	"time"
)

// Validation sentinel errors for structured error handling.
var (
	ErrInvalidWorkerCount   = errors.New("worker count must be greater than zero")
	ErrInvalidQueueSize     = errors.New("queue size must be greater than zero")
	ErrInvalidRetryInterval = errors.New("retry interval must be a positive duration")
	ErrInvalidCooldown      = errors.New("cooldown duration must be non-negative")
	ErrInvalidTimeout       = errors.New("timeout duration must be greater than zero")
	ErrInvalidAnalyzerLimit = errors.New("analyzer limit must be greater than zero")
)

// SchedulerConfig holds runtime configuration for the scheduler.
// It is validated before the scheduler starts to prevent invalid execution states.
type SchedulerConfig struct {
	// MaxWorkers limits the number of concurrent job executions.
	MaxWorkers int
	// QueueSize is the maximum number of jobs that can be buffered.
	QueueSize int
	// WorkerTimeout is the per-job execution timeout.
	WorkerTimeout time.Duration
	// RetryInterval is the minimum wait between job retries.
	RetryInterval time.Duration
	// CooldownDuration is the minimum wait between consecutive runs of the same job.
	CooldownDuration time.Duration
	// MaxAnalyzers limits the number of concurrent analyzer goroutines.
	MaxAnalyzers int
}

// DefaultSchedulerConfig returns a safe SchedulerConfig with production-ready defaults.
func DefaultSchedulerConfig() SchedulerConfig {
	return SchedulerConfig{
		MaxWorkers:       5,
		QueueSize:        100,
		WorkerTimeout:    2 * time.Minute,
		RetryInterval:    5 * time.Second,
		CooldownDuration: 0,
		MaxAnalyzers:     4,
	}
}

// Validate checks the SchedulerConfig for invalid values and returns a descriptive
// error on the first violation found. It is designed to be called at startup so
// that the scheduler never enters an inconsistent runtime state.
func (c SchedulerConfig) Validate() error {
	if c.MaxWorkers <= 0 {
		return fmt.Errorf("%w: got %d", ErrInvalidWorkerCount, c.MaxWorkers)
	}
	if c.QueueSize <= 0 {
		return fmt.Errorf("%w: got %d", ErrInvalidQueueSize, c.QueueSize)
	}
	if c.WorkerTimeout <= 0 {
		return fmt.Errorf("%w: got %v", ErrInvalidTimeout, c.WorkerTimeout)
	}
	if c.RetryInterval <= 0 {
		return fmt.Errorf("%w: got %v", ErrInvalidRetryInterval, c.RetryInterval)
	}
	if c.CooldownDuration < 0 {
		return fmt.Errorf("%w: got %v", ErrInvalidCooldown, c.CooldownDuration)
	}
	if c.MaxAnalyzers <= 0 {
		return fmt.Errorf("%w: got %d", ErrInvalidAnalyzerLimit, c.MaxAnalyzers)
	}
	return nil
}
