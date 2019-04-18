package user

import (
	"net/http"
	"time"
)

//Redirector redirector request when condition is true.
type Redirector struct {
	TargetURL string
	Cookie    *http.Cookie
	Condition func(w http.ResponseWriter, req *http.Request) bool
}

// NewRedirector create new redirector wotj govem cookie name and condition
func NewRedirector(url string, cookiename string, condition func(w http.ResponseWriter, req *http.Request) bool) *Redirector {
	return &Redirector{
		TargetURL: url,
		Cookie: &http.Cookie{
			Name:     cookiename,
			HttpOnly: false,
			Path:     "/",
		},
		Condition: condition,
	}
}

// RedirectAction  redirect action of redirector
func (r *Redirector) RedirectAction(w http.ResponseWriter, req *http.Request) {
	cookie := &http.Cookie{
		Name:     r.Cookie.Name,
		Path:     r.Cookie.Path,
		Domain:   r.Cookie.Domain,
		Value:    req.RequestURI,
		MaxAge:   r.Cookie.MaxAge,
		Secure:   r.Cookie.Secure,
		HttpOnly: r.Cookie.HttpOnly,
	}
	http.SetCookie(w, cookie)
	http.Redirect(w, req, r.TargetURL, 302)
}

//ClearSource return and clear the url before redirect.
//Return url and any error if raised.
func (r *Redirector) ClearSource(w http.ResponseWriter, req *http.Request) (string, error) {
	cookie, err := req.Cookie(r.Cookie.Name)
	if err == http.ErrNoCookie {
		return "", nil
	}
	if err != nil || cookie == nil {
		return "", err
	}
	url := cookie.Value
	newCookie := &http.Cookie{
		Name:     r.Cookie.Name,
		Path:     r.Cookie.Path,
		Domain:   r.Cookie.Domain,
		Value:    "",
		MaxAge:   r.Cookie.MaxAge,
		Secure:   r.Cookie.Secure,
		HttpOnly: r.Cookie.HttpOnly,
		Expires:  time.Unix(0, 0),
	}
	http.SetCookie(w, newCookie)
	return url, nil
}

//MustClearSource return and clear the url before redirect.
func (r *Redirector) MustClearSource(w http.ResponseWriter, req *http.Request) string {
	url, err := r.ClearSource(w, req)
	if err != nil {
		panic(err)
	}
	return url
}

//Middleware redirector middleware
func (r *Redirector) Middleware() func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		if r.Condition(w, req) {
			r.RedirectAction(w, req)
			return
		}
		next(w, req)
	}
}
