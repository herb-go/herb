package responseinfo

import (
	"net/http"
)

type Handler interface {
	Handle(r http.Request, i *Info)
}

type HandlerFunc func(r http.Request, i *Info)

func (h HandlerFunc) Handle(r http.Request, i *Info) {
	h(r, i)
}
