package httpserver

import (
	"net"
	"net/http"
	"time"
)

type Config struct {
	Net                       string
	Addr                      string
	BaseURL                   string
	ReadTimeoutInSecond       int64
	ReadHeaderTimeoutInSecond int64
	WriteTimeoutInSecond      int64
	IdleTimeoutInSecond       int64
	MaxHeaderBytes            int
}

func (c *Config) Server() *http.Server {
	server := &http.Server{
		ReadTimeout:       time.Duration(c.ReadTimeoutInSecond) * time.Second,
		ReadHeaderTimeout: time.Duration(c.ReadHeaderTimeoutInSecond) * time.Second,
		WriteTimeout:      time.Duration(c.WriteTimeoutInSecond) * time.Second,
		IdleTimeout:       time.Duration(c.IdleTimeoutInSecond) * time.Second,
		MaxHeaderBytes:    c.MaxHeaderBytes,
	}
	return server
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
