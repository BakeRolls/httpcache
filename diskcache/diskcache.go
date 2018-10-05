package diskcache

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/peterbourgon/diskv"
)

type DiskCache struct {
	diskv *diskv.Diskv
}

func New(path string) *DiskCache {
	return &DiskCache{
		diskv.New(diskv.Options{
			BasePath:     path,
			CacheSizeMax: 100 * 1024 * 1024,
		}),
	}
}

func (c *DiskCache) Get(key string) ([]byte, error) {
	return c.diskv.Read(hash(key))
}

func (c *DiskCache) Set(key string, dump []byte) error {
	return c.diskv.Write(hash(key), dump)
}

func hash(key string) string {
	sum := sha256.Sum256([]byte(key))
	return hex.EncodeToString(sum[:])
}
