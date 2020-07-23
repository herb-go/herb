package httpinfo

import (
	"net/http"
)

//Field request info field.
//Field is used to load specific information from http request
type Field interface {
	LoadInfo(r *http.Request) ([]byte, bool, error)
}

type FieldFunc func(r *http.Request) ([]byte, bool, error)

func (f FieldFunc) LoadInfo(r *http.Request) ([]byte, bool, error) {
	return f(r)
}
