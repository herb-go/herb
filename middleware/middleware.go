//Package middleware provide a app interface to use middleware easily.
//All middleware is func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc).
package middleware

import "net/http"

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
		panic("handlers cant be nil")
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
		s.handlers[current](w, r, s.Next)
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
func (a *App) HandleFunc(HandlerFunc http.HandlerFunc) {
	Use(a, WrapFunc(HandlerFunc))
}

// Handle : Use http Handler as last middleware.
func (a *App) Handle(Handler http.Handler) {
	Use(a, Wrap(Handler))
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
