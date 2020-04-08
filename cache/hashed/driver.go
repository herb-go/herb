package hashed

import (
	"encoding/binary"
	"time"

	"github.com/herb-go/herb/cache"
)

//Store cache store interface
type Store interface {
	//Close close cache store.
	//Return any error if raised.
	Close() error
	//Flush flush cache.
	//Return any error if raised.
	Flush() error
	//Hash hash given key to hashed string.
	//Return hash and any error if raised
	Hash(string) (string, error)
	//Lock lock by given hashed key.
	//Return unlock func and any error if raised
	Lock(string) (func(), error)
	//Load load hash data by given hash key.
	//Return hash data and any error if raised.
	Load(hash string) (*Hashed, error)
	//Delete delete all data by given hash key.
	//Return any error if raised.
	Delete(hash string) error
	//Save save hash data by given hash key,status and hash data.
	//Return any error if raised.
	Save(hash string, status *Status, data *Hashed) error
}

//Driver driver struct
type Driver struct {
	cache.DriverUtil
	Store        Store
	GcErrHanlder func(err error)
}

// NewDriver create new cache driver with given store
func NewDriver(Store Store) *Driver {
	return &Driver{
		Store:        Store,
		GcErrHanlder: nil,
	}
}

type context struct {
	hash     string
	unlocker func()
	data     *Hashed
}

func (d *Driver) lockAndGetData(key string) (ctx *context, err error) {
	c := &context{}
	c.hash, err = d.Store.Hash(key)
	if err != nil {
		return nil, err
	}
	c.unlocker, err = d.Store.Lock(c.hash)
	if err != nil {
		return nil, err
	}
	c.data, err = d.Store.Load(c.hash)
	if err != nil {
		c.unlocker()
		return nil, err
	}
	return c, nil
}
func (d *Driver) save(ctx *context, status *Status) error {
	if ctx.data.isEmpty() {
		return d.Store.Delete(ctx.hash)
	}
	return d.Store.Save(ctx.hash, status, ctx.data)
}

//SetBytesValue Set bytes data to cache by given key.
func (d *Driver) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return err
	}
	defer ctx.unlocker()
	status := ctx.data.set(NewData(key, now.Add(ttl).Unix(), bytes), now.Unix())
	return d.save(ctx, status)
}

//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
func (d *Driver) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return err
	}
	defer ctx.unlocker()
	status := ctx.data.update(NewData(key, now.Add(ttl).Unix(), bytes), now.Unix())
	return d.save(ctx, status)

}

//GetBytesValue Get bytes data from cache by given key.
func (d *Driver) GetBytesValue(key string) ([]byte, error) {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return nil, err
	}
	data := ctx.data.get(key, now.Unix())
	if data == nil {
		return nil, cache.ErrNotFound
	}
	return data.Data, nil
}

//Del Delete data in cache by given key.
func (d *Driver) Del(key string) error {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return err
	}
	defer ctx.unlocker()
	status := ctx.data.delete(key, now.Unix())
	return d.save(ctx, status)
}

//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
func (d *Driver) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	var v int64
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return 0, err
	}
	data := ctx.data.get(key, now.Unix())
	if data == nil {
		v = 0
	} else {
		v = int64(binary.BigEndian.Uint64(data.Data))
	}
	v = v + increment
	var bytes = make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(v))
	status := ctx.data.set(NewData(key, now.Add(ttl).Unix(), bytes), now.Unix())
	err = d.save(ctx, status)
	if err != nil {
		return 0, err
	}
	return v, nil
}

//SetCounter Set int val in cache by given key.Count cache and data cache are in two independent namespace.
func (d *Driver) SetCounter(key string, v int64, ttl time.Duration) error {
	var bytes = make([]byte, 8)
	binary.BigEndian.PutUint64(bytes, uint64(v))
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return err
	}
	defer ctx.unlocker()
	status := ctx.data.set(NewData(key, now.Add(ttl).Unix(), bytes), now.Unix())
	return d.save(ctx, status)
}

//GetCounter Get int val from cache by given key.Count cache and data cache are in two independent namespace.
func (d *Driver) GetCounter(key string) (int64, error) {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return 0, err
	}
	data := ctx.data.get(key, now.Unix())
	if data == nil {
		return 0, cache.ErrNotFound
	}
	v := binary.BigEndian.Uint64(data.Data)
	return int64(v), nil

}

//DelCounter Delete int val in cache by given key.Count cache and data cache are in two independent namespace.
func (d *Driver) DelCounter(key string) error {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return err
	}
	defer ctx.unlocker()
	status := ctx.data.delete(key, now.Unix())
	return d.save(ctx, status)
}

//SetGCErrHandler Set callback to handler error raised when gc.
func (d *Driver) SetGCErrHandler(f func(err error)) {
	d.GcErrHanlder = f
}

//Expire set item ttl by given key and ttl.
func (d *Driver) Expire(key string, ttl time.Duration) error {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return err
	}
	defer ctx.unlocker()
	status := ctx.data.expired(key, now.Add(ttl).Unix(), now.Unix())
	return d.save(ctx, status)
}

//ExpireCounter set int item ttl by given key and ttl.
func (d *Driver) ExpireCounter(key string, ttl time.Duration) error {
	now := time.Now()
	ctx, err := d.lockAndGetData(key)
	if err != nil {
		return err
	}
	defer ctx.unlocker()
	status := ctx.data.expired(key, now.Add(ttl).Unix(), now.Unix())
	return d.save(ctx, status)
}

//MGetBytesValue get multiple bytes data from cache by given keys.
//Return data bytes map and any error if raised.
func (d *Driver) MGetBytesValue(keys ...string) (map[string][]byte, error) {
	var result = map[string][]byte{}
	for k := range keys {
		bs, err := d.GetBytesValue(keys[k])
		if err != nil {
			if err == cache.ErrNotFound {
				continue
			}
			return nil, err
		}
		result[keys[k]] = bs
	}
	return result, nil
}

//MSetBytesValue set multiple bytes data to cache with given key-value map.
//Return  any error if raised.
func (d *Driver) MSetBytesValue(data map[string][]byte, ttl time.Duration) error {
	for key := range data {
		err := d.SetBytesValue(key, data[key], ttl)
		if err != nil {
			return err
		}
	}
	return nil
}

//Close close cache.
func (d *Driver) Close() error {
	return d.Store.Close()
}

//Flush Delete all data in cache.
func (d *Driver) Flush() error {
	return d.Store.Flush()
}
