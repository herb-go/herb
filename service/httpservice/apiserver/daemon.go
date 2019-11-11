package apiserver

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"sync"
	"sync/atomic"

	"github.com/herb-go/herb/server"
)

type serverMap struct {
	sync.Map
}

var serversLock = sync.Mutex{}
var servers = serverMap{}

func newAPIServer() *apiServer {
	s := &apiServer{}
	s.Handler = s.handler()
	return s
}

type apiServer struct {
	serverName   string
	listener     net.Listener
	server       *http.Server
	runningLock  sync.Mutex
	runningCount int32
	running      sync.Map
	Handler      http.Handler
	configured   bool
	config       *server.HTTPConfig
	configLock   sync.Mutex
}

func (as *apiServer) Start(name string, h func(w http.ResponseWriter, r *http.Request)) error {
	as.runningLock.Lock()
	defer as.runningLock.Unlock()
	s, ok := as.running.LoadOrStore(name, h)
	if ok == true && s != nil {
		return fmt.Errorf("simple api :\" %s\" is already running", name)
	}
	r := atomic.LoadInt32(&as.runningCount)
	atomic.AddInt32(&as.runningCount, 1)
	if r == 0 {
		return as.startServer()
	}
	return nil
}

func (as *apiServer) Stop(name string) error {
	as.runningLock.Lock()
	defer as.runningLock.Unlock()
	s, ok := as.running.Load(name)
	if ok == false || s == nil {
		return fmt.Errorf("simple api :\" %s\" is not running", name)
	}
	as.running.Delete(name)
	r := atomic.AddInt32(&as.runningCount, -1)
	if r <= 0 {
		return as.stopServer()
	}
	return nil
}

func (as *apiServer) handler() http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h, _ := as.running.Load(r.URL.Path)
		if h == nil {
			http.Error(w, http.StatusText(404), 404)
			return
		}
		h.(func(w http.ResponseWriter, r *http.Request))(w, r)
	})
}

func (as *apiServer) startServer() error {
	c := as.Config()
	as.server = c.Server()
	as.server.Handler = as.Handler
	l, err := net.Listen(c.Net, c.Addr)
	if err != nil {
		return err
	}

	as.listener = l

	go func() {
		if !c.TLS {
			as.server.Serve(l)
		} else {
			as.server.ServeTLS(l, c.TLSCertPath, c.TLSKeyPath)
		}
	}()
	return nil
}

func (as *apiServer) stopServer() error {
	var err error
	if as.listener != nil {
		err = as.listener.Close()
	}
	as.server.Close()
	as.configured = false
	return err
}

func Reset() {
	serversLock.Lock()
	defer serversLock.Unlock()
	servers = serverMap{}
}

func (as *apiServer) errConfigSetted() error {
	configcontent, err := json.Marshal(as.config)
	if err != nil {
		return err
	}
	return fmt.Errorf("apiserver : config \"%s\" has been setted as \"%s\"", as.serverName, string(configcontent))
}

func (as *apiServer) CleanConfig() {
	as.configLock.Lock()
	defer as.configLock.Unlock()
	as.config = nil
	as.configured = false
}

func (as *apiServer) SetConfig(c *server.HTTPConfig) error {
	as.configLock.Lock()
	defer as.configLock.Unlock()
	if as.configured {
		return as.errConfigSetted()
	}
	as.config = c
	as.configured = true
	return nil
}

func (as *apiServer) Config() *server.HTTPConfig {
	as.configLock.Lock()
	defer as.configLock.Unlock()
	if as.config == nil {
		as.config = defaultConfig
		as.configured = true
	}
	return as.config
}

func apiserver(name string) *apiServer {
	serversLock.Lock()
	defer serversLock.Unlock()
	v, ok := servers.Load(name)
	if v == nil || ok == false {
		s := newAPIServer()
		servers.Store(name, s)
		return s
	}
	return v.(*apiServer)
}
