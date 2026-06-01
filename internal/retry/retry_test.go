package retry

import (
	"context"
	"errors"
	"testing"
	"time"
)

var errTransient = errors.New("transient error")

// TestDo_SuccessOnFirstAttempt ensures fn is called once when it succeeds immediately.
func TestDo_SuccessOnFirstAttempt(t *testing.T) {
	calls := 0
	err := Do(context.Background(), DefaultConfig(), func(ctx context.Context) error {
		calls++
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 1 {
		t.Fatalf("expected 1 call, got %d", calls)
	}
}

// TestDo_RetriesUpToMaxAttempts verifies all attempts are used before giving up.
func TestDo_RetriesUpToMaxAttempts(t *testing.T) {
	calls := 0
	cfg := Config{MaxAttempts: 3, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond}
	err := Do(context.Background(), cfg, func(ctx context.Context) error {
		calls++
		return errTransient
	})
	if !errors.Is(err, errTransient) {
		t.Fatalf("expected errTransient, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

// TestDo_SucceedsOnRetry verifies success on a later attempt stops the loop.
func TestDo_SucceedsOnRetry(t *testing.T) {
	calls := 0
	cfg := Config{MaxAttempts: 5, BaseDelay: time.Millisecond, MaxDelay: 5 * time.Millisecond}
	err := Do(context.Background(), cfg, func(ctx context.Context) error {
		calls++
		if calls < 3 {
			return errTransient
		}
		return nil
	})
	if err != nil {
		t.Fatalf("expected nil, got %v", err)
	}
	if calls != 3 {
		t.Fatalf("expected 3 calls, got %d", calls)
	}
}

// TestDo_ContextCancellation verifies the loop stops when context is cancelled.
func TestDo_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	cfg := Config{MaxAttempts: 10, BaseDelay: 50 * time.Millisecond, MaxDelay: 100 * time.Millisecond}

	// Cancel after first failure triggers the sleep; cancel during sleep.
	err := Do(ctx, cfg, func(c context.Context) error {
		calls++
		cancel() // cancel on first attempt so the sleep is interrupted
		return errTransient
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}
	if calls > 2 {
		t.Fatalf("expected at most 2 calls after cancel, got %d", calls)
	}
}

// TestDo_ContextTimeout verifies deadline exceeded stops the loop cleanly.
func TestDo_ContextTimeout(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	cfg := Config{MaxAttempts: 10, BaseDelay: 50 * time.Millisecond, MaxDelay: 100 * time.Millisecond}
	err := Do(ctx, cfg, func(c context.Context) error {
		return errTransient
	})
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context.DeadlineExceeded, got %v", err)
	}
}

// TestJitteredDelay_BoundedByMaxDelay ensures computed delays never exceed MaxDelay.
func TestJitteredDelay_BoundedByMaxDelay(t *testing.T) {
	base := 100 * time.Millisecond
	max := 2 * time.Second
	for attempt := 0; attempt < 30; attempt++ {
		d := jitteredDelay(base, max, attempt)
		if d < 0 || d >= max {
			t.Fatalf("attempt %d: delay %v out of [0, %v)", attempt, d, max)
		}
	}
}

// TestJitteredDelay_NonDeterministic checks that repeated calls produce at least some variation.
func TestJitteredDelay_NonDeterministic(t *testing.T) {
	base := 500 * time.Millisecond
	max := 5 * time.Second
	seen := make(map[time.Duration]struct{})
	for i := 0; i < 50; i++ {
		seen[jitteredDelay(base, max, 3)] = struct{}{}
	}
	if len(seen) < 2 {
		t.Fatal("expected jitter to produce varied delays, got uniform output")
	}
}
