package middlewarefactory

import (
	"fmt"
	"sync"

	"github.com/herb-go/herb/middleware"
)

type Context struct {
	locker              sync.Mutex
	Conditionfactories  map[string]ConditionFactory
	Middlewarefactories map[string]Factory
}

func (ctx *Context) CreateMiddleware(name string, loader func(interface{}) error) (middleware.Middleware, error) {
	ctx.locker.Lock()
	defer ctx.locker.Unlock()
	f := ctx.Middlewarefactories[name]
	if f == nil {
		return nil, fmt.Errorf("middleware factory:%s %w", name, ErrFactoryNotRegistered)
	}
	return f(loader)
}
func (ctx *Context) CreateCondition(name string, loader func(interface{}) error) (Condition, error) {
	ctx.locker.Lock()
	defer ctx.locker.Unlock()
	f := ctx.Conditionfactories[name]
	if f == nil {
		return nil, fmt.Errorf("middleware factory:%s %w", name, ErrConditionFactoryNotRegistered)
	}
	return f(loader)
}
func (ctx *Context) RegisterFactory(name string, f Factory) {
	ctx.locker.Lock()
	defer ctx.locker.Unlock()
	ctx.Middlewarefactories[name] = f
}

func (ctx *Context) RegisterConditionFactory(name string, f ConditionFactory) {
	ctx.locker.Lock()
	defer ctx.locker.Unlock()
	ctx.Conditionfactories[name] = f
}
func NewContext() *Context {
	return &Context{
		Conditionfactories:  map[string]ConditionFactory{},
		Middlewarefactories: map[string]Factory{},
	}
}

var DefaultContext = NewContext()
