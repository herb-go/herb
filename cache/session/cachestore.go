//Package session is used to store user data in token based cache.
//It is normally used as user session or token.
//This package depands github.com/herb-go/herb/cache.
package session

import (
	"crypto/md5"
	"encoding/hex"
	"net/url"
	"time"

	"github.com/herb-go/herb/cache"
)

const (
	//PrefixModeRaw prefix mode which dont convert prefix
	PrefixModeRaw = "raw"
	//PrefixModeMd5 prefix mode which  convert prefix to md5 sum
	PrefixModeMd5 = "md5"
	//PrefixModeEmpty prefix mode which  convert prefix to empty string
	PrefixModeEmpty   = "empty"
	DefaultPrefixMode = PrefixModeRaw
)

const DefaultTokenLength = 64

var defaultUpdateActiveInterval = 5 * time.Minute

var defaultTokenMaxLifetime = 365 * 24 * time.Hour
var (
	defaultCookieName = "herb-session"
	defaultCookiePath = "/"
)

func defaultTokenGenerater(s *CacheDriver, prefix string) (token string, err error) {
	t, err := cache.RandMaskedBytes(cache.TokenMask, s.Length)
	if err != nil {
		return
	}
	prefixconverted, err := s.ConvertPrefix(prefix)
	if err != nil {
		return
	}
	if prefixconverted != "" {
		token = url.PathEscape(prefixconverted) + "-" + string(t)
	} else {
		token = string(t)
	}
	return
}

//NewCacheDriver create new cache driver
func NewCacheDriver() *CacheDriver {
	return &CacheDriver{
		TokenGenerater: defaultTokenGenerater,
		Length:         DefaultTokenLength,
		PrefixMode:     PrefixModeRaw,
	}
}

// MustCacheStore create new data store with given token lifetime.
//Return store created.
//Panic if any error raised.
func MustCacheStore(Cache *cache.Cache, TokenLifetime time.Duration) *Store {
	driver := NewCacheDriver()
	oc := NewCacheDriverOptionConfig()
	oc.Cache = Cache
	err := driver.Init(oc)
	if err != nil {
		panic(err)
	}
	store := New()
	store.Marshaler, err = cache.NewMarshaler(DefaultMarshaler)
	if err != nil {
		panic(err)
	}
	soc := NewOptionConfig()
	soc.Driver = driver
	soc.TokenLifetime = TokenLifetime
	err = store.Init(soc)
	if err != nil {
		panic(err)
	}
	return store
}

//CacheDriver CacheDriver is the stuct store token data in cache.
type CacheDriver struct {
	Cache          *cache.Cache //Cache which stores token data
	Length         int
	PrefixMode     string
	TokenGenerater func(s *CacheDriver, prefix string) (token string, err error) //Token name generate func
}

//Init init cache driver with given option
func (s *CacheDriver) Init(option CacheDriverOption) error {
	return option.ApplyTo(s)
}

//Close Close cachestore and return any error if raised
func (s *CacheDriver) Close() error {
	return s.Cache.Close()
}

//GenerateToken generate new token name with given prefix.
//Return the new token name and error.
func (s *CacheDriver) GenerateToken(prefix string) (token string, err error) {
	return s.TokenGenerater(s, prefix)
}

func (s *CacheDriver) ConvertPrefix(prefix string) (output string, err error) {
	switch s.PrefixMode {
	case PrefixModeEmpty:
		return "", nil
	case PrefixModeRaw:
		return prefix, nil
	case PrefixModeMd5:
		o := md5.Sum([]byte(prefix))
		return hex.EncodeToString(o[:]), nil
	}
	return "", nil
}

//Load load a given session with token from store.
func (s *CacheDriver) Load(v *Session) (err error) {
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

//Save  save given session with given ttl to store.
//Return any error if raised.
func (s *CacheDriver) Save(ts *Session, ttl time.Duration) (err error) {
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
func (s *CacheDriver) Delete(token string) error {
	return s.Cache.Del(token)
}

//GetSessionToken Get the token string from token data.
//Return token and any error raised.
func (s *CacheDriver) GetSessionToken(ts *Session) (token string, err error) {
	return ts.token, nil
}
