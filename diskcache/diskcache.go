package diskcache

import (
	"crypto/sha256"
	"encoding/hex"
	"time"

	"github.com/peterbourgon/diskv"
)

const (
	// NoExpiration stops the diskcache from removing items after a given time.
	NoExpiration time.Duration = 0
)

// DiskCache stores responses in files on disk.
type DiskCache struct {
	age   time.Duration
	diskv *diskv.Diskv
}

// New returns a new DiskCache. If age is 0, the cache won't get deleted.
func New(path string, age time.Duration) *DiskCache {
	dv := diskv.New(diskv.Options{
		BasePath:     path,
		CacheSizeMax: 100 * 1024 * 1024,
	})
	return &DiskCache{age, dv}
}

// Get a cached response dump.
func (c *DiskCache) Get(key string) ([]byte, error) {
	return c.diskv.Read(hash(key))
}

// Set a response dump.
func (c *DiskCache) Set(key string, dump []byte) error {
	key = hash(key)
	if c.age != 0 {
		time.AfterFunc(c.age, func() { c.diskv.Erase(key) })
	}
	return c.diskv.Write(key, dump)
}

func hash(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}
