package misc

import "net/http"

func If(condition bool, then http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if condition {
			then(w, r)
			return
		}
		next(w, r)
	}
}

func When(condition func() (bool, error), then http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result, err := condition()
		if err != nil {
			panic(err)
		}
		if result {
			then(w, r)
			return
		}
		next(w, r)
	}
}
func ErrorIf(condition bool, status int) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return If(condition, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(status), status)
	})
}

func ErrorWhen(condition func() (bool, error), status int) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return When(condition, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(status), status)
	})
}
