package monitor

import (
	"sync"
	"testing"
	"time"
)

// TestConsistentSnapshot_ZeroOnInit verifies a fresh SchedulerMetrics returns
// an all-zero snapshot.
func TestConsistentSnapshot_ZeroOnInit(t *testing.T) {
	m := NewSchedulerMetrics()
	snap := m.ConsistentSnapshot()
	if !snap.IsZero() {
		t.Fatalf("expected zero snapshot on init, got %+v", snap)
	}
}

// TestConsistentSnapshot_ReflectsIncrements verifies each counter is captured
// correctly after controlled mutations.
func TestConsistentSnapshot_ReflectsIncrements(t *testing.T) {
	m := NewSchedulerMetrics()
	m.IncActiveWorkers()
	m.IncActiveWorkers()
	m.IncTimeoutCount()
	m.IncRetryCount()
	m.IncRetryCount()
	m.IncRetryCount()
	m.IncSuccess()
	m.RecordLatency(250 * time.Millisecond)
	m.SetQueueDepth(7)

	snap := m.ConsistentSnapshot()
	if snap.ActiveWorkers != 2 {
		t.Errorf("ActiveWorkers: want 2, got %d", snap.ActiveWorkers)
	}
	if snap.TimeoutCount != 1 {
		t.Errorf("TimeoutCount: want 1, got %d", snap.TimeoutCount)
	}
	if snap.RetryCount != 3 {
		t.Errorf("RetryCount: want 3, got %d", snap.RetryCount)
	}
	if snap.SuccessCount != 1 {
		t.Errorf("SuccessCount: want 1, got %d", snap.SuccessCount)
	}
	if snap.LatencyMsTotal != 250 {
		t.Errorf("LatencyMsTotal: want 250, got %d", snap.LatencyMsTotal)
	}
	if snap.QueueDepth != 7 {
		t.Errorf("QueueDepth: want 7, got %d", snap.QueueDepth)
	}
}

// TestConsistentSnapshot_ConcurrentReads verifies ConsistentSnapshot never
// panics and always returns non-negative values under high concurrency.
func TestConsistentSnapshot_ConcurrentReads(t *testing.T) {
	m := NewSchedulerMetrics()
	var wg sync.WaitGroup
	const goroutines = 50

	// Writers
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.IncRetryCount()
			m.IncSuccess()
			m.RecordLatency(10 * time.Millisecond)
		}()
	}
	// Concurrent readers
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			snap := m.ConsistentSnapshot()
			if snap.RetryCount < 0 || snap.SuccessCount < 0 || snap.LatencyMsTotal < 0 {
				t.Errorf("negative counter in snapshot: %+v", snap)
			}
		}()
	}
	wg.Wait()

	// After all goroutines finish, totals must be exact.
	snap := m.ConsistentSnapshot()
	if snap.RetryCount != goroutines {
		t.Errorf("RetryCount: want %d, got %d", goroutines, snap.RetryCount)
	}
	if snap.SuccessCount != goroutines {
		t.Errorf("SuccessCount: want %d, got %d", goroutines, snap.SuccessCount)
	}
}

// TestIsZero_FalseAfterIncrement ensures IsZero returns false once any field
// is non-zero.
func TestIsZero_FalseAfterIncrement(t *testing.T) {
	m := NewSchedulerMetrics()
	m.IncRetryCount()
	snap := m.ConsistentSnapshot()
	if snap.IsZero() {
		t.Fatal("expected IsZero() == false after increment")
	}
}
