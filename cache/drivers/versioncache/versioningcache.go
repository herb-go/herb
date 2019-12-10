package versioncache

import (
	"errors"
	"strconv"
	"time"

	"github.com/herb-go/herb/cache"
)

const VersionTypeKey = byte(1)
const VersionTypeValue = byte(2)
const VersionMinLength = 256

//ErrVersionFormatWrong raised when version format wrong
var ErrVersionFormatWrong = errors.New("error version format wrong")

//Cache The redis cache Driver.
type Cache struct {
	cache.DriverUtil
	Local  *cache.Cache
	Remote *cache.Cache
}

//Flush Flush not supported.
func (c *Cache) Flush() error {
	var finalErr error
	var err error
	err = c.Remote.Flush()
	if err != nil {
		finalErr = err
	}
	err = c.Local.Flush()
	if err != nil {
		finalErr = err
	}
	return finalErr
}

//Close Close cache.
//Return any error if raised
func (c *Cache) Close() error {
	var finalErr error
	var err error
	err = c.Remote.Flush()
	if err != nil {
		finalErr = err
	}
	err = c.Local.Flush()
	if err != nil {
		finalErr = err
	}
	return finalErr
}

//SetCounter Set int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) SetCounter(key string, v int64, ttl time.Duration) error {
	return c.Remote.SetCounter(key, v, ttl)
}

//GetCounter Get int val from cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) GetCounter(key string) (int64, error) {
	return c.Remote.GetCounter(key)
}

//DelCounter Delete int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) DelCounter(key string) error {
	return c.Remote.DelCounter(key)
}

//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	return c.Remote.IncrCounter(key, increment, ttl)
}

func (c *Cache) getVersion(key string) (byte, []byte, error) {
	b, err := c.Remote.GetBytesValue(key + cache.KeyPrefix)
	if err != nil {
		return 0, nil, err
	}
	if len(b) < 2 || !(b[0] == VersionTypeKey || b[0] == VersionTypeValue) {
		return 0, nil, ErrVersionFormatWrong
	}
	return b[0], b[1:], nil
}

//SetBytesValue Set bytes data to cache by given key.
//Return any error raised.
func (c *Cache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	var err error
	if len(bytes) < VersionMinLength {
		b := make([]byte, len(bytes)+1)
		b[0] = VersionTypeValue
		copy(b[1:], bytes)
		return c.Remote.SetBytesValue(key+cache.KeyPrefix, b, ttl)
	}
	ts := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	b := make([]byte, len(ts)+1)
	b[0] = VersionTypeKey
	copy(b[1:], ts)
	err = c.Remote.SetBytesValue(key+cache.KeyPrefix, b, ttl)
	if err != nil {
		return err
	}
	var k = key + cache.KeyPrefix + string(ts)
	return c.Remote.SetBytesValue(k, bytes, ttl)
}

//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
//Return any error raised.
func (c *Cache) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {
	var err error
	if len(bytes) < VersionMinLength {
		b := make([]byte, len(bytes)+1)
		b[0] = VersionTypeValue
		copy(b[1:], bytes)
		return c.Remote.UpdateBytesValue(key+cache.KeyPrefix, b, ttl)
	}
	ts := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	b := make([]byte, len(ts)+1)
	b[0] = VersionTypeKey
	copy(b[1:], ts)
	err = c.Remote.UpdateBytesValue(key+cache.KeyPrefix, b, ttl)

	if err != nil {
		return err
	}
	var k = key + cache.KeyPrefix + string(ts)
	err = c.Remote.SetBytesValue(k, bytes, ttl)
	if err != nil {
		return err
	}
	return c.Local.SetBytesValue(k, bytes, ttl)
}

func (c *Cache) getPrefixedKeys(keys ...string) []string {
	var prefixedKeys = make([]string, len(keys))
	for k := range keys {
		prefixedKeys[k] = keys[k] + cache.KeyPrefix
	}
	return prefixedKeys
}
func (c *Cache) mGetVersions(data *map[string][]byte, keys ...string) (values map[string]string, versions []string, err error) {
	versions = make([]string, len(keys))
	values = make(map[string]string, len(keys))
	var prefixLength = len(cache.KeyPrefix)

	var versionsLength = 0
	bytesvalues, err := c.Remote.MGetBytesValue(c.getPrefixedKeys(keys...)...)
	if err != nil {
		return
	}
	for k := range bytesvalues {
		if bytesvalues[k] == nil {
			continue
		}
		if bytesvalues[k][0] == VersionTypeValue {
			(*data)[k[:len(k)-prefixLength]] = bytesvalues[k][1:]
		} else if bytesvalues[k][0] == VersionTypeKey {
			var key = k + string(bytesvalues[k][1:])
			versions[versionsLength] = key
			values[key] = k[:len(k)-prefixLength]
			versionsLength++
		} else {
			err = ErrVersionFormatWrong
			return
		}
	}
	versions = versions[:versionsLength]
	return
}
func (c *Cache) mGetLocalData(data *map[string][]byte, values map[string]string, versions []string) (RemoteVersions []string, err error) {
	var LocalData map[string][]byte
	RemoteVersions = make([]string, len(versions))
	var RemoteVersionsLength = 0
	if len(versions) > 0 {
		LocalData, err = c.Local.MGetBytesValue(versions...)
		if err != nil {
			return
		}
		for k := range versions {
			if LocalData[versions[k]] != nil {
				key, ok := values[versions[k]]
				if ok == true {
					(*data)[key] = LocalData[versions[k]]
				}
			} else {
				RemoteVersions[RemoteVersionsLength] = versions[k]
				RemoteVersionsLength++
			}
		}
		RemoteVersions = RemoteVersions[:RemoteVersionsLength]
	}
	return
}
func (c *Cache) mGetRemoteData(data *map[string][]byte, values map[string]string, versions []string) (err error) {
	var RemoteData map[string][]byte
	if len(versions) > 0 {
		RemoteData, err = c.Remote.MGetBytesValue(versions...)
		if err != nil {
			return
		}
		for k := range RemoteData {
			key, ok := values[k]
			if ok == true {
				(*data)[key] = RemoteData[k]
			}
		}
		err = c.Local.MSetBytesValue(RemoteData, 0)
	}
	return
}

//MGetBytesValue get multiple bytes data from cache by given keys.
//Return data bytes map and any error if raised.
func (c *Cache) MGetBytesValue(keys ...string) (map[string][]byte, error) {
	var err error
	var versions []string
	var values map[string]string
	var data = make(map[string][]byte, len(keys))
	values, versions, err = c.mGetVersions(&data, keys...)
	if err != nil {
		return nil, err
	}
	versions, err = c.mGetLocalData(&data, values, versions)
	if err != nil {
		return nil, err
	}
	err = c.mGetRemoteData(&data, values, versions)
	if err != nil {
		return nil, err
	}
	return data, nil
}

//MSetBytesValue set multiple bytes data to cache with given key-value map.
//Return  any error if raised.
func (c *Cache) MSetBytesValue(data map[string][]byte, ttl time.Duration) error {
	var versions = make(map[string][]byte, len(data))
	var RemoteData = make(map[string][]byte, len(data))
	ts := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	for k := range data {
		if len(data[k]) < VersionMinLength {
			b := make([]byte, len(data[k])+1)
			b[0] = VersionTypeValue
			copy(b[1:], data[k])
			versions[k+cache.KeyPrefix] = b
		} else {
			b := make([]byte, len(ts)+1)
			b[0] = VersionTypeKey
			copy(b[1:], ts)
			versions[k+cache.KeyPrefix] = b
			RemoteData[k+cache.KeyPrefix+string(ts)] = data[k]
		}
	}
	err := c.Remote.MSetBytesValue(versions, ttl)
	if err != nil {
		return err
	}
	return c.Remote.MSetBytesValue(RemoteData, ttl)
}

//Del Delete data in cache by given key.
//Return any error raised.
func (c *Cache) Del(key string) error {
	var finalErr error
	t, b, err := c.getVersion(key)
	if err != nil {
		return err
	}
	err = c.Remote.Del(key + cache.KeyPrefix)
	if err != nil {
		return err
	}
	if t == VersionTypeValue {
		return nil
	}
	k := string(b)
	finalErr = c.Remote.Del(k)
	err = c.Local.Del(k)
	if err != nil {
		finalErr = err
	}
	return finalErr
}

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (c *Cache) GetBytesValue(key string) ([]byte, error) {
	t, b, err := c.getVersion(key)
	if err != nil {
		return b, err
	}
	if t == VersionTypeValue {
		return b, nil
	}
	versionKey := key + cache.KeyPrefix + string(b)
	b, err = c.Local.GetBytesValue(versionKey)
	if err == cache.ErrNotFound {
		b, err = c.Remote.GetBytesValue(versionKey)
		if err != nil {
			return b, err
		}
		err = c.Local.SetBytesValue(versionKey, b, 0)
	}
	return b, err
}

//Expire set cache value expire duration by given key and ttl
func (c *Cache) Expire(key string, ttl time.Duration) error {
	var err error
	err = c.Remote.Expire(key+cache.KeyPrefix, ttl)
	if err != nil {
		return err
	}
	t, b, err := c.getVersion(key)
	if err == cache.ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	if t == VersionTypeValue {
		return nil
	}
	versionKey := key + cache.KeyPrefix + string(b)
	return c.Remote.Expire(versionKey, ttl)
}

//ExpireCounter set cache counter  expire duration by given key and ttl
func (c *Cache) ExpireCounter(key string, ttl time.Duration) error {
	return c.Remote.ExpireCounter(key, ttl)
}

//SetGCErrHandler Set callback to handler error raised when gc.
func (c *Cache) SetGCErrHandler(f func(err error)) {
	c.Local.SetGCErrHandler(f)
	c.Remote.SetGCErrHandler(f)
	return
}

//Config Cache driver config.
type Config struct {
	Local  cache.OptionConfig
	Remote cache.OptionConfig
}

func init() {
	cache.Register("versioncache", func(loader func(interface{}) error) (cache.Driver, error) {
		var err error
		cc := &Cache{}
		config := &Config{}
		err = loader(config)
		if err != nil {
			return nil, err
		}
		cc.Local, err = cache.NewSubCache(&config.Local)
		if err != nil {
			return nil, err
		}
		cc.Remote, err = cache.NewSubCache(&config.Remote)
		if err != nil {
			return nil, err
		}
		return cc, nil
	})
}
