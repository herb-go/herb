package httprouter

import (
	"net/http"

	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/middleware-router"
	"github.com/julienschmidt/httprouter"
)

type Router struct {
	router *httprouter.Router
}

func New() *Router {
	r := httprouter.New()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	router := Router{
		router: r,
	}
	return &router
}
func wrap(f http.Handler) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, params httprouter.Params) {
		SetParams(r, params)
		f.ServeHTTP(w, r)
	}
}

func (r *Router) Handle(method, path string) *middleware.App {
	app := middleware.New()
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

func GetParams(r *http.Request) *router.Params {
	return router.GetParams(r)
}
func SetParams(r *http.Request, params httprouter.Params) {
	p := router.GetParams(r)
	for k := range params {
		p.Set(params[k].Key, params[k].Value)
	}
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
	app := middleware.New()
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
	params := router.GetParams(r)
	r.URL.Path = params.Get("filepath")
	next(w, r)
}
