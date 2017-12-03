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
func AESTokenMarshaler(s *ClientStore, ts *Session) (err error) {
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
func AESTokenUnmarshaler(s *ClientStore, v *Session) (err error) {
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

//ClientStore ClientStore is the stuct store token data in Client side.
type ClientStore struct {
	Key              []byte                             //Crypt key
	TokenMarshaler   func(*ClientStore, *Session) error //Marshler data to Session.token
	TokenUnmarshaler func(*ClientStore, *Session) error //Unmarshler data from Session.token
}

//New New create a new client side token store with given key and token lifetime.
//Key the key used to encrpty data
//TokenLifeTime is the token initial expired tome.
//Return a new token store.
//All other property of the store can be set after creation.

func NewClientStore(key []byte, TokenLifetime time.Duration) *Store {
	return NewStore(NewClientDriver(key), TokenLifetime)
}

func NewClientDriver(key []byte) *ClientStore {
	return &ClientStore{
		Key:              key,
		TokenMarshaler:   AESTokenMarshaler,
		TokenUnmarshaler: AESTokenUnmarshaler,
	}
}

//GetSessionToken Get the token string from token data.
//Return token and any error raised.
func (s *ClientStore) GetSessionToken(ts *Session) (token string, err error) {
	err = ts.Save()
	return ts.token, err
}

//GenerateToken generate new token name with given prefix.
//Return the new token name and error.
func (s *ClientStore) GenerateToken(prefix string) (token string, err error) {
	return clientStoreNewToken, nil

}

//SearchByPrefix Search all token with given prefix.
//return all tokens start with the prefix.
//ErrFeatureNotSupported will raised if store dont support this feature.
//Return all tokens and any error if raised.
func (s *ClientStore) SearchByPrefix(prefix string) (Tokens []string, err error) {
	return nil, ErrFeatureNotSupported
}

//Load Load Session form the Session.token.
//Return any error if raised
func (s *ClientStore) Load(v *Session) (err error) {
	err = s.TokenUnmarshaler(s, v)
	if err != nil {
		return err
	}
	return
}

//Save Save Session if necessary.
//Return any error raised.
func (s *ClientStore) Save(ts *Session, ttl time.Duration) (err error) {
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
func (s *ClientStore) Delete(token string) error {
	return nil
}

//Close Close cachestore and return any error if raised
func (s *ClientStore) Close() error {
	return nil
}
