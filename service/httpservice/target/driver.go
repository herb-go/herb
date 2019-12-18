package target

import (
	"net/http"
)

type DoerFactory func(loader func(v interface{}) error) (Doer, error)

func DefaultClientDoerFactory(loader func(v interface{}) error) (Doer, error) {
	return http.DefaultClient, nil
}
