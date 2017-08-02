//Package cachegroup provides a cache drive which combined with give caches driver.
//Different cache driver will used their own default ttl to store data.
package cachegroup

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/herb-go/herb/cache"
)

const modeSet = 0
const modeUpdate = 1

//Config Cache driver list of all sub cache driver.
type Config []cache.Config

//Cache The group cache driver.
type Cache struct {
	SubCaches []*cache.Cache
}
type entry []byte

//SearchByPrefix Search All key start with given prefix.
//Return All matched key and any error raised.
func (c *Cache) SearchByPrefix(prefix string) ([]string, error) {
	return nil, cache.ErrSearchKeysNotSupported
}

//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (e *entry) Set(bytes []byte, ttl time.Duration) int64 {
	var expired int64
	var buf = make([]byte, 8)
	*e = make([]byte, len(bytes)+8)
	copy((*e)[8:], bytes)
	if ttl < 0 {
		expired = -1
	} else {
		expired = time.Now().Add(ttl).Unix()
	}
	binary.BigEndian.PutUint64(buf, uint64(expired))
	copy((*e)[0:8], buf)
	return expired
}
func (e *entry) Get() ([]byte, int64, error) {
	var b = make([]byte, len(*e))
	copy(b, *e)
	var buf []byte
	var expired int64
	if len(b) < 8 {
		return buf, expired, cache.ErrNotFound
	}
	expired = int64(binary.BigEndian.Uint64(b[0:8]))
	if expired >= 0 && expired < time.Now().Unix() {
		return buf, expired, cache.ErrNotFound
	}
	buf = make([]byte, len(b)-8)
	copy(buf, b[8:])
	return buf, expired, nil

}

//New Create new cache driver with given json bytes.
//Return new driver and any error raised.
func (c *Cache) New(bytes json.RawMessage) (cache.Driver, error) {
	config := Config{}
	err := json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	cc := Cache{}
	cc.SubCaches = make([]*cache.Cache, len(config))
	for k, v := range config {
		subcache := cache.New()
		err := subcache.OpenConfig(v)
		if err != nil {
			return &cc, err
		}
		cc.SubCaches[k] = subcache
	}
	return &cc, nil
}

//Set Set data model to cache by given key.
//Return any error raised.
func (c *Cache) Set(key string, v interface{}, ttl time.Duration) error {
	var bytes []byte
	bytes, err := cache.MarshalMsgpack(v)
	if err != nil {
		return err
	}
	return c.SetBytesValue(key, bytes, ttl)
}

//Update Update data model to cache by given key only if the cache exist.
//Return any error raised.
func (c *Cache) Update(key string, v interface{}, ttl time.Duration) error {
	var bytes []byte
	bytes, err := cache.MarshalMsgpack(v)
	if err != nil {
		return err
	}
	return c.UpdateBytesValue(key, bytes, ttl)
}

//Get Get data model from cache by given key.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raised.
func (c *Cache) Get(key string, v interface{}) error {
	bytes, err := c.GetBytesValue(key)
	if err != nil {
		return err
	}
	return cache.UnmarshalMsgpack(bytes, &v)
}

func (c *Cache) setBytesCaches(key string, caches []*cache.Cache, bytes []byte, expired int64, mode int) error {
	var finalErr error
	var err error
	var t time.Duration
	if expired < 0 {
		t = -1
	} else {
		t = time.Unix(expired, 0).Sub(time.Now())
	}
	for _, v := range caches {
		var ttl time.Duration
		if t < 0 {
			if v.TTL < 0 {
				ttl = -1
			} else {
				ttl = v.TTL
			}
		} else {
			if v.TTL < 0 {
				ttl = t
			} else {
				if v.TTL < t {
					ttl = v.TTL
				} else {
					ttl = t
				}
			}
		}
		if mode == modeSet {
			err = v.SetBytesValue(key, bytes, ttl)
		} else {
			err = v.UpdateBytesValue(key, bytes, ttl)
		}
		if err != nil && err != cache.ErrNotCacheable && err != cache.ErrEntryTooLarge {
			finalErr = err
		}
	}
	return finalErr
}

//SetBytesValue Set bytes data to cache by given key.
//Return any error raised.
func (c *Cache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	var err error
	var e entry
	expired := e.Set(bytes, ttl)
	err = c.SubCaches[len(c.SubCaches)-1].SetBytesValue(key, []byte(e), ttl)
	if err != cache.ErrNotCacheable && err != cache.ErrEntryTooLarge && err != nil {
		return err
	}
	err = c.setBytesCaches(key, c.SubCaches[0:len(c.SubCaches)-1], []byte(e), expired, modeSet)
	return err
}

//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
//Return any error raised.
func (c *Cache) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {
	var err error
	var e entry
	expired := e.Set(bytes, ttl)
	err = c.SubCaches[len(c.SubCaches)-1].UpdateBytesValue(key, []byte(e), ttl)
	if err != cache.ErrNotCacheable && err != cache.ErrEntryTooLarge && err != nil {
		return err
	}
	err = c.setBytesCaches(key, c.SubCaches[0:len(c.SubCaches)-1], []byte(e), expired, modeUpdate)
	return err
}

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (c *Cache) GetBytesValue(key string) ([]byte, error) {
	var err error
	var bytes []byte
	var buf []byte
	expiredCache := []*cache.Cache{}
	for _, v := range c.SubCaches {
		bytes, err = v.GetBytesValue(key)
		if err == cache.ErrNotFound {
			expiredCache = append(expiredCache, v)
		} else {
			break
		}
	}
	if err != nil {
		return buf, err
	}
	e := entry(bytes)

	buf, expired, err := e.Get()
	if err != nil {
		return buf, err
	}
	c.setBytesCaches(key, expiredCache, []byte(e), expired, modeSet)
	return buf, nil
}

//SetCounter Set int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) SetCounter(key string, v int64, ttl time.Duration) error {
	return c.SubCaches[len(c.SubCaches)-1].SetCounter(key, v, ttl)
}

//GetCounter Get int val from cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) GetCounter(key string) (int64, error) {
	return c.SubCaches[len(c.SubCaches)-1].GetCounter(key)
}

//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	return c.SubCaches[len(c.SubCaches)-1].IncrCounter(key, increment, ttl)
}

//Del Delete data in cache by given key.
//Return any error raised.
func (c *Cache) Del(key string) error {
	var finalErr error
	for _, v := range c.SubCaches {
		err := v.Del(key)
		if err != nil {
			finalErr = err
		}
	}
	return finalErr
}

//DelCounter Delete int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) DelCounter(key string) error {
	return c.SubCaches[len(c.SubCaches)-1].DelCounter(key)
}

//SetGCErrHandler Set callback to handler error raised when gc.
func (c *Cache) SetGCErrHandler(f func(err error)) {
	for _, v := range c.SubCaches {
		v.SetGCErrHandler(f)
	}
}

//Close Close cache.
//Return any error if raised
func (c *Cache) Close() error {
	var finalErr error
	for _, v := range c.SubCaches {
		err := v.Close()
		if err != nil {
			finalErr = err
		}
	}
	return finalErr
}

//Flush Delete all data in cache.
//Return any error if raised
func (c *Cache) Flush() error {
	var finalErr error

	for _, v := range c.SubCaches {
		err := v.Flush()
		if err != nil {
			finalErr = err
		}
	}
	return finalErr
}

func init() {
	cache.Register("cachegroup", &Cache{})
}
