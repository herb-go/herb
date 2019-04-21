package store

import (
	"fmt"
	"io"
	"sort"
	"sync"
)

//File file stored.
type File struct {
	//Store store which field stored in
	Store Store
	//ID file id
	ID  string
	url string
}

//SetURL set file url.
func (f *File) SetURL(url string) {
	f.url = url
}

//URL return file url and any error if raised.
func (f *File) URL() (url string, err error) {
	if f.url != "" {
		return f.url, nil
	}
	return f.Store.URL(f.ID)
}

//Driver store driver interface.
type Driver interface {
	//Save save data form reader to named file.
	//Return file id ,file size and any error if raised.
	Save(filename string, reader io.Reader) (id string, length int64, err error)
	//Load load file with given id.
	//Return file reader any error if raised.
	Load(id string, writer io.Writer) (io.ReadCloser, error)
	//Remove remove file by id.
	//Return any error if raised.
	Remove(id string) error
	//URL convert file id to file url.
	//Return file url and any error if raised.
	URL(id string) (string, error)
}

//Store file store.
type Store struct {
	Driver
}

//NewStore create new file store.
func New() *Store {
	return &Store{}
}

//Init applu option to store.
func (s *Store) Init(option Option) error {
	return option.ApplyTo(s)
}

// Factory store driver create factory.
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

//UnregisterAll unregister all driver
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

//NewDriver create new driver with given name,config and prefix.
//Reutrn driver created and any error if raised.
func NewDriver(name string, conf Config, prefix string) (Driver, error) {
	factorysMu.RLock()
	factoryi, ok := factories[name]
	factorysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("file: unknown driver %q (forgotten import?)", name)
	}
	return factoryi(conf, prefix)
}
