package tokenstore

import (
	"context"
	"crypto/rand"
	"net/http"
	"reflect"
	"time"

	"github.com/jarlyyn/go-utils/security"
)

const clientStoreNonceSize = 4
const clientStoreNewToken = "."

type ClientStore struct {
	Fields               map[string]TokenField //All registered field
	TokenLifetime        time.Duration         //Token initial expired time.Token life time can be update when accessed if UpdateActiveInterval is greater than 0.
	TokenMaxLifetime     time.Duration         //Token max life time.Token can't live more than TokenMaxLifetime if TokenMaxLifetime if greater than 0.
	TokenContextName     ContextKey            //Name in request context store the token  data.Default value is "token".
	CookieName           string                //Cookie name used in CookieMiddleware.Default value is "herb-session".
	CookiePath           string                //Cookie path used in cookieMiddleware.Default value is "/".
	AutoGenerate         bool                  //Whether auto generate token when guset visit.Default value is false.
	Key                  []byte
	UpdateActiveInterval time.Duration //The interval between who token active time update.If less than or equal to 0,the token life time will not be refreshed.
}

func NewClientStore(key []byte, TokenLifetime time.Duration) *ClientStore {
	return &ClientStore{
		Fields:               map[string]TokenField{},
		TokenContextName:     defaultTokenContextName,
		CookieName:           defaultCookieName,
		CookiePath:           defaultCookiePath,
		TokenLifetime:        TokenLifetime,
		UpdateActiveInterval: defaultUpdateActiveInterval,
		TokenMaxLifetime:     defaultTokenMaxLifetime,
	}
}

func (s *ClientStore) GetTokenData(token string) (td *TokenData, err error) {
	td = NewTokenData(token, s)
	err = s.LoadTokenData(td)
	return
}
func (s *ClientStore) GetTokenDataToken(td *TokenData) (token string, err error) {
	err = td.Save()
	return td.token, err
}
func (s *ClientStore) GetRequestTokenData(r *http.Request) (td *TokenData, err error) {
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
func (s *ClientStore) GenerateToken(owner string) (token string, err error) {
	return clientStoreNewToken, nil

}
func (s *ClientStore) GenerateTokenData(token string) (td *TokenData, err error) {
	td = NewTokenData(token, s)
	td.tokenChanged = true
	return td, nil

}
func (s *ClientStore) LoadTokenData(v *TokenData) (err error) {
	if v.token == clientStoreNewToken {
		return
	}
	v.notFound = true
	var data []byte
	data, err = security.AESNonceDecryptBase64(v.token, s.Key)
	if err != nil {
		return ErrDataNotFound
	}
	err = v.Unmarshal(v.token, data)
	if err != nil {
		return ErrDataNotFound
	}
	if v.ExpiredAt > 0 && v.ExpiredAt < time.Now().Unix() {
		return ErrDataNotFound
	}
	if s.TokenMaxLifetime > 0 && time.Unix(v.CreatedTime, 0).Add(s.TokenMaxLifetime).Before(time.Now()) {
		return ErrDataNotFound
	}
	v.notFound = false
	return
}
func (s *ClientStore) SaveTokenData(t *TokenData) (err error) {
	if t.token == clientStoreNewToken {
		t.updated = true
	}
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
	return nil
}
func (s *ClientStore) save(t *TokenData) (err error) {
	var data []byte
	t.Nonce = make([]byte, clientStoreNonceSize)
	_, err = rand.Read(t.Nonce)
	if err != nil {
		return
	}
	if t.ExpiredAt > 0 && t.ExpiredAt < time.Now().Unix() {
		return nil
	}
	if s.TokenMaxLifetime > 0 && time.Unix(t.CreatedTime, 0).Add(s.TokenMaxLifetime).Before(time.Now()) {
		return nil
	}
	if s.TokenLifetime >= 0 {
		t.ExpiredAt = time.Now().Add(s.TokenLifetime).Unix()
	} else {
		t.ExpiredAt = -1
	}
	data, err = t.Marshal()
	if err != nil {
		return
	}
	t.token, err = security.AESNonceEncryptBase64(data, s.Key)
	t.tokenChanged = true
	return
}
func (s *ClientStore) RegisterField(Key string, v interface{}) (*TokenField, error) {
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
func (s *ClientStore) InstallTokenToRequest(r *http.Request, token string) (td *TokenData, err error) {
	td, err = s.GetTokenData(token)
	if err == ErrDataNotFound {
		err = td.RegenerateToken("")
	}
	if err != nil {
		return
	}
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
func (s *ClientStore) CookieMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
		cw := ClientStoreResponseWriter{
			ResponseWriter: w,
			r:              r,
			store:          s,
			written:        false,
		}
		next(&cw, r)
		err = s.SaveRequestTokenData(r)
		if err != nil {
			panic(err)
		}
	}
}

//SaveRequestTokenData save the request token data.
func (s *ClientStore) SaveRequestTokenData(r *http.Request) error {
	td, err := s.GetRequestTokenData(r)
	if err != nil {
		return err
	}
	err = td.Save()
	return err
}
func (s *ClientStore) HeaderMiddleware(Name string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
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
func (s *ClientStore) LogoutMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
		v := MustGetRequestTokenData(s, r)
		v.SetToken("")
		next(w, r)
	}
}
func (s *ClientStore) Close() error {
	return nil
}

type ClientStoreResponseWriter struct {
	http.ResponseWriter
	r       *http.Request
	store   *ClientStore
	written bool
}

func (w *ClientStoreResponseWriter) WriteHeader(status int) {
	var td *TokenData
	var err error
	if w.written == false {
		w.written = true
		td, err = w.store.GetRequestTokenData(w.r)
		if err != nil {
			panic(err)
		}
		if td.tokenChanged {
			cookie := &http.Cookie{
				Name:     w.store.CookieName,
				Value:    td.token,
				Path:     w.store.CookiePath,
				Secure:   false,
				HttpOnly: true,
			}
			if w.store.TokenLifetime >= 0 {
				cookie.Expires = time.Now().Add(w.store.TokenLifetime)
			}
			http.SetCookie(w, cookie)
		}
	}
	w.ResponseWriter.WriteHeader(status)
}

func (w *ClientStoreResponseWriter) Write(data []byte) (int, error) {
	if w.written == false {
		w.WriteHeader(http.StatusOK)
	}
	return w.ResponseWriter.Write(data)
}