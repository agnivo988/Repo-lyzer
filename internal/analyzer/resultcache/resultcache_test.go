package resultcache

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

func TestSet_Get_HitPath(t *testing.T) {
	c := New(time.Minute)
	c.Set("repo:owner/name", "analysis-result")
	v, ok := c.Get("repo:owner/name")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if v.(string) != "analysis-result" {
		t.Fatalf("unexpected value: %v", v)
	}
}

func TestGet_MissOnUnknownKey(t *testing.T) {
	c := New(time.Minute)
	_, ok := c.Get("nonexistent")
	if ok {
		t.Fatal("expected cache miss for unknown key")
	}
}

func TestGet_MissAfterTTLExpiry(t *testing.T) {
	c := New(50 * time.Millisecond)
	c.Set("key", 42)
	time.Sleep(80 * time.Millisecond)
	_, ok := c.Get("key")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestInvalidate_RemovesEntry(t *testing.T) {
	c := New(time.Minute)
	c.Set("k", "v")
	c.Invalidate("k")
	_, ok := c.Get("k")
	if ok {
		t.Fatal("expected cache miss after explicit invalidation")
	}
}

func TestFlush_ClearsAll(t *testing.T) {
	c := New(time.Minute)
	for i := 0; i < 5; i++ {
		c.Set(fmt.Sprintf("key%d", i), i)
	}
	c.Flush()
	if c.Len() != 0 {
		t.Fatalf("expected empty cache after Flush, got len=%d", c.Len())
	}
}

func TestSnapshot_HitsMissesTracked(t *testing.T) {
	c := New(time.Minute)
	c.Set("a", 1)
	c.Get("a")        // hit
	c.Get("a")        // hit
	c.Get("missing")  // miss
	m := c.Snapshot()
	if m.Hits != 2 {
		t.Fatalf("expected 2 hits, got %d", m.Hits)
	}
	if m.Misses != 1 {
		t.Fatalf("expected 1 miss, got %d", m.Misses)
	}
}

func TestCache_ConcurrentSafety(t *testing.T) {
	c := New(time.Minute)
	var wg sync.WaitGroup
	const n = 100
	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			key := fmt.Sprintf("key%d", i%10)
			c.Set(key, i)
			c.Get(key)
		}(i)
	}
	wg.Wait()
	// No race or panic means the test passed.
}
