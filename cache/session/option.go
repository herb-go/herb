package session

import (
	"time"

	"github.com/herb-go/herb/cache"
)

type Option interface {
	ApplyTo(*Store) error
}

type OptionFunc func(s *Store) error

func (i OptionFunc) ApplyTo(s *Store) error {
	return i(s)
}
func OptionCommon(driver Driver, tokenLifetime time.Duration) OptionFunc {
	return func(s *Store) error {
		s.Driver = driver
		s.TokenLifetime = tokenLifetime
		return nil
	}
}

type CacheDriverOption interface {
	ApplyTo(*CacheDriver) error
}
type CacheDriverOptionFunc func(*CacheDriver) error

func (i CacheDriverOptionFunc) ApplyTo(d *CacheDriver) error {
	return i(d)
}
func CacheDriverOptionCommon(Cache *cache.Cache) CacheDriverOptionFunc {
	return func(d *CacheDriver) error {
		d.Cache = Cache
		return nil
	}
}

type ClientDriverOption interface {
	ApplyTo(*ClientDriver) error
}
type ClientDriverOptionFunc func(d *ClientDriver) error

func (i ClientDriverOptionFunc) ApplyTo(d *ClientDriver) error {
	return i(d)
}

func ClientDriverOptionCommon(key []byte) ClientDriverOptionFunc {
	return func(d *ClientDriver) error {
		d.Key = key
		return nil
	}
}
