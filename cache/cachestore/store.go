package cachestore

import "sync"

// Store interface of cache store
type Store interface {
	// Load load value with given key
	// Return value and whether load successfully
	Load(key string) (value interface{}, ok bool)
	// Store sotre value with given key
	Store(key string, value interface{})
	// LoadInterface Load load value with given key
	// Return loaded value or nil if load fail
	LoadInterface(key string) interface{}
}

// MapStore store which stores value in map.
// You should confirm safe for concurrent by yourself.
type MapStore map[string]interface{}

// Load load value with given key
// Return value and whether load successfully
func (m MapStore) Load(key string) (value interface{}, ok bool) {
	v, ok := m[key]
	return v, ok
}

// Store sotre value with given key
func (m MapStore) Store(key string, value interface{}) {
	m[key] = value
}

// LoadInterface Load load value with given key
// Return loaded value or nil if load fail
func (m MapStore) LoadInterface(key string) interface{} {
	return m[key]
}

// NewMapStore create new map store
func NewMapStore() MapStore {
	return MapStore(map[string]interface{}{})
}

// SyncMapStore store which stores value in sync.map.
type SyncMapStore struct {
	Map *sync.Map
}

// Load load value with given key
// Return value and whether load successfully
func (m SyncMapStore) Load(key string) (value interface{}, ok bool) {
	return m.Map.Load(key)
}

// Store sotre value with given key
func (m SyncMapStore) Store(key string, value interface{}) {
	m.Map.Store(key, value)
}

// LoadInterface Load load value with given key
// Return loaded value or nil if load fail
func (m SyncMapStore) LoadInterface(key string) interface{} {
	v, _ := m.Load(key)
	return v
}

// NewSyncMapStore create new sync.map store
func NewSyncMapStore() *SyncMapStore {
	return &SyncMapStore{
		Map: &sync.Map{},
	}
}
