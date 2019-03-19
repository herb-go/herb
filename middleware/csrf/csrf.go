//Package csrf provide csrf prevent middleware.
package csrf

import (
	"context"
	"net/http"
	"strconv"
	"time"
)

//ContextKey string type used in Context key
type ContextKey string

var defaultCookieName = "herb-csrf-token"
var defaultCookiePath = "/"
var defaultHeaderName = "X-CSRF-TOKEN"
var defaultFormField = "X-CSRF-TOKEN"
var defaultFailHeader = "herb-go-csrf-token-status"
var defaultFailValue = "failed"
var defaultFailStatus = http.StatusBadRequest
var defaultRequestContextKey = ContextKey("herb-csrf-token")

func (csrf *Csrf) fail(w http.ResponseWriter) {
	http.Error(w, http.StatusText(csrf.FailStatus), csrf.FailStatus)
}

//Verify Verify if the given token is equal to token value save in cookie.
//Return verification result and any error raised.
func (csrf *Csrf) Verify(r *http.Request, token string) (bool, error) {
	if !csrf.Enabled {
		return true, nil
	}
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

//ServeVerifyFormMiddleware The middleware check if the token in post form is equal to token value save in cookie
func (csrf *Csrf) ServeVerifyFormMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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

//ServeVerifyHeaderMiddleware The middleware check if the token in post form is equal to token value save in cookie
func (csrf *Csrf) ServeVerifyHeaderMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
	if !csrf.Enabled {
		return "", nil
	}
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
	} else if err != nil {
		return err
	}
	ctx := r.Context()
	r = r.WithContext(context.WithValue(ctx, csrf.RequestContextKey, c.Value))
	return nil
}

//ServeSetCsrfTokenMiddleware The middleware set a random token in cookie which is used in later verification if the cookie does not exist.
func (csrf *Csrf) ServeSetCsrfTokenMiddleware(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	if csrf.Enabled {
		err := csrf.SetCsrfToken(w, r)
		if err != nil {
			panic(err)
		}
	}
	next(w, r)
}
func (csrf *Csrf) generateToken() (string, error) {
	return csrf.TokenGenerater()
}

//DefaultTokenGenerater default csrf token generater.
//Return current timestamp string and any error if raised.
func DefaultTokenGenerater() (string, error) {
	return strconv.FormatInt(time.Now().UnixNano(), 10), nil
}

//New return a new Csrf Component with default values.
func New() *Csrf {
	c := Csrf{
		CookieName:        defaultCookieName,
		CookiePath:        defaultCookiePath,
		HeaderName:        defaultHeaderName,
		FormField:         defaultFormField,
		FailStatus:        defaultFailStatus,
		RequestContextKey: defaultRequestContextKey,
		Enabled:           true,
		TokenGenerater:    DefaultTokenGenerater,
	}
	return &c
}

//Csrf is the components provide csrf function.
//You can use Csrf.SetCsrfTokenMiddleware,Csrf.VerifyFormMiddleware,Csrf.VerifyHeaderMiddleware or Csrf.CsrfInput to protected your web app.
//All value can be change after creation.
type Csrf struct {
	CookieName        string                 //Name of cookie which the token stored in.Default value is "herb-csrf-token".
	CookiePath        string                 //Path of cookie the token stored in.Default value is "/".
	HeaderName        string                 //Name of Header which the token stroed in.Default value is "X-CSRF-TOKEN".
	FormField         string                 //Field name of post form which the token stroed in.Default value is "X-CSRF-TOKEN".
	FailStatus        int                    //Http status code returned when csrf verify failed.Default value is  http.StatusBadRequest (int 400).
	RequestContextKey ContextKey             //Context key of requst which token stored in.Default value is csrf.ContextKey("herb-csrf-token").
	Enabled           bool                   //Enabled if this middleware if enabled.
	FailHeader        string                 //FailedHeader resoponse header field send when failed
	FailValue         string                 //FailedValue resoponse header value send when failed
	TokenGenerater    func() (string, error) //TokenGenerater func to create csrf token.
}

//Config csrf config struct
type Config struct {
	CookieName        string //Name of cookie which the token stored in.Default value is "herb-csrf-token".
	CookiePath        string //Path of cookie the token stored in.Default value is "/".
	HeaderName        string //Name of Header which the token stroed in.Default value is "X-CSRF-TOKEN".
	FormField         string //Field name of post form which the token stroed in.Default value is "X-CSRF-TOKEN".
	FailStatus        int    //Http status code returned when csrf verify failed.Default value is  http.StatusBadRequest (int 400).
	RequestContextKey string //Context key of requst which token stored in.Default value is "herb-csrf-token")
	Enabled           bool   //Enabled if this middleware if enabled.
	FailHeader        string //FailedHeader resoponse header field send when failed
	FailValue         string //FailedValue resoponse header value send when failed
}

//ApplyTo apply csrf config to csrf instance.
func (c *Config) ApplyTo(csrf *Csrf) error {

	if c.CookieName != "" {
		csrf.CookieName = c.CookieName
	}
	if c.CookiePath != "" {
		csrf.CookiePath = c.CookiePath
	}
	if c.HeaderName != "" {
		csrf.HeaderName = c.HeaderName
	}
	if c.FormField != "" {
		csrf.FormField = c.FormField
	}
	if c.FailStatus != 0 {
		csrf.FailStatus = c.FailStatus
	}
	if c.RequestContextKey != "" {
		csrf.RequestContextKey = ContextKey(c.RequestContextKey)
	}
	if c.FailHeader != "" {
		csrf.FailHeader = c.FailHeader
	} else {
		csrf.FailHeader = defaultFailHeader
	}
	if c.FailValue != "" {
		csrf.FailValue = c.FailValue
	} else {
		csrf.FailValue = defaultFailValue
	}
	csrf.Enabled = c.Enabled
	return nil
}
