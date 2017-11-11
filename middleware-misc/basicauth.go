package misc

import (
	"encoding/base64"
	"net/http"
	"strings"
)

type BasicAuthConfig struct {
	Realm    string
	Username string
	Password string
}

func BasicAuth(c *BasicAuthConfig) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if c.Realm == "" || c.Username == "" || c.Password == "" {
		return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
			http.NotFound(w, r)
			return
		}
	}
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		//codes from https://gist.github.com/elithrar/9146306
		w.Header().Set("WWW-Authenticate", `Basic realm="`+c.Realm+`"`)
		s := strings.SplitN(r.Header.Get("Authorization"), " ", 2)
		if len(s) != 2 {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		b, err := base64.StdEncoding.DecodeString(s[1])
		if err != nil {
			http.Error(w, err.Error(), 401)
			return
		}

		pair := strings.SplitN(string(b), ":", 2)
		if len(pair) != 2 {
			http.Error(w, http.StatusText(401), 401)
			return
		}

		if pair[0] != c.Username || pair[1] != c.Password {
			http.Error(w, http.StatusText(401), 401)
			return
		}
		next(w, r)
	}
}
