package cache

import (
	"encoding/json"
	"time"
)

//DummyCache DummyCache dont store any data.
//Usually used in develop environment or testing
type DummyCache struct {
}

//New Create new dummy cache driver.
//Return dummy cach driver and nil.
func (c *DummyCache) New(cacheConfig json.RawMessage) (Driver, error) {
	return &DummyCache{}, nil
}

//SearchByPrefix Search All key start with given prefix.
//Return All matched key and any error raised.
func (c *DummyCache) SearchByPrefix(prefix string) ([]string, error) {
	return nil, ErrSearchKeysNotSupported
}

//Set Set data model to cache by given key.
//Return any error raised.
func (c *DummyCache) Set(key string, v interface{}, ttl time.Duration) error {
	return nil
}

//Update Update data model to cache by given key only if the cache exist.
//Return any error raised.
func (c *DummyCache) Update(key string, v interface{}, ttl time.Duration) error {
	return nil
}

//Get Get data model from cache by given key.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raised.
func (c *DummyCache) Get(key string, v interface{}) error {
	return ErrNotFound
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
	return 0, ErrNotFound
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
func init() {
	Register("dummycache", &DummyCache{})
}
