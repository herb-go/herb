package cache

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

var dummoyLoader = func(v interface{}) error {
	return nil
}

//Factory create driver with given loader.
//Reutrn driver created and any error if raised..
type Factory func(loader func(v interface{}) error) (Driver, error)

//Driver : Cache driver interface.Should Never used directly
type Driver interface {
	MinimumOperation
	//Set callback to handler error raised when gc.
	SetGCErrHandler(f func(err error))
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

//NewDriver create new dirver with given driver name and loader.
//Return driver created and any error if raised.
func NewDriver(name string, loader func(v interface{}) error) (Driver, error) {
	if loader == nil {
		loader = dummoyLoader
	}
	factorysMu.RLock()
	factoryi, ok := factories[name]
	factorysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("cache: unknown driver %q (forgotten import?)", name)
	}
	return factoryi(loader)
}

//NewSubCache  create subcache with given loader.
//Return cache created and any error if raised.
func NewSubCache(conf *OptionConfig) (*Cache, error) {
	var err error
	c := New()
	err = conf.ApplyTo(c)
	if err != nil {
		return nil, err
	}
	return c, nil
}

//MustNewDriver  create new dirver with given driver name and loader.
//Return driver created.
//Painc is any error raised.
func MustNewDriver(name string, loader func(v interface{}) error) Driver {
	d, err := NewDriver(name, loader)
	if err != nil {
		panic(err)
	}
	return d
}
