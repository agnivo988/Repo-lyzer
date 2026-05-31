package cache

import (
	"time"

	gocache "github.com/patrickmn/go-cache"
)

const (
	DefaultTTL      = 10 * time.Minute
	CleanupInterval = 15 * time.Minute
)

// New returns a new in-memory API response cache with default TTL and cleanup interval.
func New() *gocache.Cache {
	return gocache.New(DefaultTTL, CleanupInterval)
}
