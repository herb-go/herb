package apiserver

import (
	"net/http"

	"github.com/herb-go/herb/service/httpservice"
)

type Server struct {
	httpservice.Config
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
	return apiserver(o.Server.Name).SetConfig(&o.Server.Config)
}

func (o *Option) Start(handler func(w http.ResponseWriter, r *http.Request)) error {
	return apiserver(o.Server.Name).Start(o.Channel, handler)
}

func (o *Option) Stop() error {
	return apiserver(o.Server.Name).Stop(o.Channel)
}
