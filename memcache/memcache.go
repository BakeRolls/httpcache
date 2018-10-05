package memcache

import (
	"errors"
	"sync"
)

// MemCache stores responses in memory.
type MemCache struct {
	mu    sync.RWMutex
	dumps map[string][]byte
}

// New returns a new MemCache.
func New() *MemCache {
	return &MemCache{dumps: map[string][]byte{}}
}

// Get a cached response dump.
func (c *MemCache) Get(key string) ([]byte, error) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	dump, ok := c.dumps[key]
	if !ok {
		return nil, errors.New("cache not found")
	}
	return dump, nil
}

// Set a response dump.
func (c *MemCache) Set(key string, dump []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.dumps[key] = dump
	return nil
}
