package cache

import (
	"time"
)

//Node cache Collection
//Node is Permanent-able sub cache create from other cacheable.
type Node struct {
	Cache  Cacheable
	Prefix string
}

//NewNode create new cache node with given cacheable and prefix.
//Return node created.
func NewNode(c Cacheable, prefix string) *Node {
	return &Node{
		Cache:  c,
		Prefix: prefix,
	}
}

//GetCacheKey return raw cache key by given key.
//Return key and any error if raised.
func (n *Node) GetCacheKey(key string) (string, error) {
	return n.Prefix + KeyPrefix + key, nil
}

//MustGetCacheKey return raw cache key by given key.
//Return key.
//Panic if any error raised.
func (n *Node) MustGetCacheKey(key string) string {
	k, _ := n.GetCacheKey(key)
	return k
}

//Set Set data model to cache by given key.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (n *Node) Set(key string, v interface{}, ttl time.Duration) error {
	k := n.MustGetCacheKey(key)
	return n.Cache.Set(k, v, ttl)
}

//Update Update data model to cache by given key only if the cache exist.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (n *Node) Update(key string, v interface{}, TTL time.Duration) error {
	k := n.MustGetCacheKey(key)
	return n.Cache.Update(k, v, TTL)
}

//Get Get data model from cache by given key.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raised.
func (n *Node) Get(key string, v interface{}) error {
	k := n.MustGetCacheKey(key)
	return n.Cache.Get(k, v)
}

//SetBytesValue Set bytes data to cache by given key.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (n *Node) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	k := n.MustGetCacheKey(key)
	return n.Cache.SetBytesValue(k, bytes, ttl)
}

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (n *Node) GetBytesValue(key string) ([]byte, error) {
	k := n.MustGetCacheKey(key)
	return n.Cache.GetBytesValue(k)
}

//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (n *Node) UpdateBytesValue(key string, bytes []byte, TTL time.Duration) error {
	k := n.MustGetCacheKey(key)
	return n.Cache.UpdateBytesValue(k, bytes, TTL)
}

//MGetBytesValue get multiple bytes data from cache by given keys.
//Return data bytes map and any error if raised.
func (n *Node) MGetBytesValue(keys ...string) (map[string][]byte, error) {
	var result map[string][]byte
	var prefixedKeys = make([]string, len(keys))
	for k := range keys {
		prefixedKeys[k] = n.MustGetCacheKey(keys[k])
	}
	data, err := n.Cache.MGetBytesValue(prefixedKeys...)
	if err != nil {
		return result, err
	}
	result = make(map[string][]byte, len(data))
	for k := range data {
		result[k[len(n.Prefix+KeyPrefix):]] = data[k]
	}
	return result, nil
}

//MSetBytesValue set multiple bytes data to cache with given key-value map.
//Return  any error if raised.
func (n *Node) MSetBytesValue(data map[string][]byte, ttl time.Duration) error {
	var prefixed = make(map[string][]byte, len(data))
	for k := range data {
		prefixed[n.MustGetCacheKey(k)] = data[k]
	}
	return n.Cache.MSetBytesValue(prefixed, ttl)
}

//Del Delete data in cache by given name.
//Return any error raised.
func (n *Node) Del(key string) error {
	k := n.MustGetCacheKey(key)
	return n.Cache.Del(k)
}

//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return int data value and any error raised.
func (n *Node) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	k := n.MustGetCacheKey(key)
	return n.Cache.IncrCounter(k, increment, ttl)
}

//SetCounter Set int val in cache by given key.Count cache and data cache are in two independent namespace.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (n *Node) SetCounter(key string, v int64, ttl time.Duration) error {
	k := n.MustGetCacheKey(key)
	return n.Cache.SetCounter(k, v, ttl)
}

//GetCounter Get int val from cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (n *Node) GetCounter(key string) (int64, error) {
	k := n.MustGetCacheKey(key)
	return n.Cache.GetCounter(k)
}

//Load Get data model from cache by given key.If data not found,call loader to get current data value and save to cache.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (n *Node) Load(key string, v interface{}, TTL time.Duration, loader Loader) error {
	k := n.MustGetCacheKey(key)
	return n.Cache.Load(k, v, TTL, loader)
}

//Flush Delete all data in cache.
func (n *Node) Flush() error {
	return ErrFeatureNotSupported
}

//DefualtTTL return cache default ttl
func (n *Node) DefualtTTL() time.Duration {
	return n.Cache.DefualtTTL()
}

//DelCounter Delete int val in cache by given name.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (n *Node) DelCounter(key string) error {
	k, err := n.GetCacheKey(key)
	if err != nil {
		return err
	}
	return n.Cache.DelCounter(k)
}

//Expire set cache value expire duration by given key and ttl
func (n *Node) Expire(key string, ttl time.Duration) error {
	k, err := n.GetCacheKey(key)
	if err != nil {
		return err
	}
	return n.Cache.Expire(k, ttl)
}

//ExpireCounter set cache counter  expire duration by given key and ttl
func (n *Node) ExpireCounter(key string, ttl time.Duration) error {
	k, err := n.GetCacheKey(key)
	if err != nil {
		return err
	}
	return n.Cache.ExpireCounter(k, ttl)
}

// Lock lock cache value by given key.
//Return  unlock function and any error if rasied
func (n *Node) Lock(key string) (unlocker func(), err error) {
	k, err := n.GetCacheKey(key)
	if err != nil {
		return nil, err
	}
	return n.Cache.Lock(k)
}

//Wait wait any usef lock unlcok.
//Return whether waited and any error if rasied.
func (n *Node) Wait(key string) (bool, error) {
	k, err := n.GetCacheKey(key)
	if err != nil {
		return false, err
	}
	return n.Cache.Wait(k)
}

//Marshal Marshal data model to  bytes.
//Return marshaled bytes and any error rasied.
func (n *Node) Marshal(v interface{}) ([]byte, error) {
	return n.Cache.Marshal(v)
}

//Unmarshal Unmarshal bytes to data model.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raseid.
func (n *Node) Unmarshal(bytes []byte, v interface{}) error {
	return n.Cache.Unmarshal(bytes, v)
}

//Collection get a cache colletion with given prefix
func (n *Node) Collection(prefix string) *Collection {
	return NewCollection(n, prefix, n.Cache.DefualtTTL())
}

//Node get a cache node with given prefix
func (n *Node) Node(prefix string) *Node {
	return NewNode(n.Cache, n.MustGetCacheKey(prefix))
}

//Field retuan a cache field with given field name
func (n *Node) Field(fieldname string) *Field {
	return &Field{
		Cache:     n,
		FieldName: fieldname,
	}
}

//FinalKey get final key which passed to cache driver .
func (n *Node) FinalKey(key string) (string, error) {
	return n.Cache.FinalKey(n.Prefix + KeyPrefix + key)
}
