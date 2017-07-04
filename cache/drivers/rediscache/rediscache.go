package rediscache

import (
	"encoding/json"
	"sync"
	"time"

	"strconv"

	"github.com/garyburd/redigo/redis"
	"github.com/herb-go/herb/cache"
)

var DefaultGCPeriod = 30 * time.Second
var defaultGcLimit = int64(100)
var defaultMaxIdle = 200
var defaultMaxAlive = 200
var defaultIdleTimeout = 60 * time.Second
var defualtConnectTimeout = 10 * time.Second
var defualtReadTimeout = 2 * time.Second
var defualtWriteTimeout = 2 * time.Second
var defaultSepartor = string(0)
var tokenMask = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
var tokenLength = 64
var flushLua = `
	if redis.call("HEXISTS",KEYS[2],KEYS[3])==1 then return 0 end
	local v=redis.call("GET",KEYS[1]);
	if (v==false) then v="" end;
    redis.call("HSET",KEYS[2],v,0);
	redis.call("SET",KEYS[1],KEYS[3]);
	return 1;
`
var gcLua = `
	redis.replicate_commands()
	local ks=redis.call("HKEYS",KEYS[1])
	if ks ==false then return end
	local k=ks[1]
	if k ==nil then return end
	local v=redis.call("HGET",KEYS[1],k)
	local r=redis.call("SCAN",v,"MATCH",KEYS[2]..KEYS[3]..KEYS[3]..k..KEYS[3].."*","COUNT",KEYS[4])
	for _,k in ipairs(r[2]) do 
    	redis.call('DEL', k) 
	end
	if r[1]=="0" then redis.call("HDEL",KEYS[1],k) end
`

type Cache struct {
	Pool           *redis.Pool
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
	Separtor       string
}
type Config struct {
	Network     string
	Address     string
	Password    string
	Name        string
	Db          int
	MaxIdle     int
	MaxAlive    int
	IdleTimeout int
	GCPeriod    int64
	GCLimit     int64
}

func (c *Cache) Dial() (redis.Conn, error) {
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
	if err != nil {
		return err
	}
	return c.gc()
}
func (c *Cache) getKey(key string) string {
	c.versionLock.Lock()
	defer c.versionLock.Unlock()
	return c.name + c.Separtor + c.Separtor + c.version + c.Separtor + key
}
func (c *Cache) SearchByPrefix(prefix string) ([]string, error) {
	return nil, cache.ErrSearchKeysNotSupported
}
func (c *Cache) getVersionKey() string {
	return c.name + c.Separtor + "version" + c.Separtor
}
func (c *Cache) getUsedVersionsKey() string {
	return c.name + c.Separtor + "usedVersions" + c.Separtor

}
func (c *Cache) getVersionFromConn(conn redis.Conn) (string, error) {
	var version string
	vk := c.getVersionKey()
	version, err := redis.String(conn.Do("GET", vk))
	if err == redis.ErrNil {
		version = ""
	} else {
		return version, err
	}
	return version, nil
}
func (c *Cache) Flush() error {
	conn := c.Pool.Get()
	defer conn.Close()
	vk := c.getVersionKey()
	version, err := c.getVersionFromConn(conn)
	nv, err := cache.NewRandMaskedBytes(tokenMask, tokenLength, []byte(version))
	if err != nil {
		return err
	}
	vsk := c.getUsedVersionsKey()
	result, err := redis.Int64(conn.Do("EVAL", flushLua, 3, vk, vsk, string(nv)))
	if err != nil {
		return err
	}
	if result == 0 {
		return c.Flush()
	}
	return nil
}
func (c *Cache) gc() error {
	var err error
	conn := c.Pool.Get()
	defer conn.Close()
	vsk := c.getUsedVersionsKey()
	_, err = conn.Do("EVAL", gcLua, 4, vsk, c.name, c.Separtor, c.gcLimit)
	return err
}
func (c *Cache) Close() error {
	return c.Pool.Close()
}
func (c *Cache) Del(key string) error {
	k := c.getKey(key)
	conn := c.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", k)
	return err
}
func (c *Cache) Set(key string, v interface{}, ttl time.Duration) error {
	bytes, err := cache.MarshalMsgpack(&v)
	if err != nil {
		return err
	}
	return c.SetBytesValue(key, bytes, ttl)
}
func (c *Cache) setVersion(newVersion string) {
	c.versionLock.Lock()
	c.version = newVersion
	c.versionLock.Unlock()

}
func (c *Cache) SetCounter(key string, v int64, ttl time.Duration) error {
	val := strconv.FormatInt(v, 10)
	return c.SetBytesValue(key, []byte(val), ttl)
}
func (c *Cache) GetCounter(key string) (int64, error) {
	var v int64
	bytes, err := c.GetBytesValue(key)
	if err != nil {
		return v, err
	}
	return strconv.ParseInt(string(bytes), 10, 64)
}
func (c *Cache) DelCounter(key string) error {
	k := c.getKey(key)
	conn := c.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("DEL", k)
	return err
}
func (c *Cache) IncrCounter(key string, increment int64, ttl time.Duration) (int64, error) {
	var err error
	var version string
	var v int64
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)
	_, err = conn.Do("MULTI")
	if err != nil {
		return v, err
	}
	vk := c.getVersionKey()
	_, err = conn.Do("GET", vk)
	if err != nil {
		return 0, err
	}
	_, err = conn.Do("INCR", k)
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
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return v, err
	}
	values, err = redis.Scan(values, &version)
	if err == redis.ErrNil {
		version = ""
	} else if err != nil {
		return 0, err
	}
	if version != c.version {
		c.version = version
		_, err = conn.Do("DEL", k)
		if err != nil {
			return 0, err
		}
		return c.IncrCounter(key, increment, ttl)
	}
	_, err = redis.Scan(values, &v)
	return v, err
}
func (c *Cache) SetBytesValue(key string, bytes []byte, ttl time.Duration) error {
	var err error
	var version string
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)
	_, err = conn.Do("MULTI")
	if err != nil {
		return err
	}
	vk := c.getVersionKey()
	_, err = conn.Do("GET", vk)
	if err != nil {
		return err
	}
	if ttl < 0 {
		_, err = conn.Do("SET", k, bytes)
	} else {
		_, err = conn.Do("SETEX", k, int64(ttl/time.Second), bytes)
	}
	if err != nil {
		return err
	}
	values, err := redis.Values(conn.Do("EXEC"))
	if err != nil {
		return err
	}
	_, err = redis.Scan(values, &version)
	if err == redis.ErrNil {
		version = ""
	} else if err != nil {
		return err
	}
	if version != c.version {
		c.version = version
		_, err = conn.Do("DEL", k)
		if err != nil {
			return err
		}
		return c.SetBytesValue(key, bytes, ttl)
	}
	return nil
}
func (c *Cache) Get(key string, v interface{}) error {
	bytes, err := c.GetBytesValue(key)
	if err != nil {
		return err
	}
	return cache.UnmarshalMsgpack(bytes, &v)
}
func (c *Cache) GetBytesValue(key string) ([]byte, error) {
	var bytes []byte
	var version string
	conn := c.Pool.Get()
	defer conn.Close()
	k := c.getKey(key)
	v := c.getVersionKey()
	values, err := redis.Values((conn.Do("MGET", v, k)))
	b, err := redis.Scan(values, &version)
	if err == redis.ErrNil {
		version = ""
	} else {
		if err != nil {
			return bytes, err
		}
	}
	c.versionLock.Lock()
	if version != c.version {
		c.version = version
		c.versionLock.Unlock()
		return c.GetBytesValue(key)
	}
	c.versionLock.Unlock()
	_, err = redis.Scan(b, &bytes)
	if err == redis.ErrNil || bytes == nil {
		return bytes, cache.ErrNotFound
	} else if err != nil {
		return bytes, nil
	}
	return bytes, err
}
func (c *Cache) SetGCErrHandler(f func(err error)) {
	c.gcErrHandler = f
	return
}
func (_ *Cache) New(config json.RawMessage) (cache.Driver, error) {
	c := Config{}
	err := json.Unmarshal(config, &c)
	if err != nil {
		return nil, err
	}
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
	cache.Pool = redis.NewPool(cache.Dial, maxIdle)
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
	period := time.Duration(c.GCPeriod)
	if period == 0 {
		period = DefaultGCPeriod
	}
	cache.ticker = time.NewTicker(period)
	gcLimit := c.GCLimit
	if gcLimit == 0 {
		gcLimit = defaultGcLimit
	}
	cache.gcLimit = gcLimit
	go func() {
		for {
			select {
			case <-cache.ticker.C:
				err := cache.gc()
				if err != nil {
					if cache.gcErrHandler != nil {
						cache.gcErrHandler(err)
					}
				}
			case <-cache.quit:
				cache.ticker.Stop()
				return
			}
		}

	}()
	err = cache.start()
	if err != nil {
		return &cache, err
	}
	return &cache, nil
}
func init() {
	cache.Register("rediscache", &Cache{})
}
