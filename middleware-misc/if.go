package misc

import "net/http"

func If(condition bool, h ...func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if condition {
		return h
	}
	return []func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc){}
}

func StatusIfNot(condition bool, status int) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if condition {
			next(w, r)
			return
		}
		http.Error(w, http.StatusText(status), status)
	}
}
