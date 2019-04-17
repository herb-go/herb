package session

import (
	"net/http"
)

//Field session field struct
type Field struct {
	Store *Store
	Name  string
}

//Get get value from session store in  given http request and save to v.
//Return any error if raised.
func (f *Field) Get(r *http.Request, v interface{}) (err error) {
	return f.Store.Get(r, f.Name, v)
}

//Set set value to session store in given http request.
//Return any error if raised.
func (f *Field) Set(r *http.Request, v interface{}) (err error) {
	return f.Store.Set(r, f.Name, v)
}

//Flush flush session store in given http request.
//Return any error if raised.
func (f *Field) Flush(r *http.Request) (err error) {
	return f.Store.Del(r, f.Name)
}

//LoadFrom load value form given session.
//Return any error if raised.
func (f *Field) LoadFrom(ts *Session, v interface{}) (err error) {
	return ts.Get(f.Name, v)
}

//SaveTo save value to given  session.
//Return any error if raised.
func (f *Field) SaveTo(ts *Session, v interface{}) (err error) {
	return ts.Set(f.Name, v)
}

//GetSession get Session from http request.
//Return session and any error if raised.
func (f *Field) GetSession(r *http.Request) (ts *Session, err error) {
	return f.Store.GetRequestSession(r)
}

//IdentifyRequest indentify request with field.
//Return  id and any error if raised.
func (f *Field) IdentifyRequest(r *http.Request) (string, error) {
	var id = ""
	err := f.Get(r, &id)
	if err == ErrDataNotFound || err == ErrTokenNotValidated {
		return "", nil
	}
	return id, err
}

//Login login to request with given id.
//Return any error if raised.
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

//Logout  logout form request.
func (f *Field) Logout(w http.ResponseWriter, r *http.Request) error {
	s, err := f.Store.GetRequestSession(r)
	if err != nil {
		return err
	}
	s.SetToken("")
	return nil
}
