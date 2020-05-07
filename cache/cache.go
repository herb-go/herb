//Package cache provide a key-value data store with ttl interface.
package cache

import (
	"errors"
	"reflect"
	"sync/atomic"
	"time"
)

var (
	//ErrNotFound raised when the given data not found in cache.
	ErrNotFound = errors.New("Entry not found")
	//ErrNotCacheable raised if the data cannot be cached.
	ErrNotCacheable = errors.New("Not Cachable")
	//ErrEntryTooLarge raised when data is too large to store.
	ErrEntryTooLarge = errors.New("Entry too large to cache")
	//ErrKeyTooLarge raised when key is too large to store.
	ErrKeyTooLarge = errors.New("Key too large to cache")
	//ErrKeyUnavailable raised when the key is not available.For example,empty key.
	ErrKeyUnavailable = errors.New("Key Unavailable")
	//ErrFeatureNotSupported raised when calling feature on unsupported driver.
	ErrFeatureNotSupported = errors.New("Feature is not supported")
	//ErrTTLNotAvaliable raised when ttl not avaliable.
	ErrTTLNotAvaliable = errors.New("TTL not avaliable")
)

//DefaultTTL means use cache default ttl setting.
var DefaultTTL = time.Duration(0)

var (
	//KeyPrefix default key prefix
	KeyPrefix    = string([]byte{0})
	intKeyPrefix = string([]byte{69, 0})
)

//Key return cache key
func Key(key string) string {
	return KeyPrefix + key

}

//New :Create a empty cache.
func New() *Cache {
	hit := int64(0)
	miss := int64(0)
	return &Cache{
		hit:  &hit,
		miss: &miss,
	}
}

//Cache Cache stores the cache Driver and default ttl.
type Cache struct {
	Driver
	TTL  time.Duration
	hit  *int64
	miss *int64
}

//Hit return cache hit count
func (c *Cache) Hit() int64 {
	return atomic.LoadInt64(c.hit)
}

//Miss return cache miss count
func (c *Cache) Miss() int64 {
	return atomic.LoadInt64(c.miss)
}

func (c *Cache) getKey(key string) string {
	return Key(key)
}

//Init init cache with option
func (c *Cache) Init(option Option) error {
	return option.ApplyTo(c)
}

//Set Set data model to cache by given key.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (c *Cache) Set(key string, v interface{}, ttl time.Duration) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
	}
	bs, err := c.Driver.Util().Marshaler.Marshal(v)
	if err != nil {
		return err
	}
	return c.Driver.SetBytesValue(c.getKey(key), bs, ttl)
}

//Update Update data model to cache by given key only if the cache exist.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (c *Cache) Update(key string, v interface{}, ttl time.Duration) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
	}
	bs, err := c.Driver.Util().Marshaler.Marshal(v)
	if err != nil {
		return err
	}
	return c.Driver.UpdateBytesValue(c.getKey(key), bs, ttl)
}

//Get Get data model from cache by given key.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raised.
func (c *Cache) Get(key string, v interface{}) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	bs, err := c.Driver.GetBytesValue(c.getKey(key))
	if err != nil {
		return err
	}
	return c.Driver.Util().Marshaler.Unmarshal(bs, v)
}

//SetBytesValue Set bytes data to cache by given key.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (c *Cache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
	}
	if ttl < 0 {
		return ErrTTLNotAvaliable
	}
	return c.Driver.SetBytesValue(c.getKey(key), bytes, ttl)
}

//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (c *Cache) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
	}
	if ttl < 0 {
		return ErrTTLNotAvaliable
	}
	return c.Driver.UpdateBytesValue(c.getKey(key), bytes, ttl)
}

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (c *Cache) GetBytesValue(key string) ([]byte, error) {
	if key == "" {
		return nil, ErrKeyUnavailable
	}
	bs, err := c.Driver.GetBytesValue(c.getKey(key))
	if err != nil {
		atomic.AddInt64(c.hit, 1)
	} else if err == ErrNotFound {
		atomic.AddInt64(c.miss, 1)
	}
	return bs, err
}

//MGetBytesValue get multiple bytes data from cache by given keys.
//Return data bytes map and any error if raised.
func (c *Cache) MGetBytesValue(keys ...string) (map[string][]byte, error) {
	var result map[string][]byte
	var prefixedKeys = make([]string, len(keys))
	for k := range keys {
		prefixedKeys[k] = c.getKey(keys[k])
	}
	data, err := c.Driver.MGetBytesValue(prefixedKeys...)
	if err != nil {
		return result, err
	}
	result = make(map[string][]byte, len(data))
	for k := range data {
		result[k[len(KeyPrefix):]] = data[k]
	}
	atomic.AddInt64(c.hit, int64(len(data)))
	atomic.AddInt64(c.miss, int64(len(keys)-len(data)))
	return result, nil
}

//MSetBytesValue set multiple bytes data to cache with given key-value map.
//Return  any error if raised.
func (c *Cache) MSetBytesValue(data map[string][]byte, ttl time.Duration) error {
	var prefixed = make(map[string][]byte, len(data))
	for k := range data {
		prefixed[c.getKey(k)] = data[k]
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
	}
	if ttl < 0 {
		return ErrTTLNotAvaliable
	}
	return c.Driver.MSetBytesValue(prefixed, ttl)
}

//Del Delete data in cache by given name.
//Return any error raised.
func (c *Cache) Del(key string) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	return c.Driver.Del(c.getKey(key))
}

//Expire set cache value expire duration by given key and ttl
func (c *Cache) Expire(key string, ttl time.Duration) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
	}
	if ttl < 0 {
		return ErrTTLNotAvaliable
	}
	err := c.Driver.Expire(c.getKey(key), ttl)
	if err == ErrNotFound {
		err = nil
	}
	return err
}

func (c *Cache) getIntKey(key string) string {
	return intKeyPrefix + key
}

//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return int data value and any error raised.
func (c *Cache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	if key == "" {
		return 0, ErrKeyUnavailable
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
	}
	return c.Driver.IncrCounter(c.getIntKey(key), increment, ttl)
}

//SetCounter Set int val in cache by given key.Count cache and data cache are in two independent namespace.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (c *Cache) SetCounter(key string, v int64, ttl time.Duration) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
	}
	if ttl < 0 {
		return ErrTTLNotAvaliable
	}
	return c.Driver.SetCounter(c.getIntKey(key), v, ttl)
}

//GetCounter Get int val from cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) GetCounter(key string) (int64, error) {
	if key == "" {
		return 0, ErrKeyUnavailable
	}
	return c.Driver.GetCounter(c.getIntKey(key))
}

//DelCounter Delete int val in cache by given name.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) DelCounter(key string) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	err := c.Driver.DelCounter(c.getIntKey(key))
	if err == ErrNotFound {
		return nil
	}
	return err
}

//ExpireCounter set cache counter  expire duration by given key and ttl
func (c *Cache) ExpireCounter(key string, ttl time.Duration) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	if ttl < 0 {
		return ErrTTLNotAvaliable
	}
	err := c.Driver.ExpireCounter(c.getIntKey(key), ttl)
	if err == ErrNotFound {
		return nil
	}
	return err
}

//Locker create new locker with given key.
//return locker and if locker aleady locked.
func (c *Cache) Locker(key string) (*Locker, bool) {
	return c.Util().Locker(key)
}
func loadFromCache(c Cacheable, key string, v interface{}, ttl time.Duration, loader Loader) error {
	var err error
	if key == "" {
		return ErrKeyUnavailable
	}
	err = c.Get(key, v)
	if err == ErrNotFound || err == ErrKeyTooLarge {
		k := c.FinalKey(key)
		locker, ok := c.Util().Locker(k)
		if ok {
			locker.RLock()
			defer locker.RUnlock()
			err = c.Get(key, v)
			if err == nil || (err != ErrNotFound && err != ErrKeyTooLarge) {
				return err
			}
		} else {
			locker.Lock()
			defer locker.Unlock()
		}
		v2, err2 := loader(key)
		if err2 != nil {
			return err2
		}
		reflect.Indirect(reflect.ValueOf(v)).Set(reflect.Indirect(reflect.ValueOf(v2)))
		err3 := c.Set(key, v, ttl)
		if err3 == ErrNotCacheable || err3 == ErrEntryTooLarge || err3 == ErrKeyTooLarge {
			return nil
		} else if err3 != nil {
			return err3
		}
	} else if err != nil {
		return err
	}
	return nil
}

//Load Get data model from cache by given key.If data not found,call loader to get current data value and save to cache.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (c *Cache) Load(key string, v interface{}, ttl time.Duration, loader Loader) error {
	return loadFromCache(c, key, v, ttl, loader)
}

//FinalKey get final key which passed to cache driver .
func (c *Cache) FinalKey(key string) string {
	return c.getKey(key)
}

//Field retuan a cache field with given field name
func (c *Cache) Field(fieldname string) *Field {
	return &Field{
		Cache:     c,
		FieldName: fieldname,
	}
}

//DefualtTTL return cache default ttl
func (c *Cache) DefualtTTL() time.Duration {
	return c.TTL
}

//Proxy get a cache proxy with given prefix
func (c *Cache) Proxy(prefix string) *Proxy {
	return NewProxy(NewCollection(c, prefix, c.TTL))
}

//Marshal Marshal data model to  bytes.
//Return marshaled bytes and any error rasied.
func (c *Cache) Marshal(v interface{}) ([]byte, error) {
	return c.Driver.Util().Marshaler.Marshal(v)
}

//Unmarshal Unmarshal bytes to data model.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raseid.
func (c *Cache) Unmarshal(bytes []byte, v interface{}) error {
	return c.Driver.Util().Marshaler.Unmarshal(bytes, v)
}

//Loader cache value loader used in cache load method.
//Load value with given key.
//Return loaded value and any error if raised.
type Loader func(key string) (interface{}, error)
