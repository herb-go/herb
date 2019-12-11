package channel

import (
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
