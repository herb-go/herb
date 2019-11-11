package apiserver

import (
	"net/http"

	"github.com/herb-go/herb/server"
)

type Server struct {
	server.HTTPConfig
	Name string
}
type Option struct {
	Server  Server
	Method  string
	Channel string
}

func (o *Option) server() *apiServer {
	return apiserver(o.Server.Name)
}

func (o *Option) ApplyServer() error {
	if o.Server.IsEmpty() {
		return nil
	}
	return apiserver(o.Server.Name).SetConfig(&o.Server.HTTPConfig)
}

func (o *Option) Start(handler func(w http.ResponseWriter, r *http.Request)) error {
	return apiserver(o.Server.Name).Start(o.Channel, MethodMiddleware(o.Method, handler))
}

func (o *Option) Stop() error {
	return apiserver(o.Server.Name).Stop(o.Channel)
}
