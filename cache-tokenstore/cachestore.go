//Package tokenstore is used to store user data in token based cache.
//It is normally used as user session or token.
//This package depands github.com/herb-go/herb/cache.
package tokenstore

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"github.com/herb-go/herb/cache"
)

var defaultUpdateActiveInterval = 5 * time.Minute

var defaultTokenMaxLifetime = 365 * 24 * time.Hour
var (
	defaultCookieName = "herb-session"
	defaultCookiePath = "/"
)

func defaultTokenGenerater(s *CacheStore, owner string) (token string, err error) {
	t, err := cache.RandMaskedBytes(cache.TokenMask, 256)
	if err != nil {
		return
	}

	token = owner + "-" + string(t)
	return
}

//New create a new token store with given cache and token lifetime.
//Cache is the cache which dates stored in.
//TokenLifeTime is the token initial expired tome.
//Return a new token store.
//All other property of the store can be set after creation.
func New(Cache *cache.Cache, TokenLifetime time.Duration) *CacheStore {
	return &CacheStore{
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

//CacheStore CacheStore is the stuct store token data in cache.
type CacheStore struct {
	Fields               map[string]TokenField                                       //All registered field
	Cache                *cache.Cache                                                //Cache which stores token data
	TokenGenerater       func(s *CacheStore, owner string) (token string, err error) //Token name generate func
	TokenLifetime        time.Duration                                               //Token initial expired time.Token life time can be update when accessed if UpdateActiveInterval is greater than 0.
	TokenMaxLifetime     time.Duration                                               //Token max life time.Token can't live more than TokenMaxLifetime if TokenMaxLifetime if greater than 0.
	TokenContextName     ContextKey                                                  //Name in request context store the token  data.Default value is "token".
	CookieName           string                                                      //Cookie name used in CookieMiddleware.Default value is "herb-session".
	CookiePath           string                                                      //Cookie path used in cookieMiddleware.Default value is "/".
	AutoGenerate         bool                                                        //Whether auto generate token when guset visit.Default value is false.
	UpdateActiveInterval time.Duration                                               //The interval between who token active time update.If less than or equal to 0,the token life time will not be refreshed.
}

//Close Close cachestore and return any error if raised
func (s *CacheStore) Close() error {
	return s.Cache.Close()
}

//RegisterField registe filed to store.
//registered field can be used directly with request to load or save the token value.
//Parameter Key filed name.
//Parameter v should be pointer to empty data model which data filled in.
//Return a new Token field and error.
func (s *CacheStore) RegisterField(Key string, v interface{}) (*TokenField, error) {
	if v == nil {
		return nil, ErrNilPointer
	}
	t := reflect.TypeOf(v)
	if t.Kind() != reflect.Ptr {
		return nil, ErrMustRegistePtr
	}
	tp := reflect.ValueOf(v).Elem().Type()
	tf := TokenField{
		Key:   Key,
		Type:  tp,
		Store: s,
	}
	s.Fields[Key] = tf
	return &tf, nil
}

//GenerateToken generate new token name with given owner.
//Return the new token name and error.
func (s *CacheStore) GenerateToken(owner string) (token string, err error) {
	return s.TokenGenerater(s, owner)
}

//GenerateTokenData generate new token data with given token.
//Return a new TokenData.
func (s *CacheStore) GenerateTokenData(token string) (td *TokenData, err error) {
	td = NewTokenData(token, s)
	td.tokenChanged = true
	return td, nil
}

//LoadTokenData Load TokenData form the TokenData.token.
//Return any error if raised
func (s *CacheStore) LoadTokenData(v *TokenData) error {
	token := v.token
	if token == "" {
		return cache.ErrKeyUnavailable
	}

	bytes, err := s.Cache.GetBytesValue(token)
	if err == cache.ErrNotFound {
		err = ErrDataNotFound
	}
	if err != nil {
		return err
	}

	err = v.Unmarshal(token, bytes)
	if err == nil {
		v.token = token
		v.store = s
	}
	if v.ExpiredAt > 0 && v.ExpiredAt < time.Now().Unix() {
		return ErrDataNotFound
	}
	if s.TokenMaxLifetime > 0 && time.Unix(v.CreatedTime, 0).Add(s.TokenMaxLifetime).Before(time.Now()) {
		return ErrDataNotFound
	}
	return err
}

//SaveTokenData Save tokendata if necessary.
//Return any error raised.
func (s *CacheStore) SaveTokenData(t *TokenData) error {
	if s.UpdateActiveInterval > 0 {
		nextUpdateTime := time.Unix(t.LastActiveTime, 0).Add(s.UpdateActiveInterval)
		if nextUpdateTime.Before(time.Now()) {
			t.LastActiveTime = time.Now().Unix()
			t.updated = true
		}
	}
	if t.updated && t.token != "" {
		err := s.save(t)
		if err != nil {
			return err
		}
		t.updated = false
	}
	if t.tokenChanged && t.oldToken != "" {
		err := s.DeleteToken(t.oldToken)
		if err != nil {
			return err
		}
	}
	t.oldToken = t.token
	return nil
}
func (s *CacheStore) save(td *TokenData) error {
	var err error
	token := td.token
	if token == "" {
		return cache.ErrKeyUnavailable
	}
	if td.ExpiredAt > 0 && td.ExpiredAt < time.Now().Unix() {
		return nil
	}
	if s.TokenMaxLifetime > 0 && time.Unix(td.CreatedTime, 0).Add(s.TokenMaxLifetime).Before(time.Now()) {
		return nil
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
func (s *CacheStore) DeleteToken(token string) error {

	return s.Cache.Del(token)
}

//GetTokenData get the token data with give name .
//Return the TokenData
func (s *CacheStore) GetTokenData(token string) (td *TokenData) {
	td = NewTokenData(token, s)
	td.oldToken = token
	return
}

//GetTokenDataToken Get the token string from token data.
func (s *CacheStore) GetTokenDataToken(td *TokenData) (token string, err error) {
	return td.token, nil
}

//InstallTokenToRequest install the give token to request.
//Tokendata will be stored in request context which named by TokenContextName of store.
//You should use this func when use your own token binding func instead of CookieMiddleware or HeaderMiddleware
//Return the loaded TokenData and any error raised.
func (s *CacheStore) InstallTokenToRequest(r *http.Request, token string) (td *TokenData, err error) {
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
func (s *CacheStore) CookieMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
		cw := cookieResponseWriter{
			ResponseWriter: w,
			r:              r,
			store:          s,
			written:        false,
		}
		next(&cw, r)
	}
}

//HeaderMiddleware return a Middleware which install the token which special by Header with given name.
//this middleware will save token after request finished if the token changed.
func (s *CacheStore) HeaderMiddleware(Name string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var token = r.Header.Get(Name)
		_, err := s.InstallTokenToRequest(r, token)
		if err != nil {
			panic(err)
		}
		next(w, r)
		err = s.SaveRequestTokenData(r)
		if err != nil {
			panic(err)
		}
	}
}

//GetRequestTokenData get stored  token data from request.
//Return the stoed token data and any error raised.
func (s *CacheStore) GetRequestTokenData(r *http.Request) (td *TokenData, err error) {
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

//SaveRequestTokenData save the request token data.
func (s *CacheStore) SaveRequestTokenData(r *http.Request) error {
	td, err := s.GetRequestTokenData(r)
	if err != nil {
		return err
	}
	err = td.Save()
	return err
}

//LogoutMiddleware return a middleware clear the token in request.
func (s *CacheStore) LogoutMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		v := MustGetRequestTokenData(s, r)
		v.SetToken("")
		next(w, r)
	}
}
