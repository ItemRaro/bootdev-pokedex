package pokecache

import (
	"sync"
	"time"
)

type cacheEntry struct {
	val       []byte
	createdAt time.Time
}

type Cache struct {
	cacheEntry map[string]cacheEntry
	mu         sync.RWMutex
}

func NewCache(interval time.Duration) *Cache {
	cache := Cache{
		cacheEntry: make(map[string]cacheEntry),
	}

	go cache.reapLoop(interval)

	return &cache
}

func (c *Cache) Add(key string, val []byte) {
	entry := cacheEntry{
		val:       val,
		createdAt: time.Now(),
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.cacheEntry[key] = entry
}

func (c *Cache) Get(key string) ([]byte, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	entry, ok := c.cacheEntry[key]
	if !ok {
		return nil, false
	}

	return entry.val, true
}

func (c *Cache) reapLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
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
