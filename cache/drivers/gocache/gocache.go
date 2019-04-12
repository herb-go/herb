//Package gocache provides cache driver uses memory to store cache data.
//Using github.com/allegro/bigcache as driver.
package gocache

import (
	"time"

	"encoding/binary"

	"sync"

	"github.com/herb-go/herb/cache"
	gocache "github.com/patrickmn/go-cache"
)

const defaultExpirationInsecond = 60
const defaultCleanupIntervalInSecond = 60

//Cache The gocache cache Driver.
type Cache struct {
	cache.DriverUtil
	gocache      *gocache.Cache
	gcErrHandler func(err error)
	lock         sync.Mutex
}

//SetGCErrHandler Set callback to handler error raised when gc.
func (c *Cache) SetGCErrHandler(f func(err error)) {
	c.gcErrHandler = f
	return
}

//SetBytesValue Set bytes data to cache by given key.
//Return any error raised.
func (c *Cache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	c.gocache.Set(key, bytes, ttl)
	return nil
}

//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
//Return any error raised.
func (c *Cache) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {
	err := c.gocache.Replace(key, bytes, ttl)
	err = nil
	return err
}

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (c *Cache) GetBytesValue(key string) ([]byte, error) {
	bytes, found := c.gocache.Get(key)
	if found {
		return bytes.([]byte), nil
	}
	return nil, cache.ErrNotFound
}

//MGetBytesValue get multiple bytes data from cache by given keys.
//Return data bytes map and any error if raised.
func (c *Cache) MGetBytesValue(keys ...string) (map[string][]byte, error) {
	var result = make(map[string][]byte, len(keys))
	for k := range keys {
		b, err := c.GetBytesValue(keys[k])
		if err == cache.ErrNotFound {
		} else if err != nil {
			return result, err
		} else {
			result[keys[k]] = b
		}

	}
	return result, nil
}

//MSetBytesValue set multiple bytes data to cache with given key-value map.
//Return  any error if raised.
func (c *Cache) MSetBytesValue(data map[string][]byte, ttl time.Duration) error {
	for k := range data {
		c.gocache.Set(k, data[k], ttl)
	}
	return nil
}

//Flush Delete all data in cache.
//Return any error if raised
func (c *Cache) Flush() error {
	c.gocache.Flush()
	return nil
}

//Close Close cache.
//Return any error if raised
func (c *Cache) Close() error {
	c.gocache.Flush()
	return nil
}

//Del Delete data in cache by given key.
//Return any error raised.
func (c *Cache) Del(key string) error {
	c.gocache.Delete(key)
	return nil
}

//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	var v int64
	locker := c.Util().Locker(key)
	locker.Lock()
	defer locker.Unlock()
	data, found := c.gocache.Get(key)
	if found == false {
		v = 0
	} else {
		bytes := data.([]byte)
		v = int64(binary.BigEndian.Uint64(bytes[0:8]))
	}
	v = v + increment
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(v))
	c.gocache.Set(key, bytes, ttl)
	return v, nil
}

//SetCounter Set int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) SetCounter(key string, v int64, ttl time.Duration) error {
	locker := c.Util().Locker(key)
	locker.Lock()
	defer locker.Unlock()

	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(v))
	c.gocache.Set(key, bytes, ttl)
	return nil
}

//GetCounter Get int val from cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) GetCounter(key string) (int64, error) {
	var v int64
	var err error
	locker := c.Util().Locker(key)
	locker.Lock()
	defer locker.Unlock()

	data, found := c.gocache.Get(key)
	if found == false {
		err = cache.ErrNotFound
		return 0, err
	}
	bytes := data.([]byte)
	v = int64(binary.BigEndian.Uint64(bytes[0:8]))
	return v, nil
}

//DelCounter Delete int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raisegrd.
func (c *Cache) DelCounter(key string) error {
	locker := c.Util().Locker(key)
	locker.Lock()
	defer locker.Unlock()

	c.gocache.Delete(key)
	return nil
}

//Expire set cache value expire duration by given key and ttl
func (c *Cache) Expire(key string, ttl time.Duration) error {
	locker := c.Util().Locker(key)
	locker.Lock()
	defer locker.Unlock()
	data, found := c.gocache.Get(key)
	if found == false {
		return cache.ErrNotFound
	}
	bytes := data.([]byte)
	c.gocache.Set(key, bytes, ttl)
	return nil
}

//ExpireCounter set cache counter  expire duration by given key and ttl
func (c *Cache) ExpireCounter(key string, ttl time.Duration) error {
	locker := c.Util().Locker(key)
	locker.Lock()
	defer locker.Unlock()
	data, found := c.gocache.Get(key)
	if found == false {
		return cache.ErrNotFound
	}
	bytes := data.([]byte)
	c.gocache.Set(key, bytes, ttl)
	return nil
}

//Config Cache driver config.
type Config struct {
	DefaultExpirationInSecond int64 //Cache memory usage limie.
	CleanupIntervalInSecond   int64
}

//Create new cache driver.
//Return cache driver created and any error if raised.
func (config *Config) Create() (cache.Driver, error) {
	cache := Cache{
		gocache: gocache.New(time.Duration(config.DefaultExpirationInSecond)*time.Second, time.Duration(config.CleanupIntervalInSecond)*time.Second),
	}
	return &cache, nil
}

func init() {
	cache.Register("gocache", func(conf cache.Config, prefix string) (cache.Driver, error) {
		c := &Config{}
		conf.Get(prefix+"DefaultExpirationInSecond", &c.DefaultExpirationInSecond)
		if c.DefaultExpirationInSecond == 0 {
			c.DefaultExpirationInSecond = defaultExpirationInsecond
		}
		conf.Get(prefix+"CleanupIntervalInSecond", &c.CleanupIntervalInSecond)
		if c.CleanupIntervalInSecond == 0 {
			c.CleanupIntervalInSecond = defaultCleanupIntervalInSecond
		}
		return c.Create()
	})
}
