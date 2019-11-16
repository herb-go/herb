//Package middleware provide a app interface to use middleware easily.
//All middleware is func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc).
package middleware

import (
	"errors"
	"net/http"
)

// New : Create new chainable middleware app.
func New(funcs ...func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) *App {
	app := new(App)
	app.handlers = funcs
	return app
}

//ServeHTTP : Server app as http.
func ServeHTTP(app HandlerSlice, w http.ResponseWriter, r *http.Request) {
	ServeMiddleware(app, w, r, voidNextFunc)
}

// ServeMiddleware : Server  app as middleware.
func ServeMiddleware(app HandlerSlice, w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	handlers := app.Handlers()
	if len(handlers) == 0 {
		panic(errors.New("handlers can't be nil"))
	}
	s := serveWorker{
		handlers: handlers,
		final:    next,
		current:  0,
	}
	s.Next(w, r)
}

// serveWorker : Runner of middlewares.
type serveWorker struct {
	handlers []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	final    http.HandlerFunc
	current  int
}

// Next : Current middlewares running step.
func (s *serveWorker) Next(w http.ResponseWriter, r *http.Request) {
	if s.current == len(s.handlers) {
		s.final(w, r)
	} else {
		current := s.current
		s.current = s.current + 1
		handler := s.handlers[current]
		if handler == nil {
			handler = voidMiddleware
		}
		handler(w, r, s.Next)
	}
}

//Chain : Append All middlewares in src to dst.
//Dst will not change when new middleware appended to src.
func Chain(dst HandlerChain, src ...HandlerSlice) {
	if len(src) == 0 {
		return
	}
	funcs := dst.Handlers()
	funcs = append(funcs, AppToMiddlewares(src...)...)
	f := make([]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc), len(funcs))
	copy(f, funcs)
	dst.SetHandlers(f)
}

//AppToMiddlewares : Convert App to slice of middleware.
func AppToMiddlewares(app ...HandlerSlice) []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	funcs := []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){}
	for _, v := range app {
		funcs = append(funcs, v.Handlers()...)
	}
	return funcs
}

//Use : Append middlewares tp app.
func Use(app HandlerChain, middlewares ...func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {
	funcs := append(app.Handlers(), middlewares...)
	f := make([]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc), len(funcs))
	copy(f, funcs)
	app.SetHandlers(f)
}

// WrapFunc : Wrap http.HandlerFunc [func(w http.ResponseWriter, r *http.Request)] to middleware.
// Next will be called after handlerFunc finish .
func WrapFunc(handlerFunc http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		handlerFunc(w, r)
		next(w, r)
	}
}

// Wrap : Wrap http.Handler to middleware.
// Next will be called after handlerFunc finish .
func Wrap(f http.Handler) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		f.ServeHTTP(w, r)
		next(w, r)
	}
}

func voidNextFunc(w http.ResponseWriter, r *http.Request) {
}
func voidMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(w, r)
}

// HandlerSlice : interface which contains slice of middlewares.
type HandlerSlice interface {
	Handlers() []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// HandlerChain : interface which contains slice of middlewares which can be updated.
type HandlerChain interface {
	HandlerSlice
	SetHandlers([]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc))
}

// App : Slice of middlewares with tons of helpful method.
type App struct {
	handlers []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
}

// Handlers : Return all middlewares in app.
func (a *App) Handlers() []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return a.handlers
}

// SetHandlers : Set app's middlewares.
func (a *App) SetHandlers(h []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {
	a.handlers = h
}

// ServeMiddleware : Use app as a middleware.
func (a *App) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	ServeMiddleware(a, w, r, next)
}

// ServeHTTP : Use app as a http handler.
func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ServeHTTP(a, w, r)
}

// HandleFunc : Use http HandlerFunc as last middleware.
func (a *App) HandleFunc(HandlerFunc http.HandlerFunc) *App {
	Use(a, WrapFunc(HandlerFunc))
	return a
}

// Handle : Use http Handler as last middleware.
func (a *App) Handle(Handler http.Handler) *App {
	Use(a, Wrap(Handler))
	return a
}

// Use : Append middlewares to app.
func (a *App) Use(middlewares ...func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) *App {
	Use(a, middlewares...)
	return a
}

// Chain : Append All middlewares in src to app.
//App will not change when new middleware appended to src.
func (a *App) Chain(src ...HandlerSlice) *App {
	Chain(a, src...)
	return a
}

// UseApp : Append apps as middlewares to app.
// App will change when new middleware appended to apps.
func (a *App) UseApp(apps ...*App) *App {
	funcs := make([]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc), len(apps))
	for k := range apps {
		funcs[k] = apps[k].ServeMiddleware
	}
	Use(a, funcs...)
	return a
}

//Middleware middleware interface.
type Middleware func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

// Middlewares middleware list interaface.
type Middlewares []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

// Handlers : Return all middlewares in app.
func (m *Middlewares) Handlers() []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return *m
}

// SetHandlers : Set app's middlewares.
func (m *Middlewares) SetHandlers(h []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {
	ms := make([]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc), len(h))
	copy(h, ms)
	*m = ms
}

// Use : Append middlewares
func (m *Middlewares) Use(middlewares ...func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) *Middlewares {
	Use(m, middlewares...)
	return m

}

//App create app with middlewares
func (m *Middlewares) App(handler func(w http.ResponseWriter, r *http.Request)) *App {
	app := New(m.Handlers()...)
	app.HandleFunc(handler)
	return app
}

//NewMiddlewares create new middlewares with given handlers
func NewMiddlewares(middlewares ...func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) *Middlewares {
	ms := make([]func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc), len(middlewares))
	copy(middlewares, ms)
	m := Middlewares(ms)
	return &m
}
