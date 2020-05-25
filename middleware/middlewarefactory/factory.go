package middlewarefactory

import (
	"github.com/herb-go/herb/middleware"
)

//Factory middleware factory
type Factory func(loader func(v interface{}) error) (middleware.Middleware, error)
