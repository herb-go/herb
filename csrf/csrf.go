//Package csrf provide csrf prevent middleware.
package csrf

import (
	"context"
	"crypto/rand"
	"net/http"
)

//ContextKey string type used in Context key
type ContextKey string

var defaultTokenlength = 64
var defaultTokenMask = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var defaultCookieName = "herb-csrf-token"
var defaultCookiePath = "/"
var defaultHeaderName = "X-CSRF-TOKEN"
var defaultFormField = "X-CSRF-TOKEN"
var defaultFailStatus = http.StatusBadRequest
var defaultRequestContextKey = ContextKey("herb-csrf-token")

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

//Verify Verify if the given token is equal to token value save in cookie.
//Return verification result and any error raised.
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

//VerifyFormMiddleware The middleware check if the token in post form is equal to token value save in cookie
func (csrf *Csrf) VerifyFormMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	success, err := csrf.Verify(r, r.FormValue(csrf.FormField))
	if err != nil {
		panic(nil)
	}
	if !success {
		csrf.fail(w)
		return
	}
	next(w, r)

}

//VerifyHeaderMiddleware The middleware check if the token in post form is equal to token value save in cookie
func (csrf *Csrf) VerifyHeaderMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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

//CsrfInput return a html fragment that contains a csrf hidden input.
func (csrf *Csrf) CsrfInput(w http.ResponseWriter, r *http.Request) (string, error) {
	err := csrf.SetCsrfToken(w, r)
	if err != nil {
		return "", err
	}
	return `<input type="hidden" name="` + csrf.FormField + `" value="` + csrf.requestToken(r) + `"/>`, nil

}
func (csrf *Csrf) requestToken(r *http.Request) string {
	k := r.Context().Value(csrf.RequestContextKey)
	t, ok := k.(string)
	if !ok {
		return ""
	}
	return t
}

//SetCsrfToken set a random token in cookie which is used in later verification if the cookie does not exist.
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

//SetCsrfTokenMiddleware The middleware set a random token in cookie which is used in later verification if the cookie does not exist.
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

// New return a new Csrf Component with default values.
func New() *Csrf {
	c := Csrf{
		Tokenlength:       defaultTokenlength,
		TokenMask:         defaultTokenMask,
		CookieName:        defaultCookieName,
		CookiePath:        defaultCookiePath,
		HeaderName:        defaultHeaderName,
		FormField:         defaultFormField,
		FailStatus:        defaultFailStatus,
		RequestContextKey: defaultRequestContextKey,
	}
	return &c
}

//Csrf is the cimponents provide csrf function.
//You can use Csrf.SetCsrfTokenMiddleware,Csrf.VerifyFormMiddleware,Csrf.VerifyHeaderMiddleware or Csrf.CsrfInput to protected your web app.
//All value can be change after creation.
type Csrf struct {
	Tokenlength       int        //Length of csrf token.Default value is 64.
	TokenMask         []byte     //Which chars should be used in token.Default value is 0-9a-zA-z.
	CookieName        string     //Name of cookie which the token stored in.Default value is "herb-csrf-token".
	CookiePath        string     //Path of cookie the token stored in.Default value is "/".
	HeaderName        string     //Name of Header which the token stroed in.Default value is ""X-CSRF-TOKEN".
	FormField         string     //Field name of post form which the token stroed in.Default value is ""X-CSRF-TOKEN".
	FailStatus        int        //Http status code returned when csrf verify failed.Default value is  http.StatusBadRequest (int 400).
	RequestContextKey ContextKey //Context key of requst which token stored in.Default value is csrf.ContextKey("herb-csrf-token").
}
