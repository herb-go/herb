//Package freecache provides cache driver uses memory to store cache data.
//Using github.com/coocood/freecache as driver.
package freecache

import (
	"time"

	"encoding/binary"

	"sync"

	"github.com/coocood/freecache"
	"github.com/herb-go/herb/cache"
)

//Cache The freecache cache Driver.
type Cache struct {
	cache.DriverUtil
	freecache    *freecache.Cache
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
	c.lock.Lock()
	defer c.lock.Unlock()
	err := c.freecache.Set([]byte(key), bytes, int(ttl/time.Second))
	if err == freecache.ErrLargeEntry {
		return cache.ErrEntryTooLarge
	}
	return err
}

//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
//Return any error raised.
func (c *Cache) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.freecache.TTL([]byte(key))
	if err == freecache.ErrNotFound {
		return nil
	}
	err = c.freecache.Set([]byte(key), bytes, int(ttl/time.Second))
	if err == freecache.ErrLargeEntry {
		return cache.ErrEntryTooLarge
	}
	return err
}

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (c *Cache) GetBytesValue(key string) ([]byte, error) {
	bytes, err := c.freecache.Get([]byte(key))
	if err == freecache.ErrNotFound {
		err = cache.ErrNotFound
	}
	return bytes, err
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
	var err error
	var ttlsecond = int(ttl / time.Second)
	c.lock.Lock()
	defer c.lock.Unlock()
	for k := range data {
		err = c.freecache.Set([]byte(k), data[k], ttlsecond)
		if err == freecache.ErrLargeEntry {
			return cache.ErrEntryTooLarge
		}
		if err != nil {
			return err
		}
	}
	return nil
}

//Flush Delete all data in cache.
//Return any error if raised
func (c *Cache) Flush() error {
	c.freecache.Clear()
	return nil
}

//Close Close cache.
//Return any error if raised
func (c *Cache) Close() error {
	c.freecache.Clear()
	return nil
}

//Del Delete data in cache by given key.
//Return any error raised.
func (c *Cache) Del(key string) error {
	_ = c.freecache.Del([]byte(key))
	return nil
}

//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	var v int64
	var err error
	unlocker, err := c.Util().Lock(key)
	if err != nil {
		return 0, err
	}
	defer unlocker()
	bytes, err := c.freecache.Get([]byte(key))
	if err == freecache.ErrNotFound || bytes == nil || len(bytes) != 8 {
		v = 0
	} else if err != nil {
		return v, err
	} else {
		v = int64(binary.BigEndian.Uint64(bytes[0:8]))
	}
	v = v + increment
	bytes = make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(v))
	err = c.freecache.Set([]byte(key), bytes, int(ttl/time.Second))
	return v, err
}

//SetCounter Set int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) SetCounter(key string, v int64, ttl time.Duration) error {
	unlocker, err := c.Util().Lock(key)
	if err != nil {
		return err
	}
	defer unlocker()

	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(v))
	err = c.freecache.Set([]byte(key), bytes, int(ttl/time.Second))
	return err
}

//GetCounter Get int val from cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) GetCounter(key string) (int64, error) {
	var v int64
	unlocker, err := c.Util().Lock(key)
	if err != nil {
		return 0, err
	}
	defer unlocker()

	bytes, err := c.freecache.Get([]byte(key))

	if err == freecache.ErrNotFound || bytes == nil || len(bytes) != 8 {
		err = cache.ErrNotFound
	}
	if err != nil {
		return 0, err
	}
	v = int64(binary.BigEndian.Uint64(bytes[0:8]))
	return v, nil
}

//DelCounter Delete int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) DelCounter(key string) error {
	unlocker, err := c.Util().Lock(key)
	if err != nil {
		return err
	}
	defer unlocker()

	_ = c.freecache.Del([]byte(key))
	return nil
}

//Expire set cache value expire duration by given key and ttl
func (c *Cache) Expire(key string, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	b, err := c.freecache.Get([]byte(key))
	if err == freecache.ErrNotFound {
		return cache.ErrNotFound
	}
	if err != nil {
		return err
	}
	err = c.freecache.Set([]byte(key), b, int(ttl/time.Second))
	if err == freecache.ErrLargeEntry {
		return cache.ErrEntryTooLarge
	}
	return err
}

//ExpireCounter set cache counter  expire duration by given key and ttl
func (c *Cache) ExpireCounter(key string, ttl time.Duration) error {
	unlocker, err := c.Util().Lock(key)
	if err != nil {
		return err
	}
	defer unlocker()
	b, err := c.freecache.Get([]byte(key))
	if err == freecache.ErrNotFound {
		return cache.ErrNotFound
	}
	if err != nil {
		return err
	}
	err = c.freecache.Set([]byte(key), b, int(ttl/time.Second))
	if err == freecache.ErrLargeEntry {
		return cache.ErrEntryTooLarge
	}
	return err
}

//Config Cache driver config.
type Config struct {
	Size int //Cache memory usage limie.
}

//Create new cache driver.
//Return cache driver created and any error if raised.
func (config *Config) Create() (cache.Driver, error) {
	cache := Cache{
		freecache: freecache.NewCache(config.Size),
	}
	return &cache, nil
}

func init() {
	cache.Register("freecache", func(conf cache.Config, prefix string) (cache.Driver, error) {
		c := &Config{}
		conf.Get(prefix+"Size", &c.Size)
		return c.Create()
	})
}
