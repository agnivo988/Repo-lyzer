// Package queue provides a bounded, backpressure-aware job queue for Repo-lyzer's
// scheduler. It prevents unbounded memory growth during burst scheduling or retry
// storms by enforcing a hard capacity limit and exposing pressure telemetry.
//
// Design goals:
//   - Hard capacity ceiling: Enqueue returns ErrQueueFull when the limit is reached.
//   - Saturation detection: IsSaturated reports when fill level exceeds a configurable
//     high-water mark (default 80 % of capacity).
//   - Thread-safe: all operations are protected by a mutex.
//   - Zero external dependencies.
package queue

import (
	"errors"
	"sync"
	"sync/atomic"
)

// ErrQueueFull is returned by Enqueue when the queue has reached its capacity.
var ErrQueueFull = errors.New("queue is full: backpressure limit reached")

// defaultHighWaterPct is the fraction of capacity above which the queue is
// considered saturated (80 %).
const defaultHighWaterPct = 0.80

// Queue is a bounded FIFO queue with backpressure protection.
type Queue[T any] struct {
	mu           sync.Mutex
	items        []T
	cap          int
	highWater    int
	// Telemetry counters (atomic for lock-free reads).
	totalEnqueued uint64
	totalDropped  uint64
}

// New creates a Queue with the given capacity.
// The high-water mark is set to 80 % of capacity by default.
// panics if cap <= 0.
func New[T any](cap int) *Queue[T] {
	if cap <= 0 {
		panic("queue: capacity must be greater than zero")
	}
	hw := int(float64(cap) * defaultHighWaterPct)
	if hw == 0 {
		hw = 1
	}
	return &Queue[T]{
		items:     make([]T, 0, cap),
		cap:       cap,
		highWater: hw,
	}
}

// Enqueue appends item to the queue.
// Returns ErrQueueFull without blocking if the queue is at capacity.
func (q *Queue[T]) Enqueue(item T) error {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) >= q.cap {
		atomic.AddUint64(&q.totalDropped, 1)
		return ErrQueueFull
	}
	q.items = append(q.items, item)
	atomic.AddUint64(&q.totalEnqueued, 1)
	return nil
}

// Dequeue removes and returns the head item.
// Returns (zero, false) if the queue is empty.
func (q *Queue[T]) Dequeue() (T, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		var zero T
		return zero, false
	}
	item := q.items[0]
	// Shift without retaining the full backing array indefinitely.
	q.items = q.items[1:]
	return item, true
}

// Len returns the current number of items in the queue.
func (q *Queue[T]) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}

// IsSaturated returns true when the queue length has reached or exceeded the
// high-water mark. Callers can use this to apply upstream backpressure before
// the queue becomes completely full.
func (q *Queue[T]) IsSaturated() bool {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items) >= q.highWater
}

// Snapshot returns a point-in-time view of queue telemetry.
// Reads are lock-free via atomic loads so they never block Enqueue/Dequeue.
func (q *Queue[T]) Snapshot() Metrics {
	return Metrics{
		TotalEnqueued: atomic.LoadUint64(&q.totalEnqueued),
		TotalDropped:  atomic.LoadUint64(&q.totalDropped),
		CurrentLen:    q.Len(),
		Capacity:      q.cap,
	}
}

// Metrics holds a point-in-time snapshot of queue counters.
type Metrics struct {
	TotalEnqueued uint64
	TotalDropped  uint64
	CurrentLen    int
	Capacity      int
}
