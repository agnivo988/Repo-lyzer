package monitor

import (
	"sync/atomic"
	"time"
)

type SchedulerMetrics struct {
	schedulerQueueDepth           int64
	schedulerActiveWorkers        int64
	schedulerTimeoutCount         int64
	schedulerRetryCount           int64
	schedulerQuotaRecoveryCount   int64
	schedulerCooldownTriggerCount int64
	schedulerLatencyMsTotal       int64
	schedulerSuccessCount         int64
}

func NewSchedulerMetrics() *SchedulerMetrics {
	return &SchedulerMetrics{}
}

func (m *SchedulerMetrics) IncActiveWorkers() {
	atomic.AddInt64(&m.schedulerActiveWorkers, 1)
}

func (m *SchedulerMetrics) DecActiveWorkers() {
	atomic.AddInt64(&m.schedulerActiveWorkers, -1)
}

func (m *SchedulerMetrics) IncTimeoutCount() {
	atomic.AddInt64(&m.schedulerTimeoutCount, 1)
}

func (m *SchedulerMetrics) IncRetryCount() {
	atomic.AddInt64(&m.schedulerRetryCount, 1)
}

func (m *SchedulerMetrics) IncQuotaRecovery() {
	atomic.AddInt64(&m.schedulerQuotaRecoveryCount, 1)
}

func (m *SchedulerMetrics) IncCooldownTrigger() {
	atomic.AddInt64(&m.schedulerCooldownTriggerCount, 1)
}

func (m *SchedulerMetrics) RecordLatency(d time.Duration) {
	atomic.AddInt64(&m.schedulerLatencyMsTotal, d.Milliseconds())
}

func (m *SchedulerMetrics) IncSuccess() {
	atomic.AddInt64(&m.schedulerSuccessCount, 1)
}

func (m *SchedulerMetrics) SetQueueDepth(depth int64) {
	atomic.StoreInt64(&m.schedulerQueueDepth, depth)
}

func (m *SchedulerMetrics) Snapshot() map[string]int64 {
	return map[string]int64{
		"scheduler_queue_depth":            atomic.LoadInt64(&m.schedulerQueueDepth),
		"scheduler_active_workers":         atomic.LoadInt64(&m.schedulerActiveWorkers),
		"scheduler_timeout_count":          atomic.LoadInt64(&m.schedulerTimeoutCount),
		"scheduler_retry_count":            atomic.LoadInt64(&m.schedulerRetryCount),
		"scheduler_quota_recovery_count":   atomic.LoadInt64(&m.schedulerQuotaRecoveryCount),
		"scheduler_cooldown_trigger_count": atomic.LoadInt64(&m.schedulerCooldownTriggerCount),
		"scheduler_latency_ms_total":       atomic.LoadInt64(&m.schedulerLatencyMsTotal),
		"scheduler_success_count":          atomic.LoadInt64(&m.schedulerSuccessCount),
	}
}
