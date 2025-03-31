package cache

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type CacheItem struct {
	value     string
	timestamp time.Time
}

type Cache struct {
	mu            sync.RWMutex
	items         map[string]CacheItem
	maxMemoryMB   int64
	maxItemAge    time.Duration
	cleanupTicker *time.Ticker
}

type CacheConfig struct {
	MaxMemoryMB int64         // Maximum memory in MB
	MaxItemAge  time.Duration // Maximum age of items (0 for no TTL)
	CleanupInterval time.Duration // Interval for cleanup checks
}

// NewCache creates a new Cache instance
func NewCache(config CacheConfig) *Cache {
	c := &Cache{
		items:       make(map[string]CacheItem),
		maxMemoryMB: config.MaxMemoryMB,
		maxItemAge:  config.MaxItemAge,
	}

	// Start cleanup goroutine if either eviction strategy is enabled
	if config.MaxMemoryMB > 0 || config.MaxItemAge > 0 {
		c.cleanupTicker = time.NewTicker(config.CleanupInterval)
		go c.cleanupLoop()
	}

	return c
}

func (c *Cache) cleanupLoop() {
	for range c.cleanupTicker.C {
		c.cleanup()
	}
}

func (c *Cache) cleanup() {
	c.mu.Lock()
	defer c.mu.Unlock()

	now := time.Now()

	// Time-based eviction
	if c.maxItemAge > 0 {
		for key, item := range c.items {
			if now.Sub(item.timestamp) > c.maxItemAge {
				delete(c.items, key)
			}
		}
	}

	// Memory-based eviction
	if c.maxMemoryMB > 0 {
		var m runtime.MemStats
		runtime.ReadMemStats(&m)
		
		// If memory usage exceeds threshold, remove oldest items
		if m.Alloc > uint64(c.maxMemoryMB)*1024*1024 {
			// Convert map to slice for sorting
			type keyAge struct {
				key       string
				timestamp time.Time
			}
			items := make([]keyAge, 0, len(c.items))
			for k, v := range c.items {
				items = append(items, keyAge{k, v.timestamp})
			}

			// Sort by timestamp (oldest first)
			for i := 0; i < len(items)/2; i++ {
				delete(c.items, items[i].key)
			}
		}
	}
}

// Put adds or updates a key-value pair
// Returns error and a boolean indicating if the key existed
func (c *Cache) Put(key, value string) (bool, error) {
	if len(key) > 256 || len(value) > 256 {
		return false, fmt.Errorf("key or value exceeds 256 characters")
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	
	_, exists := c.items[key]
	c.items[key] = CacheItem{
		value:     value,
		timestamp: time.Now(),
	}
	return exists, nil
}

// Get retrieves a value by key
func (c *Cache) Get(key string) (string, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	
	item, exists := c.items[key]
	if !exists {
		return "", false
	}

	// Check TTL if enabled
	if c.maxItemAge > 0 && time.Since(item.timestamp) > c.maxItemAge {
		go func() {
			c.mu.Lock()
			delete(c.items, key)
			c.mu.Unlock()
		}()
		return "", false
	}

	return item.value, true
}

func (c *Cache) Stop() {
	if c.cleanupTicker != nil {
		c.cleanupTicker.Stop()
	}
} 