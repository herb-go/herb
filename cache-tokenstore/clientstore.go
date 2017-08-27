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

//AESTokenMarshaler token marshaler which crypt data with AES
//Return error if raised
func AESTokenMarshaler(s *ClientStore, td *TokenData) (err error) {
	var data []byte

	td.Nonce = make([]byte, clientStoreNonceSize)
	_, err = rand.Read(td.Nonce)
	if err != nil {
		return
	}
	data, err = td.Marshal()
	if err != nil {
		return
	}
	td.token, err = security.AESNonceEncryptBase64(data, s.Key)
	return
}

//AESTokenUnmarshaler token unmarshaler which crypt data with AES
//Return error if raised
func AESTokenUnmarshaler(s *ClientStore, v *TokenData) (err error) {
	var data []byte
	data, err = security.AESNonceDecryptBase64(v.token, s.Key)
	if err != nil {
		return ErrDataNotFound
	}
	err = v.Unmarshal(v.token, data)
	if err != nil {
		return ErrDataNotFound
	}
	return nil
}

type ClientStore struct {
	Fields               map[string]TokenField //All registered field
	TokenLifetime        time.Duration         //Token initial expired time.Token life time can be update when accessed if UpdateActiveInterval is greater than 0.
	TokenMaxLifetime     time.Duration         //Token max life time.Token can't live more than TokenMaxLifetime if TokenMaxLifetime if greater than 0.
	TokenContextName     ContextKey            //Name in request context store the token  data.Default value is "token".
	CookieName           string                //Cookie name used in CookieMiddleware.Default value is "herb-session".
	CookiePath           string                //Cookie path used in cookieMiddleware.Default value is "/".
	AutoGenerate         bool                  //Whether auto generate token when guset visit.Default value is false.
	Key                  []byte                //Crypt key
	UpdateActiveInterval time.Duration         //The interval between who token active time update.If less than or equal to 0,the token life time will not be refreshed.
	TokenMarshaler       func(*ClientStore, *TokenData) error
	TokenUnmarshaler     func(*ClientStore, *TokenData) error
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
		TokenMarshaler:       AESTokenMarshaler,
		TokenUnmarshaler:     AESTokenUnmarshaler,
	}
}

func (s *ClientStore) GetTokenData(token string) (td *TokenData) {
	td = NewTokenData(token, s)
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
func (s *ClientStore) GenerateToken(prefix string) (token string, err error) {
	return clientStoreNewToken, nil

}
func (s *ClientStore) GenerateTokenData(token string) (td *TokenData, err error) {
	td = NewTokenData(token, s)
	td.tokenChanged = true
	return td, nil
}
func (s *ClientStore) SearchByPrefix(prefix string) (Tokens []string, err error) {
	return nil, ErrFeatureNotSupported
}
func (s *ClientStore) LoadTokenData(v *TokenData) (err error) {
	if v.token == clientStoreNewToken {
		return
	}
	v.notFound = true
	err = s.TokenUnmarshaler(s, v)
	if err != nil {
		return err
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
	t.Load()
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

func (s *ClientStore) save(td *TokenData) (err error) {

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
	err = s.TokenMarshaler(s, td)
	if err != nil {
		return
	}
	td.tokenChanged = true
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
	td = s.GetTokenData(token)
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
		cw := clientStoreResponseWriter{
			ResponseWriter: w,
			r:              r,
			store:          s,
			written:        false,
		}
		next(&cw, r)
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
