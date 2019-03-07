package cachestore

import "sync"

type Store interface {
	Load(key string) (value interface{}, ok bool)
	Store(key string, value interface{})
	LoadInterface(key string) interface{}
}

type MapStore map[string]interface{}

func (m MapStore) Load(key string) (value interface{}, ok bool) {
	v, ok := m[key]
	return v, ok
}
func (m MapStore) Store(key string, value interface{}) {
	m[key] = value
}
func (m MapStore) LoadInterface(key string) interface{} {
	return m[key]
}
func NewMapStore() MapStore {
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
func (m SyncMapStore) LoadInterface(key string) interface{} {
	v, _ := m.Load(key)
	return v
}

func NewSyncMapStore() *SyncMapStore {
	return &SyncMapStore{
		Map: &sync.Map{},
	}
}
