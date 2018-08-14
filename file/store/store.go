package store

import (
	"fmt"
	"io"
	"sort"
	"sync"
)

type File struct {
	Store Store
	ID    string
	url   string
}

func (f *File) SetURL(url string) {
	f.url = url
}
func (f *File) URL() (url string, err error) {
	if f.url != "" {
		return f.url, nil
	}
	return f.Store.URL(f.ID)
}

type Driver interface {
	Save(filename string, reader io.Reader) (id string, length int64, err error)
	Load(id string, writer io.Writer) error
	Remove(id string) error
	URL(id string) (string, error)
}

type Store struct {
	Driver
}

func (s *Store) Init(option Option) error {
	return option.ApplyTo(s)
}

type Factory func(conf Config, prefix string) (Driver, error)

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
		panic("file: Register cache factory is nil")
	}
	if _, dup := factories[name]; dup {
		panic("file: Register called twice for factory " + name)
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

func NewDriver(name string, conf Config, prefix string) (Driver, error) {
	factorysMu.RLock()
	factoryi, ok := factories[name]
	factorysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("file: unknown driver %q (forgotten import?)", name)
	}
	return factoryi(conf, prefix)
}

func MustNewDriver(name string, conf Config, prefix string) Driver {
	d, err := NewDriver(name, conf, prefix)
	if err != nil {
		panic(err)
	}
	return d
}

//New :Create a empty cache.
func NewStore() *Store {
	return &Store{}
}
