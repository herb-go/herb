package requestinforeader

import (
	"fmt"
	"sync"
)

var lock sync.Mutex
var regiteredFactories = map[string]Factory{}
var defaultFactory func(string) (Factory, error)
var registered = map[string]Reader{}

func SetDefaultFactory(f func(string) (Factory, error)) {
	lock.Lock()
	defer lock.Unlock()
	defaultFactory = f
}

func RegisterFactories(name string, f Factory) {
	lock.Lock()
	defer lock.Unlock()
	regiteredFactories[name] = f
}
func Reset() {
	lock.Lock()
	defer lock.Unlock()
	registered = map[string]Reader{}
}

func Register(name string, reader Reader) {
	lock.Lock()
	defer lock.Unlock()
	registered[name] = reader
}
func GetReader(name string) (Reader, error) {
	lock.Lock()
	defer lock.Unlock()
	r := registered[name]
	if r == nil {
		return nil, fmt.Errorf("%w (%s)", ErrReaderNotFound, name)
	}
	return r, nil
}

func GetFactory(name string) (Factory, error) {

	f, ok := regiteredFactories[name]
	if ok {
		return f, nil
	}
	return defaultFactory(name)
}
