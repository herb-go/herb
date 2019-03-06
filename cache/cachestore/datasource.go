package cachestore

import (
	"github.com/herb-go/herb/cache"
)

type CachedDataSource struct {
	Cache      cache.Cacheable
	DataSource *DataSource
}

func (s *CachedDataSource) Load(m Store, keys ...string) error {
	return s.DataSource.Load(m, s.Cache, keys...)
}

func (s *CachedDataSource) NewMapStoreLoader() *Loader {
	return s.DataSource.NewMapStoreLoader(s.Cache)
}

func (s *CachedDataSource) NewSyncMapStoreLoader() *Loader {
	return s.DataSource.NewSyncMapStoreLoader(s.Cache)
}

type DataSource struct {
	Creator      func() interface{}
	SourceLoader func(...string) (map[string]interface{}, error)
}

func (s *DataSource) Load(m Store, c cache.Cacheable, keys ...string) error {
	return Load(m, c, s.SourceLoader, s.Creator, keys...)
}

func (s *DataSource) NewMapStoreLoader(c cache.Cacheable) *Loader {
	return &Loader{
		Store:      NewMapStore(),
		Cache:      c,
		DataSource: s,
	}
}

func (s *DataSource) NewSyncMapStoreLoader(c cache.Cacheable) *Loader {
	return &Loader{
		Store:      NewSyncMapStore(),
		Cache:      c,
		DataSource: s,
	}
}
