package cachestore

import "github.com/herb-go/herb/cache"

// Loader cache store loader
type Loader struct {
	Cache      cache.Cacheable
	Store      Store
	DataSource *DataSource
}

// Load load data to store by given keys.
func (l *Loader) Load(keys ...string) error {
	return l.DataSource.Load(l.Store, l.Cache, keys...)
}

// Del delete value from store and cache by given key.
func (l *Loader) Del(key string) error {
	l.Store.Delete(key)
	if l.Cache == nil {
		return nil
	}
	return l.Cache.Del(key)
}

// Flush flush all data from store and cache.
func (l *Loader) Flush() error {
	l.Store.Flush()
	if l.Cache == nil {
		return nil
	}
	return l.Cache.Flush()
}
