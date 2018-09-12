package store

import "encoding/json"

type Option interface {
	ApplyTo(*Store) error
}

type OptionFunc func(*Store) error

func (i OptionFunc) ApplyTo(store *Store) error {
	return i(store)
}

type OptionConfigJSON struct {
	Driver string
	Config ConfigJSON
}

func (o *OptionConfigJSON) ApplyTo(store *Store) error {
	driver, err := NewDriver(o.Driver, &o.Config, "")
	if err != nil {
		return err
	}
	store.Driver = driver
	return nil
}

type OptionConfigMap struct {
	Driver string
	Config ConfigMap
}

func (o *OptionConfigMap) ApplyTo(store *Store) error {
	driver, err := NewDriver(o.Driver, &o.Config, "")
	if err != nil {
		return err
	}
	store.Driver = driver
	return nil
}

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

type ConfigMap map[string]interface{}

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

func (c *ConfigMap) Set(key string, v interface{}) error {
	(*c)[key] = v
	return nil
}
