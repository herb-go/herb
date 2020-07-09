package requestinforeader

import "net/http"

type Reader func(r *http.Request) ([]byte, error)

type Factory interface {
	CreateReader(loader func(v interface{}) error) (Reader, error)
}

type FactoryFunc func(loader func(v interface{}) error) (Reader, error)

func (f FactoryFunc) CreateReader(loader func(v interface{}) error) (Reader, error) {
	return f(loader)
}
