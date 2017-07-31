//Package tokenstore is used to store user data in token based cache.
//It is normally used as user session or token.
//This package depands github.com/herb-go/herb/cache.
package tokenstore

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"errors"

	"github.com/herb-go/herb/cache"
)

//ContextKey string type used in Context key
type ContextKey string

var defaultTokenContextName = ContextKey("token")

var defaultUpdateActiveInterval = 5 * time.Minute

var defaultTokenMaxLifetime = 365 * 24 * time.Hour
var (
	defaultCookieName = "herb-session"
	defaultCookiePath = "/"
)

func defaultTokenGenerater(s *Store, owner string) (token string, err error) {
	t, err := cache.RandMaskedBytes(cache.TokenMask, 256)
	if err != nil {
		return
	}

	token = owner + "-" + string(t)
	return
}

var (
	//ErrTokenNotValidated raised when the given token is not validated(for example: token is empty string)
	ErrTokenNotValidated = errors.New("Token not validated")
	//ErrRequestTokenNotFound raised when token is not found in context.You should use cookiemiddle or headermiddle or your our function to install the token.
	ErrRequestTokenNotFound = errors.New("Request token not found.Did you forget use middleware?")
	//ErrMustRegistePtr raised when the registerd interface is not a poioter to struct.
	ErrMustRegistePtr = errors.New("Must registe struct pointer")
)

//New create a new token store with given cache and token lifetime.
//Cache is the cache which dates stored in.
//TokenLifeTime is the token initial expired tome.
//Return a new token store.
//All other property of the store can be set after creation.
func New(Cache *cache.Cache, TokenLifetime time.Duration) *Store {
	return &Store{
		Fields:               map[string]TokenField{},
		Cache:                Cache,
		TokenContextName:     defaultTokenContextName,
		CookieName:           defaultCookieName,
		CookiePath:           defaultCookiePath,
		TokenLifetime:        TokenLifetime,
		UpdateActiveInterval: defaultUpdateActiveInterval,
		TokenMaxLifetime:     defaultTokenMaxLifetime,
		TokenGenerater:       defaultTokenGenerater,
	}
}

//Store is the stuct store token data in cache.
type Store struct {
	Fields               map[string]TokenField                                  //All registered field
	Cache                *cache.Cache                                           //Cache which stores token data
	TokenGenerater       func(s *Store, owner string) (token string, err error) //Token name generate func
	TokenLifetime        time.Duration                                          //Token initial expired time.Token life time can be update when accessed if UpdateActiveInterval is greater than 0.
	TokenMaxLifetime     time.Duration                                          //Token max life time.Token can't live more than TokenMaxLifetime if TokenMaxLifetime if greater than 0.
	TokenContextName     ContextKey                                             //Name in request context store the token  data.Default value is "token".
	CookieName           string                                                 //Cookie name used in CookieMiddleware.Default value is "herb-session".
	CookiePath           string                                                 //Cookie path used in cookieMiddleware.Default value is "/".
	AutoGenerate         bool                                                   //Whether auto generate token when guset visit.Default value is false.
	UpdateActiveInterval time.Duration                                          //The interval between who token active time update.If less than or equal to 0,the token life time will not be refreshed.
}

//RegisterField registe filed to store.
//registered field can be used directly with request to load or save the token value.
//Key filed name.
//v the empty data struct pointer.
//Return a new Token field and error.
func (s *Store) RegisterField(Key string, v interface{}) (*TokenField, error) {
	tp := reflect.TypeOf(v)
	if tp.Kind() != reflect.Ptr {
		return nil, ErrMustRegistePtr
	}
	tf := TokenField{
		Key:   Key,
		Type:  tp.Elem(),
		Store: s,
	}
	s.Fields[Key] = tf
	return &tf, nil
}

//MustRegisterField registe filed to store.
//registered field can be used directly with request to load or save the token value.
//Key filed name.
//v the empty data struct pointer.
//Return a new Token field.
//Panic if any error raised.
func (s *Store) MustRegisterField(Key string, v interface{}) *TokenField {
	tf, err := s.RegisterField(Key, v)
	if err != nil {
		panic(err)
	}
	return tf
}

//GenerateToken generate new token name with given owner.
//Return the new token name and error.
func (s *Store) GenerateToken(owner string) (token string, err error) {
	return s.TokenGenerater(s, owner)
}

//GenerateTokenData generate new token data with given token.
//Return a new TokenData.
func (s *Store) GenerateTokenData(token string) *TokenData {
	td := NewTokenData(token, s)
	td.tokenChanged = true
	return td
}
func (s *Store) loadTokenData(v *TokenData) error {
	token := v.token
	if token == "" {
		return cache.ErrKeyUnavailable
	}

	bytes, err := s.Cache.GetBytesValue(token)
	if err != nil {
		return err
	}
	err = v.Unmarshal(token, bytes)
	if err == nil {
		v.token = token
		v.store = s
	}
	if s.TokenMaxLifetime > 0 && time.Unix(v.CreatedTime, 0).Add(s.TokenMaxLifetime).Before(time.Now()) {
		return cache.ErrNotFound
	}
	return err
}

func (s *Store) saveTokenData(td *TokenData) error {
	var err error
	token := td.token
	if token == "" {
		return cache.ErrKeyUnavailable
	}
	if s.TokenLifetime >= 0 {
		td.ExpiredAt = time.Now().Add(s.TokenLifetime).Unix()
	} else {
		td.ExpiredAt = -1
	}
	bytes, err := td.Marshal()
	if err != nil {
		return err
	}
	if td.oldToken == td.token {
		err = s.Cache.UpdateBytesValue(token, bytes, s.TokenLifetime)
	} else {
		err = s.Cache.SetBytesValue(token, bytes, s.TokenLifetime)
	}
	return err
}

//DeleteToken delete the token with given name.
//Return any error if raised.
func (s *Store) DeleteToken(token string) error {

	return s.Cache.Del(token)
}

//GetTokenData get the token data with give name .
//Return the TokenData
func (s *Store) GetTokenData(token string) (td *TokenData) {
	td = NewTokenData(token, s)
	td.oldToken = token
	return
}

//InstallTokenToRequest install the give token to request.
//Tokendata will be stored in request context which named by TokenContextName of store.
//You should use this func when use your own token binding func instead of CookieMiddleware or HeaderMiddleware
//Return the loaded TokenData and any error raised.
func (s *Store) InstallTokenToRequest(r *http.Request, token string) (td *TokenData, err error) {
	td = s.GetTokenData(token)
	if token == "" && s.AutoGenerate == true {
		err = td.RegenerateToken("")
		if err != nil {
			return
		}
	}

	ctx := context.WithValue(r.Context(), s.TokenContextName, td)
	*r = *r.WithContext(ctx)
	return
}

//CookieMiddleware return a Middleware which install the token which special by cookie.
//This middleware will save token after request finished if the token changed,and update cookie if necessary.
func (s *Store) CookieMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var token string
		cookie, err := r.Cookie(s.CookieName)
		if err == http.ErrNoCookie {
			err = nil
			token = ""
		} else if err != nil {
			panic(err)
		} else {
			token = cookie.Value
		}
		_, err = s.InstallTokenToRequest(r, token)
		if err != nil {
			panic(err)
		}
		cw := CookieResponseWriter{
			ResponseWriter: w,
			r:              r,
			store:          s,
			written:        false,
		}
		next(&cw, r)
		err = s.Save(r)
		if err != nil {
			panic(err)
		}
	}
}

//HeaderMiddleware return a Middleware which install the token which special by Header with given name.
//this middleware will save token after request finished if the token changed.
func (s *Store) HeaderMiddleware(Name string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var token = r.Header.Get(Name)
		_, err := s.InstallTokenToRequest(r, token)
		if err != nil {
			panic(err)
		}
		next(w, r)
		err = s.Save(r)
		if err != nil {
			panic(err)
		}
	}
}

//GetRequestTokenData get stored  token data from request.
//Return the stoed token data and any error raised.
func (s *Store) GetRequestTokenData(r *http.Request) (td *TokenData, err error) {
	var ok bool
	t := r.Context().Value(s.TokenContextName)
	if t != nil {
		td, ok = t.(*TokenData)
		if ok == false {
			return td, ErrDataTypeWrong
		}
		return td, nil
	}
	return td, ErrRequestTokenNotFound
}

//MustGetRequestTokenData get stored  token data from request.
//Panic if any error raised.
func (s *Store) MustGetRequestTokenData(r *http.Request) (v *TokenData) {
	v, err := s.GetRequestTokenData(r)
	if err != nil {
		panic(err)
	}
	return v
}

//Save save the request token data.
func (s *Store) Save(r *http.Request) error {
	td, err := s.GetRequestTokenData(r)
	if err != nil {
		return err
	}
	err = td.Save()
	return err
}

//MustRegenerateToken Regenerate the token name and data with give owner,and save to request.
//Panic if any error raised.
func (s *Store) MustRegenerateToken(r *http.Request, owner string) *TokenData {
	v := s.MustGetRequestTokenData(r)
	err := v.RegenerateToken(owner)
	if err != nil {
		panic(err)
	}
	return v
}

//LogoutMiddleware return a middleware clear the token in request.
func (s *Store) LogoutMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		v := s.MustGetRequestTokenData(r)
		v.SetToken("")
		next(w, r)
	}
}
