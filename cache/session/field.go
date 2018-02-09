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

func (f *Field) Flush(r *http.Request) (err error) {
	return f.Store.Del(r, f.Name)
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

func (f *Field) MustGetSession(r *http.Request) *Session {
	return f.Store.MustGetRequestSession(r)
}

func (f *Field) IdentifyRequest(r *http.Request) (string, error) {
	var id = ""
	err := f.Get(r, &id)
	if err == ErrDataNotFound || err == ErrTokenNotValidated {
		return "", nil
	}
	return id, err
}

func (f *Field) Login(w http.ResponseWriter, r *http.Request, id string) error {
	s, err := f.Store.GetRequestSession(r)
	if err != nil {
		return err
	}
	err = s.RegenerateToken(id)
	if err != nil {
		return err
	}
	return f.Set(r, id)
}
func (f *Field) LoginSession(id string) (*Session, error) {
	s, err := f.Store.GenerateSession("")
	if err != nil {
		return nil, err
	}
	err = s.RegenerateToken(id)
	if err != nil {
		return nil, err
	}
	err = f.SaveTo(s, id)
	if err != nil {
		return nil, err
	}
	return s, nil
}
func (f *Field) Logout(w http.ResponseWriter, r *http.Request) error {
	s, err := f.Store.GetRequestSession(r)
	if err != nil {
		return err
	}
	s.SetToken("")
	return nil
}

func (f *Field) CookieMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return f.Store.CookieMiddleware()
}

func (f *Field) HeaderMiddleware(Name string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return f.Store.HeaderMiddleware(Name)
}
