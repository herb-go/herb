package target

import (
	"net/http"
)

type DoerFactory func(loader func(v interface{}) error) (Doer, error)

type TargetFactory func(loader func(v interface{}) error) (Target, error)

func DefaultClientDoerFactory(loader func(v interface{}) error) (Doer, error) {
	return http.DefaultClient, nil
}

func URLTargetFactory(loader func(v interface{}) error) (Target, error) {
	t := NewURLTarget()
	err := loader(t)
	if err != nil {
		return nil, err
	}
	return t, nil
}
