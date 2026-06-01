package queue

import (
	"errors"
	"sync"
	"testing"
)

func TestEnqueue_AcceptsUpToCapacity(t *testing.T) {
	q := New[int](3)
	for i := 0; i < 3; i++ {
		if err := q.Enqueue(i); err != nil {
			t.Fatalf("expected nil at item %d, got %v", i, err)
		}
	}
	if q.Len() != 3 {
		t.Fatalf("expected len 3, got %d", q.Len())
	}
}

func TestEnqueue_RejectsWhenFull(t *testing.T) {
	q := New[int](2)
	_ = q.Enqueue(1)
	_ = q.Enqueue(2)
	err := q.Enqueue(3)
	if !errors.Is(err, ErrQueueFull) {
		t.Fatalf("expected ErrQueueFull, got %v", err)
	}
}

func TestDequeue_EmptyReturnsZeroFalse(t *testing.T) {
	q := New[string](5)
	v, ok := q.Dequeue()
	if ok || v != "" {
		t.Fatalf("expected ('', false) from empty queue, got (%q, %v)", v, ok)
	}
}

func TestDequeue_FIFOOrdering(t *testing.T) {
	q := New[int](5)
	for i := 1; i <= 4; i++ {
		_ = q.Enqueue(i)
	}
	for want := 1; want <= 4; want++ {
		got, ok := q.Dequeue()
		if !ok || got != want {
			t.Fatalf("expected %d, got %d (ok=%v)", want, got, ok)
		}
	}
}

func TestIsSaturated_BelowHighWater(t *testing.T) {
	q := New[int](10) // high water = 8
	for i := 0; i < 7; i++ {
		_ = q.Enqueue(i)
	}
	if q.IsSaturated() {
		t.Fatal("expected not saturated at 7/10 items")
	}
}

func TestIsSaturated_AtHighWater(t *testing.T) {
	q := New[int](10) // high water = 8
	for i := 0; i < 8; i++ {
		_ = q.Enqueue(i)
	}
	if !q.IsSaturated() {
		t.Fatal("expected saturated at high-water mark (8/10)")
	}
}

func TestSnapshot_TelemetryCounters(t *testing.T) {
	q := New[int](3)
	_ = q.Enqueue(1)
	_ = q.Enqueue(2)
	_ = q.Enqueue(3)
	_ = q.Enqueue(4) // overflow – should be dropped
	m := q.Snapshot()
	if m.TotalEnqueued != 3 {
		t.Fatalf("expected 3 enqueued, got %d", m.TotalEnqueued)
	}
	if m.TotalDropped != 1 {
		t.Fatalf("expected 1 dropped, got %d", m.TotalDropped)
	}
	if m.Capacity != 3 {
		t.Fatalf("expected capacity 3, got %d", m.Capacity)
	}
}

func TestEnqueue_ConcurrentSafety(t *testing.T) {
	const cap = 50
	q := New[int](cap)
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(v int) {
			defer wg.Done()
			_ = q.Enqueue(v) // some will fail – that's expected
		}(i)
	}
	wg.Wait()
	m := q.Snapshot()
	if m.TotalEnqueued+m.TotalDropped != 100 {
		t.Fatalf("enqueued(%d) + dropped(%d) should equal 100", m.TotalEnqueued, m.TotalDropped)
	}
	if int(m.TotalEnqueued) > cap {
		t.Fatalf("queue exceeded capacity: %d > %d", m.TotalEnqueued, cap)
	}
}
