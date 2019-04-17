package session

import "errors"

import "sync"
import "time"
import "reflect"

var (
	//ErrDataNotFound rasied when token data not found.
	ErrDataNotFound = errors.New("Data not found")
	//ErrDataTypeWrong rasied when the given model type is different form registered model type.
	ErrDataTypeWrong = errors.New("Data type wrong")
	//ErrNilPointer raised when data point to nil.
	ErrNilPointer = errors.New("Data point to nil")
)

//Flag Flag used when saving session
type Flag uint64

//FlagDefault default session flag
const FlagDefault = Flag(1)

//FlagTemporay Flag what stands for a Temporay sesson.
//For example,a login withour "remeber me".
const FlagTemporay = Flag(3)

//Session Token data in every request.
type Session struct {
	data           map[string][]byte
	ExpiredAt      int64 //Timestamp when the token expired.
	CreatedTime    int64 //Timestamp when the token created.
	LastActiveTime int64 //Timestamp when the token Last Active.
	cache          map[string]reflect.Value
	token          string
	oldToken       string
	loaded         bool
	tokenChanged   bool
	updated        bool
	notFound       bool
	Store          *Store
	Nonce          []byte
	Flag           Flag
	Mutex          *sync.RWMutex //Read write mutex.
}

//SetFlag Set a flag to session.
func (s *Session) SetFlag(flag Flag, value bool) {
	if value {
		s.Flag = s.Flag | flag
	} else {
		s.Flag = s.Flag &^ flag
	}
}

//HasFlag verify if session has given flag.
func (s *Session) HasFlag(flag Flag) bool {
	return (s.Flag & flag) != 0
}

type tokenCachedSession struct {
	Data           map[string][]byte
	CreatedTime    int64
	LastActiveTime int64
	ExpiredAt      int64
	Flag           Flag
}

//NewSession create new token data in store with given name.
//token the token name.
//s the store which token data belongs to.
//return new Session.
func NewSession(token string, s *Store) *Session {
	t := time.Now().Unix()
	return &Session{
		token:          token,
		data:           map[string][]byte{},
		cache:          map[string]reflect.Value{},
		Store:          s,
		tokenChanged:   false,
		Mutex:          &sync.RWMutex{},
		CreatedTime:    t,
		LastActiveTime: t,
		Flag:           s.DefaultSessionFlag,
		ExpiredAt:      -1,
	}

}

//Token return the toke name.
//Return any error raised.
func (s *Session) Token() (string, error) {
	return s.Store.GetSessionToken(s)
}

//MustToken return the toke name.
func (s *Session) MustToken() string {
	token, err := s.Store.GetSessionToken(s)
	if err != nil {
		panic(err)
	}
	return token
}

//SetToken update token name
func (s *Session) SetToken(newToken string) {
	s.token = newToken
	s.tokenChanged = true
	s.updated = true
}

//RegenerateToken create new token and token data with given owner.
//Return any error raised.
func (s *Session) RegenerateToken(owner string) error {
	token, err := s.Store.GenerateToken(owner)
	if err != nil {
		return err
	}
	s.SetToken(token)

	return nil
}

//Regenerate reset all session values except token
func (s *Session) Regenerate() {
	s.data = map[string][]byte{}
	s.cache = map[string]reflect.Value{}
	s.updated = false
	s.notFound = false
	s.Flag = s.Store.DefaultSessionFlag
}

//Load the token data from cache.
//Repeat call Load will only load data once.
//Return any error raised.
func (s *Session) Load() error {
	if s.token == "" {
		return ErrTokenNotValidated
	}
	if s.loaded {
		if s.notFound {
			return ErrDataNotFound
		}
		return nil
	}
	err := s.Store.LoadSession(s)
	if err == ErrDataNotFound {
		if s.tokenChanged == false {
			s.notFound = true
			s.loaded = true
			return ErrDataNotFound
		}
		err = nil
	}
	if err != nil {
		return err
	}
	return nil
}

//DeleteAndSave Delete token.
func (s *Session) DeleteAndSave() error {
	s.SetToken("")
	return s.Save()
}

//Save Save token data to cache.
//Won't do anything if token data not changed.
//You should call Save manually in your token binding func or non http request usage.
func (s *Session) Save() error {
	return s.Store.SaveSession(s)
}

//Marshal convert Session to bytes.
//Return  Converted bytes and any error raised.
func (s *Session) Marshal() ([]byte, error) {
	return s.Store.Marshaler.Marshal(
		tokenCachedSession{
			Data:           s.data,
			ExpiredAt:      s.ExpiredAt,
			CreatedTime:    s.CreatedTime,
			LastActiveTime: s.LastActiveTime,
			Flag:           s.Flag,
		})
}

//Unmarshal Unmarshal bytes to Session.
// All data in session will be overwrited.
//Return   any error raised.
func (s *Session) Unmarshal(token string, bytes []byte) error {
	var err error

	var Data = tokenCachedSession{}
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.token = token
	s.cache = map[string]reflect.Value{}
	err = s.Store.Marshaler.Unmarshal(bytes, &(Data))
	if err != nil {
		return err
	}
	s.data = Data.Data
	s.ExpiredAt = Data.ExpiredAt
	s.CreatedTime = Data.CreatedTime
	s.LastActiveTime = Data.LastActiveTime
	s.Flag = Data.Flag
	s.loaded = true
	return nil
}

//Set set value to session with given name.
//Return any error if rasied.
func (s *Session) Set(name string, v interface{}) (err error) {
	err = s.Load()
	if err == ErrDataNotFound {
		*s = *NewSession(s.token, s.Store)
		err = nil
	}
	if err != nil {
		return
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	s.SetCache(name, v)
	bytes, err := s.Store.Marshaler.Marshal(v)
	if err != nil {
		return
	}
	s.data[name] = bytes
	s.updated = true
	return
}

//Del delete value form session with given name.
func (s *Session) Del(name string) (err error) {
	err = s.Load()
	if err == ErrDataNotFound {
		*s = *NewSession(s.token, s.Store)
		err = nil
	}
	if err != nil {
		return
	}
	s.Mutex.Lock()
	defer s.Mutex.Unlock()
	delete(s.data, name)
	delete(s.cache, name)
	s.updated = true
	return
}

//Get load data model from given token data.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raised.
func (s *Session) Get(name string, v interface{}) (err error) {
	if s.token == "" {
		err = ErrTokenNotValidated
		return
	}
	err = s.Load()
	if err != nil {
		return
	}
	s.Mutex.RLock()
	defer s.Mutex.RUnlock()
	if v == nil {
		return ErrNilPointer
	}
	vt := reflect.TypeOf(v)
	if vt.Kind() != reflect.Ptr {
		return ErrNilPointer
	}

	c, ok := s.cache[name]
	if ok == true {
		dst := reflect.ValueOf(v).Elem()
		dst.Set(c)
		return
	}
	data, ok := s.data[name]
	if ok == false {
		return ErrDataNotFound
	}
	err = s.Store.Marshaler.Unmarshal(data, v)
	if err == nil {
		s.cache[name] = reflect.ValueOf(v).Elem()
	}
	return
}

//SetCache set cached value to session with given name and value.
func (s *Session) SetCache(name string, v interface{}) {
	s.cache[name] = reflect.ValueOf(v)
}
