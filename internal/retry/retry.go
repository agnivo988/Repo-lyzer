// Package retry provides a context-aware exponential back-off with full jitter.
//
// Design goals:
//   - Prevent synchronized retry storms when many goroutines fail at the same time.
//   - Keep retry growth bounded so the scheduler never stalls indefinitely.
//   - Respect context cancellation so the caller can abort cleanly.
package retry

import (
	"context"
	"math/rand"
	"time"
)

const (
	// DefaultMaxAttempts is the number of attempts before giving up (1 initial + 4 retries).
	DefaultMaxAttempts = 5
	// DefaultBaseDelay is the starting backoff window before jitter is applied.
	DefaultBaseDelay = 200 * time.Millisecond
	// DefaultMaxDelay is the ceiling for a single sleep interval.
	DefaultMaxDelay = 10 * time.Second
)

// Config holds tunable parameters for the retry loop.
type Config struct {
	// MaxAttempts is the total number of attempts (including the first call).
	MaxAttempts int
	// BaseDelay is the minimum backoff window.
	BaseDelay time.Duration
	// MaxDelay caps the computed backoff so it cannot grow unbounded.
	MaxDelay time.Duration
}

// DefaultConfig returns a production-safe retry configuration.
func DefaultConfig() Config {
	return Config{
		MaxAttempts: DefaultMaxAttempts,
		BaseDelay:   DefaultBaseDelay,
		MaxDelay:    DefaultMaxDelay,
	}
}

// Do executes fn up to cfg.MaxAttempts times.
// Between attempts it sleeps for a random duration in [0, min(BaseDelay*2^n, MaxDelay)).
// This "full jitter" strategy spreads retries uniformly across the backoff window,
// preventing thundering-herd / retry-storm behaviour.
//
// Do returns immediately if:
//   - fn returns nil (success).
//   - the context is cancelled or its deadline is exceeded.
//   - all attempts are exhausted (returns the last error).
func Do(ctx context.Context, cfg Config, fn func(ctx context.Context) error) error {
	var lastErr error
	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if err := ctx.Err(); err != nil {
			return err
		}

		lastErr = fn(ctx)
		if lastErr == nil {
			return nil
		}

		// Last attempt – no sleep needed.
		if attempt == cfg.MaxAttempts-1 {
			break
		}

		sleep := jitteredDelay(cfg.BaseDelay, cfg.MaxDelay, attempt)
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(sleep):
		}
	}
	return lastErr
}

// jitteredDelay computes a random duration in [0, cap) using full-jitter exponential backoff.
// The exponential window doubles each attempt up to maxDelay.
func jitteredDelay(base, maxDelay time.Duration, attempt int) time.Duration {
	// Shift capped at 62 to avoid int64 overflow.
	shift := attempt
	if shift > 62 {
		shift = 62
	}
	window := base * (1 << uint(shift))
	if window > maxDelay || window <= 0 {
		window = maxDelay
	}
	// Full jitter: uniform random in [0, window).
	return time.Duration(rand.Int63n(int64(window)))
}
