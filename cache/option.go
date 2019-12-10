package cache

import (
	"time"
)

//Option cache option interface.
type Option interface {
	ApplyTo(*Cache) error
}

//NewOptionConfig create new cache option.
func NewOptionConfig() *OptionConfig {
	return &OptionConfig{}
}

//OptionConfig cache option
type OptionConfig struct {
	Driver    string
	TTL       int64
	Marshaler string
	Config    func(v interface{}) error `config:", lazyload"`
}

//ApplyTo apply option to given cache.
//Return any error if raised.
func (o *OptionConfig) ApplyTo(cache *Cache) error {
	driver, err := NewDriver(o.Driver, o.Config)
	if err != nil {
		return err
	}
	cache.Driver = driver
	var mname = o.Marshaler
	if mname == "" {
		mname = DefaultMarshaler
	}
	marshaler, err := NewMarshaler(mname)
	if err != nil {
		return err
	}
	u := NewUtil()
	u.Marshaler = marshaler
	driver.SetUtil(u)
	cache.TTL = time.Duration(o.TTL * int64(time.Second))
	return nil
}
