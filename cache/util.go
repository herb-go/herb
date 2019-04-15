package cache

import "sync"

func NewUtil() *Util {
	return &Util{}
}

type Locker struct {
	sync.RWMutex
	Map *sync.Map
	Key string
}

func (l *Locker) Unlock() {
	l.RWMutex.Unlock()
	l.Map.Delete(l.Key)
}

type Util struct {
	Marshaler Marshaler
	locks     sync.Map
}

func (u *Util) Locker(key string) (*Locker, bool) {
	newlocker := &Locker{
		Map: &u.locks,
		Key: key,
	}
	v, ok := u.locks.LoadOrStore(key, newlocker)
	return v.(*Locker), ok
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
