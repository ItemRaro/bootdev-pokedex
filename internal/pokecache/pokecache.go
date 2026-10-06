// Package pokecache provides a thread-safe in-memory cache with time-based expiration.
package pokecache

import (
	"sync"
	"time"
)

// cacheEntry holds a cached value and when it was stored.
type cacheEntry struct {
	val       []byte
	createdAt time.Time
}

// Cache stores raw byte values by key and evicts stale entries in the background.
type Cache struct {
	cacheEntry map[string]cacheEntry
	mu         sync.RWMutex
}

// NewCache creates a cache whose entries expire after interval.
func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		cacheEntry: make(map[string]cacheEntry),
	}

	// Evict expired entries in the background.
	go cache.reapLoop(interval)

	return &cache
}

// Add stores val under key, replacing any existing entry.
func (c *Cache) Add(key string, val []byte) {
	entry := cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheEntry[key] = entry
}

// Get returns the value for key and whether it was found.
func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.cacheEntry[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

// reapLoop removes entries older than interval on every tick.
func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		// Hold the lock only while sweeping, not between ticks.
		c.mu.Lock()
		for key, entry := range c.cacheEntry {
			expired := entry.createdAt.Add(interval)
			if time.Now().After(expired) {
				delete(c.cacheEntry, key)
			}
		}
		c.mu.Unlock()
	}
}
