// Package resultcache provides a lightweight, TTL-based, thread-safe cache for
// analyzer results. It avoids re-running expensive analysis when the same
// repository is queried multiple times within the TTL window.
//
// Design goals:
//   - Zero external dependencies (uses only stdlib sync and time).
//   - Bounded lifecycle: entries expire automatically; no goroutine leak.
//   - Deterministic invalidation: callers can force-evict a key at any time.
//   - Generic: works with any result type via interface{} values.
package resultcache

import (
	"sync"
	"time"
)

// DefaultTTL is the duration after which a cached result is considered stale.
const DefaultTTL = 5 * time.Minute

// entry holds a cached value together with its expiry timestamp.
type entry struct {
	value     interface{}
	expiresAt time.Time
}

// Cache is a TTL-based in-memory store for analyzer results.
// The zero value is not usable; use New to construct one.
type Cache struct {
	mu      sync.RWMutex
	items   map[string]entry
	ttl     time.Duration
	// Telemetry counters (read via Metrics()).
	hits   uint64
	misses uint64
}

// New returns a Cache with the given TTL.
// If ttl <= 0 it defaults to DefaultTTL.
func New(ttl time.Duration) *Cache {
	if ttl <= 0 {
		ttl = DefaultTTL
	}
	return &Cache{
		items: make(map[string]entry),
		ttl:   ttl,
	}
}

// Set stores value under key, overwriting any existing entry.
func (c *Cache) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items[key] = entry{
		value:     value,
		expiresAt: time.Now().Add(c.ttl),
	}
}

// Get returns the cached value for key and whether it was found and unexpired.
// Expired entries are treated as cache misses and pruned on access.
func (c *Cache) Get(key string) (interface{}, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	e, ok := c.items[key]
	if !ok {
		c.misses++
		return nil, false
	}
	if time.Now().After(e.expiresAt) {
		delete(c.items, key) // lazy eviction
		c.misses++
		return nil, false
	}
	c.hits++
	return e.value, true
}

// Invalidate removes the entry for key if it exists.
func (c *Cache) Invalidate(key string) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.items, key)
}

// Flush removes all entries from the cache.
func (c *Cache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.items = make(map[string]entry)
}

// Len returns the number of entries currently stored (including expired ones
// not yet lazily evicted).
func (c *Cache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items)
}

// Metrics holds a point-in-time snapshot of cache telemetry.
type Metrics struct {
	Hits   uint64
	Misses uint64
	Len    int
}

// Snapshot returns a consistent view of cache telemetry.
func (c *Cache) Snapshot() Metrics {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return Metrics{
		Hits:   c.hits,
		Misses: c.misses,
		Len:    len(c.items),
	}
}
