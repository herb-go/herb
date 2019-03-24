//Package syncmapcache provides cache driver uses sync.map to store cache data.
package syncmapcache

import (
	"time"

	"encoding/binary"

	"sync"

	"github.com/herb-go/herb/cache"
)

const defaultCleanupIntervalInSecond = 60
const defaultSize = 50000000

type entry struct {
	NeverExpired bool
	Expired      time.Time
	Data         []byte
}

//Cache The gocache cache Driver.
type Cache struct {
	cache.DriverUtil
	datamap      *sync.Map
	Size         int64
	used         int64
	writelock    sync.Mutex
	GCInterval   time.Duration
	gcErrHandler func(err error)
	C            chan int
}

func (c *Cache) flush() {
	c.datamap.Range(func(key interface{}, value interface{}) bool {
		c.delete(key.(string))
		return true
	})
}

func (c *Cache) gc() {
	c.writelock.Lock()
	c.writelock.Unlock()
	c.datamap.Range(func(key interface{}, value interface{}) bool {
		e := value.(*entry)
		if !e.NeverExpired && e.Expired.Before(time.Now()) {
			size := int64(len(e.Data))
			c.used = c.used - size
		}
		c.datamap.Delete(key)
		return true
	})
}
func (c *Cache) get(key string) ([]byte, bool) {
	v, ok := c.datamap.Load(key)
	if ok == false || v == nil {
		return nil, false
	}
	e := v.(*entry)
	if e.NeverExpired || time.Now().Before(e.Expired) {
		return e.Data, true
	}
	return nil, false
}
func (c *Cache) delete(key string) {
	c.writelock.Lock()
	defer c.writelock.Unlock()
	defer c.datamap.Delete(key)
	v, ok := c.datamap.Load(key)
	if ok == false || v == nil {
		return
	}
	e := v.(*entry)
	size := int64(len(e.Data))
	c.used = c.used - size
}

func (c *Cache) removeAll() {
	c.datamap.Range(func(key interface{}, value interface{}) bool {
		e := value.(*entry)
		size := int64(len(e.Data))
		c.used = c.used - size
		c.datamap.Delete(key)
		return true
	})

}
func (c *Cache) set(key string, data []byte, ttl time.Duration) {
	var delta int64
	c.writelock.Lock()
	defer c.writelock.Unlock()
	if c.used+int64(len(data)) > c.Size {
		c.removeAll()
	}
	defer func() { c.used = c.used + delta }()
	v, ok := c.datamap.Load(key)
	e := &entry{
		NeverExpired: ttl < 0,
		Expired:      time.Now().Add(ttl),
		Data:         data,
	}
	c.datamap.Store(key, e)
	delta = int64(len(data))
	if ok == false || v == nil {
		return
	}
	e = v.(*entry)
	delta = delta - int64(len(e.Data))
}

func (c *Cache) replace(key string, data []byte, ttl time.Duration) {
	c.writelock.Lock()
	defer c.writelock.Unlock()
	v, ok := c.datamap.Load(key)
	if ok == false || v == nil {
		return
	}
	e := &entry{
		NeverExpired: ttl < 0,
		Expired:      time.Now().Add(ttl),
		Data:         data,
	}

	if c.used+int64(len(data))-int64(len(e.Data)) > c.Size {
		c.removeAll()
		e.Data = nil
	}
	c.datamap.Store(key, e)
	size := int64(len(data)) - int64(len(e.Data))
	c.used = c.used + size
}

//SetGCErrHandler Set callback to handler error raised when gc.
func (c *Cache) SetGCErrHandler(f func(err error)) {
	c.gcErrHandler = f
	return
}

//SetBytesValue Set bytes data to cache by given key.
//Return any error raised.
func (c *Cache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	c.set(key, bytes, ttl)
	return nil
}

//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
//Return any error raised.
func (c *Cache) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {
	c.replace(key, bytes, ttl)
	return nil
}

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (c *Cache) GetBytesValue(key string) ([]byte, error) {
	bs, found := c.get(key)
	if found {
		return bs, nil
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
			result[keys[k]] = nil
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
		c.set(k, data[k], ttl)
	}
	return nil
}

//Flush Delete all data in cache.
//Return any error if raised
func (c *Cache) Flush() error {
	c.flush()
	return nil
}

//Close Close cache.
//Return any error if raised
func (c *Cache) Close() error {
	close(c.C)
	return nil
}

//Del Delete data in cache by given key.
//Return any error raised.
func (c *Cache) Del(key string) error {
	c.delete(key)
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
	data, found := c.get(key)
	if found == false {
		v = 0
	} else {
		bs := data
		v = int64(binary.BigEndian.Uint64(bs[0:8]))
	}
	v = v + increment
	bs := make([]byte, 8)
	binary.BigEndian.PutUint64(bs, uint64(v))
	c.set(key, bs, ttl)
	return v, nil
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
	c.set(key, bytes, ttl)
	return nil
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

	bs, found := c.get(key)
	if found == false {
		err = cache.ErrNotFound
		return 0, err
	}
	v = int64(binary.BigEndian.Uint64(bs[0:8]))
	return v, nil
}

//DelCounter Delete int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raisegrd.
func (c *Cache) DelCounter(key string) error {
	unlocker, err := c.Util().Lock(key)
	if err != nil {
		return err
	}
	defer unlocker()

	c.delete(key)
	return nil
}

//Expire set cache value expire duration by given key and ttl
func (c *Cache) Expire(key string, ttl time.Duration) error {
	unlocker, err := c.Util().Lock(key)
	if err != nil {
		return err
	}
	defer unlocker()
	bs, found := c.get(key)
	if found == false {
		return cache.ErrNotFound
	}
	c.set(key, bs, ttl)
	return nil
}

//ExpireCounter set cache counter  expire duration by given key and ttl
func (c *Cache) ExpireCounter(key string, ttl time.Duration) error {
	unlocker, err := c.Util().Lock(key)
	if err != nil {
		return err
	}
	defer unlocker()
	bs, found := c.get(key)
	if found == false {
		return cache.ErrNotFound
	}
	c.set(key, bs, ttl)
	return nil
}

//Config Cache driver config.
type Config struct {
	CleanupIntervalInSecond int64
	Size                    int64
}

//Create new cache driver.
//Return cache driver created and any error if raised.
func (config *Config) Create() (cache.Driver, error) {
	cache := Cache{
		datamap: &sync.Map{},
		Size:    config.Size,
		C:       make(chan int),
	}
	gctick := time.Tick(time.Duration(config.CleanupIntervalInSecond) * time.Second)
	go func() {
		for {
			select {
			case <-gctick:
				cache.gc()
			case <-cache.C:
				return
			}
		}
	}()
	return &cache, nil
}

func init() {
	cache.Register("syncmapcache", func(conf cache.Config, prefix string) (cache.Driver, error) {
		c := &Config{}

		conf.Get(prefix+"CleanupIntervalInSecond", &c.CleanupIntervalInSecond)
		if c.CleanupIntervalInSecond == 0 {
			c.CleanupIntervalInSecond = defaultCleanupIntervalInSecond
		}
		conf.Get(prefix+"Size", &c.Size)
		if c.Size <= 0 {
			c.Size = defaultSize
		}
		return c.Create()
	})
}