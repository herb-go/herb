package session

import (
	"net/http"
	"time"
)

type cookieResponseWriter struct {
	http.ResponseWriter
	r       *http.Request
	store   *Store
	written bool
}

func (w *cookieResponseWriter) WriteHeader(status int) {
	var td *Session
	var err error
	if w.written == false {
		w.written = true
		td, err = w.store.GetRequestSession(w.r)
		if err != nil {
			panic(err)
		}
		err = w.store.SaveRequestSession(w.r)
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
			if td.token != "" {
				if w.store.TokenLifetime >= 0 {
					cookie.Expires = time.Now().Add(w.store.TokenLifetime)
				}
			} else {
				cookie.Expires = time.Unix(0, 0)
			}
			if td.HasFlag(FlagTemporay) {
				cookie.Expires = time.Unix(0, 0)
			}
			http.SetCookie(w, cookie)
		}
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *cookieResponseWriter) Write(data []byte) (int, error) {
	if w.written == false {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}
