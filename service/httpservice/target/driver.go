package target

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync"
)

type DoerFactory func(loader func(v interface{}) error) (Doer, error)

func DefaultClientDoerFactory(loader func(v interface{}) error) (Doer, error) {
	return http.DefaultClient, nil
}

var (
	doerFactorysMu sync.RWMutex
	doerFactories  = make(map[string]DoerFactory)
)

// RegisterDoer makes a driver creator available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func RegisterDoer(name string, f DoerFactory) {
	doerFactorysMu.Lock()
	defer doerFactorysMu.Unlock()
	if f == nil {
		panic(errors.New("target: Register doer factory is nil"))
	}
	if _, dup := doerFactories[name]; dup {
		panic(errors.New("target: Register called twice for factory " + name))
	}
	doerFactories[name] = f
}

//UnregisterAllDoer unregister all driver
func UnregisterAllDoer() {
	doerFactorysMu.Lock()
	defer doerFactorysMu.Unlock()
	// For tests.
	doerFactories = make(map[string]DoerFactory)
}

//DoerFactories returns a sorted list of the names of the registered factories.
func DoerFactories() []string {
	doerFactorysMu.RLock()
	defer doerFactorysMu.RUnlock()
	var list []string
	for name := range doerFactories {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

//NewDoer create new doer with given name,loader.
//Return driver created and any error if raised.
func NewDoer(name string, loader func(interface{}) error) (Doer, error) {
	doerFactorysMu.RLock()
	factoryi, ok := doerFactories[name]
	doerFactorysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("target: unknown driver %q (forgotten import?)", name)
	}
	return factoryi(loader)
}
