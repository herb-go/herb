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
		Values:               map[string]TokenValue{},
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
	Values               map[string]TokenValue
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

func (s *Store) RegisterField(Key string, v interface{}) (*TokenValue, error) {
	tp := reflect.TypeOf(v)
	if tp.Kind() != reflect.Ptr {
		return nil, ErrMustRegisterPtr
	}
	tv := TokenValue{
		Key:   Key,
		Type:  tp.Elem(),
		store: s,
	}
	s.Values[Key] = tv
	return &tv, nil
}
func (s *Store) MustRegisterField(Key string, v interface{}) *TokenValue {
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
func (s *Store) AssignTokenValues(token string) *TokenValues {
	tv := newTokenValues(token, s)
	t := time.Now().Unix()
	tv.CreatedTime = t
	tv.LastActiveTime = t
	tv.tokenChanged = true
	return tv
}
func (s *Store) GetTokenValues(v *TokenValues) error {
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

func (s *Store) SetTokenValues(v *TokenValues) error {
	token := v.token
	if token == "" {
		return cache.ErrKeyUnavailable
	}
	if s.TokenLifetime >= 0 {
		v.ExpiredAt = time.Now().Add(s.TokenLifetime).Unix()
	} else {
		v.ExpiredAt = -1
	}
	bytes, err := v.Marshal()
	if err != nil {
		return err
	}
	return s.Cache.SetBytesValue(token, bytes, s.TokenLifetime)
}
func (s *Store) DeleteToken(token string) error {

	return s.Cache.Del(token)
}

func (s *Store) InstallTokenToRequest(r *http.Request, token string) (v *TokenValues, err error) {
	v = newTokenValues(token, s)
	v.oldToken = token
	if token == "" && s.AutoGenerate == true {
		err = v.RegenerateToken("")
		if err != nil {
			return
		}
	}

	ctx := context.WithValue(r.Context(), s.TokenContextName, v)
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
func (s *Store) GetRequestTokenValues(r *http.Request) (v *TokenValues, err error) {
	var ok bool
	tv := r.Context().Value(s.TokenContextName)
	if tv != nil {
		v, ok = tv.(*TokenValues)
		if ok == false {
			return v, ErrDataTypeWrong
		}
		return v, nil
	}
	return v, ErrRequestTokenNotFound
}
func (s *Store) Save(r *http.Request) error {
	tv, err := s.GetRequestTokenValues(r)
	if err != nil {
		return err
	}
	err = tv.Save()
	return err
}
func (s *Store) MustGetRequestTokenValues(r *http.Request) (v *TokenValues) {
	v, err := s.GetRequestTokenValues(r)
	if err != nil {
		panic(err)
	}
	return v
}
func (s *Store) MustRegenerateToken(r *http.Request, owner string) {
	v := s.MustGetRequestTokenValues(r)
	err := v.RegenerateToken(owner)
	if err != nil {
		panic(err)
	}
}
func (s *Store) LogoutMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		v := s.MustGetRequestTokenValues(r)
		v.SetToken("")
		next(w, r)
	}
}
