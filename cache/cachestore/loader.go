package cachestore

type Loader struct {
	Store      Store
	DataSource *DataSource
}

func (l *Loader) Load(keys ...string) error {
	return l.DataSource.Load(l.Store, keys...)
}
