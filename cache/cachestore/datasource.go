package cachestore

import (
	"github.com/herb-go/herb/cache"
)

type DataSource struct {
	Cache        cache.Cacheable
	Creator      func() interface{}
	SourceLoader func(...string) (map[string]interface{}, error)
}

func (s *DataSource) Load(m Store, keys ...string) error {
	return Load(m, s.Cache, s.SourceLoader, s.Creator, keys...)
}

func (s *DataSource) NewMapStoreLoader() *Loader {
	return &Loader{
		Store:      NewMapStore(),
		DataSource: s,
	}
}

func (s *DataSource) NewSyncMapStoreLoader() *Loader {
	return &Loader{
		Store:      NewSyncMapStore(),
		DataSource: s,
	}
}
