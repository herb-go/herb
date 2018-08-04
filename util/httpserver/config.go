package httpserver

import (
	"net"
)

type Config struct {
	Net     string
	Addr    string
	BaseURL string
}

func (c *Config) Listen() (net.Listener, error) {
	return net.Listen(c.Net, c.Addr)
}
func (c *Config) MustListen() net.Listener {
	l, err := net.Listen(c.Net, c.Addr)
	if err != nil {
		panic(err)
	}
	return l
}
