package spilter

import (
	"net/http"
	"strings"

	"github.com/herb-go/herb/middleware/router"
)

func SplitFirstFolderInto(name string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return SplitFirstInto("/", name)
}

func SplitFirstInto(sep string, name string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		p := r.URL.Path
		if p[0] == '/' {
			p = p[1:]
		}
		pathlist := strings.SplitN(p, sep, 2)
		if len(pathlist) < 2 {
			pathlist = append(pathlist, "")
		}
		router.GetParams(r).Set(name, pathlist[0])
		r.URL.Path = "/" + pathlist[1]
		next(w, r)
	}
}

func DropAfterFirst(sep string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		p := r.URL.Path
		if p[0] == '/' {
			p = p[1:]
		}
		pathlist := strings.SplitN(p, sep, 2)
		r.URL.Path = "/" + pathlist[0]
		next(w, r)
	}
}
