package freecache

import (
	"time"

	"encoding/binary"
	"encoding/json"

	"sync"

	"github.com/coocood/freecache"
	"github.com/herb-go/herb/cache"
)

type Cache struct {
	freecache    *freecache.Cache
	gcErrHandler func(err error)
	lock         sync.Mutex
}

func (c *Cache) SetGCErrHandler(f func(err error)) {
	c.gcErrHandler = f
	return
}
func (c *Cache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	err := c.freecache.Set([]byte(key), bytes, int(ttl/time.Second))
	if err == freecache.ErrLargeEntry {
		return cache.ErrEntryTooLarge
	}
	return err
}
func (c *Cache) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {
	c.lock.Lock()
	defer c.lock.Unlock()
	_, err := c.freecache.Get([]byte(key))
	if err == freecache.ErrNotFound {
		return nil
	}
	err = c.freecache.Set([]byte(key), bytes, int(ttl/time.Second))
	if err == freecache.ErrLargeEntry {
		return cache.ErrEntryTooLarge
	}
	return err
}
func (c *Cache) GetBytesValue(key string) ([]byte, error) {
	bytes, err := c.freecache.Get([]byte(key))
	if err == freecache.ErrNotFound {
		err = cache.ErrNotFound
	}
	return bytes, err
}

func (c *Cache) Set(key string, v interface{}, ttl time.Duration) error {
	bytes, err := cache.MarshalMsgpack(&v)
	if err != nil {
		return err
	}
	return c.SetBytesValue(key, bytes, ttl)
}
func (c *Cache) Update(key string, v interface{}, ttl time.Duration) error {
	bytes, err := cache.MarshalMsgpack(&v)
	if err != nil {
		return err
	}
	return c.UpdateBytesValue(key, bytes, ttl)
}
func (c *Cache) Get(key string, v interface{}) error {
	bytes, err := c.GetBytesValue(key)
	if err != nil {
		return err
	}
	return cache.UnmarshalMsgpack(bytes, &v)
}

func (c *Cache) Flush() error {
	c.freecache.Clear()
	return nil
}
func (c *Cache) Close() error {
	c.freecache.Clear()
	return nil
}
func (c *Cache) Del(key string) error {
	_ = c.freecache.Del([]byte(key))
	return nil
}
func (c *Cache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	var v int64
	var err error
	c.lock.Lock()
	defer c.lock.Unlock()
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
func (c *Cache) SetCounter(key string, v int64, ttl time.Duration) error {
	bytes := make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(v))
	err := c.freecache.Set([]byte(key), bytes, int(ttl/time.Second))
	return err
}
func (c *Cache) GetCounter(key string) (int64, error) {
	var v int64
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
func (c *Cache) DelCounter(key string) error {
	_ = c.freecache.Del([]byte(key))
	return nil
}
func (c *Cache) SearchByPrefix(prefix string) ([]string, error) {
	return nil, cache.ErrSearchKeysNotSupported
}

type Config struct {
	Size int
}

func (_ *Cache) New(config json.RawMessage) (cache.Driver, error) {
	c := Config{}
	err := json.Unmarshal(config, &c)
	if err != nil {
		return nil, err
	}
	cache := Cache{
		freecache: freecache.NewCache(c.Size),
	}
	return &cache, nil
}

func init() {
	cache.Register("freecache", &Cache{})
}
