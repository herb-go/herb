package httprouter

import (
	"errors"
	"net/http"

	"context"

	"github.com/herb-go/herb/middleware"
	"github.com/julienschmidt/httprouter"
)

type keyType string

const ContextParamsKey = keyType("httprouterParams")

type Router struct {
	handlers []func(http.ResponseWriter, *http.Request, http.HandlerFunc)
	router   *httprouter.Router
}

func New(funcs ...func(http.ResponseWriter, *http.Request, http.HandlerFunc)) *Router {
	r := httprouter.New()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	router := Router{
		router: r,
	}
	router.handlers = funcs
	return &router
}
func wrap(f http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		f.ServeHTTP(w, SetParams(r, params))
	}
}

func (r *Router) Handle(method, path string) *middleware.App {
	app := middleware.New(r.handlers...)
	r.router.Handle(method, path, wrap(app))
	return app
}
func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.router.ServeHTTP(w, r)
}
func (router *Router) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	router.router.ServeHTTP(w, r)
}
func (router *Router) Handlers() []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){router.ServeMiddleware}
}

func (r *Router) Use(funcs ...func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) *Router {
	r.handlers = append(r.handlers, funcs...)
	return r
}
func (r *Router) Chain(apps ...middleware.HandlerSlice) *Router {
	r.Use(middleware.AppToMiddlewares(apps...)...)
	return r
}
func GetParams(r *http.Request) *httprouter.Params {
	p := r.Context().Value(ContextParamsKey)
	params, _ := p.(httprouter.Params)
	return &params
}
func SetParams(r *http.Request, params httprouter.Params) *http.Request {
	ctx := context.WithValue(r.Context(), ContextParamsKey, params)
	r = r.WithContext(ctx)
	return r
}

func (r *Router) GET(path string) *middleware.App {
	return r.Handle("GET", path)
}

func (r *Router) HEAD(path string) *middleware.App {
	return r.Handle("HEAD", path)
}

func (r *Router) OPTIONS(path string) *middleware.App {
	return r.Handle("OPTIONS", path)
}

func (r *Router) POST(path string) *middleware.App {
	return r.Handle("POST", path)
}

func (r *Router) PUT(path string) *middleware.App {
	return r.Handle("PUT", path)
}

func (r *Router) PATCH(path string) *middleware.App {
	return r.Handle("PATCH", path)
}

func (r *Router) DELETE(path string) *middleware.App {
	return r.Handle("DELETE", path)
}
func (r *Router) ALL(path string) *middleware.App {
	app := middleware.New(r.handlers...)
	handler := wrap(app)
	r.router.GET(path, handler)
	r.router.POST(path, handler)
	r.router.PUT(path, handler)
	r.router.DELETE(path, handler)
	r.router.PATCH(path, handler)
	r.router.OPTIONS(path, handler)
	r.router.HEAD(path, handler)
	return app
}
func (r *Router) StripPrefix(path string) *middleware.App {
	app := middleware.New(stripPrefixfunc)
	app.Use(r.handlers...)
	handler := wrap(app)
	p := path + "/*filepath"
	r.router.GET(p, handler)
	r.router.POST(p, handler)
	r.router.PUT(p, handler)
	r.router.DELETE(p, handler)
	r.router.PATCH(p, handler)
	r.router.OPTIONS(p, handler)
	r.router.HEAD(p, handler)
	return app
}

func stripPrefixfunc(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	params := r.Context().Value(ContextParamsKey)
	p, ok := params.(httprouter.Params)
	if ok {
		r.URL.Path = p.ByName("filepath")
		next(w, r)
		return
	}
	panic(errors.New("Strip prefix Error"))
}
