package cache

import (
	"time"
)

//DummyCache DummyCache dont store any data.
//Usually used in develop environment or testing
type DummyCache struct {
	DriverUtil
}

//SetBytesValue Set bytes data to cache by given key.
//Return any error raised.
func (c *DummyCache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	return nil
}

//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
//Return any error raised.
func (c *DummyCache) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {
	return nil
}

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (c *DummyCache) GetBytesValue(key string) ([]byte, error) {
	return nil, ErrNotFound
}

//MGetBytesValue get multiple bytes data from cache by given keys.
//Return data bytes map and any error if raised.
func (c *DummyCache) MGetBytesValue(keys ...string) (map[string][]byte, error) {
	return map[string][]byte{}, nil
}

//MSetBytesValue set multiple bytes data to cache with given key-value map.
//Return  any error if raised.
func (c *DummyCache) MSetBytesValue(data map[string][]byte, ttl time.Duration) error {
	return nil
}

//Del Delete data in cache by given key.
//Return any error raised.
func (c *DummyCache) Del(key string) error {
	return nil
}

//SetGCErrHandler Set callback to handler error raised when gc.
func (c *DummyCache) SetGCErrHandler(f func(err error)) {
	return
}

//Close Close cache.
//Return any error if raised
func (c *DummyCache) Close() error {
	return nil
}

//Flush Delete all data in cache.
//Return any error if raised
func (c *DummyCache) Flush() error {
	return nil
}

//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *DummyCache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	return 0, nil
}

//SetCounter Set int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *DummyCache) SetCounter(key string, v int64, ttl time.Duration) error {
	return nil
}

//GetCounter Get int val from cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *DummyCache) GetCounter(key string) (int64, error) {
	return 0, ErrNotFound
}

//DelCounter Delete int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *DummyCache) DelCounter(key string) error {
	return nil
}

//Expire set cache value expire duration by given key and ttl
func (c *DummyCache) Expire(key string, ttl time.Duration) error {
	return nil
}

//ExpireCounter set cache counter  expire duration by given key and ttl
func (c *DummyCache) ExpireCounter(key string, ttl time.Duration) error {
	return nil
}

func init() {
	Register("dummycache", func(conf Config, prefix string) (Driver, error) {
		return &DummyCache{}, nil
	})
}
