//Package rediscache provides cache driver uses redis to store cache data.
//Using github.com/garyburd/redigo/redis as driver.
package rediscache

import (
	"sync"
	"time"

	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/herb-go/herb/cache"
)

var defaultMaxIdle = 200
var defaultMaxAlive = 200
var defaultIdleTimeout = 60 * time.Second
var defualtConnectTimeout = 10 * time.Second
var defualtReadTimeout = 2 * time.Second
var defualtWriteTimeout = 2 * time.Second
var defaultSepartor = string(0)

const modeSet = 0
const modeUpdate = 1

//Cache The redis cache Driver.
type Cache struct {
	Pool           *redis.Pool //Redis pool.
	ticker         *time.Ticker
	name           string
	quit           chan int
	gcErrHandler   func(err error)
	gcLimit        int64
	network        string
	address        string
	password       string
	version        string
	versionLock    sync.Mutex
	db             int
	connectTimeout time.Duration
	readTimeout    time.Duration
	writeTimeout   time.Duration
	Separtor       string //Separtor in redis key.
}

func (c *Cache) dial() (redis.Conn, error) {
	conn, err := redis.DialTimeout(c.network, c.address, c.connectTimeout, c.readTimeout, c.writeTimeout)
	if err != nil {
		return nil, err
	}
	if c.password != "" {
		_, err = conn.Do("auth", c.password)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}
	_, err = conn.Do("SELECT", c.db)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}
func (c *Cache) start() error {
	conn := c.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("PING")
	return err
}
func (c *Cache) getKey(key string) string {
	return c.name + c.Separtor + key
}

//Flush Flush not supported.
func (c *Cache) Flush() error {
	return cache.ErrFeatureNotSupported
}

//Close Close cache.
//Return any error if raised
func (c *Cache) Close() error {
	return c.Pool.Close()
}

//Del Delete data in cache by given key.
//Return any error raised.
func (c *Cache) Del(key string) error {
	k := c.getKey(key)
	conn := c.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", k)
	return err
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
	val := strconv.FormatInt(v, 10)
	return c.SetBytesValue(key, []byte(val), ttl)
}

//GetCounter Get int val from cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) GetCounter(key string) (int64, error) {
	var v int64
	bytes, err := c.GetBytesValue(key)
	if err != nil {
		return v, err
	}
	return strconv.ParseInt(string(bytes), 10, 64)
}

//DelCounter Delete int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return any error raised.
func (c *Cache) DelCounter(key string) error {
	k := c.getKey(key)
	conn := c.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", k)
	return err
}

//IncrCounter Increase int val in cache by given key.Count cache and data cache are in two independent namespace.
//Return int data value and any error raised.
func (c *Cache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	var err error
	var v int64
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)

	v, err = redis.Int64(conn.Do("INCRBY", k, increment))
	if err != nil {
		return v, err
	}
	if ttl < 0 {
		_, err = conn.Do("PERSIST", k)
	} else {
		_, err = conn.Do("EXPIRE", k, int64(ttl/time.Second))
	}
	if err != nil {
		return v, err
	}

	return v, err
}
func (c *Cache) doSet(key string, bytes []byte, ttl time.Duration, mode int) error {
	var err error
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)
	if ttl < 0 {
		if mode == modeUpdate {
			_, err = conn.Do("SET", k, bytes, "XX")
		} else {
			_, err = conn.Do("SET", k, bytes)
		}
	} else {
		if mode == modeUpdate {
			_, err = conn.Do("SET", k, bytes, "EX", int64(ttl/time.Second), "XX")

		} else {
			_, err = conn.Do("SET", k, bytes, "EX", int64(ttl/time.Second))

		}
	}
	return err
}

//SetBytesValue Set bytes data to cache by given key.
//Return any error raised.
func (c *Cache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	return c.doSet(key, bytes, ttl, modeSet)
}

//UpdateBytesValue Update bytes data to cache by given key only if the cache exist.
//Return any error raised.
func (c *Cache) UpdateBytesValue(key string, bytes []byte, ttl time.Duration) error {
	return c.doSet(key, bytes, ttl, modeUpdate)
}

func (c *Cache) MGetBytesValue(keys ...string) (map[string][]byte, error) {
	var data = make(map[string][]byte, len(keys))
	var err error
	conn := c.Pool.Get()
	defer conn.Close()
	for key := range keys {
		k := c.getKey(keys[key])
		err := (conn.Send("GET", k))
		if err != nil {
			return nil, err
		}
	}

	err = conn.Flush()
	if err != nil {
		return nil, err
	}
	for key := range keys {
		bs, err := redis.Bytes((conn.Receive()))
		if err == redis.ErrNil {
			data[keys[key]] = nil
			continue
		}
		if err != nil {
			return nil, err
		}
		data[keys[key]] = bs
	}

	return data, nil
}
func (c *Cache) MSetBytesValue(data map[string][]byte, ttl time.Duration) error {
	var err error
	conn := c.Pool.Get()
	defer conn.Close()
	var ttlInSecond = int64(ttl / time.Second)
	for key := range data {
		k := c.getKey(key)
		if ttl < 0 {
			err = conn.Send("SET", k, data[key])
		} else {
			err = conn.Send("SET", k, data[key], "EX", ttlInSecond)

		}
		if err != nil {
			return err
		}
	}
	return conn.Flush()

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

//GetBytesValue Get bytes data from cache by given key.
//Return data bytes and any error raised.
func (c *Cache) GetBytesValue(key string) ([]byte, error) {
	var bs []byte
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)
	bs, err := redis.Bytes((conn.Do("GET", k)))
	if err == redis.ErrNil {
		return nil, cache.ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	if bs == nil {
		return nil, cache.ErrNotFound
	}
	return bs, err
}

func (c *Cache) Expire(key string, ttl time.Duration) error {
	var err error
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)
	if ttl < 0 {
		_, err = conn.Do("PERSIST", k)
	} else {
		_, err = conn.Do("EXPIRE", k, int64(ttl/time.Second))
	}
	return err
}

func (c *Cache) ExpireCounter(key string, ttl time.Duration) error {
	var err error
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)
	if ttl < 0 {
		_, err = conn.Do("PERSIST", k)
	} else {
		_, err = conn.Do("EXPIRE", k, int64(ttl/time.Second))
	}
	return err
}

//SetGCErrHandler Set callback to handler error raised when gc.
func (c *Cache) SetGCErrHandler(f func(err error)) {
	return
}

//Config Cache driver config.
type Config struct {
	Network     string //Network string of redis conn.
	Address     string //Redis server address.
	Name        string ////Redis server username.
	Password    string //Redis server password.
	Db          int    //Redis server database id.
	MaxIdle     int    //Max idle conn in redis pool.
	MaxAlive    int    //Max Alive conn in redis pool.
	IdleTimeout int    //Idel comm time.
	GCPeriod    int64  //Period of gc.Default value is 30 second.
	GCLimit     int64  //Max delete limit in every gc call.Default value is 100.
}

func (c *Config) Create() (cache.Driver, error) {
	cache := Cache{}
	cache.name = c.Name
	cache.network = c.Network
	cache.address = c.Address
	cache.password = c.Password
	cache.db = c.Db
	cache.connectTimeout = defualtConnectTimeout
	cache.readTimeout = defualtReadTimeout
	cache.writeTimeout = defualtWriteTimeout
	cache.Separtor = defaultSepartor
	maxIdle := c.MaxIdle
	if maxIdle == 0 {
		maxIdle = defaultMaxIdle
	}
	cache.Pool = redis.NewPool(cache.dial, maxIdle)
	cache.Pool.MaxActive = c.MaxAlive
	if cache.Pool.MaxActive == 0 {
		cache.Pool.MaxActive = defaultMaxAlive
	}
	cache.Pool.IdleTimeout = time.Duration(c.IdleTimeout) * time.Second
	if cache.Pool.IdleTimeout == 0 {
		cache.Pool.IdleTimeout = defaultIdleTimeout
	}
	cache.Pool.Wait = true
	cache.quit = make(chan int)
	err := cache.start()
	if err != nil {
		return &cache, err
	}
	return &cache, nil
}
func init() {
	cache.Register("rediscache", func() cache.DriverConfig {
		return &Config{}
	})
}
