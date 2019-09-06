package cache

import (
	"encoding/json"
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
	Config    Config
}

//ApplyTo apply option to given cache.
//Return any error if raised.
func (o *OptionConfig) ApplyTo(cache *Cache) error {
	driver, err := NewDriver(o.Driver, o.Config, "")
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

//OptionConfigMap option config in map format.
type OptionConfigMap struct {
	Driver    string
	Marshaler string
	TTL       int64
	Config    ConfigMap
}

//ApplyTo apply config map option to cache.
//Return any error if raised.
func (o *OptionConfigMap) ApplyTo(c *Cache) error {
	oc := NewOptionConfig()
	oc.Driver = o.Driver
	oc.Config = &o.Config
	oc.Marshaler = o.Marshaler
	oc.TTL = o.TTL
	return oc.ApplyTo(c)
}

//Config cache config interface.
type Config interface {
	Get(key string, v interface{}) error
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
