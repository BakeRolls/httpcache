package memcache

import (
	"errors"
	"sync"
)

type MemCache struct {
	mu    sync.RWMutex
	dumps map[string][]byte
}

func New() *MemCache {
	return &MemCache{dumps: map[string][]byte{}}
}

func (c *MemCache) Get(key string) ([]byte, error) {
	c.mu.RLock()
	dump, ok := c.dumps[key]
	if !ok {
		return nil, errors.New("cache not found")
	}
	c.mu.RUnlock()
	return dump, nil
}

func (c *MemCache) Set(key string, dump []byte) error {
	c.mu.Lock()
	c.dumps[key] = dump
	c.mu.Unlock()
	return nil
}
