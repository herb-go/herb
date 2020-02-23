package channel

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/herb-go/herb/service"
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

func SetConfig(c *httpservice.Config) {
	if c == nil {
		return

	}
	host := convertListenerToString(getListener(&c.ListenerConfig))
	setConfig(host, c)
}

var DefaultConfig = PresetDefaultConfig()

func PresetDefaultConfig() *httpservice.Config {
	return &httpservice.Config{}
}

type Server struct {
	running  int
	handlers map[string]*Handler
	server   *http.Server
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
	server, err := s.config.Serve(s.mux)
	if err != nil {
		return err
	}
	s.server = server
	s.server.Handler = s.mux
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
	if s.server != nil {
		err = s.server.Close()
	}
	s.server = nil
	return err
}
func (s *Server) handle(path string, h http.Handler) (err error) {
	defer func() {
		r := recover()
		if r == nil {
			return
		}
		e, ok := r.(error)
		if ok == false {
			s, ok := r.(string)
			if ok == false {
				panic(r)
			}
			e = errors.New(s)
		}
		err = e
	}()
	_, ok := s.handlers[path]
	if ok {
		return fmt.Errorf("channel: %s %w", path, ErrChannelUsed)
	}
	handler := NewHandler(http.StripPrefix(path, h))
	s.handlers[path] = handler
	s.mux.Handle(path, handler)
	return nil
}
func newServer(l *service.ListenerConfig) *Server {
	host := convertListenerToString(getListener(l))
	c := getConfig(host)
	if c == nil {
		c = DefaultConfig.Clone()
	}
	c.ListenerConfig = *l
	m := http.NewServeMux()
	return &Server{
		running:  0,
		config:   c,
		handlers: map[string]*Handler{},
		mux:      m,
	}
}

func GetServer(l *service.ListenerConfig) *Server {
	locker.Lock()
	defer locker.Unlock()
	return getServer(l)
}
func getServer(l *service.ListenerConfig) *Server {
	host := convertListenerToString(getListener(l))
	s := servers[host]
	if s != nil {
		return s
	}
	s = newServer(l)
	servers[host] = s
	return s
}

func ResetServers() {
	locker.Lock()
	defer locker.Unlock()
	for _, s := range servers {
		s.stopServer()
	}
	servers = map[string]*Server{}
}

func ResetConfigs() {
	locker.Lock()
	defer locker.Unlock()
	configs = sync.Map{}
}