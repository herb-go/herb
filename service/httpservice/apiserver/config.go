package apiserver

import (
	"github.com/herb-go/herb/service"
	"github.com/herb-go/herb/service/httpservice"
)

var defaultConfig = &httpservice.Config{
	ListenerConfig: service.ListenerConfig{
		Net:  "tcp",
		Addr: ":6789",
	},
}
