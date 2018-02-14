package cache

import (
	"encoding/json"
	"fmt"
	"time"
)

type Initializer interface {
	Init(*Cache) error
}

type InitializerFunc func(*Cache) error

func (i InitializerFunc) Init(cache *Cache) error {
	return i(cache)
}
func Option(driverName string, cacheConfig json.RawMessage, ttlInSecond int64) InitializerFunc {
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

func (c ConfigJSON) Init(cache *Cache) error {
	var config Config
	err := json.Unmarshal(c, &config)
	if err != nil {
		return err
	}
	return Option(config.Driver, config.Config, config.TTL)(cache)
}
