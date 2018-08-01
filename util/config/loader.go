package config

import (
	"sync"
)

var registeredLoader = []func(){}

var Lock sync.RWMutex

func RegisterLoader(loader func()) {
	registeredLoader = append(registeredLoader, loader)
}

func LoadAll() {
	defer Lock.RUnlock()
	Lock.RLock()
	for _, k := range registeredLoader {
		k()
	}
}
