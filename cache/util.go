package cache

import "sync"

//NewUtil create new util
func NewUtil() *Util {
	return &Util{}
}

//Locker cache locker
type Locker struct {
	sync.RWMutex
	Map *sync.Map
	Key string
}

//Unlock unlock and delete locker
func (l *Locker) Unlock() {
	l.RWMutex.Unlock()
	l.Map.Delete(l.Key)
}

//Util cache util
type Util struct {
	Marshaler Marshaler
	locks     sync.Map
}

//Locker create new locker with given key.
//Return locker and if locker is locked.
func (u *Util) Locker(key string) (*Locker, bool) {
	newlocker := &Locker{
		Map: &u.locks,
		Key: key,
	}
	v, ok := u.locks.LoadOrStore(key, newlocker)
	return v.(*Locker), ok
}

//DriverUtil drive util struct.
type DriverUtil struct {
	util *Util
}

//Util return cache util
func (u *DriverUtil) Util() *Util {
	return u.util
}

//SetUtil set util to cache
func (u *DriverUtil) SetUtil(util *Util) {
	u.util = util
}
