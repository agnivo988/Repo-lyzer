package config

import (
	"errors"
	"testing"
	"time"
)

func TestSchedulerConfig_Validate_ValidDefaults(t *testing.T) {
	cfg := DefaultSchedulerConfig()
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected valid default config, got error: %v", err)
	}
}

func TestSchedulerConfig_Validate_NegativeWorkerCount(t *testing.T) {
	cfg := DefaultSchedulerConfig()
	cfg.MaxWorkers = -1
	err := cfg.Validate()
	if err == nil {
		t.Fatal("expected error for negative MaxWorkers")
	}
	if !errors.Is(err, ErrInvalidWorkerCount) {
		t.Fatalf("expected ErrInvalidWorkerCount, got: %v", err)
	}
}

func TestSchedulerConfig_Validate_ZeroWorkerCount(t *testing.T) {
	cfg := DefaultSchedulerConfig()
	cfg.MaxWorkers = 0
	if err := cfg.Validate(); !errors.Is(err, ErrInvalidWorkerCount) {
		t.Fatalf("expected ErrInvalidWorkerCount, got: %v", err)
	}
}

func TestSchedulerConfig_Validate_NegativeQueueSize(t *testing.T) {
	cfg := DefaultSchedulerConfig()
	cfg.QueueSize = -5
	if err := cfg.Validate(); !errors.Is(err, ErrInvalidQueueSize) {
		t.Fatalf("expected ErrInvalidQueueSize, got: %v", err)
	}
}

func TestSchedulerConfig_Validate_ZeroTimeout(t *testing.T) {
	cfg := DefaultSchedulerConfig()
	cfg.WorkerTimeout = 0
	if err := cfg.Validate(); !errors.Is(err, ErrInvalidTimeout) {
		t.Fatalf("expected ErrInvalidTimeout, got: %v", err)
	}
}

func TestSchedulerConfig_Validate_NegativeRetryInterval(t *testing.T) {
	cfg := DefaultSchedulerConfig()
	cfg.RetryInterval = -1 * time.Second
	if err := cfg.Validate(); !errors.Is(err, ErrInvalidRetryInterval) {
		t.Fatalf("expected ErrInvalidRetryInterval, got: %v", err)
	}
}

func TestSchedulerConfig_Validate_NegativeCooldown(t *testing.T) {
	cfg := DefaultSchedulerConfig()
	cfg.CooldownDuration = -30 * time.Second
	if err := cfg.Validate(); !errors.Is(err, ErrInvalidCooldown) {
		t.Fatalf("expected ErrInvalidCooldown, got: %v", err)
	}
}

func TestSchedulerConfig_Validate_ZeroAnalyzerLimit(t *testing.T) {
	cfg := DefaultSchedulerConfig()
	cfg.MaxAnalyzers = 0
	if err := cfg.Validate(); !errors.Is(err, ErrInvalidAnalyzerLimit) {
		t.Fatalf("expected ErrInvalidAnalyzerLimit, got: %v", err)
	}
}

func TestSchedulerConfig_Validate_ZeroCooldownAllowed(t *testing.T) {
	// Zero cooldown is explicitly permitted (means no cooldown enforcement).
	cfg := DefaultSchedulerConfig()
	cfg.CooldownDuration = 0
	if err := cfg.Validate(); err != nil {
		t.Fatalf("expected zero cooldown to be valid, got: %v", err)
	}
}
