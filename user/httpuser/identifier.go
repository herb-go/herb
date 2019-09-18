package httpuser

import (
	"net/http"
	"time"
)

//Identifier http request identifier
type Identifier interface {
	//IdentifyRequest identify http request
	//return identification and any error if rasied.
	IdentifyRequest(r *http.Request) (string, error)
}

//LogoutProvider Logout provider interface
type LogoutProvider interface {
	Logout(w http.ResponseWriter, r *http.Request) error
}

//LoginProvider Login provider interface
type LoginProvider interface {
	Login(w http.ResponseWriter, r *http.Request, id string) error
}

//LoginRedirector login redirector struct
type LoginRedirector struct {
	//LoginURL redirector will redirect user to this url if user didnot log in.
	LoginURL string
	//Cookie cookie settings.
	Cookie *http.Cookie
}

//NewLoginRedirector create new login redirector with given login url and cookie name.
func NewLoginRedirector(loginurl string, cookiename string) *LoginRedirector {
	return &LoginRedirector{
		LoginURL: loginurl,
		Cookie: &http.Cookie{
			Name:     cookiename,
			HttpOnly: false,
			Path:     "/",
		},
	}
}

//RedirectAction action which set cookie and redirect user.
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

//ClearSource return and clear the url before redirect.
//Return url and any error if raised.
func (lr *LoginRedirector) ClearSource(w http.ResponseWriter, r *http.Request) (string, error) {
	cookie, err := r.Cookie(lr.Cookie.Name)
	if err == http.ErrNoCookie {
		return "", nil
	}
	if err != nil || cookie == nil {
		return "", err
	}
	url := cookie.Value
	newCookie := &http.Cookie{
		Name:     lr.Cookie.Name,
		Path:     lr.Cookie.Path,
		Domain:   lr.Cookie.Domain,
		Value:    "",
		MaxAge:   lr.Cookie.MaxAge,
		Secure:   lr.Cookie.Secure,
		HttpOnly: lr.Cookie.HttpOnly,
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, newCookie)
	return url, nil
}

//MustClearSource return and clear the url before redirect.
func (lr *LoginRedirector) MustClearSource(w http.ResponseWriter, r *http.Request) string {
	url, err := lr.ClearSource(w, r)
	if err != nil {
		panic(err)
	}
	return url
}

//Middleware redirector middleware
func (lr *LoginRedirector) Middleware(s Identifier) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return LoginRequiredMiddleware(s, lr.RedirectAction)
}

//LoginRequiredMiddleware middleware which indentify user with identifier.
//If indentify fail param unauthorizedAction will be executed.
func LoginRequiredMiddleware(identifier Identifier, unauthorizedAction http.HandlerFunc) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		id, err := identifier.IdentifyRequest(r)
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

//LogoutMiddleware middleware which will logout user.
func LogoutMiddleware(s LogoutProvider) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		err := s.Logout(w, r)
		if err != nil {
			panic(err)
		}
		next(w, r)
	}
}

//MiddlewareForbiddenExceptForUsers middleware which identify user with identifier,and return http 403 error if user is not in users list.
func MiddlewareForbiddenExceptForUsers(identifier Identifier, users []string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		id, err := identifier.IdentifyRequest(r)
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
