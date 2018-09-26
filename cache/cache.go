//Package cache provide a key-value data store with ttl interface.
package cache

import (
	"errors"
	"fmt"
	"reflect"
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
	//ErrPermanentCacheNotSupport raised when cache driver not support permanent ttl.
	ErrPermanentCacheNotSupport = errors.New("Permanent cache is not supported.can use ttl <0 on this cache")
)

//DefualtTTL means use cache default ttl setting.
var DefualtTTL = time.Duration(0)

//TTLForever When cache ttl sets to TTLForever,the cache will never expire.
var TTLForever = time.Duration(-1)

var (
	//KeyPrefix default key prefix
	KeyPrefix    = string([]byte{0})
	intKeyPrefix = string([]byte{69, 0})
)

//Factory create driver with given config and prefix
//Reutrn driver created and any error if raised..
type Factory func(conf Config, prefix string) (Driver, error)

//Driver : Cache driver interface.Should Never used directly
type Driver interface {
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
	factorysMu sync.RWMutex
	factories  = make(map[string]Factory)
)

// Register makes a driver creator available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, f Factory) {
	factorysMu.Lock()
	defer factorysMu.Unlock()
	if f == nil {
		panic("cache: Register cache factory is nil")
	}
	if _, dup := factories[name]; dup {
		panic("cache: Register called twice for factory " + name)
	}
	factories[name] = f
}
func unregisterAll() {
	factorysMu.Lock()
	defer factorysMu.Unlock()
	// For tests.
	factories = make(map[string]Factory)
}

//Factories returns a sorted list of the names of the registered factories.
func Factories() []string {
	factorysMu.RLock()
	defer factorysMu.RUnlock()
	var list []string
	for name := range factories {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

//NewDriver create new dirver with given driver name,config and prefix.
//Return driver created and any error if raised.
func NewDriver(name string, conf Config, prefix string) (Driver, error) {
	factorysMu.RLock()
	factoryi, ok := factories[name]
	factorysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("cache: unknown driver %q (forgotten import?)", name)
	}
	return factoryi(conf, prefix)
}

//NewSubCache  create subcache with given config and prefix.
//Return cache created and any error if raised.
func NewSubCache(conf Config, prefix string) (*Cache, error) {
	var err error
	c := New()
	var TTL int64
	var DriverName string
	var d Driver
	err = conf.Get(prefix+"TTL", &TTL)
	if err != nil {
		return nil, err
	}
	err = conf.Get(prefix+"Driver", &DriverName)
	if err != nil {
		return nil, err
	}
	d, err = NewDriver(DriverName, conf, prefix+"Config.")
	if err != nil {
		return nil, err
	}
	c.Driver = d
	return c, nil
}

//MustNewDriver  create new dirver with given driver name,config and prefix.
//Return driver created.
//Painc is any error raised.
func MustNewDriver(name string, conf Config, prefix string) Driver {
	d, err := NewDriver(name, conf, prefix)
	if err != nil {
		panic(err)
	}
	return d
}

//New :Create a empty cache.
func New() *Cache {
	return &Cache{}
}

//Cache Cache stores the cache Driver and default ttl.
type Cache struct {
	Driver
	TTL   time.Duration
	locks sync.Map
}

func (c *Cache) getKey(key string) string {
	return KeyPrefix + key
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

//ExpireCounter set cache counter  expire duration by given key and ttl
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

//Load Get data model from cache by given key.If data not found,call loader to get current data value and save to cache.
//If ttl is DefualtTTL(0),use default ttl in config instead.
//Return any error raised.
func (c *Cache) Load(key string, v interface{}, ttl time.Duration, loader Loader) error {
	if key == "" {
		return ErrKeyUnavailable
	}
	l, _ := c.locks.Load(key)
	if l != nil {
		lock := l.(*sync.RWMutex)
		lock.RLock()
		lock.RUnlock()
	}
	err := c.Get(key, v)
	if err == ErrNotFound {
		l, _ := c.locks.Load(key)
		if l != nil {
			lock := l.(*sync.RWMutex)
			lock.RLock()
			lock.RUnlock()
			return c.Load(key, v, ttl, loader)
		}
		lock := &sync.RWMutex{}
		c.locks.Store(key, lock)
		lock.Lock()
		defer func() {
			lock.Unlock()
			c.locks.Delete(key)
		}()
		v2, err2 := loader(key)
		if err2 != nil {
			return err2
		}
		reflect.Indirect(reflect.ValueOf(v)).Set(reflect.Indirect(reflect.ValueOf(v2)))
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

//FinalKey get final key which passed to cache driver .
func (c *Cache) FinalKey(key string) (string, error) {
	return c.getKey(key), nil
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

//Collection get a cache colletion with given prefix
func (c *Cache) Collection(prefix string) *Collection {
	return NewCollection(c, prefix, c.TTL)
}

//Node get a cache node with given prefix
func (c *Cache) Node(prefix string) *Node {
	return NewNode(c, prefix)
}

//Loader cache value loader used in cache load method.
//Load value with given key.
//Return loaded value and any error if raised.
type Loader func(key string) (interface{}, error)
