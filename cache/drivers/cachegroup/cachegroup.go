package cachegroup

import (
	"encoding/binary"
	"encoding/json"
	"time"

	"github.com/herb-go/herb/cache"
)

type Config []cache.Config
type Cache struct {
	SubCaches []*cache.Cache
}
type entry []byte

func (c *Cache) SearchByPrefix(prefix string) ([]string, error) {
	return nil, cache.ErrSearchKeysNotSupported
}
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
func (ca *Cache) New(bytes json.RawMessage) (cache.Driver, error) {
	config := Config{}
	err := json.Unmarshal(bytes, &config)
	if err != nil {
		return nil, err
	}
	c := Cache{}
	c.SubCaches = make([]*cache.Cache, len(config))
	for k, v := range config {
		subcache := cache.New()
		err := subcache.OpenConfig(v)
		if err != nil {
			return &c, err
		}
		c.SubCaches[k] = subcache
	}
	return &c, nil
}
func (c *Cache) Set(key string, v interface{}, ttl time.Duration) error {
	var bytes []byte
	bytes, err := cache.MarshalMsgpack(v)
	if err != nil {
		return err
	}
	return c.SetBytesValue(key, bytes, ttl)
}
func (c *Cache) Get(key string, v interface{}) error {
	bytes, err := c.GetBytesValue(key)
	if err != nil {
		return err
	}
	return cache.UnmarshalMsgpack(bytes, &v)
}
func (c *Cache) setBytesCaches(key string, caches []*cache.Cache, bytes []byte, expired int64) error {
	var finalErr error
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
		err := v.SetBytesValue(key, bytes, ttl)
		if err != nil && err != cache.ErrNotCacheable && err != cache.ErrEntryTooLarge {
			finalErr = err
		}
	}
	return finalErr
}

func (c *Cache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	var err error
	var e entry
	expired := e.Set(bytes, ttl)
	c.SubCaches[len(c.SubCaches)-1].SetBytesValue(key, []byte(e), ttl)
	if err != cache.ErrNotCacheable && err != cache.ErrEntryTooLarge && err != nil {
		return err
	}
	err = c.setBytesCaches(key, c.SubCaches[0:len(c.SubCaches)-1], []byte(e), expired)
	return err
}
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
	c.setBytesCaches(key, expiredCache, []byte(e), expired)
	return buf, nil
}
func (c *Cache) SetCounter(key string, v int64, ttl time.Duration) error {
	return c.SubCaches[len(c.SubCaches)-1].SetCounter(key, v, ttl)
}
func (c *Cache) GetCounter(key string) (int64, error) {
	return c.SubCaches[len(c.SubCaches)-1].GetCounter(key)
}

func (c *Cache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	return c.SubCaches[len(c.SubCaches)-1].IncrCounter(key, increment, ttl)
}
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
func (c *Cache) DelCounter(key string) error {
	return c.SubCaches[len(c.SubCaches)-1].DelCounter(key)
}
func (c *Cache) SetGCErrHandler(f func(err error)) {
	for _, v := range c.SubCaches {
		v.SetGCErrHandler(f)
	}
}
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
