package cache

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"

	"gopkg.in/vmihailenco/msgpack.v2"
)

var ErrNotFound = errors.New("Entry not found")
var ErrNotCacheable = errors.New("Not Cachable")
var ErrEntryTooLarge = errors.New("Entry too large to cache")
var ErrKeyTooLarge = errors.New("Key too large to cache")
var ErrKeyUnavailable = errors.New("Key Unavailable")
var ErrSearchKeysNotSupported = errors.New("Search Keys by Prefix is not supported")
var DefualtTTL = time.Duration(0)
var TTLForever = time.Duration(-1)
var keyPrefix = string([]byte{0})
var intKeyPrefix = string([]byte{69, 0})

type Config struct {
	Driver string
	Config json.RawMessage
	TTL    int64
}

//Driver : Cache driver interface.Should Never not be used directly/
type Driver interface {
	New(cacheConfig json.RawMessage) (Driver, error)
	Set(key string, v interface{}, ttl time.Duration) error
	Update(key string, v interface{}, ttl time.Duration) error
	Get(key string, v interface{}) error
	SetBytesValue(key string, bytes []byte, ttl time.Duration) error
	UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error
	GetBytesValue(key string) ([]byte, error)
	Del(key string) error
	SearchByPrefix(prefix string) ([]string, error)
	IncrCounter(key string, increment int64, ttl time.Duration) (int64, error)
	SetCounter(key string, v int64, ttl time.Duration) error
	GetCounter(key string) (int64, error)
	DelCounter(key string) error
	SetGCErrHandler(f func(err error))
	Close() error
	Flush() error
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
func New() *Cache {
	return &Cache{}
}

type Cache struct {
	Driver
	TTL time.Duration
}

func (c *Cache) Open(driverName string, cacheConfig json.RawMessage, ttlInSecond int64) error {
	driversMu.RLock()
	driveri, ok := drivers[driverName]
	driversMu.RUnlock()
	if !ok {
		return fmt.Errorf("cache: unknown driver %q (forgotten import?)", driverName)
	}
	driver, err := driveri.New(cacheConfig)
	if err != nil {
		return err
	}
	c.Driver = driver
	c.TTL = time.Duration(ttlInSecond * int64(time.Second))
	return nil
}
func (c *Cache) OpenConfig(config Config) error {
	return c.Open(config.Driver, config.Config, config.TTL)
}
func (c *Cache) OpenJson(data []byte) error {
	var config Config
	err := json.Unmarshal(data, &config)
	if err != nil {
		return err
	}
	return c.OpenConfig(config)
}

func (c *Cache) getKey(key string) string {
	return keyPrefix + key
}
func (c *Cache) Set(key string, v interface{}, ttl time.Duration) error {
	if key == "" {
		return ErrNotFound
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
	}
	return c.Driver.Set(c.getKey(key), v, ttl)
}
func (c *Cache) Update(key string, v interface{}, ttl time.Duration) error {
	if key == "" {
		return ErrNotFound
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
	}
	return c.Driver.Update(c.getKey(key), v, ttl)
}
func (c *Cache) Get(key string, v interface{}) error {
	if key == "" {
		return ErrNotFound
	}
	return c.Driver.Get(c.getKey(key), &v)
}
func (c *Cache) SearchByPrefix(prefix string) ([]string, error) {

	return c.Driver.SearchByPrefix(prefix)
}
func (c *Cache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	if key == "" {
		return ErrNotFound
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
	}

	return c.Driver.SetBytesValue(c.getKey(key), bytes, ttl)
}
func (c *Cache) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {
	if key == "" {
		return ErrNotFound
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
	}

	return c.Driver.UpdateBytesValue(c.getKey(key), bytes, ttl)
}
func (c *Cache) GetBytesValue(key string) ([]byte, error) {
	if key == "" {
		return nil, ErrNotFound
	}
	return c.Driver.GetBytesValue(c.getKey(key))
}
func (c *Cache) Del(key string) error {
	if key == "" {
		return ErrNotFound
	}
	return c.Driver.Del(c.getKey(key))
}
func (c *Cache) getIntKey(key string) string {
	return intKeyPrefix + key
}
func (c *Cache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	if key == "" {
		return 0, ErrNotFound
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
	}
	return c.Driver.IncrCounter(c.getIntKey(key), increment, ttl)
}
func (c *Cache) SetCounter(key string, v int64, ttl time.Duration) error {
	if key == "" {
		return ErrNotFound
	}
	if ttl == DefualtTTL {
		ttl = c.TTL
	}
	return c.Driver.SetCounter(c.getIntKey(key), v, ttl)
}
func (c *Cache) SearchCounterByPrefix(prefix string) ([]string, error) {
	return c.Driver.SearchByPrefix(c.getIntKey(prefix))
}
func (c *Cache) GetCounter(key string) (int64, error) {
	if key == "" {
		return 0, ErrNotFound
	}
	return c.Driver.GetCounter(c.getIntKey(key))
}
func (c *Cache) DelCounter(key string) error {
	if key == "" {
		return ErrNotFound
	}
	return c.Driver.DelCounter(c.getIntKey(key))
}

func (c *Cache) Load(key string, v interface{}, ttl time.Duration, closure func(v interface{}) error) error {
	if key == "" {
		return ErrNotFound
	}
	err := c.Get(key, &v)
	if err == ErrNotFound {
		err2 := closure(&v)
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

func MarshalMsgpack(v interface{}) ([]byte, error) {
	return msgpack.Marshal(&v)
}

func UnmarshalMsgpack(bytes []byte, v interface{}) error {
	return msgpack.Unmarshal(bytes, &v)
}
