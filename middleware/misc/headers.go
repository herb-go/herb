package misc

import (
	"net/http"
)

type Headers map[string]string

func (h *Headers) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	for k, v := range *h {
		w.Header().Set(k, v)
	}
	next(w, r)
}
