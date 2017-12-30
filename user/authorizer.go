package user

import (
	"net/http"
)

type Authorizer interface {
	Authorize(*http.Request) (bool, error)
}

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
func AuthorizeOrForbiddenMiddleware(authorizer Authorizer) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return AuthorizeMiddleware(authorizer, nil)
}
