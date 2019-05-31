package basicauth

import (
	"context"
	"net/http"
)

//ContextName basicauth context name type.
type ContextName string

//IdentifyRequest indentify request by basic auth user.
//Return uesrname and any error if raised.
func (c ContextName) IdentifyRequest(r *http.Request) (string, error) {
	var username string
	u := r.Context().Value(c)
	if u != nil {
		username = (u).(string)
	}
	return username, nil
}

//Username context name for username.
var Username = ContextName("Username")

//Authorizer interface for basicauth authorizer.
type Authorizer interface {
	GetRealm() (string, error)
	Authorize(Username string, Password string) (bool, error)
}

//Users  user map authorizer.
type Users struct {
	Realm string
	Users map[string]string
}

//GetRealm return basic auth realm.
//return realm and any error if raised.
func (u *Users) GetRealm() (string, error) {
	return u.Realm, nil
}

//Authorize authorize user with username and password.
//return authorize result and any error if raised.
func (u *Users) Authorize(Username string, Password string) (bool, error) {
	if u.Users == nil || u.Users[Username] == "" || u.Users[Username] != Password {
		return false, nil
	}
	return true, nil
}

//SingleUser single user authorizer.
type SingleUser struct {
	//Realm basic auth realm.
	Realm string
	//Username basic auth username.
	Username string
	//Password basic auth password.
	Password string
}

//GetRealm return basic auth realm.
//return realm and any error if raised.
func (u *SingleUser) GetRealm() (string, error) {
	return u.Realm, nil
}

//Authorize authorize user with username and password.
//return authorize result and any error if raised.
func (u *SingleUser) Authorize(Username string, Password string) (bool, error) {
	if u.Username != Username || u.Password != Password {
		return false, nil
	}
	return true, nil
}

//GetUsername get basic auth user name from request.
func GetUsername(r *http.Request) string {
	var username string
	u := r.Context().Value(Username)
	if u != nil {
		username = (u).(string)
	}
	return username
}

//SetUsername set username to request context.
func SetUsername(r *http.Request, username string) {
	ctx := context.WithValue(r.Context(), Username, username)
	*r = *r.WithContext(ctx)
}

//Middleware use authorizer as basic auth middleware.
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
