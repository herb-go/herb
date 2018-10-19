package session

import (
	"crypto/rand"
	"time"

	"github.com/jarlyyn/go-utils/security"
)

const clientStoreNonceSize = 4
const clientStoreNewToken = "."

//AESTokenMarshaler token marshaler which crypt data with AES
//Return error if raised
func AESTokenMarshaler(s *ClientDriver, ts *Session) (err error) {
	var data []byte

	ts.Nonce = make([]byte, clientStoreNonceSize)
	_, err = rand.Read(ts.Nonce)
	if err != nil {
		return
	}
	data, err = ts.Marshal()
	if err != nil {
		return
	}
	ts.token, err = security.AESNonceEncryptBase64(data, s.Key)
	return
}

//AESTokenUnmarshaler token unmarshaler which crypt data with AES
//Return error if raised
func AESTokenUnmarshaler(s *ClientDriver, v *Session) (err error) {
	var data []byte
	data, err = security.AESNonceDecryptBase64(v.token, s.Key)
	if err != nil {
		return ErrDataNotFound
	}
	err = v.Unmarshal(v.token, data)
	if err != nil {
		return ErrDataNotFound
	}
	return nil
}

//ClientDriver ClientDriver is the stuct store token data in Client side.
type ClientDriver struct {
	Key              []byte                              //Crypt key
	TokenMarshaler   func(*ClientDriver, *Session) error //Marshler data to Session.token
	TokenUnmarshaler func(*ClientDriver, *Session) error //Unmarshler data from Session.token
}

//NewClientDriver New create a new client side token store with given key and token lifetime.
//Key the key used to encrpty data
//TokenLifeTime is the token initial expired tome.
//Return a new token store.
//All other property of the store can be set after creation.
func NewClientDriver() *ClientDriver {
	return &ClientDriver{
		TokenMarshaler:   AESTokenMarshaler,
		TokenUnmarshaler: AESTokenUnmarshaler,
	}
}

//MustClientStore create new client store with given  key and ttl.
//Return store created.
//Panic if any error raised.
func MustClientStore(key []byte, TokenLifetime time.Duration) *Store {
	driver := NewClientDriver()
	oc := NewClientDriverOptionConfig()
	oc.Key = key
	err := driver.Init(oc)
	if err != nil {
		panic(err)
	}
	store := New()
	soc := NewOptionConfig()
	soc.Driver = driver
	soc.TokenLifetime = TokenLifetime
	err = store.Init(soc)
	if err != nil {
		panic(err)
	}
	return store
}

//GetSessionToken Get the token string from token data.
//Return token and any error raised.
func (s *ClientDriver) GetSessionToken(ts *Session) (token string, err error) {
	err = ts.Save()
	return ts.token, err
}

//Init init client driver with given option.
//Return any error if raised.
func (s *ClientDriver) Init(option ClientDriverOption) error {
	return option.ApplyTo(s)
}

//GenerateToken generate new token name with given prefix.
//Return the new token name and error.
func (s *ClientDriver) GenerateToken(prefix string) (token string, err error) {
	return clientStoreNewToken, nil

}

//Load Load Session form the Session.token.
//Return any error if raised
func (s *ClientDriver) Load(v *Session) (err error) {
	err = s.TokenUnmarshaler(s, v)
	if err != nil {
		return err
	}
	return
}

//Save Save Session if necessary.
//Return any error raised.
func (s *ClientDriver) Save(ts *Session, ttl time.Duration) (err error) {
	ts.oldToken = ts.token
	err = s.TokenMarshaler(s, ts)
	if err != nil {
		return
	}
	if ts.oldToken != ts.token {
		ts.tokenChanged = true
	}
	return
}

//Delete delete the token with given name.
//Return any error if raised.
func (s *ClientDriver) Delete(token string) error {
	return nil
}

//Close Close cachestore and return any error if raised
func (s *ClientDriver) Close() error {
	return nil
}
