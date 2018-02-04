package user

import (
	"net/http"
)

//Authorizer user role authorizer interface
type Authorizer interface {
	//Authorize authorize http request.
	//Return authorize result and any error raised.
	Authorize(*http.Request) (bool, error)
}

//AuthorizeMiddleware middleware which authorize http request with authorizer.
//Params unauthorizedAction will be executed if authorize fail.
//If authorize fail and params unauthorizedAction is nil,http error 403 will be execute.
func AuthorizeMiddleware(authorizer Authorizer, unauthorizedAction http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		result, err := authorizer.Authorize(r)
		if err != nil {
			panic(err)
		}
		if result != true {
			if unauthorizedAction == nil {
				http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
			} else {
				unauthorizedAction(w, r)
			}
			return
		}
		next(w, r)
	}
}

//AuthorizeOrForbiddenMiddleware middleware which authorize http request with authorizer.
//http error 403  will be executed if authorize fail.
func AuthorizeOrForbiddenMiddleware(authorizer Authorizer) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return AuthorizeMiddleware(authorizer, nil)
}
