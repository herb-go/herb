package session

import (
	"time"

	"github.com/herb-go/herb/cache"
)

// Option store init option interface.
type Option interface {
	ApplyTo(*Store) error
}

//NewOptionConfig create new store init config
func NewOptionConfig() *OptionConfig {
	return &OptionConfig{}
}

//OptionConfig store init config
type OptionConfig struct {
	Driver        Driver
	TokenLifetime time.Duration
}

//ApplyTo apply config to sessoion store
func (o *OptionConfig) ApplyTo(s *Store) error {
	s.Driver = o.Driver
	s.TokenLifetime = o.TokenLifetime
	return nil
}

//CacheDriverOption cache driver init option interface.
type CacheDriverOption interface {
	ApplyTo(*CacheDriver) error
}

//NewCacheDriverOptionConfig create new cache driver init option
func NewCacheDriverOptionConfig() *CacheDriverOptionConfig {
	return &CacheDriverOptionConfig{}
}

//CacheDriverOptionConfig cache driver init option
type CacheDriverOptionConfig struct {
	Cache      *cache.Cache
	Length     int
	PrefixMode string
}

//ApplyTo apply cache driver option config to cache driver.
//return any error if raised.
func (o *CacheDriverOptionConfig) ApplyTo(d *CacheDriver) error {
	d.Cache = o.Cache
	if o.Length != 0 {
		d.Length = o.Length
	}
	if o.PrefixMode != "" {
		d.PrefixMode = o.PrefixMode
	}
	return nil
}

//ClientDriverOption client driver init option interface.
type ClientDriverOption interface {
	ApplyTo(*ClientDriver) error
}

//NewClientDriverOptionConfig create new client driver init option.
func NewClientDriverOptionConfig() *ClientDriverOptionConfig {
	return &ClientDriverOptionConfig{}
}

//ClientDriverOptionConfig client driver init option.
type ClientDriverOptionConfig struct {
	Key []byte
}

//ApplyTo apply client driver option config to cache driver.
//return any error if raised.
func (o *ClientDriverOptionConfig) ApplyTo(d *ClientDriver) error {
	d.Key = o.Key
	return nil
}
