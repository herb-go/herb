package cache

import "time"

//Cacheable interface which can be used as cache.
type Cacheable interface {
	//Set Set data model to cache by given key.
	//If ttl is DefualtTTL(0),use default ttl in config instead.
	//Return any error raised.
	Set(key string, v interface{}, ttl time.Duration) error
	//Get Get data model from cache by given key.
	//Parameter v should be pointer to empty data model which data filled in.
	//Return any error raised.
	Get(key string, v interface{}) error
	//Update Update data model to cache by given key only if the cache exist.
	//If ttl is DefualtTTL(0),use default ttl in config instead.
	//Return any error raised.
	Update(key string, v interface{}, ttl time.Duration) error
	//SetBytesValue Set bytes data to cache by given key.
	//If ttl is DefualtTTL(0),use default ttl in config instead.
	//Return any error raised.
	SetBytesValue(key string, bytes []byte, ttl time.Duration) error
	//GetBytesValue Get bytes data from cache by given key.
	//Return data bytes and any error raised.
	GetBytesValue(key string) ([]byte, error)
	//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
	//If ttl is DefualtTTL(0),use default ttl in config instead.
	//Return any error raised.
	UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error
	//Del Delete data in cache by given name.
	//Return any error raised.
	Del(key string) error
	//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
	//If ttl is DefualtTTL(0),use default ttl in config instead.
	//Return int data value and any error raised.
	IncrCounter(key string, increment int64, ttl time.Duration) (int64, error)
	//SetCounter Set int val in cache by given key.Count cache and data cache are in two independent namespace.
	//If ttl is DefualtTTL(0),use default ttl in config instead.
	//Return any error raised.
	SetCounter(key string, v int64, ttl time.Duration) error
	//GetCounter Get int val from cache by given key.Count cache and data cache are in two independent namespace.
	//Return int data value and any error raised.
	GetCounter(key string) (int64, error)
	//DelCounter Delete int val in cache by given name.Count cache and data cache are in two independent namespace.
	//Return any error raised.
	DelCounter(key string) error
	//Expire set cache value expire duration by given key and ttl
	Expire(key string, ttl time.Duration) error
	//ExpireCounter set cache counter  expire duration by given key and ttl
	ExpireCounter(key string, ttl time.Duration) error
	//Load Get data model from cache by given key.If data not found,call loader to get current data value and save to cache.
	//If ttl is DefualtTTL(0),use default ttl in config instead.
	//Return any error raised.
	Load(key string, v interface{}, ttl time.Duration, loader Loader) error
	//MGetBytesValue get multiple bytes data from cache by given keys.
	//Return data bytes map and any error if raised.
	MGetBytesValue(keys ...string) (map[string][]byte, error)
	//MSetBytesValue set multiple bytes data to cache with given key-value map.
	//Return  any error if raised.
	MSetBytesValue(map[string][]byte, time.Duration) error
	//FinalKey get final key which passed to cache driver .
	FinalKey(string) (string, error)
	//Flush Delete all data in cache.
	Flush() error
	//DefualtTTL return cache default ttl
	DefualtTTL() time.Duration
	Lock(key string) (unlocker func(), err error)
	Wait(key string) (bool, error)
}
