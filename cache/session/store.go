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

const (
	//StoreModeCookie which store session token in cookie.
	StoreModeCookie = "cookie"
	//StoreModeHeader which store session token in header.
	StoreModeHeader = "header"
)

//Driver store driver.
type Driver interface {
	GetSessionToken(ts *Session) (token string, err error)
	GenerateToken(owner string) (token string, err error)
	Load(v *Session) error
	Save(t *Session, ttl time.Duration) error
	Delete(token string) error
	Close() error
}

//Store Basic token store interface
type Store struct {
	Driver               Driver
	Marshaler            cache.Marshaler
	TokenLifetime        time.Duration //Token initial expired time.Token life time can be update when accessed if UpdateActiveInterval is greater than 0.
	TokenMaxLifetime     time.Duration //Token max life time.Token can't live more than TokenMaxLifetime if TokenMaxLifetime if greater than 0.
	TokenContextName     ContextKey    //Name in request context store the token  data.Default Session is "token".
	CookieName           string        //Cookie name used in CookieMiddleware.Default Session is "herb-session".
	CookiePath           string        //Cookie path used in cookieMiddleware.Default Session is "/".
	CookieSecure         bool          //Cookie secure value used in cookie middleware.
	AutoGenerate         bool          //Whether auto generate token when guset visit.Default Session is false.
	Mode                 string        //Mode used in auto install middleware.
	UpdateActiveInterval time.Duration //The interval between who token active time update.If less than or equal to 0,the token life time will not be refreshed.
	DefaultSessionFlag   Flag          //Default flag when creating session.
}

// New create empty session store.
func New() *Store {
	return &Store{
		TokenContextName:     defaultTokenContextName,
		CookieName:           defaultCookieName,
		CookiePath:           defaultCookiePath,
		UpdateActiveInterval: defaultUpdateActiveInterval,
		TokenMaxLifetime:     defaultTokenMaxLifetime,
	}
}

//Init  init store with given option.
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
		return ErrTokenNotValidated
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
	return nil
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

//RegenerateToken regenerate session token with given prefix.
//Return session and any error if raised.
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

//AutoGenerateMiddleware middleware that auto generate session.
func (s *Store) AutoGenerateMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		var ts, err = s.GetRequestSession(r)
		if err != nil {
			panic(err)
		}
		if ts.token == "" || ts.token == clientStoreNewToken {
			err := ts.RegenerateToken("")
			if err != nil {
				panic(err)
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
		writer := w.(ResponseWriter)
		cw := cookieResponseWriter{
			ResponseWriter: writer,
			r:              r,
			store:          s,
			written:        false,
		}
		next(&cw, r)
	}
}

//HeaderMiddleware return a Middleware which install the token which special by Header with given name.
//This middleware will save token after request finished if the token changed.
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

//InstallMiddleware middleware which auto install session depand on store mode.
//Cookie middleware will be installed if no valid store mode given.
func (s *Store) InstallMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	switch s.Mode {
	case StoreModeHeader:
		return s.HeaderMiddleware(s.CookieName)
	}
	return s.CookieMiddleware()
}

//GetRequestSession get stored  token data from request.
//Return the stored token data and any error raised.
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

//MustRegenerateRequsetSession rregenerate request session with given prefix.
//Return session regenerated.
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

//Set set value to request with given field name.
//Return any error if raised.
func (s *Store) Set(r *http.Request, fieldName string, v interface{}) (err error) {
	ts, err := s.GetRequestSession(r)
	if err != nil {
		return
	}
	return ts.Set(fieldName, v)
}

//Get get value form request with given field name.
//Return any error if raised.
func (s *Store) Get(r *http.Request, fieldName string, v interface{}) (err error) {
	ts, err := s.GetRequestSession(r)
	if err != nil {
		return
	}
	return ts.Get(fieldName, v)
}

//Del delete value from request with given field name.
//Return any error if raised.
func (s *Store) Del(r *http.Request, fieldName string) (err error) {
	ts, err := s.GetRequestSession(r)
	if err != nil {
		return
	}
	return ts.Del(fieldName)
}

//ExpiredAt get session expired timestamp from rerquest.
//Return  expired timestamp and any error if raised.
func (s *Store) ExpiredAt(r *http.Request) (ExpiredAt int64, err error) {
	ts, err := s.GetRequestSession(r)
	if err != nil {
		return
	}
	return ts.ExpiredAt, nil
}

//Field create store field with given name.
//Return field created.
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
		v, err := s.GetRequestSession(r)
		if err != nil {
			panic(err)
		}
		v.SetToken("")
		next(w, r)
	}
}

//IsNotFoundError return if given error if a not found error.
func (s Store) IsNotFoundError(err error) bool {
	return err == ErrDataNotFound
}
