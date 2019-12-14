package channel

import (
	"fmt"
	"net"
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
	running  int
	handlers map[string]*Handler
	server   *http.Server
	listener net.Listener
	mux      *http.ServeMux
	channels sync.Map
	config   *httpservice.Config
}

func (s *Server) start(path string) error {
	var err error
	h, ok := s.handlers[path]
	if ok == false {
		return fmt.Errorf("channel: %s %w", path, ErrChannelNotRegistered)
	}
	if !h.Stoped() {
		return fmt.Errorf("channel: %s %w", path, ErrChannelStarted)
	}
	h.Start()
	if s.running == 0 {
		err = s.startServer()
		if err != nil {
			return err
		}
	}
	s.running++
	return nil
}

func (s *Server) startServer() error {
	l, err := net.Listen(s.config.Net, s.config.Addr)
	if err != nil {
		return err
	}
	s.listener = l
	s.server = s.config.Server()
	go func() {
		if !s.config.TLS {
			s.server.Serve(l)
		} else {
			s.server.ServeTLS(l, s.config.TLSCertPath, s.config.TLSKeyPath)
		}
	}()
	return nil
}
func (s *Server) stop(path string) error {
	var err error
	h, ok := s.handlers[path]
	if ok == false {
		return fmt.Errorf("channel: %s %w", path, ErrChannelNotRegistered)
	}
	if h.Stoped() {
		return fmt.Errorf("channel: %s %w", path, ErrChannelStopped)
	}
	h.Stop()

	if s.running == 1 {
		err = s.stopServer()
		if err != nil {
			return err
		}
	}
	s.running--
	return nil
}

func (s *Server) stopServer() error {
	var err error
	err = s.server.Close()
	if err != nil {
		return err
	}
	return s.listener.Close()
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
	s := c.Server()
	m := http.NewServeMux()
	s.Handler = m
	return &Server{
		running:  0,
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
