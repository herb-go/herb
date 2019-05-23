package datastore

import "github.com/herb-go/herb/cache"

// Loader data store loader
type Loader struct {
	Cache           cache.Cacheable
	Store           Store
	BatchDataLoader BatchDataLoader
}

// Load load data to store by given keys.
func (l *Loader) Load(keys ...string) error {
	return LoadWithBatchLoader(l.Store, l.Cache, l.BatchDataLoader, keys...)
}

// Delete delete value from store and cache by given key.
func (l *Loader) Delete(key string) error {
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

// NewLoader create new loader with given store,cache and batchlaoder
func NewLoader(s Store, c cache.Cacheable, l BatchDataLoader) *Loader {
	return &Loader{
		Cache:           c,
		Store:           s,
		BatchDataLoader: l,
	}
}

// NewMapLoader create new loader with new MapStore,cache and batchlaoder
func NewMapLoader(c cache.Cacheable, l BatchDataLoader) *Loader {
	return NewLoader(NewMapStore(), c, l)
}

//NewSyncMapLoader create new loader with new SyncMapStore,cache and batchlaoder
func NewSyncMapLoader(c cache.Cacheable, l BatchDataLoader) *Loader {
	return NewLoader(NewSyncMapStore(), c, l)
}
