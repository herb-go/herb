package tokenstore

import (
	"net/http"
	"reflect"
	"sync"

	"github.com/herb-go/herb/cache"
)

//TokenField is registered token field which stand for a data model.
type TokenField struct {
	Key   string       //Token field name.
	Type  reflect.Type //Registered token struct type.
	Store *Store       //Token store which this field belongs to.
}

//LoadFrom load data model from given token data.
//Return any error raised.
func (f *TokenField) LoadFrom(td *TokenData, v interface{}) (err error) {
	if td.token == "" {
		err = ErrDataNotFound
		return
	}
	err = td.Load()
	if err != nil {
		return
	}
	td.Mutex.RLock()
	defer td.Mutex.RUnlock()
	key := f.Key
	typ := reflect.TypeOf(v)
	if typ.Elem() != f.Type {
		return ErrDataTypeWrong
	}
	if reflect.ValueOf(v).IsNil() {
		return ErrNilPoint
	}
	c, ok := td.cache[key]
	if ok == true {
		dst := reflect.ValueOf(v).Elem()
		src := reflect.ValueOf(c).Elem()
		dst.Set(src)
	}
	data, ok := td.data[key]
	if ok == false {
		return ErrDataNotFound
	}
	err = cache.UnmarshalMsgpack(data, v)
	if err == nil {
		td.cache[key] = v
	}
	return
}

//GetFromToken get data model from given token
//Return any error raised.
func (f *TokenField) GetFromToken(token string, v interface{}) (err error) {
	var td *TokenData
	td = f.Store.GetTokenData(token)
	return f.LoadFrom(td, v)
}

//Get get data model form request.
//Return any error raised.
func (f *TokenField) Get(r *http.Request, v interface{}) error {
	var m, err = f.Store.GetRequestTokenData(r)
	if err != nil {
		return err
	}
	return f.LoadFrom(m, v)
}

//RwMutex return the RwMutex of request token data and any error if raised.
func (f *TokenField) RwMutex(r *http.Request) (*sync.RWMutex, error) {
	var td, err = f.Store.GetRequestTokenData(r)
	if err != nil {
		return nil, err
	}
	return td.Mutex, nil
}

//ExpiredAt return the timestamp when token will expired at and any error rasied.
func (f *TokenField) ExpiredAt(r *http.Request) (int64, error) {
	var td, err = f.Store.GetRequestTokenData(r)
	if err != nil {
		return 0, err
	}
	return td.ExpiredAt, nil
}

//GetToken return the token name of request.
func (f *TokenField) GetToken(r *http.Request) (string, error) {
	var td, err = f.Store.GetRequestTokenData(r)
	if err != nil {
		return "", err
	}
	return td.token, nil
}

//SaveTo save datamodel to given token data.
//Return any error is raised.
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

//Set set data model to request.
//Return any error raised.
func (f *TokenField) Set(r *http.Request, v interface{}) error {
	var td, err = f.Store.GetRequestTokenData(r)
	if err != nil {
		return err
	}
	err = f.SaveTo(td, v)
	return err
}

//MustLogin quick regenerate token and save data model to request.
//Usually used to login user.
//Return the new token data.
//Panic is any error raised.
func (f *TokenField) MustLogin(r *http.Request, owner string, v interface{}) (td *TokenData) {
	var err error
	td, err = f.Store.GetRequestTokenData(r)
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

//MustLoginTokenData quick regenerate token and save data model.
//Usually used to login user.
//Return the new token data.
//Panic is any error raised.
func (f *TokenField) MustLoginTokenData(owner string, v interface{}) (td *TokenData) {
	var err error
	td = f.Store.GenerateTokenData(owner)
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
