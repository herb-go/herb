package misc

import "net/http"

func NoCache(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	h := w.Header()
	h.Set("Cache-Control", "no-cache, no-store, must-revalidate")
	h.Set("Pragma", "no-cache")
	h.Set("Expires", "0")
	next(w, r)
}
