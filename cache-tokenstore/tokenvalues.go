package tokenstore

import "errors"
import "reflect"
import "github.com/herb-go/herb/cache"
import "net/http"
import "sync"
import "time"

var ErrDataNotFound = errors.New("Data not found")
var ErrDataTypeWrong = errors.New("Data type wrong")
var ErrNilPoint = errors.New("Data point to nil")
var ErrDataTypeNotRegister = errors.New("Data type not register")

type TokenValues struct {
	data           map[string][]byte
	ExpiredAt      int64
	CreatedTime    int64
	LastActiveTime int64
	cache          map[string]interface{}
	token          string
	oldToken       string
	loaded         bool
	tokenChanged   bool
	updated        bool
	store          *Store
	Mutex          *sync.RWMutex
}
type tokenValuesData struct {
	Data           map[string][]byte
	CreatedTime    int64
	LastActiveTime int64
	ExpiredAt      int64
}

func newTokenValues(token string, s *Store) *TokenValues {
	t := time.Now().Unix()
	return &TokenValues{
		token:          token,
		data:           map[string][]byte{},
		cache:          map[string]interface{}{},
		store:          s,
		tokenChanged:   false,
		Mutex:          &sync.RWMutex{},
		CreatedTime:    t,
		LastActiveTime: t,
		ExpiredAt:      -1,
	}

}
func (t *TokenValues) Token() string {
	return t.token
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
func (t *TokenValues) Save() error {
	nextUpdateTime := time.Unix(t.LastActiveTime, 0).Add(t.store.UpdateActiveInterval)
	if nextUpdateTime.Before(time.Now()) {
		t.LastActiveTime = time.Now().Unix()
		t.updated = true
	}
	if t.updated && t.token != "" {
		err := t.store.SetTokenValues(t)
		if err != nil {
			return err
		}
	}
	if t.tokenChanged && t.oldToken != "" {
		err := t.store.DeleteToken(t.oldToken)
		if err != nil {
			return err
		}
	}
	return nil
}
func (t *TokenValues) Marshal() ([]byte, error) {
	return cache.MarshalMsgpack(
		tokenValuesData{
			Data:           t.data,
			ExpiredAt:      t.ExpiredAt,
			CreatedTime:    t.CreatedTime,
			LastActiveTime: t.LastActiveTime,
		})
}
func (t *TokenValues) Unmarshal(token string, bytes []byte) error {
	var err error
	var Data = tokenValuesData{}
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	t.token = token
	t.cache = map[string]interface{}{}
	err = cache.UnmarshalMsgpack(bytes, &(Data))
	if err != nil {
		return err
	}
	t.data = Data.Data
	t.ExpiredAt = Data.ExpiredAt
	t.CreatedTime = Data.CreatedTime
	t.LastActiveTime = Data.LastActiveTime
	return nil
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
	m.Mutex.RLock()
	defer m.Mutex.RUnlock()
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

func (f *TokenValue) RwMutex(r *http.Request) (*sync.RWMutex, error) {
	var m, err = f.store.GetRequestTokenValues(r)
	if err != nil {
		return nil, err
	}
	return m.Mutex, nil
}
func (f *TokenValue) ExpiredAt(r *http.Request) (int64, error) {
	var m, err = f.store.GetRequestTokenValues(r)
	if err != nil {
		return 0, err
	}
	return m.ExpiredAt, nil
}
func (f *TokenValue) GetToken(r *http.Request) (string, error) {
	var m, err = f.store.GetRequestTokenValues(r)
	if err != nil {
		return "", err
	}
	return m.token, nil
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
	m.Mutex.Lock()
	defer m.Mutex.Unlock()
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
