package misc

import (
	"net/http"
)

//Headers headers middleware which add headers to each response
type Headers map[string]string

//ServeMiddleware serve headers settings as middleware
func (h *Headers) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	for k, v := range *h {
		w.Header().Set(k, v)
	}
	next(w, r)
}
