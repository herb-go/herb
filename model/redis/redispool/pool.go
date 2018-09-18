package redispool

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var defaultMaxIdle = 200
var defaultMaxAlive = 200
var defaultIdleTimeout = 60 * time.Second
var defualtConnectTimeout = 10 * time.Second
var defualtReadTimeout = 2 * time.Second
var defualtWriteTimeout = 2 * time.Second

func New() *Pool {
	return &Pool{}
}

type Pool struct {
	*redis.Pool
	Network        string //Network string of redis conn.
	Address        string //Redis server address.
	Name           string ////Redis server username.
	Password       string //Redis server password.
	ConnectTimeout time.Duration
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	Db             int           //Redis server database id.
	MaxIdle        int           //Max idle conn in redis pool.
	MaxAlive       int           //Max Alive conn in redis pool.
	IdleTimeout    time.Duration //Idel conn time.
}

func (p *Pool) dial() (redis.Conn, error) {
	conn, err := redis.DialTimeout(p.Network, p.Address, p.ConnectTimeout, p.ReadTimeout, p.WriteTimeout)
	if err != nil {
		return nil, err
	}
	if p.Password != "" {
		_, err = conn.Do("auth", p.Password)
		if err != nil {
			conn.Close()
			return nil, err
		}
	}
	_, err = conn.Do("SELECT", p.Db)
	if err != nil {
		conn.Close()
		return nil, err
	}
	return conn, nil
}

func (p *Pool) Open() *redis.Pool {
	maxIdle := p.MaxIdle
	if maxIdle == 0 {
		maxIdle = defaultMaxIdle
	}
	p.Pool = redis.NewPool(p.dial, maxIdle)
	p.Pool.MaxActive = p.MaxAlive
	if p.Pool.MaxActive == 0 {
		p.Pool.MaxActive = defaultMaxAlive
	}
	p.Pool.IdleTimeout = p.IdleTimeout
	if p.Pool.IdleTimeout == 0 {
		p.Pool.IdleTimeout = defaultIdleTimeout
	}
	p.Pool.Wait = true
	return p.Pool
}

type Config struct {
	Network                string //Network string of redis conn.
	Address                string //Redis server address.
	Password               string //Redis server password.
	ConnectTimeoutInSecond int64
	ReadTimeoutInSecond    int64
	WriteTimeoutInSecond   int64
	Db                     int //Redis server database id.
	MaxIdle                int //Max idle conn in redis pool.
	MaxAlive               int //Max Alive conn in redis pool.
	IdleTimeoutInSecond    int64
}

func (c *Config) ApplyTo(p *Pool) error {
	p.Network = c.Network
	p.Address = c.Address
	p.Password = c.Password
	p.ConnectTimeout = time.Duration(c.ConnectTimeoutInSecond) * time.Second
	p.ReadTimeout = time.Duration(c.ReadTimeoutInSecond) * time.Second
	p.WriteTimeout = time.Duration(c.WriteTimeoutInSecond) * time.Second
	p.Db = c.Db
	p.MaxIdle = c.MaxIdle
	p.MaxAlive = c.MaxAlive
	p.IdleTimeout = time.Duration(c.IdleTimeoutInSecond) * time.Second
	return nil
}
