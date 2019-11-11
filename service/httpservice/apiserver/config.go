package apiserver

import "github.com/herb-go/herb/server"

var defaultConfig = &server.HTTPConfig{
	ListenerConfig: server.ListenerConfig{
		Net:  "tcp",
		Addr: ":6789",
	},
}
