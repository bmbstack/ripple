package cache

import (
	"sync"
	"time"
)

// MemoryItem represents a memory cache item.
type MemoryItem struct {
	val        string
	created    int64
	expiration time.Duration
}

// MemoryCache represents a memory cache adapter implementation.
type MemoryCache struct {
	lock  sync.RWMutex
	items map[string]*MemoryItem
}

// NewMemoryCache creates and returns a new memory cache.
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{items: make(map[string]*MemoryItem)}
}

// put value into cache with key forever save
func (c *MemoryCache) Forever(key, val string) error {
	c.lock.Lock()
	defer c.lock.Unlock()

	c.items[key] = &MemoryItem{
		val:        val,
		created:    time.Now().Unix(),
		expiration: 0,
	}
	return nil
}

func (c *MemoryCache) Connect(opt Options) error {
	return nil
}

// Set puts value into cache with key and expire time.
func (c *MemoryCache) Set(key, val string, expiration time.Duration) {
}

// Get gets cached value by given key.
func (c *MemoryCache) Get(key string) string {
	return ""
}

// Delete deletes cached value by given key.
func (c *MemoryCache) Delete(key string) {
}

// Delete deletes cached value by given prefix.
func (c *MemoryCache) DeleteByPrefix(prefix string) {
}

// Incr increases cached int-type value by given key as a counter.
func (c *MemoryCache) Incr(key string) int64 {
	return 0
}

// Decr decreases cached int-type value by given key as a counter.
func (c *MemoryCache) Decr(key string) int64 {
	return 0
}

// IsExist returns true if cached value exists.
func (c *MemoryCache) IsExist(key string) bool {
	return false
}

// update expire time
func (c *MemoryCache) Touch(key string) {
}

// Flush deletes all cached data.
func (c *MemoryCache) Flush() {
}
