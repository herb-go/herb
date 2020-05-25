package middlewarefactory

import (
	"fmt"
	"sync"

	"github.com/herb-go/herb/middleware"
)

type Context interface {
	ConditionCreator
	MiddlewareCreator
}

type PlainContext struct {
	locker              sync.Mutex
	Conditionfactories  map[string]ConditionFactory
	Middlewarefactories map[string]Factory
}

func (ctx *PlainContext) CreateMiddleware(name string, loader func(interface{}) error) (middleware.Middleware, error) {
	ctx.locker.Lock()
	defer ctx.locker.Unlock()
	f := ctx.Middlewarefactories[name]
	if f == nil {
		return nil, fmt.Errorf("middleware factory:%s %w", name, ErrFactoryNotRegistered)
	}
	return f(loader)
}
func (ctx *PlainContext) CreateCondition(name string, loader func(interface{}) error) (Condition, error) {
	ctx.locker.Lock()
	defer ctx.locker.Unlock()
	f := ctx.Conditionfactories[name]
	if f == nil {
		return nil, fmt.Errorf("middleware factory:%s %w", name, ErrConditionFactoryNotRegistered)
	}
	return f(loader)
}
func (ctx *PlainContext) MustRegisterFactory(name string, f Factory) {
	ctx.locker.Lock()
	defer ctx.locker.Unlock()
	_, ok := ctx.Middlewarefactories[name]
	if ok {
		panic(fmt.Errorf("middleware factory:%s %w", name, ErrFactoryRegistered))
	}
	ctx.Middlewarefactories[name] = f
}

func (ctx *PlainContext) MustRegisterConditionFactory(name string, f ConditionFactory) {
	ctx.locker.Lock()
	defer ctx.locker.Unlock()
	_, ok := ctx.Conditionfactories[name]
	if ok {
		panic(fmt.Errorf("middleware factory:%s %w", name, ErrConditionFactoryRegistered))
	}
	ctx.Conditionfactories[name] = f
}
func NewPlainContext() *PlainContext {
	return &PlainContext{
		Conditionfactories:  map[string]ConditionFactory{},
		Middlewarefactories: map[string]Factory{},
	}
}

var DefaultContext = NewPlainContext()

var MustRegisterFactory = func(name string, f Factory) {
	DefaultContext.MustRegisterFactory(name, f)
}
var MustRegisterConditionFactory = func(name string, f ConditionFactory) {
	DefaultContext.MustRegisterConditionFactory(name, f)
}

func RegisterBuildin() {
	MustRegisterFactory("respone", NewResponseFactory())
	MustRegisterConditionFactory("time", NewTimeConditionFactory())
}
