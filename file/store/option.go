package store

import "encoding/json"

type Option interface {
	ApplyTo(*Store) error
}

type OptionFunc func(*Store) error

func (i OptionFunc) ApplyTo(store *Store) error {
	return i(store)
}

type OptionConfig struct {
	Driver string
	Config ConfigJSON
}

func (o *OptionConfig) ApplyTo(store *Store) error {
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
