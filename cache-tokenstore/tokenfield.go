package tokenstore

import (
	"net/http"
	"reflect"
	"sync"

	"github.com/herb-go/herb/cache"
)

type TokenField struct {
	Key   string
	Type  reflect.Type
	store *Store
}

func (f *TokenField) LoadFrom(m *TokenData, v interface{}) (err error) {
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
func (f *TokenField) GetTokenData(token string, v interface{}) (err error) {
	var td *TokenData
	td = f.store.GetTokenData(token)
	return f.LoadFrom(td, v)
}
func (f *TokenField) Get(r *http.Request, v interface{}) error {
	var m, err = f.store.GetRequestTokenData(r)
	if err != nil {
		return err
	}
	return f.LoadFrom(m, v)
}

func (f *TokenField) RwMutex(r *http.Request) (*sync.RWMutex, error) {
	var td, err = f.store.GetRequestTokenData(r)
	if err != nil {
		return nil, err
	}
	return td.Mutex, nil
}
func (f *TokenField) ExpiredAt(r *http.Request) (int64, error) {
	var td, err = f.store.GetRequestTokenData(r)
	if err != nil {
		return 0, err
	}
	return td.ExpiredAt, nil
}
func (f *TokenField) GetToken(r *http.Request) (string, error) {
	var td, err = f.store.GetRequestTokenData(r)
	if err != nil {
		return "", err
	}
	return td.token, nil
}

func (f *TokenField) SaveTo(td *TokenData, v interface{}) (err error) {
	if td.token == "" {
		err = ErrDataNotFound
		return
	}
	err = td.Load()
	if err != nil {
		return
	}
	key := f.Key
	if reflect.TypeOf(v) != f.Type {
		return ErrDataTypeWrong
	}
	td.Mutex.Lock()
	defer td.Mutex.Unlock()
	td.cache[key] = v
	bytes, err := cache.MarshalMsgpack(v)
	if err != nil {
		return
	}
	td.data[key] = bytes
	err = nil
	td.updated = true
	return
}

func (f *TokenField) Set(r *http.Request, v interface{}) error {
	var td, err = f.store.GetRequestTokenData(r)
	if err != nil {
		return err
	}
	err = f.SaveTo(td, v)
	return err
}

func (f *TokenField) MustGenerate(r *http.Request, owner string, v interface{}) (td *TokenData) {
	var err error
	td, err = f.store.GetRequestTokenData(r)
	if err != nil {
		panic(err)
	}
	err = td.RegenerateToken(owner)
	if err != nil {
		panic(err)
	}
	err = f.SaveTo(td, v)
	if err != nil {
		panic(err)
	}
	return
}

func (f *TokenField) MustGenerateTokenData(owner string, v interface{}) (td *TokenData) {
	var err error
	td = f.store.GenerateTokenData(owner)
	err = f.SaveTo(td, v)
	if err != nil {
		panic(err)
	}
	err = td.Save()
	if err != nil {
		panic(err)
	}
	return
}
