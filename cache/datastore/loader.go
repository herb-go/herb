package datastore

import "github.com/herb-go/herb/cache"

// Loader data store loader
type Loader struct {
	Cache  cache.Cacheable
	Store  Store
	Loader BatchLoader
}

// Load load data to store by given keys.
func (l *Loader) Load(keys ...string) error {
	return Load(l.Store, l.Cache, l.Loader.BatchLoadData, l.Loader.NewDataElement, keys...)
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
