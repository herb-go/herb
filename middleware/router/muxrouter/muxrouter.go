package muxrouter

import (
	"errors"
	"net/http"

	"github.com/herb-go/herb/middleware/router"

	"github.com/herb-go/herb/middleware"
)

type Router struct {
	mux             *http.ServeMux
	notfoundHanlder http.Handler
	homepage        *middleware.App
}

func (muxrouter *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		if muxrouter.homepage != nil {
			muxrouter.homepage.ServeHTTP(w, r)
			return
		}
	} else {
		h, p := muxrouter.mux.Handler(r)
		if p != "" {
			h.ServeHTTP(w, r)
			return
		}
	}
	muxrouter.NotFound(w, r)

}
func (muxrouter *Router) SetNotFoundHandler(h http.Handler) {
	muxrouter.notfoundHanlder = h
}
func (muxrouter *Router) HandleHomepage() *middleware.App {
	if muxrouter.homepage != nil {
		panic(errors.New("homepage handled"))
	}
	muxrouter.homepage = middleware.New()
	return muxrouter.homepage
}
func (muxrouter *Router) StripPrefix(prefix string) *middleware.App {
	if prefix == "" {
		panic(errors.New("empty prefix"))
	}
	if prefix[len(prefix)-1] == '/' {
		panic(errors.New("prefix ends with '/'"))
	}
	a := muxrouter.Handle(prefix + "/")
	a.Use(router.NewStripPrefixMiddleware(prefix))
	muxrouter.mux.Handle(prefix, a)
	return a
}
func (muxrouter *Router) Handle(pattern string) *middleware.App {
	a := middleware.New()
	muxrouter.mux.Handle(pattern, a)
	return a
}
func (muxrouter *Router) NotFound(w http.ResponseWriter, r *http.Request) {
	muxrouter.notfoundHanlder.ServeHTTP(w, r)
}
func New() *Router {
	notfound := http.HandlerFunc(http.NotFound)
	r := &Router{
		mux:             &http.ServeMux{},
		notfoundHanlder: notfound,
	}
	return r
}
