package cache

import (
	"encoding/json"
	"time"
)

type DummyCache struct {
}

func (c *DummyCache) New(cacheConfig json.RawMessage) (Driver, error) {
	return &DummyCache{}, nil
}
func (c *DummyCache) SearchByPrefix(prefix string) ([]string, error) {
	return nil, ErrSearchKeysNotSupported
}
func (c *DummyCache) Set(key string, v interface{}, ttl time.Duration) error {
	return nil
}
func (c *DummyCache) Get(key string, v interface{}) error {
	return ErrNotFound
}
func (c *DummyCache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	return nil
}
func (c *DummyCache) GetBytesValue(key string) ([]byte, error) {
	return nil, ErrNotFound
}

func (c *DummyCache) Del(key string) error {
	return nil
}
func (c *DummyCache) SetGCErrHandler(f func(err error)) {
	return
}
func (c *DummyCache) Close() error {
	return nil
}
func (c *DummyCache) Flush() error {
	return nil
}
func (c *DummyCache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	return 0, ErrNotFound
}
func (c *DummyCache) SetCounter(key string, v int64, ttl time.Duration) error {
	return nil
}
func (c *DummyCache) GetCounter(key string) (int64, error) {
	return 0, ErrNotFound
}
func (c *DummyCache) DelCounter(key string) error {
	return nil
}
func init() {
	Register("dummycache", &DummyCache{})
}
