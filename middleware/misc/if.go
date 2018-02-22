package misc

import "net/http"

//If middleware checks condition.
//If condition is true,parans then will be executed.
func If(condition bool, then http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if condition {
			then(w, r)
			return
		}
		next(w, r)
	}
}

//When middleware checks condition.
//If result of condition function is true,parans then will be executed.
//Panic if any error raised.
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

//ErrorIf middleware return http error if condition is true.
func ErrorIf(condition bool, status int) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return If(condition, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(status), status)
	})
}

//ErrorWhen middleware return http error if result of condition function is true.
//Panic if any error raised.
func ErrorWhen(condition func() (bool, error), status int) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return When(condition, func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, http.StatusText(status), status)
	})
}
