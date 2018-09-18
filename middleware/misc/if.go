package misc

import "net/http"

//MiddlewareIf return middleware after condition checked.
//If condition is true,mirddleware 'then' will be returned.
//Only checks condition when init.
func MiddlewareIf(condition bool, then func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if condition {
		return then
	}
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		next(w, r)
	}
}

//MiddlewareWhen return middleware after condition checked.
//If result of condition function is true,params 'then' will be executed.
//Panic if any error raised.
//Checks condition every time.
func MiddlewareWhen(condition func() (bool, error), then func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result, err := condition()
		if err != nil {
			panic(err)
		}
		if result {
			then(w, r, next)
			return
		}
		next(w, r)
	}
}

//If middleware checks condition.
//If condition is true,params 'then' will be executed.
//Only checks condition when init.
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
//If result of condition function is true,params 'then' will be executed.
//Panic if any error raised.
//Checks condition every time.
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
