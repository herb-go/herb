package httpinfo

import (
	"net/http"
)

type Validator interface {
	Validate(*http.Request, *Response) (bool, error)
}

type ValidatorFunc func(*http.Request, *Response) (bool, error)

func (f ValidatorFunc) Validate(r *http.Request, resp *Response) (bool, error) {
	return f(r, resp)
}

var ValidatorNever = ValidatorFunc(func(*http.Request, *Response) (bool, error) {
	return false, nil
})

var ValidatorAlways = ValidatorFunc(func(*http.Request, *Response) (bool, error) {
	return true, nil
})
