package router

import (
	"net/http"

	"github.com/herb-go/herb/middleware"
)

type Proxy struct {
	Router
	App *middleware.App
}

func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.App.ServeMiddleware(w, r, p.Router.ServeHTTP)
}

type Factory struct {
	creator func() Router
	app     *middleware.App
}

func (f *Factory) CreateRouter() Router {
	return &Proxy{
		Router: f.creator(),
		App:    f.app,
	}
}
func (f *Factory) Middlewares() *middleware.App {
	return f.app
}
func NewFactory(creator func() Router) *Factory {
	return &Factory{
		creator: creator,
		app:     middleware.New(),
	}
}
