//Package session is used to store user data in token based cache.
//It is normally used as user session or token.
//This package depands github.com/herb-go/herb/cache.
package session

import (
	"net/url"
	"time"

	"github.com/herb-go/herb/cache"
)

var defaultUpdateActiveInterval = 5 * time.Minute

var defaultTokenMaxLifetime = 365 * 24 * time.Hour
var (
	defaultCookieName = "herb-session"
	defaultCookiePath = "/"
)

func defaultTokenGenerater(s *CacheStore, prefix string) (token string, err error) {
	t, err := cache.RandMaskedBytes(cache.TokenMask, 256)
	if err != nil {
		return
	}

	token = url.PathEscape(prefix) + "-" + string(t)
	return
}

//NewCacheStore New create a new token store with given cache and token lifetime.
//Cache is the cache which dates stored in.
//TokenLifeTime is the token initial expired tome.
//Return a new token store.
//All other property of the store can be set after creation.
func NewCacheStore(Cache *cache.Cache, TokenLifetime time.Duration) *Store {
	s := New()
	s.Init(NewCacheDriver(Cache), TokenLifetime)
	return s
}

func NewCacheDriver(Cache *cache.Cache) *CacheStore {
	return &CacheStore{
		Cache:          Cache,
		TokenGenerater: defaultTokenGenerater,
	}
}

//CacheStore CacheStore is the stuct store token data in cache.
type CacheStore struct {
	Cache          *cache.Cache                                                 //Cache which stores token data
	TokenGenerater func(s *CacheStore, prefix string) (token string, err error) //Token name generate func
}

//Close Close cachestore and return any error if raised
func (s *CacheStore) Close() error {
	return s.Cache.Close()
}

//GenerateToken generate new token name with given prefix.
//Return the new token name and error.
func (s *CacheStore) GenerateToken(prefix string) (token string, err error) {
	return s.TokenGenerater(s, prefix)
}

//Load load a given session with token from store.
func (s *CacheStore) Load(v *Session) (err error) {
	token := v.token
	bytes, err := s.Cache.GetBytesValue(token)
	if err == cache.ErrNotFound {
		err = ErrDataNotFound
	}
	if err != nil {
		return err
	}

	err = v.Unmarshal(token, bytes)
	if err == nil {
		v.token = token
	}
	return
}

func (s *CacheStore) Save(ts *Session, ttl time.Duration) (err error) {
	bytes, err := ts.Marshal()
	if err != nil {
		return err
	}
	if ts.oldToken == ts.token {
		err = s.Cache.UpdateBytesValue(ts.token, bytes, ttl)
	} else {
		err = s.Cache.SetBytesValue(ts.token, bytes, ttl)
	}
	return
}

//Delete delete the token with given name.
//Return any error if raised.
func (s *CacheStore) Delete(token string) error {
	return s.Cache.Del(token)
}

//GetSessionToken Get the token string from token data.
//Return token and any error raised.
func (s *CacheStore) GetSessionToken(ts *Session) (token string, err error) {
	return ts.token, nil
}
