package session

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/herb-go/herb/cache"
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

//MustRegisterField registe filed to store.
//registered field can be used directly with request to load or save the token Session.
//Parameter Key filed name.
//Parameter v should be pointer to empty data model which data filled in.
//Return a new Token field.
//Panic if any error raised.
// func MustRegisterField(s Store, Key string, v interface{}) *TokenField {
// 	tf, err := s.RegisterField(Key, v)
// 	if err != nil {
// 		panic(err)
// 	}
// 	return tf
// }

//Store Basic token storce interface
type Driver interface {
	GetSessionToken(ts *Session) (token string, err error)
	GenerateToken(owner string) (token string, err error)
	Load(v *Session) error
	Save(t *Session, ttl time.Duration) error
	Delete(token string) error
	Close() error
}
type Store struct {
	Driver               Driver
	TokenLifetime        time.Duration //Token initial expired time.Token life time can be update when accessed if UpdateActiveInterval is greater than 0.
	TokenMaxLifetime     time.Duration //Token max life time.Token can't live more than TokenMaxLifetime if TokenMaxLifetime if greater than 0.
	TokenContextName     ContextKey    //Name in request context store the token  data.Default Session is "token".
	CookieName           string        //Cookie name used in CookieMiddleware.Default Session is "herb-session".
	CookiePath           string        //Cookie path used in cookieMiddleware.Default Session is "/".
	CookieSecure         bool          //Cookie secure value used in cookie middleware.
	AutoGenerate         bool          //Whether auto generate token when guset visit.Default Session is false.
	UpdateActiveInterval time.Duration //The interval between who token active time update.If less than or equal to 0,the token life time will not be refreshed.
	DefaultSessionFlag   Flag          //Default flag when creating session.
}

func New() *Store {
	return &Store{
		TokenContextName:     defaultTokenContextName,
		CookieName:           defaultCookieName,
		CookiePath:           defaultCookiePath,
		UpdateActiveInterval: defaultUpdateActiveInterval,
		TokenMaxLifetime:     defaultTokenMaxLifetime,
	}
}

func (s *Store) Init(option Option) error {
	return option.ApplyTo(s)
}

//Close Close cachestore and return any error if raised
func (s *Store) Close() error {
	return s.Driver.Close()
}

//GenerateToken generate new token name with given prefix.
//Return the new token name and error.
func (s *Store) GenerateToken(prefix string) (token string, err error) {
	return s.Driver.GenerateToken(prefix)
}

//GenerateSession generate new token data with given token.
//Return a new Session and error.
func (s *Store) GenerateSession(token string) (ts *Session, err error) {
	ts = NewSession(token, s)
	ts.tokenChanged = true
	return ts, nil
}

//LoadSession Load Session form the Session.token.
//Return any error if raised
func (s *Store) LoadSession(v *Session) error {
	token := v.token
	if token == "" {
		return cache.ErrKeyUnavailable
	}
	err := s.Driver.Load(v)
	if err != nil {
		return err
	}
	v.Store = s
	if v.ExpiredAt > 0 && v.ExpiredAt < time.Now().Unix() {
		return ErrDataNotFound
	}
	if s.TokenMaxLifetime > 0 && time.Unix(v.CreatedTime, 0).Add(s.TokenMaxLifetime).Before(time.Now()) {
		return ErrDataNotFound
	}
	return err
}

//SaveSession Save Session if necessary.
//Return any error raised.
func (s *Store) SaveSession(t *Session) error {
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
func (s *Store) save(ts *Session) error {
	token := ts.token
	if token == "" {
		return cache.ErrKeyUnavailable
	}
	if ts.ExpiredAt > 0 && ts.ExpiredAt < time.Now().Unix() {
		return nil
	}
	if s.TokenMaxLifetime > 0 && time.Unix(ts.CreatedTime, 0).Add(s.TokenMaxLifetime).Before(time.Now()) {
		return nil
	}
	if s.TokenLifetime >= 0 {
		ts.ExpiredAt = time.Now().Add(s.TokenLifetime).Unix()
	} else {
		ts.ExpiredAt = -1
	}

	err := s.Driver.Save(ts, s.TokenLifetime)
	if err != nil {
		return err
	}
	ts.loaded = true
	return nil
}

//DeleteToken delete the token with given name.
//Return any error if raised.
func (s *Store) DeleteToken(token string) error {

	return s.Driver.Delete(token)
}

//GetSession get the token data with give token .
//Return the Session
func (s *Store) GetSession(token string) (ts *Session) {
	ts = NewSession(token, s)
	ts.oldToken = token
	return
}

//GetSessionToken Get the token string from token data.
//Return token and any error raised.
func (s *Store) GetSessionToken(ts *Session) (token string, err error) {
	return s.Driver.GetSessionToken(ts)
}
func (s *Store) MustGetSessionToken(ts *Session) (token string) {
	token, err := s.Driver.GetSessionToken(ts)
	if err != nil {
		panic(err)
	}
	return
}
func (s *Store) RegenerateToken(prefix string) (ts *Session, err error) {
	ts = NewSession("", s)
	err = ts.RegenerateToken(prefix)
	return
}

//Install install the give token to request.
//Session will be stored in request context which named by TokenContextName of store.
//You should use this func when use your own token binding func instead of CookieMiddleware or HeaderMiddleware
//Return the loaded Session and any error raised.
func (s *Store) Install(r *http.Request, token string) (ts *Session, err error) {
	ts = s.GetSession(token)

	if (token == "" || token == clientStoreNewToken) && s.AutoGenerate == true {
		err = ts.RegenerateToken("")
		if err != nil {
			return
		}
	}

	ctx := context.WithValue(r.Context(), s.TokenContextName, ts)
	*r = *r.WithContext(ctx)
	return
}

func (s *Store) AutoGenerateMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var ts = s.MustGetRequestSession(r)
		if ts.token == "" || ts.token == clientStoreNewToken {
			err := ts.RegenerateToken("")
			if err != nil {
				return
			}
			ctx := context.WithValue(r.Context(), s.TokenContextName, ts)
			*r = *r.WithContext(ctx)
		}
		next(w, r)
	}
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
		_, err = s.Install(r, token)
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
func (s *Store) HeaderMiddleware(Name string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var token = r.Header.Get(Name)
		_, err := s.Install(r, token)
		if err != nil {
			panic(err)
		}
		next(w, r)
		err = s.SaveRequestSession(r)
		if err != nil {
			panic(err)
		}
	}
}

//GetRequestSession get stored  token data from request.
//Return the stoed token data and any error raised.
func (s *Store) GetRequestSession(r *http.Request) (ts *Session, err error) {
	var ok bool
	t := r.Context().Value(s.TokenContextName)
	if t != nil {
		ts, ok = t.(*Session)
		if ok == false {
			return ts, ErrDataTypeWrong
		}
		return ts, nil
	}
	return ts, ErrRequestTokenNotFound
}
func (s *Store) MustRegenerateRequsetSession(r *http.Request, prefix string) (ts *Session) {
	var err error
	ts, err = s.GetRequestSession(r)
	if err != nil {
		panic(err)
	}
	err = ts.RegenerateToken(prefix)
	if err != nil {
		panic(err)
	}
	return ts
}
func (s *Store) Set(r *http.Request, fieldName string, v interface{}) (err error) {
	ts, err := s.GetRequestSession(r)
	if err != nil {
		return
	}
	return ts.Set(fieldName, v)
}

func (s *Store) Get(r *http.Request, fieldName string, v interface{}) (err error) {
	ts, err := s.GetRequestSession(r)
	if err != nil {
		return
	}
	return ts.Get(fieldName, v)
}

func (s *Store) Del(r *http.Request, fieldName string) (err error) {
	ts, err := s.GetRequestSession(r)
	if err != nil {
		return
	}
	return ts.Del(fieldName)
}
func (s *Store) ExpiredAt(r *http.Request) (ExpiredAt int64, err error) {
	ts, err := s.GetRequestSession(r)
	if err != nil {
		return
	}
	return ts.ExpiredAt, nil
}

func (s *Store) Field(name string) *Field {
	return &Field{Name: name, Store: s}
}

//SaveRequestSession save the request token data.
func (s *Store) SaveRequestSession(r *http.Request) error {
	ts, err := s.GetRequestSession(r)
	if err != nil {
		return err
	}
	err = ts.Save()
	return err
}

//DestoryMiddleware return a middleware clear the token in request.
func (s *Store) DestoryMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		v := s.MustGetRequestSession(r)
		v.SetToken("")
		next(w, r)
	}
}

//MustGetRequestSession get stored  token data from request.
//Panic if any error raised.
func (s Store) MustGetRequestSession(r *http.Request) (v *Session) {
	v, err := s.GetRequestSession(r)
	if err != nil {
		panic(err)
	}
	return v
}

//MustRegenerateRequestToken Regenerate the token name and data with give owner,and save to request.
//Panic if any error raised.
func (s Store) MustRegenerateRequestToken(r *http.Request, owner string) *Session {
	v := s.MustGetRequestSession(r)
	err := v.RegenerateToken(owner)
	if err != nil {
		panic(err)
	}
	return v
}

func (s Store) IsNotFound(err error) bool {
	return err == ErrDataNotFound
}
