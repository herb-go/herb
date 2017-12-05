package cache

import "time"
import "strconv"

type Collection struct {
	Cache  Cacheable
	Prefix string
	ttl    time.Duration
}

func NewCollection(cache Cacheable, prefix string, ttlInSecond int64) *Collection {
	return &Collection{
		Cache:  cache,
		Prefix: prefix,
		ttl:    time.Duration(ttlInSecond) * time.Second,
	}
}
func (c *Collection) GetCacheKey(key string) (string, error) {
	var ts string
	err := c.Cache.Load(c.Prefix, &ts, c.ttl, func(v interface{}) error {
		ts = strconv.FormatInt(time.Now().UnixNano(), 10)
		return nil
	})
	if err != nil {
		return "", err
	}
	return c.Prefix + keyPrefix + ts + keyPrefix + key, nil
}
func (c *Collection) MustGetCacheKey(key string) string {
	k, err := c.GetCacheKey(key)
	if err != nil {
		panic(err)
	}
	return k
}

func (c *Collection) Set(key string, v interface{}, ttl time.Duration) error {
	if ttl < 0 || (ttl == 0 && c.Cache.DefualtTTL() < 0) {
		return ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.Set(k, v, ttl)
}

func (c *Collection) Get(key string, v interface{}) error {
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.Get(k, &v)
}
func (c *Collection) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	if ttl < 0 || (ttl == 0 && c.Cache.DefualtTTL() < 0) {
		return ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.SetBytesValue(k, bytes, ttl)

}
func (c *Collection) GetBytesValue(key string) ([]byte, error) {
	k, err := c.GetCacheKey(key)
	if err != nil {
		return nil, err
	}
	return c.GetBytesValue(k)
}
func (c *Collection) Del(key string) error {
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.Del(k)
}
func (c *Collection) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	if ttl < 0 || (ttl == 0 && c.Cache.DefualtTTL() < 0) {
		return 0, ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return 0, err
	}
	return c.Cache.IncrCounter(k, increment, ttl)

}
func (c *Collection) SetCounter(key string, v int64, ttl time.Duration) error {
	if ttl < 0 || (ttl == 0 && c.Cache.DefualtTTL() < 0) {
		return ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.SetCounter(k, v, ttl)

}
func (c *Collection) GetCounter(key string) (int64, error) {
	k, err := c.GetCacheKey(key)
	if err != nil {
		return 0, err
	}
	return c.Cache.GetCounter(k)

}
func (c *Collection) Load(key string, v interface{}, ttl time.Duration, closure func(v interface{}) error) error {
	if ttl < 0 || (ttl == 0 && c.Cache.DefualtTTL() < 0) {
		return ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.Load(k, &v, ttl, closure)
}

func (c *Collection) Flush() error {
	return c.Cache.Del(c.Prefix)
}

func (c *Collection) DefualtTTL() time.Duration {
	return c.Cache.DefualtTTL()
}

func (c *Collection) SubCollection(prefix string) *Collection {
	return &Collection{
		Cache:  c,
		Prefix: prefix,
		ttl:    c.ttl,
	}
}
