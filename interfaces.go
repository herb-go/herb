package herb

import (
	"net/http"
)

type Middleware func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)

type DataMap interface {
	Get(name string, v interface{}) error
	Set(name string, v interface{}) error
}
type DataField interface {
	Get(v interface{}) error
	Set(v interface{}) error
}
type RequestItem interface {
	LoadFromRequest(r *http.Request) error
	SaveToRequest(r *http.Request) error
}

type Requestvalidator interface {
	ValidateRequest(r *http.Request) (bool, error)
}
