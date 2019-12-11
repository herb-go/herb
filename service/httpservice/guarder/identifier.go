package guarder

import (
	"errors"
	"fmt"
	"sort"
	"sync"
)

type Identifier interface {
	IdentifyParams(p *Params) (string, error)
}

//IdentifierFactory guarder factory
type IdentifierFactory func(loader func(interface{}) error) (Identifier, error)

var (
	identifierFactorysMu sync.RWMutex
	identifierFactories  = make(map[string]IdentifierFactory)
)

// RegisterIdentifier makes a driver creator available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func RegisterIdentifier(name string, f IdentifierFactory) {
	identifierFactorysMu.Lock()
	defer identifierFactorysMu.Unlock()
	if f == nil {
		panic(errors.New("guarder: Register identifier factory is nil"))
	}
	if _, dup := identifierFactories[name]; dup {
		panic(errors.New("guarder: Register called twice for identifier factory " + name))
	}
	identifierFactories[name] = f
}

//UnregisterAllIdentifier unregister all driver
func UnregisterAllIdentifier() {
	identifierFactorysMu.Lock()
	defer identifierFactorysMu.Unlock()
	// For tests.
	identifierFactories = make(map[string]IdentifierFactory)
}

//IdentifierFactories returns a sorted list of the names of the registered factories.
func IdentifierFactories() []string {
	identifierFactorysMu.RLock()
	defer identifierFactorysMu.RUnlock()
	var list []string
	for name := range identifierFactories {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

//NewIdentifierDriver create new driver with given name and loader.
//Reutrn driver created and any error if raised.
func NewIdentifierDriver(name string, loader func(interface{}) error) (Identifier, error) {
	identifierFactorysMu.RLock()
	factoryi, ok := identifierFactories[name]
	identifierFactorysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("guarder: unknown identifier driver %q (forgotten import?)", name)
	}
	return factoryi(loader)
}
