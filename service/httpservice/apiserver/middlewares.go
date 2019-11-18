package apiserver

import (
	"net/http"
)

//MiddlewareList http middleware list interface
type MiddlewareList interface {
	ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	Append(...func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc))
}
