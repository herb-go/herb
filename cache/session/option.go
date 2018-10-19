package session

import (
	"time"

	"github.com/herb-go/herb/cache"
)

// Option store init option.
type Option interface {
	ApplyTo(*Store) error
}

func NewOptionConfig() *OptionConfig {
	return &OptionConfig{}
}

type OptionConfig struct {
	Driver        Driver
	TokenLifetime time.Duration
}

func (o *OptionConfig) ApplyTo(s *Store) error {
	s.Driver = o.Driver
	s.TokenLifetime = o.TokenLifetime
	return nil
}

//CacheDriverOption cache driver init option.
type CacheDriverOption interface {
	ApplyTo(*CacheDriver) error
}

func NewCacheDriverOptionConfig() *CacheDriverOptionConfig {
	return &CacheDriverOptionConfig{}
}

type CacheDriverOptionConfig struct {
	Cache *cache.Cache
}

func (o *CacheDriverOptionConfig) ApplyTo(d *CacheDriver) error {
	d.Cache = o.Cache
	return nil
}

//ClientDriverOption client driver init option.
type ClientDriverOption interface {
	ApplyTo(*ClientDriver) error
}

func NewClientDriverOptionConfig() *ClientDriverOptionConfig {
	return &ClientDriverOptionConfig{}
}

type ClientDriverOptionConfig struct {
	Key []byte
}

func (o *ClientDriverOptionConfig) ApplyTo(d *ClientDriver) error {
	d.Key = o.Key
	return nil
}
