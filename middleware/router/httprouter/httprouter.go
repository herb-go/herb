package httprouter

import (
	"net/http"

	"github.com/herb-go/herb/middleware"
	"github.com/herb-go/herb/middleware/router"
	"github.com/julienschmidt/httprouter"
)

//Router router main struct.
type Router struct {
	app    *middleware.App
	router *httprouter.Router
}

//New create new router.
func New() *Router {
	r := httprouter.New()
	r.RedirectTrailingSlash = false
	r.RedirectFixedPath = false
	router := Router{
		app:    middleware.New(),
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

//Handle return app which will response to given method and path.
func (r *Router) Handle(method, path string) *middleware.App {
	app := middleware.New()
	r.router.Handle(method, path, wrap(app))
	return app
}

//Middlewares return router middlewares.
func (r *Router) Middlewares() *middleware.App {
	return r.app
}

//ServeHTTP serve router as http.handler.
func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.app.ServeMiddleware(w, req, r.router.ServeHTTP)
}

//GetParams get router params from request.
func GetParams(r *http.Request) *router.Params {
	return router.GetParams(r)
}

//SetParams set router to request.
func SetParams(r *http.Request, params httprouter.Params) {
	p := router.GetParams(r)
	for k := range params {
		p.Set(params[k].Key, params[k].Value)
	}
}

//GET return app which will response to GET method and path.
func (r *Router) GET(path string) *middleware.App {
	return r.Handle("GET", path)
}

//HEAD return app which will response to HEAD method and path.
func (r *Router) HEAD(path string) *middleware.App {
	return r.Handle("HEAD", path)
}

//OPTIONS return app which will response to HEAD method and path.
//Request called to path which any handle by OPTIONS method will return 404 instead of 405 error due to httprouter.
func (r *Router) OPTIONS(path string) *middleware.App {
	return r.Handle("OPTIONS", path)
}

//POST return app which will response to POST method and path.
func (r *Router) POST(path string) *middleware.App {
	return r.Handle("POST", path)
}

//PUT return app which will response to PUT method and path.
func (r *Router) PUT(path string) *middleware.App {
	return r.Handle("PUT", path)
}

//PATCH return app which will response to PATCH method and path.
func (r *Router) PATCH(path string) *middleware.App {
	return r.Handle("PATCH", path)
}

//DELETE return app which will response to DELETE method and path.
func (r *Router) DELETE(path string) *middleware.App {
	return r.Handle("DELETE", path)
}

//ALL return app which will response to all method and path.
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

//StripPrefix strip request prefix and server as a middleware app
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
