package channel

import (
	"fmt"
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
	running  *int
	handlers map[string]*Handler
	server   *http.Server
	mux      *http.ServeMux
	channels sync.Map
}

func (s *Server) handle(path string, h http.Handler) error {
	_, ok := s.handlers[path]
	if ok {
		return fmt.Errorf("channel: %s %w", path, ErrChannelUsed)
	}
	s.handlers[path] = NewHandler(h)
	return nil
}
func newServer(host string) *Server {
	c := getConfig(host)
	if c == nil {
		c = DefaultConfig.Clone()
	}
	running := 0
	s := c.Server()
	m := http.NewServeMux()
	s.Handler = m
	return &Server{
		running:  &running,
		handlers: map[string]*Handler{},
		server:   s,
		mux:      m,
	}
}

func GetServer(host string) *Server {
	locker.Lock()
	defer locker.Unlock()
	return getServer(host)
}
func getServer(host string) *Server {
	s := servers[host]
	if s != nil {
		return s
	}
	s = newServer(host)
	servers[host] = s
	return s
}
