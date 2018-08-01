package cache

import (
	"encoding/json"
	"time"
)

type Option interface {
	ApplyTo(*Cache) error
}

type OptionFunc func(*Cache) error

func (i OptionFunc) ApplyTo(cache *Cache) error {
	return i(cache)
}
func OptionJSON(driverName string, creatorjson json.RawMessage, ttlInSecond int64) OptionFunc {
	return func(cache *Cache) error {
		config, err := NewDriverConfig(driverName)
		if err != nil {
			return err
		}
		err = json.Unmarshal(creatorjson, config)
		if err != nil {
			return err
		}
		driver, err := config.Create()
		if err != nil {
			return err
		}
		cache.Driver = driver
		cache.TTL = time.Duration(ttlInSecond * int64(time.Second))
		return nil

	}
}

//Config :The cache config json format struct
type Config struct {
	Driver string
	Config json.RawMessage
	TTL    int64
}

func (c Config) ApplyTo(cache *Cache) error {
	if len(c.Config) == 0 || c.Config == nil {
		c.Config = json.RawMessage("{}")
	}
	return OptionJSON(c.Driver, c.Config, c.TTL)(cache)
}

type ConfigString struct {
	Driver string
	Config string
	TTL    int64
}

func (c ConfigString) ApplyTo(cache *Cache) error {
	config := json.RawMessage(c.Config)
	if len(c.Config) == 0 {
		config = json.RawMessage("{}")
	}
	return OptionJSON(c.Driver, config, c.TTL)(cache)
}
