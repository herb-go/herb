package misc

import "net/http"

func PoweredBy(poweredBy string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		if poweredBy != "" {
			w.Header().Set("Powered-By", poweredBy)
		}
		next(w, r)
	}
}
