package tokenstore

import "errors"
import "reflect"
import "github.com/herb-go/herb/cache"
import "net/http"

var ErrDataNotFound = errors.New("Data not found")
var ErrDataTypeWrong = errors.New("Data type wrong")
var ErrNilPoint = errors.New("Data point to nil")
var ErrDataTypeNotRegister = errors.New("Data type not register")

type TokenValues struct {
	data         map[string][]byte
	cache        map[string]interface{}
	token        string
	oldToken     string
	loaded       bool
	tokenChanged bool
	updated      bool
	store        *Store
}

func (t *TokenValues) SetToken(newToken string) {
	t.token = newToken
	t.tokenChanged = true
	t.updated = true
}
func (t *TokenValues) RegenerateToken(owner string) error {
	token, err := t.store.GenerateToken(owner)
	if err != nil {
		return err
	}
	t.data = map[string][]byte{}
	t.cache = map[string]interface{}{}
	t.SetToken(token)
	return nil
}

func (t *TokenValues) Load() error {
	if t.token == "" {
		return ErrTokenNotValidated
	}
	if t.loaded {
		return nil
	}
	err := t.store.GetTokenValues(t)
	if err == cache.ErrNotFound {
		if t.tokenChanged == false {
			return ErrDataNotFound
		} else {
			err = nil
		}
	}
	if err != nil {
		return err
	}
	return nil
}
func (t *TokenValues) Marshal() ([]byte, error) {
	return cache.MarshalMsgpack(t.data)
}
func (t *TokenValues) Unmarshal(token string, bytes []byte) error {
	t.token = token
	t.cache = map[string]interface{}{}
	return cache.UnmarshalMsgpack(bytes, &(t.data))
}

type TokenValue struct {
	Key   string
	Type  reflect.Type
	store *Store
}

func (f *TokenValue) GetTokenValuesData(m *TokenValues, v interface{}) (err error) {
	if m.token == "" {
		err = ErrDataNotFound
		return
	}
	err = m.Load()
	if err != nil {
		return
	}
	key := f.Key
	typ := reflect.TypeOf(v)
	if typ.Elem() != f.Type {
		return ErrDataTypeWrong
	}
	if reflect.ValueOf(v).IsNil() {
		return ErrNilPoint
	}
	c, ok := m.cache[key]
	if ok == true {
		dst := reflect.ValueOf(v).Elem()
		src := reflect.ValueOf(c).Elem()
		dst.Set(src)
	}
	data, ok := m.data[key]
	if ok == false {
		return ErrDataNotFound
	}
	err = cache.UnmarshalMsgpack(data, v)
	if err == nil {
		m.cache[key] = v
	}
	return
}
func (f *TokenValue) Get(r *http.Request, v interface{}) error {
	var m, err = f.store.GetRequestTokenValues(r)
	if err != nil {
		return err
	}
	return f.GetTokenValuesData(m, v)
}

func (f *TokenValue) SetTokenValuesData(m *TokenValues, v interface{}) (err error) {
	if m.token == "" {
		err = ErrDataNotFound
		return
	}
	err = m.Load()
	if err != nil {
		return
	}

	key := f.Key
	if reflect.TypeOf(v) != f.Type {
		return ErrDataTypeWrong
	}
	m.cache[key] = v
	bytes, err := cache.MarshalMsgpack(v)
	if err != nil {
		return
	}
	m.data[key] = bytes
	err = nil
	m.updated = true
	return
}

func (f *TokenValue) Set(r *http.Request, v interface{}) error {
	var m, err = f.store.GetRequestTokenValues(r)
	if err != nil {
		return err
	}
	err = f.SetTokenValuesData(m, v)
	return err
}
