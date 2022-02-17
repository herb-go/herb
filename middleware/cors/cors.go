package cors

import (
	"net/http"
	"path"
	"strings"
)

const (
	HeaderOrigin                        = "origin"
	HeaderAccessControlAllowOrigin      = "Access-Control-Allow-Origin"
	HeaderMaxAge                        = "Access-Control-Max-Age"
	HeaderVary                          = "Vary"
	HeaderAccessControlAllowCredentials = "Access-Control-Allow-Credentials"
	HeaderAccessControlAllowMethods     = "Access-Control-Allow-Methods"
	HeaderAccessControlAllowHeaders     = "Access-Control-Allow-Headers"
	HeaderAccessControlExposeHeaders    = "Access-Control-Expose-Headers"
)

const DefaultMaxAge = "86400"

func DefaultOriginValidator(c *CORS, r *http.Request) (string, error) {
	origin := r.Header.Get(HeaderOrigin)
	if origin == "" {
		return "", nil
	}
	for k := range c.Origins {
		ok, err := path.Match(c.Origins[k], origin)
		if err == nil && ok {
			return origin, nil
		}
	}
	return "", nil
}

type CORS struct {
	Enabled          bool
	AllowedHeaders   []string
	AllowedMethods   []string
	ExposeHeaders    []string
	MaxAge           string
	Origins          []string
	originValidator  func(c *CORS, r *http.Request) (string, error)
	AllowCredentials bool
}

func New() *CORS {
	return &CORS{
		AllowedMethods: []string{"POST", "GET", "OPTIONS"},
	}
}
func (c *CORS) Preflight(w http.ResponseWriter, r *http.Request) {
	if c.Enabled {
		v := c.originValidator
		if v == nil {
			v = DefaultOriginValidator
		}
		origin, err := v(c, r)
		if err != nil {
			panic(err)
		}
		if origin != "" {
			w.Header().Set(HeaderAccessControlAllowOrigin, origin)
			w.Header().Set(HeaderVary, "Origin")
			if c.MaxAge != "" {
				w.Header().Set(HeaderMaxAge, c.MaxAge)
			} else {
				w.Header().Set(HeaderMaxAge, DefaultMaxAge)
			}
			if len(c.AllowedHeaders) > 0 {
				w.Header().Set(HeaderAccessControlAllowHeaders, strings.Join(c.AllowedHeaders, ", "))
			}
			if len(c.AllowedMethods) > 0 {
				w.Header().Set(HeaderAccessControlAllowMethods, strings.Join(c.AllowedMethods, ", "))
			}
			if len(c.ExposeHeaders) > 0 {
				w.Header().Set(HeaderAccessControlExposeHeaders, strings.Join(c.ExposeHeaders, ", "))
			}
			if c.AllowCredentials {
				w.Header().Set(HeaderAccessControlAllowCredentials, "true")
			}
		}
	}
	w.WriteHeader(200)
}
func (c *CORS) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if c.Enabled {
		v := c.originValidator
		if v == nil {
			v = DefaultOriginValidator
		}
		origin, err := v(c, r)
		if err != nil {
			panic(err)
		}
		if origin != "" {
			w.Header().Set(HeaderAccessControlAllowOrigin, origin)
			w.Header().Set(HeaderVary, "Origin")
		}
	}
	next(w, r)
}
