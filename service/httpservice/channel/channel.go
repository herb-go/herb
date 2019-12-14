package channel

import (
	"net/http"
	"net/url"
	"sync"

	"github.com/herb-go/herb/service"
)

var channels = sync.Map{}

type Channel struct {
	service.ListenerConfig
	Path string
}

func (c *Channel) Host() string {
	u := &url.URL{
		Scheme: c.Net,
		Host:   c.Addr,
	}
	return u.String()
}
func (c *Channel) getServer() *Server {
	return getServer(c.Host())
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
