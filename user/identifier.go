package user

import (
	"net/http"
	"time"
)

type Identifier interface {
	IdentifyRequest(r *http.Request) (string, error)
}

type LogoutService interface {
	Logout(r *http.Request) error
}

type LoginRedirector struct {
	LoginURL string
	Cookie   *http.Cookie
}

func NewLoginRedirector(loginurl string, cookiename string) *LoginRedirector {
	return &LoginRedirector{
		LoginURL: loginurl,
		Cookie: &http.Cookie{
			Name:     cookiename,
			HttpOnly: true,
			Path:     "/",
		},
	}
}
func (lr *LoginRedirector) RedirectAction(w http.ResponseWriter, r *http.Request) {
	cookie := &http.Cookie{
		Name:     lr.Cookie.Name,
		Path:     lr.Cookie.Path,
		Domain:   lr.Cookie.Domain,
		Value:    r.RequestURI,
		MaxAge:   lr.Cookie.MaxAge,
		Secure:   lr.Cookie.Secure,
		HttpOnly: lr.Cookie.HttpOnly,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, r, lr.LoginURL, 302)
}
func (lr *LoginRedirector) ClearSource(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(lr.Cookie.Name)
	if err != nil || cookie == nil {
		return "", err
	}
	url := cookie.Value
	cookie.Value = ""
	cookie.Expires = time.Unix(0, 0)
	http.SetCookie(w, cookie)
	return url, nil
}
func (lr *LoginRedirector) MustClearSource(w http.ResponseWriter, r *http.Request) string {
	url, err := lr.ClearSource(w, r)
	if err != nil {
		panic(err)
	}
	return url
}
func (lr *LoginRedirector) Middleware(s Identifier) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return LoginRequiredMiddleware(s, lr.RedirectAction)
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
