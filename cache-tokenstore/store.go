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
)

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
func MustGetRequestTokenData(s *CacheStore, r *http.Request) (v *TokenData) {
	v, err := s.GetRequestTokenData(r)
	if err != nil {
		panic(err)
	}
	return v
}

type Store interface {
	GetTokenData(token string) (td *TokenData)
	GetRequestTokenData(r *http.Request) (td *TokenData, err error)
	GenerateToken(owner string) (token string, err error)
	GenerateTokenData(token string) *TokenData
	LoadTokenData(v *TokenData) error
	SaveTokenData(t *TokenData) error
	RegisterField(Key string, v interface{}) (*TokenField, error)
	CookieMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	HeaderMiddleware(Name string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	LogoutMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc)
	Close() error
}
