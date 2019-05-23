package datastore

import "github.com/herb-go/herb/cache"

// BatchDataLoader batch data loader interface
type BatchDataLoader interface {
	// NewDataElement create empty data element.
	// Must return pointer of data.
	NewDataElement() interface{}
	//BatchLoadData   load data by given keys to map of data pointers.
	// Return map of data pointers and any error if raised.
	BatchLoadData(...string) (map[string]interface{}, error)
}

//LoadWithBatchLoader load data by given key list into data store.
//Param s target data store.
//Param c map cache
//Param l BatchDataLoader
//Param keys key list to load
//Return any error if raised.
func LoadWithBatchLoader(s Store, c cache.Cacheable, l BatchDataLoader, keys ...string) error {
	return Load(s, c, l.BatchLoadData, l.NewDataElement, keys...)
}
