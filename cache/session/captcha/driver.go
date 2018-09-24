package captcha

import (
	"fmt"
	"net/http"
	"sort"

	"github.com/herb-go/herb/cache"
)

//Driver captca driver interface.
type Driver interface {
	//Name return driver name.
	Name() string
	//MustCaptcha execute captcha to giver http request and response and scene or reset value.
	//Panic if any error rasied.
	MustCaptcha(scene string, reset bool, w http.ResponseWriter, r *http.Request)
	//Verify verify if token is validated with given http rquest and scene.
	//return verify result and any error raised.
	Verify(r *http.Request, scene string, token string) (bool, error)
}

//Factory driver createor with given config and prefix.
//Return driver and any error raised.
type Factory func(conf cache.Config, prefix string) (Driver, error)

// Register makes a driver creator available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, f Factory) {
	factorysMu.Lock()
	defer factorysMu.Unlock()
	if f == nil {
		panic("captcha: Register captcha factory is nil")
	}
	if _, dup := factories[name]; dup {
		panic("captcha: Register called twice for factory " + name)
	}
	factories[name] = f
}

//unregisterAll unregister all drivers.
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

//NewDriver create new driver with given name,config and prefix.
//Return driver created and any error if raised.
func NewDriver(name string, conf cache.Config, prefix string) (Driver, error) {
	factorysMu.RLock()
	factoryi, ok := factories[name]
	factorysMu.RUnlock()
	if !ok {
		return nil, fmt.Errorf("captcha: unknown driver %q (forgotten import?)", name)
	}
	return factoryi(conf, prefix)
}
