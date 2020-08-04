package misc

import "net/http"

//MethodMiddleware middleware which check request in method list.
//Status 405 will return if not match.
func MethodMiddleware(methods ...string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	var methodMap = map[string]bool{}
	for k := range methods {
		methodMap[methods[k]] = true
	}
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if methodMap[r.Method] {
			next(w, r)
			return
		}
		http.Error(w, http.StatusText(405), 405)
	}
}

//Method http request method middleware type
type Method string

//ServeMiddleware serve as middleware
func (m Method) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if r.Method != string(m) {
		http.Error(w, http.StatusText(405), 405)
		return
	}
	next(w, r)
}

//MethodPOST method post only middleare
var MethodPOST = Method("POST")

//MethodGET method get only middleare
var MethodGET = Method("GET")

//MethodPUT method put only middleare
var MethodPUT = Method("PUT")

//MethodDELETE method delete only middleare
var MethodDELETE = Method("DELETE")

//MethodOPTIONS method options only middleare
var MethodOPTIONS = Method("OPTIONS")
