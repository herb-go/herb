package tokenstore

import (
	"context"
	"net/http"
	"reflect"
	"time"

	"errors"

	"github.com/herb-go/herb/cache"
)

var DefaultTokenSepartor = "-"
var DefaultTokenLength = 256
var DefaultTokenContextName = "token"
var DefaultUpdateActiveInterval = 5 * time.Minute
var DefaultTokenMaxLifetime = 365 * 24 * time.Hour

const DefaultCookieName = "herb-session"
const DefaultCookiePath = "/"

var ErrTokenNotValidated = errors.New("Token not validated.")
var ErrRequestTokenNotFound = errors.New("Request token not found.Did you forget use middleware?")
var ErrMustRegisterPtr = errors.New("Must register struct pointer.")

func GenerateToken(owner string) (token string, err error) {
	t, err := cache.RandMaskedBytes(cache.TokenMask, DefaultTokenLength)
	if err != nil {
		return
	}
	if owner == "" {
		token = string(t)
	} else {
		token = owner + DefaultTokenSepartor + string(t)
	}
	return
}

func New(Cache *cache.Cache, TokenLifetime time.Duration) *Store {
	return &Store{
		Fields:               map[string]TokenField{},
		Cache:                Cache,
		TokenLength:          DefaultTokenLength,
		TokenSepartor:        DefaultTokenSepartor,
		TokenContextName:     DefaultTokenContextName,
		CookieName:           DefaultCookieName,
		CookiePath:           DefaultCookiePath,
		TokenLifetime:        TokenLifetime,
		UpdateActiveInterval: DefaultUpdateActiveInterval,
		TokenMaxLifetime:     DefaultTokenMaxLifetime,
	}
}
func NewWithContextName(Cache *cache.Cache, TokenLifetime time.Duration, ContenxtName string) *Store {
	s := New(Cache, TokenLifetime)
	s.TokenContextName = ContenxtName
	return s
}

type Store struct {
	Fields               map[string]TokenField
	Cache                *cache.Cache
	TokenLength          int
	TokenSepartor        string
	TokenLifetime        time.Duration
	TokenMaxLifetime     time.Duration
	TokenContextName     string
	CookieName           string
	CookiePath           string
	AutoGenerate         bool
	UpdateActiveInterval time.Duration
}

func (s *Store) RegisterField(Key string, v interface{}) (*TokenField, error) {
	tp := reflect.TypeOf(v)
	if tp.Kind() != reflect.Ptr {
		return nil, ErrMustRegisterPtr
	}
	tv := TokenField{
		Key:   Key,
		Type:  tp.Elem(),
		store: s,
	}
	s.Fields[Key] = tv
	return &tv, nil
}
func (s *Store) MustRegisterField(Key string, v interface{}) *TokenField {
	tv, err := s.RegisterField(Key, v)
	if err != nil {
		panic(err)
	}
	return tv
}
func (s *Store) GenerateToken(owner string) (token string, err error) {
	t, err := cache.RandMaskedBytes(cache.TokenMask, DefaultTokenLength)
	if err != nil {
		return
	}
	token = owner + s.TokenSepartor + string(t)
	return
}
func (s *Store) SearchTokensByOwner(owner string) ([]string, error) {
	return s.Cache.SearchByPrefix(owner + s.TokenSepartor)
}
func (s *Store) Assign(token string) *TokenData {
	tv := newTokenData(token, s)
	tv.tokenChanged = true
	return tv
}
func (s *Store) GetTokenData(v *TokenData) error {
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

func (s *Store) SetTokenData(td *TokenData) error {
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
	return s.Cache.SetBytesValue(token, bytes, s.TokenLifetime)
}
func (s *Store) DeleteToken(token string) error {

	return s.Cache.Del(token)
}

func (s *Store) InstallTokenToRequest(r *http.Request, token string) (td *TokenData, err error) {
	td = newTokenData(token, s)
	td.oldToken = token
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
func (s *Store) GetRequestTokenData(r *http.Request) (td *TokenData, err error) {
	var ok bool
	tv := r.Context().Value(s.TokenContextName)
	if tv != nil {
		td, ok = tv.(*TokenData)
		if ok == false {
			return td, ErrDataTypeWrong
		}
		return td, nil
	}
	return td, ErrRequestTokenNotFound
}
func (s *Store) Save(r *http.Request) error {
	tv, err := s.GetRequestTokenData(r)
	if err != nil {
		return err
	}
	err = tv.Save()
	return err
}
func (s *Store) MustGetRequestTokenData(r *http.Request) (v *TokenData) {
	v, err := s.GetRequestTokenData(r)
	if err != nil {
		panic(err)
	}
	return v
}
func (s *Store) MustRegenerateToken(r *http.Request, owner string) *TokenData {
	v := s.MustGetRequestTokenData(r)
	err := v.RegenerateToken(owner)
	if err != nil {
		panic(err)
	}
	return v
}
func (s *Store) LogoutMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		v := s.MustGetRequestTokenData(r)
		v.SetToken("")
		next(w, r)
	}
}
