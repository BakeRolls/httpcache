package memcache

import (
	"errors"
	"fmt"
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
	defer c.mu.RUnlock()
	dump, ok := c.dumps[key]
	if !ok {
		return nil, errors.New("cache not found")
	}
	return dump, nil
}

func (c *MemCache) Set(key string, dump []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	fmt.Println("set", key, len(dump))
	c.dumps[key] = dump
	return nil
}
