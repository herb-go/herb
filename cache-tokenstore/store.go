package tokenstore

import "net/http"

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
	Close() error
}
