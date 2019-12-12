package service

import "net"

//ListenerConfig listener config struct
type ListenerConfig struct {
	//Net net interface,"tcp" for example.
	Net string
	//Addr network addr.
	Addr string
}

func (c *ListenerConfig) Clone() *ListenerConfig {
	return &ListenerConfig{
		Net:  c.Net,
		Addr: c.Addr,
	}
}

//Listen listen net and addr in config.
//Return net listener and any error if raised.
func (c *ListenerConfig) Listen() (net.Listener, error) {
	return net.Listen(c.Net, c.Addr)
}
