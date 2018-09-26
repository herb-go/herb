package cache

import (
	"encoding/json"
	"time"
)

//Option cache option
type Option interface {
	ApplyTo(*Cache) error
}

//OptionFunc cache option function interface.
type OptionFunc func(*Cache) error

//ApplyTo apply option finction to given cache.
//Return any error if raised.
func (i OptionFunc) ApplyTo(cache *Cache) error {
	return i(cache)
}

//OptionConfigJSON option config in json format
type OptionConfigJSON struct {
	Driver string
	TTL    int64
	Config ConfigJSON
}

//ApplyTo apply config json option to cache.
//Return any error if raised.
func (o *OptionConfigJSON) ApplyTo(c *Cache) error {
	return OptionConfig(o.Driver, &o.Config, o.TTL).ApplyTo(c)
}

//OptionConfigMap option config in map format.
type OptionConfigMap struct {
	Driver string
	TTL    int64
	Config ConfigMap
}

//ApplyTo apply config map option to cache.
//Return any error if raised.
func (o *OptionConfigMap) ApplyTo(c *Cache) error {
	return OptionConfig(o.Driver, &o.Config, o.TTL).ApplyTo(c)
}

//OptionConfig option config return option function with given drivername,config,ttl.
//Return option function
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

//Config cache config interface.
type Config interface {
	Get(key string, v interface{}) error
}

//ConfigJSON config in json format.
type ConfigJSON map[string]string

//Get get value from config json.
//Return any error if raised.
func (c *ConfigJSON) Get(key string, v interface{}) error {
	s, ok := (*c)[key]
	if !ok {
		return nil
	}
	return json.Unmarshal([]byte(s), v)
}

//Set set value to config json.
//Return any error if raised.
func (c *ConfigJSON) Set(key string, v interface{}) error {
	s, err := json.Marshal(v)
	if err != nil {
		return nil
	}
	(*c)[key] = string(s)
	return nil
}

//ConfigMap config in map format.
type ConfigMap map[string]interface{}

//Get get value from config map.
//Return any error if raised.
func (c *ConfigMap) Get(key string, v interface{}) error {
	i, ok := (*c)[key]
	if !ok {
		return nil
	}
	bs, err := json.Marshal(i)
	if err != nil {
		return err
	}
	return json.Unmarshal(bs, v)
}

//Set set value to config map.
//Return any error if raised.
func (c *ConfigMap) Set(key string, v interface{}) error {
	(*c)[key] = v
	return nil
}
