package cache

import (
	"encoding/json"
	"fmt"
	"time"
)

type Option interface {
	ApplyTo(*Cache) error
}

type OptionFunc func(*Cache) error

func (i OptionFunc) ApplyTo(cache *Cache) error {
	return i(cache)
}
func OptionCommon(driverName string, cacheConfig json.RawMessage, ttlInSecond int64) OptionFunc {
	return func(cache *Cache) error {
		driversMu.RLock()
		driveri, ok := drivers[driverName]
		driversMu.RUnlock()
		if !ok {
			return fmt.Errorf("cache: unknown driver %q (forgotten import?)", driverName)
		}
		driver, err := driveri.New(cacheConfig)
		if err != nil {
			return err
		}
		cache.Driver = driver
		cache.TTL = time.Duration(ttlInSecond * int64(time.Second))
		return nil

	}
}

type ConfigJSON []byte

func (c ConfigJSON) ApplyTo(cache *Cache) error {
	var config Config
	err := json.Unmarshal(c, &config)
	if err != nil {
		return err
	}
	return OptionCommon(config.Driver, config.Config, config.TTL)(cache)
}
