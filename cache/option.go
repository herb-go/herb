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

type OptionConfigJSON struct {
	Driver string
	TTL    int64
	Config ConfigJSON
}

func (o *OptionConfigJSON) ApplyTo(c *Cache) error {
	return OptionConfig(o.Driver, &o.Config, o.TTL).ApplyTo(c)
}
func OptionConfig(driverName string, conf Config, ttlInSecond int64) OptionFunc {
	return func(cache *Cache) error {
		driver, err := NewDriver(driverName, conf, "")
		if err != nil {
			return err
		}
		cache.Driver = driver
		cache.TTL = time.Duration(ttlInSecond * int64(time.Second))
		return nil

	}
}

// func OptionJSON(driverName string, creatorjson json.RawMessage, ttlInSecond int64) OptionFunc {
// 	return func(cache *Cache) error {
// 		config, err := NewDriverConfig(driverName)
// 		if err != nil {
// 			return err
// 		}
// 		err = json.Unmarshal(creatorjson, config)
// 		if err != nil {
// 			return err
// 		}
// 		driver, err := config.Create()
// 		if err != nil {
// 			return err
// 		}
// 		cache.Driver = driver
// 		cache.TTL = time.Duration(ttlInSecond * int64(time.Second))
// 		return nil

// 	}
// }

type Config interface {
	Get(key string, v interface{}) error
}
type ConfigJSON map[string]string

func (c *ConfigJSON) Get(key string, v interface{}) error {
	s, ok := (*c)[key]
	if !ok {
		return nil
	}
	return json.Unmarshal([]byte(s), v)
}
func (c *ConfigJSON) Set(key string, v interface{}) error {
	s, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	(*c)[key] = string(s)
	return nil
}

//Config :The cache config json format struct
// type Config struct {
// 	Driver string
// 	Config json.RawMessage
// 	TTL    int64
// }

// func (c Config) ApplyTo(cache *Cache) error {
// 	if len(c.Config) == 0 || c.Config == nil {
// 		c.Config = json.RawMessage("{}")
// 	}
// 	return OptionJSON(c.Driver, c.Config, c.TTL)(cache)
// }

// type ConfigLoader struct {
// 	Driver string
// 	loader func(v interface{}) error
// 	TTL    int64
// }

// func (c *ConfigLoader) SetLoader(loader func(v interface{}) error) {
// 	c.loader = loader
// }

// func (c *ConfigLoader) ApplyTo(cache *Cache) error {
// 	return OptionLoader(c.Driver, c.loader, c.TTL)(cache)
// }
