package session

import "errors"
import "github.com/herb-go/herb/cache"
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

//SessionFlagDefault default session flag
const FlagDefault = Flag(0)

//SessionFlagTemporay Flag what stands for a Temporay sesson.
//For example,a login withour "remeber me".
const FlagTemporay = Flag(1)

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
	Nonce          []byte
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
func (ts *Session) Token() (string, error) {
	return ts.Store.GetSessionToken(ts)
}

//MustToken return the toke name.
func (ts *Session) MustToken() string {
	token, err := ts.Store.GetSessionToken(ts)
	if err != nil {
		panic(err)
	}
	return token
}

//SetToken update token name
func (ts *Session) SetToken(newToken string) {
	ts.token = newToken
	ts.tokenChanged = true
	ts.updated = true
}

//RegenerateToken create new token and token data with given owner.
//Return any error raised.
func (ts *Session) RegenerateToken(owner string) error {
	token, err := ts.Store.GenerateToken(owner)
	if err != nil {
		return err
	}
	ts.SetToken(token)

	return nil
}

//Regenerate reset all session values except token
func (ts *Session) Regenerate() {
	ts.data = map[string][]byte{}
	ts.cache = map[string]reflect.Value{}
	ts.updated = false
	ts.notFound = false
	ts.Flag = ts.Store.DefaultSessionFlag
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
	return cache.MarshalMsgpack(
		tokenCachedSession{
			Data:           s.data,
			ExpiredAt:      s.ExpiredAt,
			CreatedTime:    s.CreatedTime,
			Nonce:          s.Nonce,
			LastActiveTime: s.LastActiveTime,
			Flag:           s.Flag,
		})
}

//Unmarshal Unmarshal bytes to Session.
//Return   any error raised.
func (t *Session) Unmarshal(token string, bytes []byte) error {
	var err error
	var Data = tokenCachedSession{}
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	t.token = token
	t.cache = map[string]reflect.Value{}
	err = cache.UnmarshalMsgpack(bytes, &(Data))
	if err != nil {
		return err
	}
	t.data = Data.Data
	t.ExpiredAt = Data.ExpiredAt
	t.CreatedTime = Data.CreatedTime
	t.LastActiveTime = Data.LastActiveTime
	t.Flag = Data.Flag
	t.loaded = true
	return nil
}

func (t *Session) Set(name string, v interface{}) (err error) {
	err = t.Load()
	if err == ErrDataNotFound {
		*t = *NewSession(t.token, t.Store)
		err = nil
	}
	if err != nil {
		return
	}

	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	t.SetCache(name, v)
	bytes, err := cache.MarshalMsgpack(v)
	if err != nil {
		return
	}
	t.data[name] = bytes
	t.updated = true
	return
}

func (t *Session) Del(name string) (err error) {
	err = t.Load()
	if err == ErrDataNotFound {
		*t = *NewSession(t.token, t.Store)
		err = nil
	}
	if err != nil {
		return
	}
	t.Mutex.Lock()
	defer t.Mutex.Unlock()
	delete(t.data, name)
	delete(t.cache, name)
	t.updated = true
	return
}

//LoadFrom load data model from given token data.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raised.
func (t *Session) Get(name string, v interface{}) (err error) {
	if t.token == "" {
		err = ErrTokenNotValidated
		return
	}
	err = t.Load()
	if err != nil {
		return
	}
	t.Mutex.RLock()
	defer t.Mutex.RUnlock()
	vt := reflect.TypeOf(v)
	if vt.Kind() != reflect.Ptr {
		return ErrNilPointer
	}
	if v == nil || reflect.ValueOf(v).IsNil() {
		return ErrNilPointer
	}

	c, ok := t.cache[name]
	if ok == true {
		dst := reflect.ValueOf(v).Elem()
		dst.Set(c)
		return
	}
	data, ok := t.data[name]
	if ok == false {
		return ErrDataNotFound
	}
	err = cache.UnmarshalMsgpack(data, v)
	if err == nil {
		t.cache[name] = reflect.ValueOf(v).Elem()
	}
	return
}
func (t *Session) SetCache(name string, v interface{}) {
	t.cache[name] = reflect.ValueOf(v)
}
