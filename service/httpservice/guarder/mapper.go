package guarder

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"

	"github.com/herb-go/herb/service"
)

type Mapper interface {
	ReadParamsFromRequest(*http.Request) (*Params, error)
	WriteParamsToRequest(*http.Request, *Params) error
}

//MapperFactory guarder factory
type MapperFactory func(conf service.Config, prefix string) (Mapper, error)

var (
	mapperFactorysMu sync.RWMutex
	mapperFactories  = make(map[string]MapperFactory)
)

// RegisterMapper makes a driver creator available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func RegisterMapper(name string, f MapperFactory) {
	mapperFactorysMu.Lock()
	defer mapperFactorysMu.Unlock()
	if f == nil {
		panic(errors.New("guarder: Register mapper factory is nil"))
	}
	if _, dup := mapperFactories[name]; dup {
		panic(errors.New("guarder: Register called twice for mapper factory " + name))
	}
	mapperFactories[name] = f
}

//UnregisterAllMapper unregister all driver
func UnregisterAllMapper() {
	mapperFactorysMu.Lock()
	defer mapperFactorysMu.Unlock()
	// For tests.
	mapperFactories = make(map[string]MapperFactory)
}

//MapperFactories returns a sorted list of the names of the registered factories.
func MapperFactories() []string {
	mapperFactorysMu.RLock()
	defer mapperFactorysMu.RUnlock()
	var list []string
	for name := range mapperFactories {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

//NewMapperDriver create new driver with given name,config and prefix.
//Reutrn driver created and any error if raised.
func NewMapperDriver(name string, conf service.Config, prefix string) (Mapper, error) {
	mapperFactorysMu.RLock()
	factoryi, ok := mapperFactories[name]
	mapperFactorysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("guarder: unknown mapper driver %q (forgotten import?)", name)
	}
	return factoryi(conf, prefix)
}
