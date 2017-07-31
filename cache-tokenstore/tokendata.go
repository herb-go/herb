package tokenstore

import "errors"
import "github.com/herb-go/herb/cache"
import "sync"
import "time"

var (
	//ErrDataNotFound rasied when token data not found.
	ErrDataNotFound = errors.New("Data not found")
	//ErrDataTypeWrong rasied when the given model type is different form registered model type.
	ErrDataTypeWrong = errors.New("Data type wrong")
	//ErrNilPoint raised when data point to nil.
	ErrNilPoint = errors.New("Data point to nil")
)

//Token data in every request.
type TokenData struct {
	data           map[string][]byte
	ExpiredAt      int64 //Timestamp when the token expired.
	CreatedTime    int64 //Timestamp when the token created.
	LastActiveTime int64 //Timestamp when the token Last Active.
	cache          map[string]interface{}
	token          string
	oldToken       string
	loaded         bool
	tokenChanged   bool
	updated        bool
	store          *Store
	Mutex          *sync.RWMutex //Read write mutex.
}
type tokenCachedData struct {
	Data           map[string][]byte
	CreatedTime    int64
	LastActiveTime int64
	ExpiredAt      int64
}

//NewTokenData create new token data in store with given name.
//token the token name.
//s the store which token data belongs to.
//return new TokenData.
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

//Token return the toke name.
func (t *TokenData) Token() string {
	return t.token
}

//SetToken update token name
func (t *TokenData) SetToken(newToken string) {
	t.token = newToken
	t.tokenChanged = true
	t.updated = true
}

//RegenerateToken create new token and token data with given owner.
//Return any error raised.
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

//Load the token data from cache.
//Repeat call Load will only load data once.
//Return any error raised.
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
		}
		err = nil
	}
	if err != nil {
		return err
	}
	return nil
}

//DeleteAndSave delete token.
func (t *TokenData) DeleteAndSave() error {
	t.SetToken("")
	return t.Save()
}

//Save save token data to cache.
//Won't do anything if token data not changed.
//You should call Save manually in your token binding func or non http request usage.
func (t *TokenData) Save() error {
	if t.store.UpdateActiveInterval > 0 {
		nextUpdateTime := time.Unix(t.LastActiveTime, 0).Add(t.store.UpdateActiveInterval)
		if nextUpdateTime.Before(time.Now()) {
			t.LastActiveTime = time.Now().Unix()
			t.updated = true
		}
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

//Marshal convert TokenData to bytes.
//Return  Converted bytes and any error raised.
func (t *TokenData) Marshal() ([]byte, error) {
	return cache.MarshalMsgpack(
		tokenCachedData{
			Data:           t.data,
			ExpiredAt:      t.ExpiredAt,
			CreatedTime:    t.CreatedTime,
			LastActiveTime: t.LastActiveTime,
		})
}

//Unmarshal Unmarshal bytes to TokenData.
//Return   any error raised.
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
