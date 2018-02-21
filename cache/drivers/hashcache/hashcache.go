package hashcache

import (
	"errors"
	"strconv"
	"time"

	"github.com/herb-go/herb/cache"
)

const hashTypeKey = byte(1)
const hashTypeValue = byte(2)
const hashMinLength = 256

var ErrHashFormatWrong = errors.New("error hash format wrong")

//Cache The redis cache Driver.
type Cache struct {
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

//Get Get data model from cache by given key.
//Parameter v should be pointer to empty data model which data filled in.
//Return any error raised.
func (c *Cache) Get(key string, v interface{}) error {
	bytes, err := c.GetBytesValue(key)
	if err != nil {
		return err
	}
	return cache.UnmarshalMsgpack(bytes, v)
}

//Set Set data model to cache by given key.
//Return any error raised.
func (c *Cache) Set(key string, v interface{}, ttl time.Duration) error {
	bytes, err := cache.MarshalMsgpack(v)
	if err != nil {
		return err
	}
	return c.SetBytesValue(key, bytes, ttl)
}

//Update Update data model to cache by given key only if the cache exist.
//Return any error raised.
func (c *Cache) Update(key string, v interface{}, ttl time.Duration) error {
	bytes, err := cache.MarshalMsgpack(v)
	if err != nil {
		return err
	}
	return c.UpdateBytesValue(key, bytes, ttl)
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

func (c *Cache) getHash(key string) (byte, []byte, error) {
	b, err := c.Remote.GetBytesValue(key + cache.KeyPrefix)
	if err != nil {
		return 0, nil, err
	}
	if len(b) < 2 || !(b[0] == hashTypeKey || b[0] == hashTypeValue) {
		return 0, nil, ErrHashFormatWrong
	}
	return b[0], b[1:], nil
}

//SetBytesValue Set bytes data to cache by given key.
//Return any error raised.
func (c *Cache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	var err error
	if len(bytes) < hashMinLength {
		b := make([]byte, len(bytes)+1)
		b[0] = hashTypeValue
		copy(b[1:], bytes)
		return c.Remote.SetBytesValue(key+cache.KeyPrefix, b, ttl)
	}
	ts := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	b := make([]byte, len(ts)+1)
	b[0] = hashTypeKey
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
	if len(bytes) < hashMinLength {
		b := make([]byte, len(bytes)+1)
		b[0] = hashTypeValue
		copy(b[1:], bytes)
		return c.Remote.UpdateBytesValue(key+cache.KeyPrefix, b, ttl)
	}
	ts := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	b := make([]byte, len(ts)+1)
	b[0] = hashTypeKey
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
func (c *Cache) mGetHashs(data *map[string][]byte, keys ...string) (hashedKeys map[string]string, hashedValueKeys []string, err error) {
	hashedValueKeys = make([]string, len(keys))
	hashedKeys = make(map[string]string, len(keys))
	var prefixLength = len(cache.KeyPrefix)

	var hashedValueKeysLength = 0
	hashs, err := c.Remote.MGetBytesValue(c.getPrefixedKeys(keys...)...)
	if err != nil {
		return
	}
	for k := range hashs {
		if hashs[k] == nil {
			continue
		}
		if hashs[k][0] == hashTypeValue {
			(*data)[k[:len(k)-prefixLength]] = hashs[k][1:]
		} else if hashs[k][0] == hashTypeKey {
			var key = k + string(hashs[k][1:])
			hashedValueKeys[hashedValueKeysLength] = key
			hashedKeys[key] = k[:len(k)-prefixLength]
			hashedValueKeysLength++
		} else {
			err = ErrHashFormatWrong
			return
		}
	}
	hashedValueKeys = hashedValueKeys[:hashedValueKeysLength]
	return
}
func (c *Cache) mGetLocalData(data *map[string][]byte, hashedKeys map[string]string, hashedValueKeys []string) (RemoteValueKeys []string, err error) {
	var LocalData map[string][]byte
	RemoteValueKeys = make([]string, len(hashedValueKeys))
	var RemoteValueKeysLength = 0
	if len(hashedValueKeys) > 0 {
		LocalData, err = c.Local.MGetBytesValue(hashedValueKeys...)
		if err != nil {
			return
		}
		for k := range LocalData {
			if LocalData[k] != nil {
				key, ok := hashedKeys[k]
				if ok == true {
					(*data)[key] = LocalData[k]
				}
			} else {
				RemoteValueKeys[RemoteValueKeysLength] = k
				RemoteValueKeysLength++
			}
		}
		RemoteValueKeys = RemoteValueKeys[:RemoteValueKeysLength]
	}
	return
}
func (c *Cache) mGetRemoteData(data *map[string][]byte, hashedKeys map[string]string, RemoteValueKeys []string) (err error) {
	var RemoteData map[string][]byte
	if len(RemoteValueKeys) > 0 {
		RemoteData, err = c.Remote.MGetBytesValue(RemoteValueKeys...)
		if err != nil {
			return
		}
		for k := range RemoteData {
			key, ok := hashedKeys[k]
			if ok == true {
				(*data)[key] = RemoteData[k]
			}
		}
		err = c.Local.MSetBytesValue(RemoteData, 0)
	}
	return
}
func (c *Cache) MGetBytesValue(keys ...string) (map[string][]byte, error) {
	var err error
	var unfinishedKeys []string
	var hashedKeys map[string]string
	var data = make(map[string][]byte, len(keys))
	hashedKeys, unfinishedKeys, err = c.mGetHashs(&data, keys...)
	if err != nil {
		return nil, err
	}
	unfinishedKeys, err = c.mGetLocalData(&data, hashedKeys, unfinishedKeys)
	if err != nil {
		return nil, err
	}
	err = c.mGetRemoteData(&data, hashedKeys, unfinishedKeys)
	if err != nil {
		return nil, err
	}
	return data, nil
}
func (c *Cache) MSetBytesValue(data map[string][]byte, ttl time.Duration) error {
	var hashs = make(map[string][]byte, len(data))
	var RemoteData = make(map[string][]byte, len(data))
	ts := []byte(strconv.FormatInt(time.Now().UnixNano(), 10))
	for k := range data {
		if len(data[k]) < hashMinLength {
			b := make([]byte, len(data[k])+1)
			b[0] = hashTypeValue
			copy(b[1:], data[k])
			hashs[k+cache.KeyPrefix] = b
		} else {
			b := make([]byte, len(ts)+1)
			b[0] = hashTypeKey
			copy(b[1:], ts)
			hashs[k+cache.KeyPrefix] = b
			RemoteData[k+cache.KeyPrefix+string(ts)] = data[k]
		}
	}
	err := c.Remote.MSetBytesValue(hashs, ttl)
	if err != nil {
		return err
	}
	return c.Remote.MSetBytesValue(RemoteData, ttl)
}

//Del Delete data in cache by given key.
//Return any error raised.
func (c *Cache) Del(key string) error {
	var finalErr error
	t, b, err := c.getHash(key)
	if err != nil {
		return err
	}
	err = c.Remote.Del(key + cache.KeyPrefix)
	if err != nil {
		return err
	}
	if t == hashTypeValue {
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
	t, b, err := c.getHash(key)
	if err != nil {
		return b, err
	}
	if t == hashTypeValue {
		return b, nil
	}
	hashKey := key + cache.KeyPrefix + string(b)
	b, err = c.Local.GetBytesValue(hashKey)
	if err == cache.ErrNotFound {
		b, err = c.Remote.GetBytesValue(hashKey)
		if err != nil {
			return b, err
		}
		err = c.Local.SetBytesValue(hashKey, b, 0)
	}
	return b, err
}

func (c *Cache) Expire(key string, ttl time.Duration) error {
	var err error
	err = c.Remote.Expire(key+cache.KeyPrefix, ttl)
	if err != nil {
		return err
	}
	t, b, err := c.getHash(key)
	if err == cache.ErrNotFound {
		return nil
	}
	if err != nil {
		return err
	}
	if t == hashTypeValue {
		return nil
	}
	hashKey := key + cache.KeyPrefix + string(b)
	return c.Remote.Expire(hashKey, ttl)
}

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
	Local  cache.Config
	Remote cache.Config
}

//New Create new cache driver with given json bytes.
//Return new driver and any error raised.
func (config *Config) Create() (cache.Driver, error) {
	cc := Cache{}
	localcache := cache.New()
	err := localcache.Init(config.Local)
	if err != nil {
		return &cc, err
	}
	cc.Local = localcache
	remotecache := cache.New()
	err = remotecache.Init(config.Remote)
	if err != nil {
		return &cc, err
	}
	cc.Remote = remotecache
	return &cc, nil
}
func init() {
	cache.Register("hashcache", func() cache.DriverConfig {
		return &Config{}
	})
}
