type ClientStoreResponseWriter struct {
	http.ResponseWriter
	r       *http.Request
	store   *ClientStore
	written bool
}

func (w *ClientStoreResponseWriter) WriteHeader(status int) {
	var td *TokenData
	var err error
	if w.written == false {
		w.written = true
		td, err = w.store.GetRequestTokenData(w.r)
		if err != nil {
			panic(err)
		}
		err = w.store.SaveRequestTokenData(w.r)
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

func (w *ClientStoreResponseWriter) Write(data []byte) (int, error) {
	if w.written == false {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}
