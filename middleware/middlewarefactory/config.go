package middlewarefactory

import (
	"net/http"

	"github.com/herb-go/herb/middleware"
)

type ConditionConfig struct {
	Type       string
	Config     func(v interface{}) error `config:", lazyload"`
	Not        bool
	Or         bool
	Disabled   bool
	Conditions []*ConditionConfig
}

func (c *ConditionConfig) Create(creator ConditionCreator) (Condition, error) {
	var err error
	pc := NewPlainCondition()
	if c.Type != "" {
		pc.Condition, err = creator.CreateCondition(c.Type, c.Config)
		if err != nil {
			return nil, err
		}
	}
	pc.Not = c.Not
	pc.Or = c.Or
	pc.Disabled = c.Disabled
	for k := range c.Conditions {
		condition, err := c.Conditions[k].Create(creator)
		if err != nil {
			return nil, err
		}
		pc.Conditions = append(pc.Conditions, condition)
	}
	return pc, nil
}

type ConditionCreator interface {
	CreateCondition(string, func(interface{}) error) (Condition, error)
}

type MiddlewareConfig struct {
	Type   string
	Config func(v interface{}) error `config:", lazyload"`
}

type MiddlewareCreator interface {
	CreateMiddleware(string, func(interface{}) error) (middleware.Middleware, error)
}

type Config struct {
	Condition   *ConditionConfig
	Middlewares []*MiddlewareConfig
}

func (c *Config) Middleware(ctx *Context) (middleware.Middleware, error) {
	condition, err := c.Condition.Create(ctx)
	if err != nil {
		return nil, err
	}
	middlewares := make(middleware.Middlewares, len(c.Middlewares))
	for k := range c.Middlewares {
		m, err := ctx.CreateMiddleware(c.Middlewares[k].Type, c.Middlewares[k].Config)
		if err != nil {
			return nil, err
		}
		middlewares[k] = m
	}
	m := func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result, err := condition.MatchRequest(r)
		if err != nil {
			panic(err)
		}
		if result {
			middleware.ServeMiddleware(&middlewares, w, r, next)
		} else {
			next(w, r)
		}
	}
	return m, nil
}

type ConfigList []*Config

func (c *ConfigList) Middleware(ctx *Context) (middleware.Middleware, error) {
	result := middleware.Middlewares{}
	for k := range *c {
		m, err := (*c)[k].Middleware(ctx)
		if err != nil {
			return nil, err
		}
		result = append(result, m)
	}
	return result.ServeMiddleware, nil
}
