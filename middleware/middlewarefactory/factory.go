package middlewarefactory

import (
	"github.com/herb-go/herb/middleware"
)

//Factory middleware factory
type Factory interface {
	CreateMiddleware(name string, loader func(v interface{}) error) (middleware.Middleware, error)
}

//FactoryFunc factory func type
type FactoryFunc (func(name string, loader func(v interface{}) error) (middleware.Middleware, error))

//CreateMiddleware create new middleware with given name and loader
func (f FactoryFunc) CreateMiddleware(name string, loader func(v interface{}) error) (middleware.Middleware, error) {
	return f(name, loader)
}
