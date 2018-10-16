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

//MustGetSession get Session from http request.
//Return session.
//Panic if any error raised.
func (f *Field) MustGetSession(r *http.Request) *Session {
	return f.Store.MustGetRequestSession(r)
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

//LoginSession login into session with given id.
//Return session and any error if raised.
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

//Logout  logout form request.
func (f *Field) Logout(w http.ResponseWriter, r *http.Request) error {
	s, err := f.Store.GetRequestSession(r)
	if err != nil {
		return err
	}
	s.SetToken("")
	return nil
}

//CookieMiddleware return a Middleware which install the token which special by cookie.
//This middleware will save token after request finished if the token changed,and update cookie if necessary.
func (f *Field) CookieMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return f.Store.CookieMiddleware()
}

//HeaderMiddleware return a Middleware which install the token which special by Header with given name.
//This middleware will save token after request finished if the token changed.
func (f *Field) HeaderMiddleware(Name string) func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return f.Store.HeaderMiddleware(Name)
}

//InstallMiddleware middleware which auto install session depand on store mode.
//Cookie middleware will be installed if no valid store mode given.
func (f *Field) InstallMiddleware() func(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	return f.Store.InstallMiddleware()
}
