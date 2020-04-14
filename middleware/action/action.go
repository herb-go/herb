package action

import (
	"net/http"
	"sync/atomic"
)

//Action action struct
type Action struct {
	visit      *int64
	middleware func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	handler    func(w http.ResponseWriter, r *http.Request)
}

//ServeHTTP serve as http server
func (a *Action) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	atomic.AddInt64(a.visit, 1)
	if a.middleware == nil {
		a.handler(w, r)
		return
	}
	a.middleware(w, r, a.handler)
}

//SetMiddleware set action middleware
func (a *Action) SetMiddleware(f func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) {
	a.middleware = f
}

//WithHandler set action handler
func (a *Action) WithHandler(f func(w http.ResponseWriter, r *http.Request)) *Action {
	a.handler = f
	return a
}

//Handler return action handler
func (a *Action) Handler() func(w http.ResponseWriter, r *http.Request) {
	return a.handler
}

//Visit action visit count
func (a *Action) Visit() int64 {
	return atomic.LoadInt64(a.visit)
}

//New create new cation
func New() *Action {
	return &Action{}
}
