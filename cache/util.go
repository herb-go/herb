package cache

import "sync"

func NewUtil() *Util {
	return &Util{}
}

type Util struct {
	Marshaler Marshaler
	locks     sync.Map
}

// Lock lock cache value by given key.
//Return  unlock function and any error if rasied
func (u *Util) Lock(key string) (func(), error) {
	lock := &sync.RWMutex{}
	u.locks.Store(key, lock)
	lock.Lock()
	return func() {
		lock.Unlock()
		u.locks.Delete(key)
	}, nil
}

//Wait wait any usef lock unlcok.
//Return whether waited and any error if rasied.
func (u *Util) Wait(key string) (bool, error) {
	l, _ := u.locks.Load(key)
	if l != nil {
		lock := l.(*sync.RWMutex)
		lock.RLock()
		lock.RUnlock()
		return true, nil
	}
	return false, nil
}

type DriverUtil struct {
	util Util
}

func (u *DriverUtil) Util() *Util {
	return &u.util
}

func (u *DriverUtil) SetUtil(util *Util) {
	u.util = *util
}
