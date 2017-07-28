package tokenstore

import "errors"
import "github.com/herb-go/herb/cache"
import "sync"
import "time"

var ErrDataNotFound = errors.New("Data not found")
var ErrDataTypeWrong = errors.New("Data type wrong")
var ErrNilPoint = errors.New("Data point to nil")
var ErrDataTypeNotRegister = errors.New("Data type not register")

type TokenData struct {
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
type tokenCachedData struct {
	Data           map[string][]byte
	CreatedTime    int64
	LastActiveTime int64
	ExpiredAt      int64
}

func NewTokenData(token string, s *Store) *TokenData {
	t := time.Now().Unix()
	return &TokenData{
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
func (t *TokenData) Token() string {
	return t.token
}
func (t *TokenData) SetToken(newToken string) {
	t.token = newToken
	t.tokenChanged = true
	t.updated = true
}
func (t *TokenData) RegenerateToken(owner string) error {
	token, err := t.store.GenerateToken(owner)
	if err != nil {
		return err
	}
	t.data = map[string][]byte{}
	t.cache = map[string]interface{}{}
	t.SetToken(token)
	return nil
}

func (t *TokenData) Load() error {
	if t.token == "" {
		return ErrTokenNotValidated
	}
	if t.loaded {
		return nil
	}
	err := t.store.loadTokenData(t)
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
func (t *TokenData) DeleteAndSave() error {
	t.SetToken("")
	return t.Save()
}
func (t *TokenData) Save() error {
	nextUpdateTime := time.Unix(t.LastActiveTime, 0).Add(t.store.UpdateActiveInterval)
	if nextUpdateTime.Before(time.Now()) {
		t.LastActiveTime = time.Now().Unix()
		t.updated = true
	}
	if t.updated && t.token != "" {
		err := t.store.saveTokenData(t)
		if err != nil {
			return err
		}
		t.updated = false
	}
	if t.tokenChanged && t.oldToken != "" {
		err := t.store.DeleteToken(t.oldToken)
		if err != nil {
			return err
		}
	}
	return nil
}
func (t *TokenData) Marshal() ([]byte, error) {
	return cache.MarshalMsgpack(
		tokenCachedData{
			Data:           t.data,
			ExpiredAt:      t.ExpiredAt,
			CreatedTime:    t.CreatedTime,
			LastActiveTime: t.LastActiveTime,
		})
}
func (t *TokenData) Unmarshal(token string, bytes []byte) error {
	var err error
	var Data = tokenCachedData{}
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
