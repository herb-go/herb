package httphook

import (
	"net/http"

	"github.com/herb-go/herb/middleware/httpinfo"
)

type Handler interface {
	Handle(req *http.Request, resp *httpinfo.Response)
}

type HandlerFunc func(req *http.Request, resp *httpinfo.Response)

func (h HandlerFunc) Handle(req *http.Request, resp *httpinfo.Response) {
	h(req, resp)
}

var NopHandler = HandlerFunc(func(req *http.Request, resp *httpinfo.Response) {})
