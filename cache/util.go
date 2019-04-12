package cache

import "sync"

func NewUtil() *Util {
	return &Util{}
}

type Locker struct {
	Map      *sync.Map
	Key      string
	rwlocker *sync.RWMutex
	locker   *sync.RWMutex
}

func (l *Locker) RLock() {
	lock, ok := l.Map.Load(l.Key)
	if ok {
		l.rwlocker = lock.(*sync.RWMutex)
		l.rwlocker.RLock()
	}
}

func (l *Locker) RUnlock() {
	if l.rwlocker != nil {
		l.rwlocker.RUnlock()
	}
}
func (l *Locker) Lock() {
	var locker *sync.RWMutex
	v, _ := l.Map.LoadOrStore(l.Key, &sync.RWMutex{})
	locker = v.(*sync.RWMutex)
	l.locker = locker
	l.locker.Lock()
}
func (l *Locker) Unlock() {
	if l.locker != nil {
		l.locker.Unlock()
		l.Map.Delete(l.Key)
	}
}

type Util struct {
	Marshaler Marshaler
	locks     sync.Map
}

func (u *Util) Locker(key string) *Locker {
	return &Locker{
		Map: &u.locks,
		Key: key,
	}
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
