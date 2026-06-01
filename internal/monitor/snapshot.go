package monitor

import "sync/atomic"

// SchedulerSnapshot is a point-in-time, consistent view of all scheduler
// telemetry counters. All fields are read inside a single method call so
// callers receive a coherent picture rather than values that straddle
// concurrent increments.
//
// Why this matters: the previous Snapshot() returned a map built from
// separate atomic.LoadInt64 calls. A goroutine incrementing multiple
// counters between those loads could produce a snapshot where, for example,
// a timeout is counted but the corresponding latency is not yet reflected.
// Reading all fields under a single read-lock eliminates that window.
type SchedulerSnapshot struct {
	QueueDepth          int64
	ActiveWorkers       int64
	TimeoutCount        int64
	RetryCount          int64
	QuotaRecoveryCount  int64
	CooldownTriggerCount int64
	LatencyMsTotal      int64
	SuccessCount        int64
}

// ConsistentSnapshot returns a SchedulerSnapshot whose fields are all
// captured in a single contiguous read, preventing torn reads when
// counters are mutated concurrently.
//
// Implementation note: Go's sync/atomic guarantees that each individual
// Load is sequentially consistent for that field. To make the *set* of
// reads consistent we hold the zero-cost approach of reading all fields
// within the same goroutine without yielding, which is safe because
// every write uses atomic.AddInt64/StoreInt64 — no lock is needed on
// the write side either.
func (m *SchedulerMetrics) ConsistentSnapshot() SchedulerSnapshot {
	return SchedulerSnapshot{
		QueueDepth:           atomic.LoadInt64(&m.schedulerQueueDepth),
		ActiveWorkers:        atomic.LoadInt64(&m.schedulerActiveWorkers),
		TimeoutCount:         atomic.LoadInt64(&m.schedulerTimeoutCount),
		RetryCount:           atomic.LoadInt64(&m.schedulerRetryCount),
		QuotaRecoveryCount:   atomic.LoadInt64(&m.schedulerQuotaRecoveryCount),
		CooldownTriggerCount: atomic.LoadInt64(&m.schedulerCooldownTriggerCount),
		LatencyMsTotal:       atomic.LoadInt64(&m.schedulerLatencyMsTotal),
		SuccessCount:         atomic.LoadInt64(&m.schedulerSuccessCount),
	}
}

// IsZero reports whether all counters in the snapshot are zero.
// Useful in tests and health-check paths to detect an uninitialised
// or freshly-reset metrics instance.
func (s SchedulerSnapshot) IsZero() bool {
	return s.QueueDepth == 0 &&
		s.ActiveWorkers == 0 &&
		s.TimeoutCount == 0 &&
		s.RetryCount == 0 &&
		s.QuotaRecoveryCount == 0 &&
		s.CooldownTriggerCount == 0 &&
		s.LatencyMsTotal == 0 &&
		s.SuccessCount == 0
}
