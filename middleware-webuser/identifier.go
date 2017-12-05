package webuser

import (
	"net/http"
)

type Identifier interface {
	IdentifyRequest(r *http.Request) (string, error)
}

type LogoutService interface {
	Logout(r *http.Request) error
}

func LoginRequiredMiddleware(s Identifier, unauthorizedAction http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		id, err := s.IdentifyRequest(r)
		if err != nil {
			panic(err)
		}
		if id == "" {
			if unauthorizedAction == nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
			} else {
				unauthorizedAction(w, r)
			}
			return
		}
		next(w, r)
	}
}
func LogoutMiddleware(s LogoutService) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		err := s.Logout(r)
		if err != nil {
			panic(err)
		}
		next(w, r)
	}
}

func ForbiddenExceptForUsers(s Identifier, users []string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		id, err := s.IdentifyRequest(r)
		if err != nil {
			panic(err)
		}
		if id != "" && users != nil {
			for _, v := range users {
				if v == id {
					next(w, r)
					return
				}
			}
		}
		http.Error(w, http.StatusText(http.StatusForbidden), http.StatusForbidden)
	}
}
