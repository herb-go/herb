package session

import (
	"time"

	"github.com/herb-go/herb/cache"
)

// Option store init option.
type Option interface {
	ApplyTo(*Store) error
}

//OptionFunc store init option func type.
type OptionFunc func(s *Store) error

//ApplyTo apply option func to given store.
//Return any error if raised.
func (i OptionFunc) ApplyTo(s *Store) error {
	return i(s)
}

// OptionCommon create option function with given driver and lifetime.
func OptionCommon(driver Driver, tokenLifetime time.Duration) OptionFunc {
	return func(s *Store) error {
		s.Driver = driver
		s.TokenLifetime = tokenLifetime
		return nil
	}
}

//CacheDriverOption cache driver init option.
type CacheDriverOption interface {
	ApplyTo(*CacheDriver) error
}

//CacheDriverOptionFunc cache driver init option function type.
type CacheDriverOptionFunc func(*CacheDriver) error

//ApplyTo appky init option function to cache driver.
func (i CacheDriverOptionFunc) ApplyTo(d *CacheDriver) error {
	return i(d)
}

//CacheDriverOptionCommon create cache driver init option function with given cache.
func CacheDriverOptionCommon(Cache *cache.Cache) CacheDriverOptionFunc {
	return func(d *CacheDriver) error {
		d.Cache = Cache
		return nil
	}
}

//ClientDriverOption client driver init option.
type ClientDriverOption interface {
	ApplyTo(*ClientDriver) error
}

//ClientDriverOptionFunc client driver init option function type.
type ClientDriverOptionFunc func(d *ClientDriver) error

//ApplyTo apply  client driver init option function to client driver.
func (i ClientDriverOptionFunc) ApplyTo(d *ClientDriver) error {
	return i(d)
}

// ClientDriverOptionCommon create client driver init option function with given key.
func ClientDriverOptionCommon(key []byte) ClientDriverOptionFunc {
	return func(d *ClientDriver) error {
		d.Key = key
		return nil
	}
}
