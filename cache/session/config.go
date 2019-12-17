package session

import (
	"time"

	"github.com/herb-go/herb/cache"
)

//DefaultMarshaler default session Marshaler
var DefaultMarshaler = "msgpack"

//DriverNameCacheStore driver name for data store
const DriverNameCacheStore = "cache"

//DriverNameClientStore driver name for client store
const DriverNameClientStore = "cookie"

//StoreConfig store config struct.
type StoreConfig struct {
	DriverName                   string
	Marshaler                    string
	Mode                         string
	TokenLifetime                string //Token initial expired time.Token life time can be update when accessed if UpdateActiveInterval is greater than 0.
	TokenMaxLifetime             string //Token max life time.Token can't live more than TokenMaxLifetime if TokenMaxLifetime if greater than 0.
	TokenContextName             string //Name in request context store the token  data.Default Session is "token".
	CookieName                   string //Cookie name used in CookieMiddleware.Default Session is "herb-session".
	CookiePath                   string //Cookie path used in cookieMiddleware.Default Session is "/".
	CookieSecure                 bool   //Cookie secure value used in cookie middleware.
	AutoGenerate                 bool   //Whether auto generate token when guset visit.Default Session is false.
	UpdateActiveIntervalInSecond int64  //The interval between which token active time update.If less than or equal to 0,the token life time will not be refreshed.
	DefaultSessionFlag           Flag   //Default flag when creating session.
	ClientStoreKey               string
	TokenPrefixMode              string
	TokenLength                  int
	Cache                        cache.OptionConfig
}

//ApplyTo apply config to store.
//Return any error if raised.
func (s *StoreConfig) ApplyTo(store *Store) error {
	var err error
	if s.TokenLifetime != "" {
		store.TokenLifetime, err = time.ParseDuration(s.TokenLifetime)
		if err != nil {
			return err
		}
	}
	if s.TokenMaxLifetime != "" {
		store.TokenMaxLifetime, err = time.ParseDuration(s.TokenMaxLifetime)
		if err != nil {
			return err
		}
	}
	if s.TokenContextName != "" {
		store.TokenContextName = ContextKey(s.TokenContextName)
	}
	if s.CookieName != "" {
		store.CookieName = s.CookieName
	}
	if s.CookiePath != "" {
		store.CookiePath = s.CookiePath
	}
	if s.CookieSecure {
		store.CookieSecure = s.CookieSecure
	}
	store.AutoGenerate = s.AutoGenerate
	if s.UpdateActiveIntervalInSecond != 0 {
		store.UpdateActiveInterval = time.Duration(s.UpdateActiveIntervalInSecond) * time.Second
	}
	store.Mode = s.Mode
	store.DefaultSessionFlag = s.DefaultSessionFlag
	var marshaler string
	marshaler = s.Marshaler
	if marshaler == "" {
		marshaler = DefaultMarshaler
	}
	m, err := cache.NewMarshaler(marshaler)
	if err != nil {
		return err
	}
	store.Marshaler = m
	switch s.DriverName {
	case DriverNameCacheStore:
		c := cache.New()
		err := c.Init(&s.Cache)
		if err != nil {
			return err
		}
		driver := NewCacheDriver()
		coc := NewCacheDriverOptionConfig()
		coc.Cache = c
		coc.Length = s.TokenLength
		coc.PrefixMode = s.TokenPrefixMode
		err = driver.Init(coc)
		if err != nil {
			return err
		}
		soc := NewOptionConfig()
		soc.Driver = driver
		soc.TokenLifetime = store.TokenLifetime
		return store.Init(soc)
	case DriverNameClientStore:
		driver := NewClientDriver()
		coc := NewClientDriverOptionConfig()
		coc.Key = []byte(s.ClientStoreKey)
		err := driver.Init(coc)
		if err != nil {
			return err
		}
		soc := NewOptionConfig()
		soc.Driver = driver
		soc.TokenLifetime = store.TokenLifetime
		return store.Init(soc)
	}
	return nil
}
