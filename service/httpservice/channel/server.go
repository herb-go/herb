package channel

import (
	"net/http"
	"sync"

	"github.com/herb-go/herb/service/httpservice"
)

var servers = sync.Map{}

var configs = sync.Map{}

var locker = sync.Mutex{}

func GetConfig(host string) *httpservice.Config {
	v, ok := configs.Load(host)
	if v == nil || ok == false {
		return nil
	}
	return v.(*httpservice.Config)
}

func SetConfig(host string, c *httpservice.Config) {
	configs.Store(host, c)
}

var DefaultConfig = &httpservice.Config{}

type Server struct {
	server   *http.Server
	channels sync.Map
}
