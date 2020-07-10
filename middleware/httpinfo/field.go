package httpinfo

import (
	"net/http"
)

//Field request info field.
//Field is used to load specific information from http request
type Field interface {
	LoadInfo(r *http.Request) ([]byte, bool, error)
}
