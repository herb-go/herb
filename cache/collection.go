package cache

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"time"
)

//Collection cache Collection
//Collection is flushable sub cache create from other cacheable.
type Collection struct {
	//Cache raw cache
	Cache Cacheable
	//Prefix cache key prefix
	Prefix string
	// default ttl
	TTL time.Duration
}

//CollectionTTLMultiple default collection ttl multiple
var CollectionTTLMultiple = 10

//NewCollection create new cache collection with given cache,prefix and ttl.
//Return collection created.
func NewCollection(cache Cacheable, prefix string, TTL time.Duration) *Collection {
	return &Collection{
		Cache:  cache,
		Prefix: prefix,
		TTL:    TTL,
	}
}

//GetCacheKey return raw cache key by given key.
//Return key and any error if raised.
func (c *Collection) GetCacheKey(key string) (string, error) {
	var ts string
	var data int64
	err := c.Cache.Get(c.Prefix, &ts)
	if err == ErrNotFound {
		data = time.Now().UnixNano()
		buf := bytes.NewBuffer(nil)
		err = binary.Write(buf, binary.BigEndian, data)
		if err != nil {
			return "", err
		}
		ts = hex.EncodeToString(buf.Bytes())
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

//Set Set data model to cache by given key.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
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

//Update Update data model to cache by given key only if the cache exist.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (c *Collection) Update(key string, v interface{}, TTL time.Duration) error {
	if TTL < 0 || (TTL == 0 && c.Cache.DefualtTTL() < 0) {
		return ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.Update(k, v, TTL)
}

//Get Get data model from cache by given key.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raised.
func (c *Collection) Get(key string, v interface{}) error {
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.Get(k, v)
}

//SetBytesValue Set bytes data to cache by given key.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
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

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (c *Collection) GetBytesValue(key string) ([]byte, error) {
	k, err := c.GetCacheKey(key)
	if err != nil {
		return nil, err
	}
	return c.Cache.GetBytesValue(k)
}

//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (c *Collection) UpdateBytesValue(key string, bytes []byte, TTL time.Duration) error {
	if TTL < 0 || (TTL == 0 && c.Cache.DefualtTTL() < 0) {
		return ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.UpdateBytesValue(k, bytes, TTL)
}

//MGetBytesValue get multiple bytes data from cache by given keys.
//Return data bytes map and any error if raised.
func (c *Collection) MGetBytesValue(keys ...string) (map[string][]byte, error) {
	prefix, err := c.GetCacheKey("")
	var result map[string][]byte
	var prefixedKeys = make([]string, len(keys))
	for k := range keys {
		prefixedKeys[k] = prefix + keys[k]
	}
	data, err := c.Cache.MGetBytesValue(prefixedKeys...)
	if err != nil {
		return result, err
	}
	result = make(map[string][]byte, len(data))
	for k := range data {
		result[k[len(prefix):]] = data[k]
	}
	return result, nil

}

//MSetBytesValue set multiple bytes data to cache with given key-value map.
//Return  any error if raised.
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

//Del Delete data in cache by given name.
//Return any error raised.
func (c *Collection) Del(key string) error {
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.Del(k)
}

//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return int data value and any error raised.
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

//SetCounter Set int val in cache by given key.Count cache and data cache are in two independent namespace.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
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

//DelCounter Delete int val in cache by given name.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Collection) DelCounter(key string) error {
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.DelCounter(k)
}

//GetCounter Get int val from cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Collection) GetCounter(key string) (int64, error) {
	k, err := c.GetCacheKey(key)
	if err != nil {
		return 0, err
	}
	return c.Cache.GetCounter(k)

}

//Load Get data model from cache by given key.If data not found,call loader to get current data value and save to cache.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (c *Collection) Load(key string, v interface{}, TTL time.Duration, loader Loader) error {
	if TTL < 0 || (TTL == 0 && c.Cache.DefualtTTL() < 0) {
		return ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return loadFromCache(c, k, v, TTL, loader)
}

//Flush Delete all data in cache.
func (c *Collection) Flush() error {
	return c.Cache.Del(c.Prefix)
}

//DefualtTTL return cache default ttl
func (c *Collection) DefualtTTL() time.Duration {
	return c.Cache.DefualtTTL()
}

//Expire set cache value expire duration by given key and ttl
func (c *Collection) Expire(key string, TTL time.Duration) error {
	if TTL < 0 || (TTL == 0 && c.Cache.DefualtTTL() < 0) {
		return ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.Expire(k, TTL)
}

//ExpireCounter set cache counter  expire duration by given key and ttl
func (c *Collection) ExpireCounter(key string, TTL time.Duration) error {
	if TTL < 0 || (TTL == 0 && c.Cache.DefualtTTL() < 0) {
		return ErrPermanentCacheNotSupport
	}
	k, err := c.GetCacheKey(key)
	if err != nil {
		return err
	}
	return c.Cache.ExpireCounter(k, TTL)
}

//Locker create new locker with given key.
//return locker and if locker aleady locked
func (c *Collection) Locker(key string) (*Locker, bool) {
	return c.Cache.Locker(key)
}

//Marshal Marshal data model to  bytes.
//Return marshaled bytes and any error rasied.
func (c *Collection) Marshal(v interface{}) ([]byte, error) {
	return c.Cache.Marshal(v)
}

//Unmarshal Unmarshal bytes to data model.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raseid.
func (c *Collection) Unmarshal(bytes []byte, v interface{}) error {
	return c.Cache.Unmarshal(bytes, v)
}

//Collection get a cache colletion with given prefix
func (c *Collection) Collection(prefix string) *Collection {
	return NewCollection(c, prefix, c.TTL)
}

//Node get a cache node with given prefix
func (c *Collection) Node(prefix string) *Node {
	return NewNode(c, prefix)
}

//Field retuan a cache field with given field name
func (c *Collection) Field(fieldname string) *Field {
	return &Field{
		Cache:     c,
		FieldName: fieldname,
	}
}

//FinalKey get final key which passed to cache driver .
func (c *Collection) FinalKey(key string) (string, error) {
	return c.Cache.FinalKey(c.Prefix + KeyPrefix + key)
}
