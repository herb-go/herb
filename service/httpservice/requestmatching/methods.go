package requestmatching

import (
	"net/http"
	"strings"
)

//Methods methods pattern
type Methods map[string]bool

//Add add method to pattern.
//Method will be converted to upper.
func (m *Methods) Add(method string) {
	(*m)[strings.ToUpper(method)] = true
}

//MatchRequest match request.
//Return result and any error if raised.
func (m *Methods) MatchRequest(r *http.Request) (bool, error) {
	if len(*m) == 0 {
		return true, nil
	}
	return (*m)[r.Method], nil
}

//NewMethods create new methods pattern.
func NewMethods() *Methods {
	return &Methods{}
}
