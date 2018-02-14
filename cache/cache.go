//Package cache provide a key-value data store with ttl interface.
package cache

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"
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

	ErrPermanentCacheNotSupport = errors.New("Permanent cache is not supported.can use ttl <0 on this cache")
)

//DefualtTTL means use cache default ttl setting.
var DefualtTTL = time.Duration(0)

//TTLForever When cache ttl sets to TTLForever,the cache will never expire.
var TTLForever = time.Duration(-1)

var (
	KeyPrefix    = string([]byte{0})
	intKeyPrefix = string([]byte{69, 0})
)

//Config :The cache config struct
type Config struct {
	Driver string
	Config json.RawMessage
	TTL    int64
}

func (config *Config) Init(cache *Cache) error {
	return Option(config.Driver, config.Config, config.TTL)(cache)
}

//Driver : Cache driver interface.Should Never used directly
type Driver interface {
	New(cacheConfig json.RawMessage) (Driver, error)                           //Create new cache with given config.
	Set(key string, v interface{}, ttl time.Duration) error                    //Set data model to cache by given key.
	Update(key string, v interface{}, ttl time.Duration) error                 //Update data model to cache by given key only if the cache exist.
	Get(key string, v interface{}) error                                       //Get data model from cache by given key.
	SetBytesValue(key string, bytes []byte, ttl time.Duration) error           //Set bytes data to cache by given key.
	UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error        //Update bytes data to cache by given key only if the cache exist.
	GetBytesValue(key string) ([]byte, error)                                  //Get bytes data from cache by given key.
	Del(key string) error                                                      //Delete data in cache by given key.
	IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) //Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
	SetCounter(key string, v int64, ttl time.Duration) error                   //Set int val in cache by given key.Count cache and data cache are in two independent namespace.
	GetCounter(key string) (int64, error)                                      //Get int val from cache by given key.Count cache and data cache are in two independent namespace.
	DelCounter(key string) error                                               //Delete int val in cache by given key.Count cache and data cache are in two independent namespace.
	SetGCErrHandler(f func(err error))                                         //Set callback to handler error raised when gc.
	Expire(key string, ttl time.Duration) error
	ExpireCounter(key string, ttl time.Duration) error
	MGetBytesValue(keys ...string) (map[string][]byte, error)
	MSetBytesValue(map[string][]byte, time.Duration) error
	Close() error //Close cache.
	Flush() error //Delete all data in cache.
}

var (
	driversMu sync.RWMutex
	drivers   = make(map[string]Driver)
)

// Register makes a database driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, driver Driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if driver == nil {
		panic("cache: Register driver is nil")
	}
	if _, dup := drivers[name]; dup {
		panic("cache: Register called twice for driver " + name)
	}
	drivers[name] = driver
}
func unregisterAllDrivers() {
	driversMu.Lock()
	defer driversMu.Unlock()
	// For tests.
	drivers = make(map[string]Driver)
}

// Drivers returns a sorted list of the names of the registered drivers.
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	var list []string
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

//New :Create a empty cache.
func New() *Cache {
	return &Cache{}
}

func Init(cache *Cache, i Initializer) error {
	return i.Init(cache)
}

func NewAndInit(i Initializer) (*Cache, error) {
	cache := New()
	return cache, i.Init(cache)
}

//Cache Cache stores the cache Driver and default ttl.
type Cache struct {
	Driver
	TTL time.Duration
}

func (c *Cache) getKey(key string) string {
	return KeyPrefix + key
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
	return c.Driver.Set(c.getKey(key), v, ttl)
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
	return c.Driver.Update(c.getKey(key), v, ttl)
}

//Get Get data model from cache by given key.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raised.
func (c *Cache) Get(key string, v interface{}) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	return c.Driver.Get(c.getKey(key), v)
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

	return c.Driver.UpdateBytesValue(c.getKey(key), bytes, ttl)
}

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (c *Cache) GetBytesValue(key string) ([]byte, error) {
	if key == "" {
		return nil, ErrKeyUnavailable
	}
	return c.Driver.GetBytesValue(c.getKey(key))
}
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
	return result, nil
}
func (c *Cache) MSetBytesValue(data map[string][]byte, ttl time.Duration) error {
	var prefixed = make(map[string][]byte, len(data))
	for k := range data {
		prefixed[c.getKey(k)] = data[k]
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
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

func (c *Cache) Expire(key string, ttl time.Duration) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
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

func (c *Cache) ExpireCounter(key string, ttl time.Duration) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	err := c.Driver.ExpireCounter(c.getIntKey(key), ttl)
	if err == ErrNotFound {
		return nil
	}
	return err
}

//Load Get data model from cache by given key.If data not found,call closure to get current data value and save to cache.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (c *Cache) Load(key string, v interface{}, ttl time.Duration, closure func(v interface{}) error) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	err := c.Get(key, v)
	if err == ErrNotFound {
		err2 := closure(v)
		if err2 != nil {
			return err2
		}
		err3 := c.Set(key, v, ttl)
		if err3 == ErrNotCacheable {
			return nil
		} else if err3 != nil {
			return err3
		}
	} else if err != nil {
		return err
	}
	return nil
}

func (c *Cache) Field(fieldname string) *Field {
	return &Field{
		Cache:     c,
		FieldName: fieldname,
	}
}

func (c *Cache) DefualtTTL() time.Duration {
	return c.TTL
}

func (c *Cache) Collection(prefix string) *Collection {
	return NewCollection(c, prefix, c.TTL)
}
func (c *Cache) Node(prefix string) *Node {
	return NewNode(c, prefix)
}
