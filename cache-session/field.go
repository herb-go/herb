package session

import (
	"net/http"
)

type Field struct {
	Store *Store
	Name  string
}

func (f *Field) Get(r *http.Request, v interface{}) (err error) {
	return f.Store.Get(r, f.Name, v)
}

func (f *Field) Set(r *http.Request, v interface{}) (err error) {
	return f.Store.Set(r, f.Name, v)
}

func (f *Field) LoadFrom(ts *Session, v interface{}) (err error) {
	return ts.Get(f.Name, v)
}

func (f *Field) SaveTo(ts *Session, v interface{}) (err error) {
	return ts.Set(f.Name, v)
}

func (f *Field) GetSession(r *http.Request) (ts *Session, err error) {
	return f.Store.GetRequestSession(r)
}
