package diff

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sync"
)

// CacheEntry holds a cached Report keyed by a hash of the two state inputs.
type CacheEntry struct {
	Key    string  `json:"key"`
	Report *Report `json:"report"`
}

// DiffCache is a thread-safe in-memory cache for diff reports.
type DiffCache struct {
	mu      sync.RWMutex
	entries map[string]*CacheEntry
}

// NewDiffCache creates an empty DiffCache.
func NewDiffCache() *DiffCache {
	return &DiffCache{
		entries: make(map[string]*CacheEntry),
	}
}

// CacheKey derives a deterministic cache key from two serialisable state
// values. It returns an error if either value cannot be marshalled.
func CacheKey(a, b interface{}) (string, error) {
	ab, err := json.Marshal(a)
	if err != nil {
		return "", fmt.Errorf("cache key marshal a: %w", err)
	}
	bb, err := json.Marshal(b)
	if err != nil {
		return "", fmt.Errorf("cache key marshal b: %w", err)
	}
	h := sha256.New()
	h.Write(ab)
	h.Write([]byte("|"))
	h.Write(bb)
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Get returns the cached Report for key, or nil if not present.
func (c *DiffCache) Get(key string) *Report {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if e, ok := c.entries[key]; ok {
		return e.Report
	}
	return nil
}

// Set stores a Report under the given key.
func (c *DiffCache) Set(key string, r *Report) {
	if r == nil {
		return
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries[key] = &CacheEntry{Key: key, Report: r}
}

// Len returns the number of cached entries.
func (c *DiffCache) Len() int {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.entries)
}

// Flush removes all cached entries.
func (c *DiffCache) Flush() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.entries = make(map[string]*CacheEntry)
}
