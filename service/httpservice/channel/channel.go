package channel

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/herb-go/herb/service"
)

var channels = sync.Map{}

var DefaultNet = "tcp"

var DefaultAddr = "127.0.0.1:2531"

type Channel struct {
	service.ListenerConfig
	Path string
}

func getListener(l *service.ListenerConfig) (string, string) {
	net := DefaultNet
	addr := DefaultAddr
	if l != nil {
		if l.Net != "" {
			net = l.Net
		}
		if l.Addr != "" {
			addr = l.Addr
		}
	}
	return net, addr
}

func convertListenerToString(net, addr string) string {
	u := &url.URL{
		Scheme: net,
		Host:   addr,
	}
	return u.String()
}

func (c *Channel) getServer() *Server {
	return getServer(&c.ListenerConfig)
}
func (c *Channel) Handle(h http.Handler) error {
	locker.Lock()
	defer locker.Unlock()
	return c.getServer().handle(c.Path, h)
}

func (c *Channel) Start() error {
	locker.Lock()
	defer locker.Unlock()
	return c.getServer().start(c.Path)
}

func (c *Channel) Stop() error {
	locker.Lock()
	defer locker.Unlock()
	return c.getServer().stop(c.Path)
}

func NewChannel() *Channel {
	return &Channel{}
}
