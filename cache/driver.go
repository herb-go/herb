package cache

import (
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

//Factory create driver with given config and prefix
//Reutrn driver created and any error if raised..
type Factory func(conf Config, prefix string) (Driver, error)

//Driver : Cache driver interface.Should Never used directly
type Driver interface {
	Util() *Util
	SetUtil(*Util)
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
		panic(errors.New("cache: Register cache factory is nil"))
	}
	if _, dup := factories[name]; dup {
		panic(errors.New("cache: Register called twice for factory " + name))
	}
	factories[name] = f
}

//UnregisterAll all factorys
func UnregisterAll() {
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
	c.TTL = time.Duration(TTL) * time.Second
	err = conf.Get(prefix+"Driver", &DriverName)
	if err != nil {
		return nil, err
	}
	d, err = NewDriver(DriverName, conf, prefix+"Config.")
	if err != nil {
		return nil, err
	}
	var mname = ""
	if mname == "" {
		mname = DefaultMarshaler
	}
	err = conf.Get(prefix+"Marshaler", &mname)
	if err != nil {
		return nil, err
	}
	marshaler, err := NewMarshaler(mname)
	if err != nil {
		return nil, err
	}
	u := NewUtil()
	u.Marshaler = marshaler
	d.SetUtil(u)

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
