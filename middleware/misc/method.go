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
