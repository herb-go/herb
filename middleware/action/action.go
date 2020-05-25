package action

import (
	"net/http"

	"github.com/herb-go/herb/middleware"
)

//Action action struct
type Action struct {
	app     *middleware.App
	handler func(w http.ResponseWriter, r *http.Request)
}

//ServeHTTP serve as http server
func (a *Action) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.app.ServeMiddleware(w, r, a.handler)
}

//Middlewares return middlewares of aciton
func (a *Action) Middlewares() *middleware.App {
	return a.app
}

//New create new cation
func New(f func(w http.ResponseWriter, r *http.Request)) *Action {
	return &Action{
		app:     middleware.New(),
		handler: f,
	}
}
