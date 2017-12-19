package cache

import (
	"strconv"
	"time"
)

type Collection struct {
	Cache  Cacheable
	Prefix string
	TTL    time.Duration
}

var CollectionTTLMultiple = 10

func NewCollection(cache Cacheable, prefix string, TTL time.Duration) *Collection {
	return &Collection{
		Cache:  cache,
		Prefix: prefix,
		TTL:    TTL,
	}
}
func (c *Collection) GetCacheKey(key string) (string, error) {
	var ts string
	var data int64
	err := c.Cache.Get(c.Prefix, &ts)
	if err == ErrNotFound {
		data = time.Now().UnixNano()
		ts = strconv.FormatInt(data, 10)
		err = nil
		ttl := c.TTL
		if !c.persist() {
			ttl = ttl * time.Duration(CollectionTTLMultiple)
		}
		err2 := c.Cache.Set(c.Prefix, ts, ttl)
		if err2 == ErrNotCacheable {
			err2 = nil
		}
		if err2 != nil {
			return "", err2
		}
	}
	if err != nil {
		return "", err
	}
	return c.Prefix + KeyPrefix + ts + KeyPrefix + key, nil
}
func (c *Collection) persist() bool {
	return c.TTL < 0 || (c.TTL == 0 && c.Cache.DefualtTTL() < 0)
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

func (c *Collection) MGetBytesValue(keys ...string) (map[string][]byte, error) {
	prefix, err := c.GetCacheKey("")
	if err != nil {
		return map[string][]byte{}, err
	}
	var prefixedKeys = make([]string, len(keys))
	for k := range keys {
		prefixedKeys[k] = prefix + keys[k]
	}
	return c.Cache.MGetBytesValue(prefixedKeys...)
}
func (c *Collection) MSetBytesValue(data map[string][]byte, ttl time.Duration) error {
	prefix, err := c.GetCacheKey("")
	if err != nil {
		return err
	}
	var prefixed = make(map[string][]byte, len(data))
	for k := range data {
		prefixed[prefix+k] = data[k]
	}
	return c.Cache.MSetBytesValue(prefixed, ttl)
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
func (n *Collection) Expire(key string, ttl time.Duration) error {
	k, err := n.GetCacheKey(key)
	if err != nil {
		return err
	}
	return n.Cache.Expire(k, ttl)
}
func (n *Collection) ExpireCounter(key string, ttl time.Duration) error {
	k, err := n.GetCacheKey(key)
	if err != nil {
		return err
	}
	return n.Cache.ExpireCounter(k, ttl)
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
