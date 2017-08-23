package tokenstore

import (
	"errors"
	"net/http"
)

var (
	//ErrTokenNotValidated raised when the given token is not validated(for example: token is empty string)
	ErrTokenNotValidated = errors.New("Token not validated")
	//ErrRequestTokenNotFound raised when token is not found in context.You should use cookiemiddle or headermiddle or your our function to install the token.
	ErrRequestTokenNotFound = errors.New("Request token not found.Did you forget use middleware?")
	//ErrMustRegistePtr raised when the registerd interface is not a poioter to struct.
	ErrMustRegistePtr = errors.New("Must registe struct pointer")
	//ErrFeatureNotSupported raised when fearture is not supoprted.
	ErrFeatureNotSupported = errors.New("Feature is not supported")
)

//ContextKey string type used in Context key
type ContextKey string

var defaultTokenContextName = ContextKey("token")

//MustRegenerateToken Regenerate the token name and data with give owner,and save to request.
//Panic if any error raised.
func MustRegenerateToken(s *CacheStore, r *http.Request, owner string) *TokenData {
	v := MustGetRequestTokenData(s, r)
	err := v.RegenerateToken(owner)
	if err != nil {
		panic(err)
	}
	return v
}

//MustGetRequestTokenData get stored  token data from request.
//Panic if any error raised.
func MustGetRequestTokenData(s Store, r *http.Request) (v *TokenData) {
	v, err := s.GetRequestTokenData(r)
	if err != nil {
		panic(err)
	}
	return v
}

//MustRegisterField registe filed to store.
//registered field can be used directly with request to load or save the token value.
//Parameter Key filed name.
//Parameter v should be pointer to empty data model which data filled in.
//Return a new Token field.
//Panic if any error raised.
func MustRegisterField(s Store, Key string, v interface{}) *TokenField {
	tf, err := s.RegisterField(Key, v)
	if err != nil {
		panic(err)
	}
	return tf
}

type Store interface {
	GetTokenData(token string) (td *TokenData, err error)
	GetTokenDataToken(td *TokenData) (token string, err error)
	GetRequestTokenData(r *http.Request) (td *TokenData, err error)
	GenerateToken(owner string) (token string, err error)
	GenerateTokenData(token string) (td *TokenData, err error)
	LoadTokenData(v *TokenData) error
	SaveTokenData(t *TokenData) error
	RegisterField(Key string, v interface{}) (*TokenField, error)
	InstallTokenToRequest(r *http.Request, token string) (td *TokenData, err error)
	CookieMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	HeaderMiddleware(Name string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	LogoutMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	Close() error
}
