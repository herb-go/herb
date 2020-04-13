package action

import "net/http"

//Action action struct
type Action struct {
	middleware func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	handler    func(w http.ResponseWriter, r *http.Request)
}

//ServeHTTP serve as http server
func (a *Action) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if a.middleware == nil {
		a.handler(w, r)
		return
	}
	a.middleware(w, r, a.handler)
}

//SetHandler set action handler
func (a *Action) SetHandler(f func(w http.ResponseWriter, r *http.Request)) *Action {
	a.handler = f
	return a
}

//Handler return action handler
func (a *Action) Handler() func(w http.ResponseWriter, r *http.Request) {
	return a.handler
}

//New create new cation
func New() *Action {
	return &Action{}
}
