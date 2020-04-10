package action

import "net/http"

//Action action struct
type Action struct {
	handler func(w http.ResponseWriter, r *http.Request)
}

//ServeHTTP serve as http server
func (a *Action) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler(w, r)
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
