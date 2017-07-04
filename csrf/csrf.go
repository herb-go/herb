package csrf

import (
	"context"
	"crypto/rand"
	"net/http"
)

type keyType string

var DefaultTokenlength = 64
var DefaultTokenMask = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var DefaultCookieName = "herb-csrf-token"
var DefaultCookiePath = "/"
var DefaultHeaderName = "X-CSRF-TOKEN"
var DefaultPostName = "X-CSRF-TOKEN"
var DefaultFailStatus = http.StatusBadRequest
var RequestContextKey = keyType("herb-csrf-token")

func randMaskedBytes(mask []byte, length int) ([]byte, error) {
	token := make([]byte, length)
	masked := make([]byte, length)
	_, err := rand.Read(token)
	if err != nil {
		return masked, err
	}
	l := len(mask)
	for k, v := range token {
		index := int(v) % l
		masked[k] = mask[index]
	}
	return masked, nil
}
func (csrf *Csrf) fail(w http.ResponseWriter) {
	http.Error(w, http.StatusText(csrf.FailStatus), csrf.FailStatus)
}
func (csrf *Csrf) Verify(r *http.Request, token string) (bool, error) {
	c, err := r.Cookie(csrf.CookieName)
	if err == http.ErrNoCookie {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	if c.Value == "" {
		return false, nil
	}
	return c.Value == token, nil
}
func (csrf *Csrf) VerifyPost(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	success, err := csrf.Verify(r, r.FormValue(csrf.HeaderName))
	if err != nil {
		panic(nil)
	}
	if !success {
		csrf.fail(w)
		return
	}
	next(w, r)

}
func (csrf *Csrf) VerifyHeader(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	success, err := csrf.Verify(r, r.Header.Get(csrf.HeaderName))

	if err != nil {
		panic(nil)
	}
	if !success {
		csrf.fail(w)
		return
	}
	next(w, r)
}
func (csrf *Csrf) CsrfInput(w http.ResponseWriter, r *http.Request) (string, error) {
	err := csrf.SetCsrfToken(w, r)
	if err != nil {
		return "", err
	}
	return `<input type="hidden" name="` + csrf.PostName + `" value="` + csrf.requestToken(r) + `"/>`, nil

}
func (csrf *Csrf) requestToken(r *http.Request) string {
	k := r.Context().Value(csrf.RequestContextKey)
	t, ok := k.(string)
	if !ok {
		return ""
	}
	return t
}
func (csrf *Csrf) SetCsrfToken(w http.ResponseWriter, r *http.Request) error {
	rt := csrf.requestToken(r)
	if rt != "" {
		return nil
	}
	c, err := r.Cookie(csrf.CookieName)
	if err == http.ErrNoCookie {
		t, err := csrf.generateToken()
		if err != nil {
			return err
		}
		c = &http.Cookie{
			Name:     csrf.CookieName,
			Value:    t,
			Path:     csrf.CookiePath,
			Secure:   false,
			HttpOnly: false,
		}
		http.SetCookie(w, c)
	} else if err == nil {
		return err
	}
	ctx := r.Context()
	r = r.WithContext(context.WithValue(ctx, csrf.RequestContextKey, c.Value))
	return nil
}
func (csrf *Csrf) SetCsrfTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	err := csrf.SetCsrfToken(w, r)
	if err != nil {
		panic(err)
	}
	next(w, r)
}
func (csrf *Csrf) generateToken() (string, error) {
	token, err := randMaskedBytes(csrf.TokenMask, csrf.Tokenlength)
	if err != nil {
		return "", err
	}
	return string(token), err
}
func New() *Csrf {
	c := Csrf{
		Tokenlength:       DefaultTokenlength,
		TokenMask:         DefaultTokenMask,
		CookieName:        DefaultCookieName,
		CookiePath:        DefaultCookiePath,
		HeaderName:        DefaultHeaderName,
		PostName:          DefaultPostName,
		FailStatus:        DefaultFailStatus,
		RequestContextKey: RequestContextKey,
	}
	return &c
}

type Csrf struct {
	Tokenlength       int
	TokenMask         []byte
	CookieName        string
	CookiePath        string
	HeaderName        string
	PostName          string
	FailStatus        int
	RequestContextKey keyType
}
