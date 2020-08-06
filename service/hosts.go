package service

import (
	"net/http"
	"strings"
)

//HostPattern hostname pattern type
type HostPattern string

//Match check if given hostname match host
func (p HostPattern) Match(v string) bool {
	if p == "" || v == "" {
		return false
	}
	if p[0] == '.' {
		l := strings.SplitN(v, ".", 2)
		return len(l) == 2 && l[1] == string(p[1:])
	} else if p[0] == '*' {
		return strings.HasSuffix(v, string(p[1:]))
	}
	return string(p) == v
}

//Hosts hostname list middleware struct
type Hosts struct {
	Patterns []HostPattern
}

//ServeMiddleware server hostname list middleware
func (h Hosts) ServeMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if len(h.Patterns) != 0 {
		host := strings.SplitN(r.Host, ":", 2)
		for k := range h.Patterns {
			if h.Patterns[k].Match(host[0]) {
				next(w, r)
				return
			}
		}
		http.NotFound(w, r)
		return
	}
	next(w, r)
}
