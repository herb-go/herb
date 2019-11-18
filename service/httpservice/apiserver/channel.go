package apiserver

import (
	"net/http"

	"github.com/herb-go/herb/service/httpservice"
)

type Server struct {
	httpservice.Config
	Name string
}

type Channel struct {
	Server
	Middlewares MiddlewareList
	Channel     string
}

func (c *Channel) server() *apiServer {
	return apiserver(c.Server.Name)
}

func (c *Channel) ApplyServer() error {
	if c.Server.IsEmpty() {
		return nil
	}
	return apiserver(c.Server.Name).SetConfig(&c.Server.Config)
}

func (c *Channel) Start(handler func(w http.ResponseWriter, r *http.Request)) error {
	var h func(w http.ResponseWriter, r *http.Request)
	if c.Middlewares == nil {
		h = handler
	} else {
		h = func(w http.ResponseWriter, r *http.Request) {
			c.Middlewares.ServeMiddleware(w, r, handler)
		}
	}
	return apiserver(c.Server.Name).Start(c.Channel, h)
}

func (c *Channel) Stop() error {
	return apiserver(c.Server.Name).Stop(c.Channel)
}
