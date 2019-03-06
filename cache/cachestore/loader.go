package cachestore

import "github.com/herb-go/herb/cache"

type Loader struct {
	Cache      cache.Cacheable
	Store      Store
	DataSource *DataSource
}

func (l *Loader) Load(keys ...string) error {
	return l.DataSource.Load(l.Store, l.Cache, keys...)
}
