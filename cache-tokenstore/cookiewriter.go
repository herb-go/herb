package tokenstore

import (
	"net/http"
	"time"
)

//CookieResponseWriter ResponseWriter that update cookie if token changed.
type CookieResponseWriter struct {
	http.ResponseWriter
	r       *http.Request
	store   *Store
	written bool
}

//WriteHeader Worked as Response Writer WriteHeader.
func (w *CookieResponseWriter) WriteHeader(status int) {
	var td *TokenData
	var err error
	if w.written == false {
		w.written = true
		td, err = w.store.GetRequestTokenData(w.r)
		if err != nil {
			panic(err)
		}
		if td.tokenChanged {
			cookie := &http.Cookie{
				Name:     w.store.CookieName,
				Value:    td.token,
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

//WriteHeader Worked as Response Writer Write.
func (w *CookieResponseWriter) Write(data []byte) (int, error) {
	if w.written == false {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}
