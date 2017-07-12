package tokenstore

import (
	"net/http"
	"time"
)

type CookieWriter struct {
	http.ResponseWriter
	r       *http.Request
	store   *Store
	written bool
}

func (w *CookieWriter) WriteHeader(status int) {
	var v *TokenValues
	var err error
	if w.written == false {
		w.written = true
		v, err = w.store.GetRequestTokenValues(w.r)
		if err != nil {
			panic(err)
		}
		if v.tokenChanged {
			cookie := &http.Cookie{
				Name:     w.store.CookieName,
				Value:    v.token,
				Path:     w.store.CookiePath,
				Secure:   false,
				HttpOnly: true,
			}
			if w.store.TokenLifetime >= 0 {
				cookie.Expires = time.Now().Add(w.store.TokenLifetime)
			}
			http.SetCookie(w, cookie)
		}
	}
	w.ResponseWriter.WriteHeader(status)
}
func (w *CookieWriter) Write(data []byte) (int, error) {
	if w.written == false {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}
