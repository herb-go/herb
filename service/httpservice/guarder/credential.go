package guarder

import (
	"errors"
	"fmt"
	"sort"
	"sync"

	"github.com/herb-go/herb/service"
)

type Credential interface {
	CredentialParams() (*Params, error)
}

//CredentialFactory guarder factory
type CredentialFactory func(conf service.Config, prefix string) (Credential, error)

var (
	credentialFactorysMu sync.RWMutex
	credentialFactories  = make(map[string]CredentialFactory)
)

// RegisterCredential makes a driver creator available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func RegisterCredential(name string, f CredentialFactory) {
	credentialFactorysMu.Lock()
	defer credentialFactorysMu.Unlock()
	if f == nil {
		panic(errors.New("guarder: Register credential factory is nil"))
	}
	if _, dup := credentialFactories[name]; dup {
		panic(errors.New("guarder: Register called twice for credential factory " + name))
	}
	credentialFactories[name] = f
}

//UnregisterAllCredential unregister all driver
func UnregisterAllCredential() {
	credentialFactorysMu.Lock()
	defer credentialFactorysMu.Unlock()
	// For tests.
	credentialFactories = make(map[string]CredentialFactory)
}

//CredentialFactories returns a sorted list of the names of the registered factories.
func CredentialFactories() []string {
	credentialFactorysMu.RLock()
	defer credentialFactorysMu.RUnlock()
	var list []string
	for name := range credentialFactories {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

//NewCredentialDriver create new driver with given name,config and prefix.
//Reutrn driver created and any error if raised.
func NewCredentialDriver(name string, conf service.Config, prefix string) (Credential, error) {
	credentialFactorysMu.RLock()
	factoryi, ok := credentialFactories[name]
	credentialFactorysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("guarder: unknown credential driver %q (forgotten import?)", name)
	}
	return factoryi(conf, prefix)
}
