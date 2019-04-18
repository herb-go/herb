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
	TokenLifetimeInHour          int64  //Token initial expired time in second.Token life time can be update when accessed if UpdateActiveInterval is greater than 0.
	TokenLifetimeInDay           int64  //Token initial expired time in day.Token life time can be update when accessed if UpdateActiveInterval is greater than 0.Skipped if  TokenLifetimeInHour is set.
	TokenMaxLifetimeInDay        int64  //Token max life time.Token can't live more than TokenMaxLifetime if TokenMaxLifetime if greater than 0.
	TokenContextName             string //Name in request context store the token  data.Default Session is "token".
	CookieName                   string //Cookie name used in CookieMiddleware.Default Session is "herb-session".
	CookiePath                   string //Cookie path used in cookieMiddleware.Default Session is "/".
	CookieSecure                 bool   //Cookie secure value used in cookie middleware.
	AutoGenerate                 bool   //Whether auto generate token when guset visit.Default Session is false.
	UpdateActiveIntervalInSecond int64  //The interval between who token active time update.If less than or equal to 0,the token life time will not be refreshed.
	DefaultSessionFlag           Flag   //Default flag when creating session.
	ClientStoreKey               string
	TokenPrefixMode              string
	TokenLength                  int
	Cache                        cache.OptionConfigMap
}

//ApplyTo apply config to store.
//Return any error if raised.
func (s *StoreConfig) ApplyTo(store *Store) error {
	if s.TokenLifetimeInHour != 0 {
		store.TokenLifetime = time.Duration(s.TokenLifetimeInHour) * time.Hour
	} else if s.TokenLifetimeInDay != 0 {
		store.TokenLifetime = time.Duration(s.TokenLifetimeInDay) * 24 * time.Hour
	}
	if s.TokenMaxLifetimeInDay != 0 {
		store.TokenMaxLifetime = time.Duration(s.TokenMaxLifetimeInDay) * time.Hour * 24
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
