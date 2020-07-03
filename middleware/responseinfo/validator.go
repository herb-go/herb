package responseinfo

import "net/http"

type Validator interface {
	Validate(*http.Request, *Info) (bool, error)
}

type ValidatorFunc func(*http.Request, *Info) (bool, error)

func (f ValidatorFunc) Validate(r *http.Request, i *Info) (bool, error) {
	return f(r, i)
}
