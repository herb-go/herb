package cache

import (
	"sync"
	"time"
)

//NewUtil create new util
func NewUtil() *Util {
	return &Util{
		locks: &sync.Map{},
	}
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
	Marshaler         Marshaler
	locks             *sync.Map
	CollectionFactory func(Cacheable, string, time.Duration) *Collection
	NodeFactory       func(Cacheable, string) *Node
}

//Clone clone util
func (u *Util) Clone() *Util {
	return &Util{
		Marshaler:         u.Marshaler,
		locks:             u.locks,
		CollectionFactory: u.CollectionFactory,
		NodeFactory:       u.NodeFactory,
	}
}

//Locker create new locker with given key.
//Return locker and if locker is locked.
func (u *Util) Locker(key string) (*Locker, bool) {
	newlocker := &Locker{
		Map: u.locks,
		Key: key,
	}
	v, ok := u.locks.LoadOrStore(key, newlocker)
	return v.(*Locker), ok
}

//Marshal Marshal data model to  bytes.
//Return marshaled bytes and any error rasied.
func (u *Util) Marshal(v interface{}) ([]byte, error) {
	return u.Marshaler.Marshal(v)
}

//Unmarshal Unmarshal bytes to data model.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raseid.
func (u *Util) Unmarshal(bytes []byte, v interface{}) error {
	return u.Marshaler.Unmarshal(bytes, v)
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
