package cache

import "sync"

func NewUtil() *Util {
	return &Util{}
}

type Locker struct {
	Map     *sync.Map
	Key     string
	rlocker *sync.Mutex
	locker  *sync.Mutex
}

func (l *Locker) RLock() {
	lock, ok := l.Map.Load(l.Key)
	if ok {
		l.rlocker = lock.(*sync.Mutex)
		l.rlocker.Lock()
	}
}

func (l *Locker) Lock() {
	var locker *sync.Mutex
	if l.rlocker != nil {
		l.locker = l.rlocker
	} else {
		v, _ := l.Map.LoadOrStore(l.Key, &sync.Mutex{})
		locker = v.(*sync.Mutex)
		l.locker = locker
		l.locker.Lock()
	}
}
func (l *Locker) Unlock() {
	if l.locker != nil {
		l.locker.Unlock()
		l.Map.Delete(l.Key)
		l.locker = nil
		l.rlocker = nil
	} else if l.rlocker != nil {
		l.rlocker.Unlock()
		l.rlocker = nil
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
