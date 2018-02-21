package basicauth

import (
	"context"
	"net/http"
)

type ContextName string

func (c ContextName) IdentifyRequest(r *http.Request) (string, error) {
	var username string
	u := r.Context().Value(c)
	if u != nil {
		username = (u).(string)
	}
	return username, nil
}

var Username = ContextName("Username")

type Authorizer interface {
	GetRealm() (string, error)
	Authorize(Username string, Password string) (bool, error)
}
type SingleUser struct {
	Realm    string
	Username string
	Password string
}

func GetUsername(r *http.Request) string {
	var username string
	u := r.Context().Value(Username)
	if u != nil {
		username = (u).(string)
	}
	return username
}
func SetUsername(r *http.Request, username string) {
	ctx := context.WithValue(r.Context(), Username, username)
	*r = *r.WithContext(ctx)
}
func (c *SingleUser) GetRealm() (string, error) {
	return c.Realm, nil
}
func (c *SingleUser) Authorize(Username string, Password string) (bool, error) {
	if c.Username != Username || c.Password != Password {
		return false, nil
	}
	return true, nil
}
func Middleware(c Authorizer) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	realm, err := c.GetRealm()
	if err != nil {
		panic(err)
	}
	if realm == "" {
		return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			http.NotFound(w, r)
			return
		}
	}
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		username, password, _ := r.BasicAuth()
		result, err := c.Authorize(username, password)
		if err != nil {
			panic(err)
		}
		if !result {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			http.Error(w, http.StatusText(401), 401)
			return
		}
		SetUsername(r, username)
		next(w, r)
	}
}
