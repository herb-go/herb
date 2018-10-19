package session

import (
	"time"

	"github.com/herb-go/herb/cache"
)

var DefaultMarshaler = "msgpack"

//DriverNameCacheStore driver name for cache store
const DriverNameCacheStore = "cache"

//DriverNameClientStore driver name for client store
const DriverNameClientStore = "cookie"

//StoreConfig store config struct.
type StoreConfig struct {
	DriverName                   string
	Marshaler                    string
	TokenLifetimeInDay           int64  //Token initial expired time.Token life time can be update when accessed if UpdateActiveInterval is greater than 0.
	TokenMaxLifetimeInDay        int64  //Token max life time.Token can't live more than TokenMaxLifetime if TokenMaxLifetime if greater than 0.
	TokenContextName             string //Name in request context store the token  data.Default Session is "token".
	CookieName                   string //Cookie name used in CookieMiddleware.Default Session is "herb-session".
	CookiePath                   string //Cookie path used in cookieMiddleware.Default Session is "/".
	CookieSecure                 bool   //Cookie secure value used in cookie middleware.
	AutoGenerate                 bool   //Whether auto generate token when guset visit.Default Session is false.
	UpdateActiveIntervalInSecond int64  //The interval between who token active time update.If less than or equal to 0,the token life time will not be refreshed.
	DefaultSessionFlag           Flag   //Default flag when creating session.
	ClientStoreKey               string
	Cache                        cache.OptionConfigMap
}

//ApplyTo apply config to store.
//Return any error if raised.
func (s *StoreConfig) ApplyTo(store *Store) error {
	if s.TokenLifetimeInDay != 0 {
		store.TokenLifetime = time.Duration(s.TokenLifetimeInDay) * time.Hour * 24
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
