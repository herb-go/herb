package channel

import (
	"net/http"
	"sync"

	"github.com/herb-go/herb/service/httpservice"
)

var servers = map[string]*Server{}

var configs = sync.Map{}

var locker = sync.Mutex{}

func getConfig(host string) *httpservice.Config {
	v, ok := configs.Load(host)
	if v == nil || ok == false {
		return nil
	}
	return v.(*httpservice.Config)
}

func setConfig(host string, c *httpservice.Config) {
	configs.Store(host, c)
}

var DefaultConfig = &httpservice.Config{}

type Server struct {
	handelrs map[string]Handler
	server   *http.Server
	channels sync.Map
}

func newServer(host string) *Server {
	c := getConfig(host)
	if c == nil {
		c = DefaultConfig.Clone()
	}
	return &Server{
		handelrs: map[string]Handler{},
		server:   c.Server(),
	}
}
func GetServer(host string) *Server {
	locker.Lock()
	defer locker.Unlock()
	s := servers[host]
	if s != nil {
		return s
	}
	s = newServer(host)
	servers[host] = s
	return s
}
