package service

import "net"

//ListenerOption listener option
type ListenerOption interface {
	Listen() (net.Listener, error)
}

//TLSOption server tls option interface
type TLSOption interface {
	//ServerIsTLSEnabeld return is server tls enabled
	ServerTLSEnabeld() bool
	//ServerTLSCertPath return server tls cert path
	ServerTLSCertPath() string
	//ServerTLSKeyPath resturn serve tls key path
	ServerTLSKeyPath() string
}
