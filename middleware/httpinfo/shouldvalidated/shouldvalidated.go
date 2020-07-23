package shouldvalidated

import (
	"net/http"

	"github.com/herb-go/herb/middleware/httpinfo"
)

type Validator func([]byte) (bool, error)

func New(field httpinfo.Field, v Validator, OnFail http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		data, ok, err := field.LoadInfo(r)
		if err != nil {
			panic(err)
		}
		if !ok {
			http.NotFound(w, r)
			return
		}
		ok, err = v(data)
		if err != nil {
			panic(err)
		}
		if !ok {
			OnFail(w, r)
			return
		}
		next(w, r)
	}
}

func NewNotFoundMiddleware(field httpinfo.Field, v Validator) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return New(field, v, http.NotFound)
}
