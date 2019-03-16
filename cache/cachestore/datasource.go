package cachestore

import (
	"github.com/herb-go/herb/cache"
)

// CachedDataSource data source with cache
type CachedDataSource struct {
	Cache      cache.Cacheable
	DataSource *DataSource
}

//Load load data by give keys to store.
// Return any error if raised.
func (s *CachedDataSource) Load(m Store, keys ...string) error {
	return s.DataSource.Load(m, s.Cache, keys...)
}

// NewMapStoreLoader create map store lodaer.
func (s *CachedDataSource) NewMapStoreLoader() *Loader {
	return s.DataSource.NewMapStoreLoader(s.Cache)
}

// NewSyncMapStoreLoader create sync map store lodaer.
func (s *CachedDataSource) NewSyncMapStoreLoader() *Loader {
	return s.DataSource.NewSyncMapStoreLoader(s.Cache)
}

//DataSource cache store datasource
type DataSource struct {
	// Creator empty createor.
	// Must return pointer of data.
	Creator func() interface{}
	//SourceLoader  source loader which load data by given keys to map of data pointers.
	// Return map of data pointers and any error if raised.
	SourceLoader func(...string) (map[string]interface{}, error)
}

//Load load data bygiven keys to store wich given cache
//Return any error if rasied.
func (s *DataSource) Load(m Store, c cache.Cacheable, keys ...string) error {
	return Load(m, c, s.SourceLoader, s.Creator, keys...)
}

// NewMapStoreLoader create map store loader with given cache
func (s *DataSource) NewMapStoreLoader(c cache.Cacheable) *Loader {
	return &Loader{
		Store:      NewMapStore(),
		Cache:      c,
		DataSource: s,
	}
}

// NewSyncMapStoreLoader create sync map store loader with given cache
func (s *DataSource) NewSyncMapStoreLoader(c cache.Cacheable) *Loader {
	return &Loader{
		Store:      NewSyncMapStore(),
		Cache:      c,
		DataSource: s,
	}
}

// NewDataSource create new data source.
func NewDataSource() *DataSource {
	return &DataSource{}
}
