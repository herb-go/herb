package cache

import "time"
import "strconv"

type Collection struct {
	Cache  Cacheable
	Prefix string
	TTL    time.Duration
}

func NewCollection(cache Cacheable, prefix string, TTL time.Duration) *Collection {
	return &Collection{
		Cache:  cache,
		Prefix: prefix,
		TTL:    TTL,
	}
}
func (c *Collection) GetCacheKey(key string) (string, error) {
	var ts string
	data, err := c.Cache.GetBytesValue(c.Prefix)
	if err == nil {
		ts = string(data)
	} else if err == ErrNotFound {
		err = nil
		ts = strconv.FormatInt(time.Now().UnixNano(), 10)
		err2 := c.Cache.SetBytesValue(c.Prefix, []byte(ts), c.TTL)
		if err2 == ErrNotCacheable {
			err2 = nil
		}
		if err2 != nil {
			return "", err2
		}
	} else {
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

func (c *Collection) Set(key string, v interface{}, TTL time.Duration) error {
	if TTL < 0 || (TTL == 0 && c.Cache.DefualtTTL() < 0) {
		return ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.Set(k, v, TTL)
}

func (c *Collection) Get(key string, v interface{}) error {
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.Get(k, &v)
}
func (c *Collection) SetBytesValue(key string, bytes []byte, TTL time.Duration) error {
	if TTL < 0 || (TTL == 0 && c.Cache.DefualtTTL() < 0) {
		return ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.SetBytesValue(k, bytes, TTL)

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
func (c *Collection) IncrCounter(key string, increment int64, TTL time.Duration) (int64, error) {
	if TTL < 0 || (TTL == 0 && c.Cache.DefualtTTL() < 0) {
		return 0, ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return 0, err
	}
	return c.Cache.IncrCounter(k, increment, TTL)

}
func (c *Collection) SetCounter(key string, v int64, TTL time.Duration) error {
	if TTL < 0 || (TTL == 0 && c.Cache.DefualtTTL() < 0) {
		return ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.SetCounter(k, v, TTL)

}
func (c *Collection) DelCounter(key string) error {
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.DelCounter(k)
}
func (c *Collection) GetCounter(key string) (int64, error) {
	k, err := c.GetCacheKey(key)
	if err != nil {
		return 0, err
	}
	return c.Cache.GetCounter(k)

}
func (c *Collection) Load(key string, v interface{}, TTL time.Duration, closure func(v interface{}) error) error {
	if TTL < 0 || (TTL == 0 && c.Cache.DefualtTTL() < 0) {
		return ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.Load(k, &v, TTL, closure)
}

func (c *Collection) Flush() error {
	return c.Cache.Del(c.Prefix)
}

func (c *Collection) DefualtTTL() time.Duration {
	return c.Cache.DefualtTTL()
}

func (c *Collection) Collection(prefix string) *Collection {
	return NewCollection(c, prefix, c.TTL)
}
func (c *Collection) Node(prefix string) *Node {
	return NewNode(c, prefix)
}
func (c *Collection) Field(fieldname string) *Field {
	return &Field{
		Cache:     c,
		FieldName: fieldname,
	}
}
