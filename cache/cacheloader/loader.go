package cacheloader

import (
	"sync"

	"github.com/herb-go/herb/cache"
)

type Store interface {
	Load(key string) (value interface{}, ok bool)
	Store(key string, value interface{})
}

type MapStore map[string]interface{}

func (m MapStore) Load(key string) (value interface{}, ok bool) {
	v, ok := m[key]
	return v, ok
}
func (m MapStore) Store(key string, value interface{}) {
	m[key] = value
}

func NewMapStore() Store {
	return MapStore(map[string]interface{}{})
}

type SyncMapStore struct {
	Map *sync.Map
}

func (m SyncMapStore) Load(key string) (value interface{}, ok bool) {
	return m.Map.Load(key)
}
func (m SyncMapStore) Store(key string, value interface{}) {
	m.Map.Store(key, value)
}

func NewSyncMapStore() Store {
	return SyncMapStore{
		Map: &sync.Map{},
	}
}

func unmarshalMapElement(cm Store, creator func() interface{}, key string, data []byte, c cache.Cacheable) (err error) {
	v := creator()
	err = c.Unmarshal(data, v)
	if err != nil {
		return err
	}
	cm.Store(key, v)
	return nil
}

//Load load data by given key list into cache map.
//Param cm target data map pointer.
//Param c map cache
//Param loader func by which data load.
//Param creator map element creator.
//Param keys key list to load
//Return any error if raised.
func Load(cm Store, c cache.Cacheable, loader func(...string) (map[string]interface{}, error), creator func() interface{}, keys ...string) error {
	var keysmap = make(map[string]bool, len(keys))
	var filteredKeys = make([]string, len(keys))
	var filteredKeysLength = 0
	if len(keys) == 0 {
		return nil
	}
	for k := range keys {
		if keysmap[keys[k]] == true {
			continue
		}
		keysmap[keys[k]] = true
		if _, ok := cm.Load(keys[k]); !ok {

			filteredKeys[filteredKeysLength] = keys[k]
			filteredKeysLength++
		}
	}
	filteredKeys = filteredKeys[:filteredKeysLength]
	if filteredKeysLength == 0 {
		return nil
	}
	lockKey := filteredKeys[0]
	_, err := c.Wait(lockKey)
	if err != nil {
		return err
	}
	results, err := c.MGetBytesValue(filteredKeys...)
	if err != nil {
		return err
	}
	var uncachedKeys = make([]string, len(filteredKeys))
	var uncachedKeysLength = 0
	for i := range filteredKeys {
		k := filteredKeys[i]
		if results[k] == nil {
			cm.Store(k, creator())
			uncachedKeys[uncachedKeysLength] = k
			uncachedKeysLength++
		} else {
			err = unmarshalMapElement(cm, creator, k, results[k], c)
			if err != nil {
				return err
			}
		}
	}
	uncachedKeys = uncachedKeys[:uncachedKeysLength]
	if uncachedKeysLength == 0 {
		return nil
	}
	unlocker, err := c.Lock(lockKey)
	if err != nil {
		return err
	}
	defer unlocker()
	loaded, err := loader(uncachedKeys...)
	if err != nil {
		return err
	}
	var data = make(map[string][]byte, len(loaded))
	for k := range loaded {
		v := loaded[k]
		cm.Store(k, v)
		data[k], err = c.Marshal(v)
		if err != nil {
			return err
		}
	}
	return c.MSetBytesValue(data, 0)
}

type Loader struct {
	Cache   cache.Cacheable
	Creator func() interface{}
	Loader  func(...string) (map[string]interface{}, error)
}

func (l *Loader) Load(m Store, keys ...string) error {
	return Load(m, l.Cache, l.Loader, l.Creator, keys...)
}
