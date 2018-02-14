package session

import (
	"time"

	"github.com/herb-go/herb/cache"
)

type Initializer interface {
	Init(*Store) error
}

type InitializerFunc func(s *Store) error

func (i InitializerFunc) Init(s *Store) error {
	return i(s)
}
func Option(driver Driver, tokenLifetime time.Duration) InitializerFunc {
	return func(s *Store) error {
		s.Driver = driver
		s.TokenLifetime = tokenLifetime
		return nil
	}
}

type CacheDriverInitializer interface {
	Init(*CacheDriver) error
}
type CacheDriverInitializerFunc func(*CacheDriver) error

func (i CacheDriverInitializerFunc) Init(d *CacheDriver) error {
	return i(d)
}
func CacheDriverOption(Cache *cache.Cache) CacheDriverInitializerFunc {
	return func(d *CacheDriver) error {
		d.Cache = Cache
		return nil
	}
}

type ClientDriverInitializer interface {
	Init(*ClientDriver) error
}
type ClientDriverInitializerFunc func(d *ClientDriver) error

func (i ClientDriverInitializerFunc) Init(d *ClientDriver) error {
	return i(d)
}

func ClientDriverOption(key []byte) ClientDriverInitializerFunc {
	return func(d *ClientDriver) error {
		d.Key = key
		return nil
	}
}
